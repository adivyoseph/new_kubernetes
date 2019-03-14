/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package expand

import (
	"fmt"
	"net"
	"time"

	"k8s.io/klog"

	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	kcache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/kubernetes/pkg/controller"
	"k8s.io/kubernetes/pkg/controller/volume/events"
	"k8s.io/kubernetes/pkg/util/mount"
	"k8s.io/kubernetes/pkg/volume"
	"k8s.io/kubernetes/pkg/volume/util"
	"k8s.io/kubernetes/pkg/volume/util/operationexecutor"
	"k8s.io/kubernetes/pkg/volume/util/subpath"
	"k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
)

// ExpandController expands the pvs
type ExpandController interface {
	Run(stopCh <-chan struct{})
}

type expandController struct {
	// kubeClient is the kube API client used by volumehost to communicate with
	// the API server.
	kubeClient clientset.Interface

	// pvcLister is the shared PVC lister used to fetch and store PVC
	// objects from the API server. It is shared with other controllers and
	// therefore the PVC objects in its store should be treated as immutable.
	pvcLister  corelisters.PersistentVolumeClaimLister
	pvcsSynced kcache.InformerSynced

	pvLister corelisters.PersistentVolumeLister
	pvSynced kcache.InformerSynced

	// cloud provider used by volume host
	cloud cloudprovider.Interface

	// volumePluginMgr used to initialize and fetch volume plugins
	volumePluginMgr volume.VolumePluginMgr

	// recorder is used to record events in the API server
	recorder record.EventRecorder

	operationGenerator operationexecutor.OperationGenerator

	queue workqueue.RateLimitingInterface
}

func NewExpandController(
	kubeClient clientset.Interface,
	pvcInformer coreinformers.PersistentVolumeClaimInformer,
	pvInformer coreinformers.PersistentVolumeInformer,
	cloud cloudprovider.Interface,
	plugins []volume.VolumePlugin) (ExpandController, error) {

	expc := &expandController{
		kubeClient: kubeClient,
		cloud:      cloud,
		pvcLister:  pvcInformer.Lister(),
		pvcsSynced: pvcInformer.Informer().HasSynced,
		pvLister:   pvInformer.Lister(),
		pvSynced:   pvInformer.Informer().HasSynced,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "expand"),
	}

	if err := expc.volumePluginMgr.InitPlugins(plugins, nil, expc); err != nil {
		return nil, fmt.Errorf("Could not initialize volume plugins for Expand Controller : %+v", err)
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	expc.recorder = eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "volume_expand"})
	blkutil := volumepathhandler.NewBlockVolumePathHandler()

	expc.operationGenerator = operationexecutor.NewOperationGenerator(
		kubeClient,
		&expc.volumePluginMgr,
		expc.recorder,
		false,
		blkutil)

	pvcInformer.Informer().AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc: expc.enqueuePVC,
		UpdateFunc: func(old, new interface{}) {
			oldPVC := old.(*v1.PersistentVolumeClaim)
			oldSize := oldPVC.Spec.Resources.Requests[v1.ResourceStorage]

			newPVC := new.(*v1.PersistentVolumeClaim)
			newSize := newPVC.Spec.Resources.Requests[v1.ResourceStorage]
			if newSize.Cmp(oldSize) > 0 {
				expc.enqueuePVC(new)
			}
		},
		DeleteFunc: expc.enqueuePVC,
	})

	return expc, nil
}

func (expc *expandController) enqueuePVC(obj interface{}) {
	pvc := obj.(*v1.PersistentVolumeClaim)
	size := pvc.Spec.Resources.Requests[v1.ResourceStorage]
	statusSize := pvc.Status.Capacity[v1.ResourceStorage]

	if pvc.Status.Phase == v1.ClaimBound && size.Cmp(statusSize) > 0 {
		key, err := kcache.DeletionHandlingMetaNamespaceKeyFunc(pvc)
		if err != nil {
			runtime.HandleError(fmt.Errorf("Couldn't get key for object %#v: %v", pvc, err))
			return
		}
		expc.queue.Add(key)
	}
}

func (expc *expandController) processNextWorkItem() bool {
	key, shutdown := expc.queue.Get()
	if shutdown {
		return false
	}
	defer expc.queue.Done(key)

	err := expc.syncHandler(key.(string))
	if err == nil {
		expc.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("%v failed with : %v", key, err))
	expc.queue.AddRateLimited(key)

	return true
}

func (expc *expandController) syncHandler(key string) error {
	namespace, name, err := kcache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	pvc, err := expc.pvcLister.PersistentVolumeClaims(namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		klog.V(5).Infof("Error getting PVC %q (uid: %q) from informer : %v", util.GetPersistentVolumeClaimQualifiedName(pvc), pvc.UID, err)
		return err
	}

	pv, err := getPersistentVolume(pvc, expc.pvLister)
	if err != nil {
		klog.V(5).Infof("Error getting Persistent Volume for PVC %q (uid: %q) from informer : %v", util.GetPersistentVolumeClaimQualifiedName(pvc), pvc.UID, err)
		return err
	}
	if pv.Spec.ClaimRef == nil || pvc.Namespace != pv.Spec.ClaimRef.Namespace || pvc.Name != pv.Spec.ClaimRef.Name {
		err := fmt.Errorf("Persistent Volume is not bound to PVC being updated : %s", util.ClaimToClaimKey(pvc))
		klog.V(4).Infof("%v", err)
		return err
	}

	volumeSpec := volume.NewSpecFromPersistentVolume(pv, false)
	volumePlugin, err := expc.volumePluginMgr.FindExpandablePluginBySpec(volumeSpec)
	if err != nil || volumePlugin == nil {
		msg := fmt.Errorf("didn't find a plugin capable of expanding the volume; " +
			"waiting for an external controller to process this PVC")
		expc.recorder.Event(pvc, v1.EventTypeNormal, events.ExternalExpanding,
			fmt.Sprintf("Ignoring the PVC: %v.", msg))
		klog.V(3).Infof("Ignoring the PVC %q (uid: %q) : %v.",
			util.GetPersistentVolumeClaimQualifiedName(pvc), pvc.UID, msg)
		return err
	}

	return expc.expand(pvc, pv)
}

func (expc *expandController) expand(pvc *v1.PersistentVolumeClaim, pv *v1.PersistentVolume) error {
	pvc, err := util.MarkResizeInProgress(pvc, expc.kubeClient)
	if err != nil {
		klog.V(5).Infof("Error setting PVC %s in progress with error : %v", util.GetPersistentVolumeClaimQualifiedName(pvc), err)
		return err
	}

	generatedOperations, err := expc.operationGenerator.GenerateExpandVolumeFunc(pvc, pv)
	if err != nil {
		klog.Errorf("Error starting ExpandVolume for pvc %s with %v", util.GetPersistentVolumeClaimQualifiedName(pvc), err)
		return err
	}
	klog.V(5).Infof("Starting ExpandVolume for volume %s", util.GetPersistentVolumeClaimQualifiedName(pvc))
	_, detailedErr := generatedOperations.Run()

	return detailedErr
}

// TODO make concurrency configurable (workers/threadiness argument). previously, nestedpendingoperations spawned unlimited goroutines
func (expc *expandController) Run(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()
	defer expc.queue.ShutDown()

	klog.Infof("Starting expand controller")
	defer klog.Infof("Shutting down expand controller")

	if !controller.WaitForCacheSync("expand", stopCh, expc.pvcsSynced, expc.pvSynced) {
		return
	}

	go wait.Until(expc.runWorker, time.Second, stopCh)

	<-stopCh
}

func (expc *expandController) runWorker() {
	for expc.processNextWorkItem() {
	}
}

func getPersistentVolume(pvc *v1.PersistentVolumeClaim, pvLister corelisters.PersistentVolumeLister) (*v1.PersistentVolume, error) {
	volumeName := pvc.Spec.VolumeName
	pv, err := pvLister.Get(volumeName)

	if err != nil {
		return nil, fmt.Errorf("failed to find PV %q in PV informer cache with error : %v", volumeName, err)
	}

	return pv.DeepCopy(), nil
}

// Implementing VolumeHost interface
func (expc *expandController) GetPluginDir(pluginName string) string {
	return ""
}

func (expc *expandController) GetVolumeDevicePluginDir(pluginName string) string {
	return ""
}

func (expc *expandController) GetPodsDir() string {
	return ""
}

func (expc *expandController) GetPodVolumeDir(podUID types.UID, pluginName string, volumeName string) string {
	return ""
}

func (expc *expandController) GetPodVolumeDeviceDir(podUID types.UID, pluginName string) string {
	return ""
}

func (expc *expandController) GetPodPluginDir(podUID types.UID, pluginName string) string {
	return ""
}

func (expc *expandController) GetKubeClient() clientset.Interface {
	return expc.kubeClient
}

func (expc *expandController) NewWrapperMounter(volName string, spec volume.Spec, pod *v1.Pod, opts volume.VolumeOptions) (volume.Mounter, error) {
	return nil, fmt.Errorf("NewWrapperMounter not supported by expand controller's VolumeHost implementation")
}

func (expc *expandController) NewWrapperUnmounter(volName string, spec volume.Spec, podUID types.UID) (volume.Unmounter, error) {
	return nil, fmt.Errorf("NewWrapperUnmounter not supported by expand controller's VolumeHost implementation")
}

func (expc *expandController) GetCloudProvider() cloudprovider.Interface {
	return expc.cloud
}

func (expc *expandController) GetMounter(pluginName string) mount.Interface {
	return nil
}

func (expc *expandController) GetExec(pluginName string) mount.Exec {
	return mount.NewOsExec()
}

func (expc *expandController) GetHostName() string {
	return ""
}

func (expc *expandController) GetHostIP() (net.IP, error) {
	return nil, fmt.Errorf("GetHostIP not supported by expand controller's VolumeHost implementation")
}

func (expc *expandController) GetNodeAllocatable() (v1.ResourceList, error) {
	return v1.ResourceList{}, nil
}

func (expc *expandController) GetSecretFunc() func(namespace, name string) (*v1.Secret, error) {
	return func(_, _ string) (*v1.Secret, error) {
		return nil, fmt.Errorf("GetSecret unsupported in expandController")
	}
}

func (expc *expandController) GetConfigMapFunc() func(namespace, name string) (*v1.ConfigMap, error) {
	return func(_, _ string) (*v1.ConfigMap, error) {
		return nil, fmt.Errorf("GetConfigMap unsupported in expandController")
	}
}

func (expc *expandController) GetServiceAccountTokenFunc() func(_, _ string, _ *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
	return func(_, _ string, _ *authenticationv1.TokenRequest) (*authenticationv1.TokenRequest, error) {
		return nil, fmt.Errorf("GetServiceAccountToken unsupported in expandController")
	}
}

func (expc *expandController) DeleteServiceAccountTokenFunc() func(types.UID) {
	return func(types.UID) {
		klog.Errorf("DeleteServiceAccountToken unsupported in expandController")
	}
}

func (expc *expandController) GetNodeLabels() (map[string]string, error) {
	return nil, fmt.Errorf("GetNodeLabels unsupported in expandController")
}

func (expc *expandController) GetNodeName() types.NodeName {
	return ""
}

func (expc *expandController) GetEventRecorder() record.EventRecorder {
	return expc.recorder
}

func (expc *expandController) GetSubpather() subpath.Interface {
	// not needed for expand controller
	return nil
}
