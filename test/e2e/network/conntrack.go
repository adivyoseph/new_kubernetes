/*
Copyright 2020 The Kubernetes Authors.

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

package network

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	e2enode "k8s.io/kubernetes/test/e2e/framework/node"
	e2epod "k8s.io/kubernetes/test/e2e/framework/pod"
	e2eservice "k8s.io/kubernetes/test/e2e/framework/service"
	e2eskipper "k8s.io/kubernetes/test/e2e/framework/skipper"
	"k8s.io/kubernetes/test/e2e/network/common"
	imageutils "k8s.io/kubernetes/test/utils/image"
)

const (
	serviceName = "svc-udp"
	podClient   = "pod-client"
	podBackend1 = "pod-server-1"
	podBackend2 = "pod-server-2"
	srcPort     = 12345
)

// Linux NAT uses conntrack to perform NAT, everytime a new
// flow is seen, a connection is created in the conntrack table, and it
// is being used by the NAT module.
// Each entry in the conntrack table has associated a timeout, that removes
// the connection once it expires.
// UDP is a connectionless protocol, so the conntrack module tracking functions
// are not very advanced.
// It uses a short timeout (30 sec by default) that is renewed if there are new flows
// matching the connection. Otherwise it expires the entry.
// This behaviour can cause issues in Kubernetes when one entry on the conntrack table
// is never expired because the sender does not stop sending traffic, but the pods or
// endpoints were deleted, blackholing the traffic
// In order to mitigate this problem, Kubernetes delete the stale entries:
// - when an endpoint is removed
// - when a service goes from no endpoints to new endpoint

// Ref: https://api.semanticscholar.org/CorpusID:198903401
// Boye, Magnus. "Netfilter Connection Tracking and NAT Implementation." (2012).

var _ = common.SIGDescribe("Conntrack", func() {

	fr := framework.NewDefaultFramework("conntrack")

	type nodeInfo struct {
		name   string
		nodeIP string
	}

	var (
		cs                             clientset.Interface
		ns                             string
		clientNodeInfo, serverNodeInfo nodeInfo
	)

	logContainsFn := func(text, podName string) wait.ConditionFunc {
		return func() (bool, error) {
			logs, err := e2epod.GetPodLogs(cs, ns, podName, podName)
			if err != nil {
				// Retry the error next time.
				return false, nil
			}
			if !strings.Contains(string(logs), text) {
				return false, nil
			}
			return true, nil
		}
	}

	ginkgo.BeforeEach(func() {
		cs = fr.ClientSet
		ns = fr.Namespace.Name

		nodes, err := e2enode.GetBoundedReadySchedulableNodes(cs, 2)
		framework.ExpectNoError(err)
		if len(nodes.Items) < 2 {
			e2eskipper.Skipf(
				"Test requires >= 2 Ready nodes, but there are only %v nodes",
				len(nodes.Items))
		}

		var family v1.IPFamily
		if framework.TestContext.ClusterIsIPv6() {
			family = v1.IPv6Protocol
		} else {
			family = v1.IPv4Protocol
		}

		ips := e2enode.GetAddressesByTypeAndFamily(&nodes.Items[0], v1.NodeInternalIP, family)

		clientNodeInfo = nodeInfo{
			name:   nodes.Items[0].Name,
			nodeIP: ips[0],
		}

		ips = e2enode.GetAddressesByTypeAndFamily(&nodes.Items[1], v1.NodeInternalIP, family)

		serverNodeInfo = nodeInfo{
			name:   nodes.Items[1].Name,
			nodeIP: ips[0],
		}
	})

	ginkgo.It("should be able to preserve UDP traffic when server pod cycles for a NodePort service", func() {

		// Create a NodePort service
		udpJig := e2eservice.NewTestJig(cs, ns, serviceName)
		ginkgo.By("creating a UDP service " + serviceName + " with type=NodePort in " + ns)
		udpService, err := udpJig.CreateUDPService(func(svc *v1.Service) {
			svc.Spec.Type = v1.ServiceTypeNodePort
			svc.Spec.Ports = []v1.ServicePort{
				{Port: 80, Name: "udp", Protocol: v1.ProtocolUDP, TargetPort: intstr.FromInt(80)},
			}
		})
		framework.ExpectNoError(err)

		// Create a pod in one node to create the UDP traffic against the NodePort service every 5 seconds
		ginkgo.By("creating a client pod for probing the service " + serviceName)
		clientPod := e2epod.NewAgnhostPod(ns, podClient, nil, nil, nil)
		nodeSelection := e2epod.NodeSelection{Name: clientNodeInfo.name}
		e2epod.SetNodeSelection(&clientPod.Spec, nodeSelection)
		cmd := fmt.Sprintf(`date; for i in $(seq 1 3000); do echo "$(date) Try: ${i}"; echo hostname | nc -u -w 5 -p %d %s %d; echo; done`, srcPort, serverNodeInfo.nodeIP, udpService.Spec.Ports[0].NodePort)
		clientPod.Spec.Containers[0].Command = []string{"/bin/sh", "-c", cmd}
		clientPod.Spec.Containers[0].Name = podClient
		fr.PodClient().CreateSync(clientPod)

		// Read the client pod logs
		logs, err := e2epod.GetPodLogs(cs, ns, podClient, podClient)
		framework.ExpectNoError(err)
		framework.Logf("Pod client logs: %s", logs)

		// Add a backend pod to the service in the other node
		ginkgo.By("creating a backend pod " + podBackend1 + " for the service " + serviceName)
		serverPod1 := e2epod.NewAgnhostPod(ns, podBackend1, nil, nil, nil, "netexec", fmt.Sprintf("--udp-port=%d", 80))
		serverPod1.Labels = udpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod1.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod1)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend1: {80}})

		// Note that the fact that Endpoints object already exists, does NOT mean
		// that iptables (or whatever else is used) was already programmed.
		// Additionally take into account that UDP conntract entries timeout is
		// 30 seconds by default.
		// Based on the above check if the pod receives the traffic.
		ginkgo.By("checking client pod connected to the backend 1 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend1, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 1")
		}

		// Create a second pod
		ginkgo.By("creating a second backend pod " + podBackend2 + " for the service " + serviceName)
		serverPod2 := e2epod.NewAgnhostPod(ns, podBackend2, nil, nil, nil, "netexec", fmt.Sprintf("--udp-port=%d", 80))
		serverPod2.Labels = udpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod2.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod2)

		// and delete the first pod
		framework.Logf("Cleaning up %s pod", podBackend1)
		fr.PodClient().DeleteSync(podBackend1, metav1.DeleteOptions{}, framework.DefaultPodDeletionTimeout)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend2: {80}})

		// Check that the second pod keeps receiving traffic
		// UDP conntrack entries timeout is 30 sec by default
		ginkgo.By("checking client pod connected to the backend 2 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend2, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 2")
		}
	})

	ginkgo.It("should be able to preserve UDP traffic when server pod cycles for a ClusterIP service", func() {

		// Create a ClusterIP service
		udpJig := e2eservice.NewTestJig(cs, ns, serviceName)
		ginkgo.By("creating a UDP service " + serviceName + " with type=ClusterIP in " + ns)
		udpService, err := udpJig.CreateUDPService(func(svc *v1.Service) {
			svc.Spec.Type = v1.ServiceTypeClusterIP
			svc.Spec.Ports = []v1.ServicePort{
				{Port: 80, Name: "udp", Protocol: v1.ProtocolUDP, TargetPort: intstr.FromInt(80)},
			}
		})
		framework.ExpectNoError(err)

		// Create a pod in one node to create the UDP traffic against the ClusterIP service every 5 seconds
		ginkgo.By("creating a client pod for probing the service " + serviceName)
		clientPod := e2epod.NewAgnhostPod(ns, podClient, nil, nil, nil)
		nodeSelection := e2epod.NodeSelection{Name: clientNodeInfo.name}
		e2epod.SetNodeSelection(&clientPod.Spec, nodeSelection)
		cmd := fmt.Sprintf(`date; for i in $(seq 1 3000); do echo "$(date) Try: ${i}"; echo hostname | nc -u -w 5 -p %d %s %d; echo; done`, srcPort, udpService.Spec.ClusterIP, udpService.Spec.Ports[0].Port)
		clientPod.Spec.Containers[0].Command = []string{"/bin/sh", "-c", cmd}
		clientPod.Spec.Containers[0].Name = podClient
		fr.PodClient().CreateSync(clientPod)

		// Read the client pod logs
		logs, err := e2epod.GetPodLogs(cs, ns, podClient, podClient)
		framework.ExpectNoError(err)
		framework.Logf("Pod client logs: %s", logs)

		// Add a backend pod to the service in the other node
		ginkgo.By("creating a backend pod " + podBackend1 + " for the service " + serviceName)
		serverPod1 := e2epod.NewAgnhostPod(ns, podBackend1, nil, nil, nil, "netexec", fmt.Sprintf("--udp-port=%d", 80))
		serverPod1.Labels = udpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod1.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod1)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend1: {80}})

		// Note that the fact that Endpoints object already exists, does NOT mean
		// that iptables (or whatever else is used) was already programmed.
		// Additionally take into account that UDP conntract entries timeout is
		// 30 seconds by default.
		// Based on the above check if the pod receives the traffic.
		ginkgo.By("checking client pod connected to the backend 1 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend1, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 1")
		}

		// Create a second pod
		ginkgo.By("creating a second backend pod " + podBackend2 + " for the service " + serviceName)
		serverPod2 := e2epod.NewAgnhostPod(ns, podBackend2, nil, nil, nil, "netexec", fmt.Sprintf("--udp-port=%d", 80))
		serverPod2.Labels = udpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod2.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod2)

		// and delete the first pod
		framework.Logf("Cleaning up %s pod", podBackend1)
		fr.PodClient().DeleteSync(podBackend1, metav1.DeleteOptions{}, framework.DefaultPodDeletionTimeout)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend2: {80}})

		// Check that the second pod keeps receiving traffic
		// UDP conntrack entries timeout is 30 sec by default
		ginkgo.By("checking client pod connected to the backend 2 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend2, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 2")
		}
	})

	ginkgo.It("should be able to preserve TCP traffic when server pod cycles for a NodePort service", func() {

		// Create a NodePort service
		tcpJig := e2eservice.NewTestJig(cs, ns, serviceName)
		ginkgo.By("creating a TCP service " + serviceName + " with type=NodePort in " + ns)
		tcpService, err := tcpJig.CreateTCPService(func(svc *v1.Service) {
			svc.Spec.Type = v1.ServiceTypeNodePort
			svc.Spec.Ports = []v1.ServicePort{
				{Port: 80, Name: "tcp", Protocol: v1.ProtocolTCP, TargetPort: intstr.FromInt(80)},
			}
		})
		framework.ExpectNoError(err)

		// Create a pod in one node to create the UDP traffic against the NodePort service every 5 seconds
		ginkgo.By("creating a client pod for probing the service " + serviceName)
		clientPod := e2epod.NewAgnhostPod(ns, podClient, nil, nil, nil)
		nodeSelection := e2epod.NodeSelection{Name: clientNodeInfo.name}
		e2epod.SetNodeSelection(&clientPod.Spec, nodeSelection)
		cmd := fmt.Sprintf(`date; for i in $(seq 1 3000); do echo "$(date) Try: ${i}"; echo hostname | nc -w 5 -p %d %s %d; echo; done`, srcPort, serverNodeInfo.nodeIP, tcpService.Spec.Ports[0].NodePort)
		clientPod.Spec.Containers[0].Command = []string{"/bin/sh", "-c", cmd}
		clientPod.Spec.Containers[0].Name = podClient
		fr.PodClient().CreateSync(clientPod)

		// Read the client pod logs
		logs, err := e2epod.GetPodLogs(cs, ns, podClient, podClient)
		framework.ExpectNoError(err)
		framework.Logf("Pod client logs: %s", logs)

		// Add a backend pod to the service in the other node
		ginkgo.By("creating a backend pod " + podBackend1 + " for the service " + serviceName)
		serverPod1 := e2epod.NewAgnhostPod(ns, podBackend1, nil, nil, nil, "netexec", fmt.Sprintf("--tcp-port=%d", 80))
		serverPod1.Labels = tcpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod1.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod1)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend1: {80}})

		// Note that the fact that Endpoints object already exists, does NOT mean
		// that iptables (or whatever else is used) was already programmed.
		// Additionally take into account that TCP SYN conntrakt entries timeout is
		// net.netfilter.nf_conntrack_tcp_timeout_syn_sent = 120
		// Based on the above check if the pod receives the traffic.
		ginkgo.By("checking client pod connected to the backend 1 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend1, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 1")
		}

		// Create a second pod
		ginkgo.By("creating a second backend pod " + podBackend2 + " for the service " + serviceName)
		serverPod2 := e2epod.NewAgnhostPod(ns, podBackend2, nil, nil, nil, "netexec", fmt.Sprintf("--tcp-port=%d", 80))
		serverPod2.Labels = tcpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod2.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod2)

		// and delete the first pod
		framework.Logf("Cleaning up %s pod", podBackend1)
		fr.PodClient().DeleteSync(podBackend1, metav1.DeleteOptions{}, framework.DefaultPodDeletionTimeout)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend2: {80}})

		// Check that the second pod keeps receiving traffic
		// UDP conntrack entries timeout is 30 sec by default
		ginkgo.By("checking client pod connected to the backend 2 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend2, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 2")
		}
	})

	ginkgo.It("should be able to preserve TCP traffic when server pod cycles for a ClusterIP service", func() {

		// Create a ClusterIP service
		tcpJig := e2eservice.NewTestJig(cs, ns, serviceName)
		ginkgo.By("creating a TCP service " + serviceName + " with type=ClusterIP in " + ns)
		tcpService, err := tcpJig.CreateTCPService(func(svc *v1.Service) {
			svc.Spec.Type = v1.ServiceTypeClusterIP
			svc.Spec.Ports = []v1.ServicePort{
				{Port: 80, Name: "tcp", Protocol: v1.ProtocolTCP, TargetPort: intstr.FromInt(80)},
			}
		})
		framework.ExpectNoError(err)

		// Create a pod in one node to create the UDP traffic against the ClusterIP service every 5 seconds
		ginkgo.By("creating a client pod for probing the service " + serviceName)
		clientPod := e2epod.NewAgnhostPod(ns, podClient, nil, nil, nil)
		nodeSelection := e2epod.NodeSelection{Name: clientNodeInfo.name}
		e2epod.SetNodeSelection(&clientPod.Spec, nodeSelection)
		cmd := fmt.Sprintf(`date; for i in $(seq 1 3000); do echo "$(date) Try: ${i}"; echo hostname | nc -w 5 -p %d %s %d; echo; done`, srcPort, tcpService.Spec.ClusterIP, tcpService.Spec.Ports[0].Port)
		clientPod.Spec.Containers[0].Command = []string{"/bin/sh", "-c", cmd}
		clientPod.Spec.Containers[0].Name = podClient
		fr.PodClient().CreateSync(clientPod)

		// Read the client pod logs
		logs, err := e2epod.GetPodLogs(cs, ns, podClient, podClient)
		framework.ExpectNoError(err)
		framework.Logf("Pod client logs: %s", logs)

		// Add a backend pod to the service in the other node
		ginkgo.By("creating a backend pod " + podBackend1 + " for the service " + serviceName)
		serverPod1 := e2epod.NewAgnhostPod(ns, podBackend1, nil, nil, nil, "netexec", fmt.Sprintf("--tcp-port=%d", 80))
		serverPod1.Labels = tcpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod1.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod1)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend1: {80}})

		// Note that the fact that Endpoints object already exists, does NOT mean
		// that iptables (or whatever else is used) was already programmed.
		// Additionally take into account that TCP SYN conntrakt entries timeout is
		// net.netfilter.nf_conntrack_tcp_timeout_syn_sent = 120
		// Based on the above check if the pod receives the traffic.
		ginkgo.By("checking client pod connected to the backend 1 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend1, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 1")
		}

		// Create a second pod
		ginkgo.By("creating a second backend pod " + podBackend2 + " for the service " + serviceName)
		serverPod2 := e2epod.NewAgnhostPod(ns, podBackend2, nil, nil, nil, "netexec", fmt.Sprintf("--tcp-port=%d", 80))
		serverPod2.Labels = tcpJig.Labels
		nodeSelection = e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod2.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod2)

		// and delete the first pod
		framework.Logf("Cleaning up %s pod", podBackend1)
		fr.PodClient().DeleteSync(podBackend1, metav1.DeleteOptions{}, framework.DefaultPodDeletionTimeout)

		validateEndpointsPortsOrFail(cs, ns, serviceName, portsByPodName{podBackend2: {80}})

		// Check that the second pod keeps receiving traffic
		// UDP conntrack entries timeout is 30 sec by default
		ginkgo.By("checking client pod connected to the backend 2 on Node IP " + serverNodeInfo.nodeIP)
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn(podBackend2, podClient)); err != nil {
			logs, err = e2epod.GetPodLogs(cs, ns, podClient, podClient)
			framework.ExpectNoError(err)
			framework.Logf("Pod client logs: %s", logs)
			framework.Failf("Failed to connect to backend 2")
		}
	})

	// Regression test for #74839, where:
	// Packets considered INVALID by conntrack are now dropped. In particular, this fixes
	// a problem where spurious retransmits in a long-running TCP connection to a service
	// IP could result in the connection being closed with the error "Connection reset by
	// peer"
	// xref: https://kubernetes.io/blog/2019/03/29/kube-proxy-subtleties-debugging-an-intermittent-connection-reset/
	ginkgo.It("should drop INVALID conntrack entries", func() {
		serverLabel := map[string]string{
			"app": "boom-server",
		}

		serverPod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "boom-server",
				Labels: serverLabel,
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "boom-server",
						Image: imageutils.GetE2EImage(imageutils.RegressionIssue74839),
						Ports: []v1.ContainerPort{
							{
								ContainerPort: 9000, // Default port exposed by boom-server
							},
						},
						Env: []v1.EnvVar{
							{
								Name: "POD_IP",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "status.podIP",
									},
								},
							},
							{
								Name: "POD_IPS",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										APIVersion: "v1",
										FieldPath:  "status.podIPs",
									},
								},
							},
						},
					},
				},
			},
		}
		nodeSelection := e2epod.NodeSelection{Name: serverNodeInfo.name}
		e2epod.SetNodeSelection(&serverPod.Spec, nodeSelection)
		fr.PodClient().CreateSync(serverPod)
		ginkgo.By("Server pod created on node " + serverNodeInfo.name)

		svc := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "boom-server",
			},
			Spec: v1.ServiceSpec{
				Selector: serverLabel,
				Ports: []v1.ServicePort{
					{
						Protocol: v1.ProtocolTCP,
						Port:     9000,
					},
				},
			},
		}
		_, err := fr.ClientSet.CoreV1().Services(fr.Namespace.Name).Create(context.TODO(), svc, metav1.CreateOptions{})
		framework.ExpectNoError(err)

		ginkgo.By("Server service created")

		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "startup-script",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "startup-script",
						Image: imageutils.GetE2EImage(imageutils.BusyBox),
						Command: []string{
							"sh", "-c", "while true; do sleep 2; nc boom-server 9000& done",
						},
					},
				},
				RestartPolicy: v1.RestartPolicyNever,
			},
		}
		nodeSelection = e2epod.NodeSelection{Name: clientNodeInfo.name}
		e2epod.SetNodeSelection(&pod.Spec, nodeSelection)

		fr.PodClient().CreateSync(pod)
		ginkgo.By("Client pod created")

		// The client will open connections against the server
		// The server will inject invalid packets
		// if conntrack does not drop the invalid packets it will go through without NAT
		// so the client will receive an unexpected TCP connection and RST the connection
		// the server will log ERROR if that happens
		ginkgo.By("checking client pod does not RST the TCP connection because it receives and INVALID packet")
		if err := wait.PollImmediate(5*time.Second, time.Minute, logContainsFn("ERROR", "boom-server")); err == nil {
			logs, err := e2epod.GetPodLogs(cs, ns, "boom-server", "boom-server")
			framework.ExpectNoError(err)
			framework.Logf("boom-server pod logs: %s", logs)
			framework.Failf("Boom server pod received a RST from the client")
		}

		logs, err := e2epod.GetPodLogs(cs, ns, "boom-server", "boom-server")
		framework.ExpectNoError(err)
		if !strings.Contains(string(logs), "connection established") {
			framework.Logf("boom-server pod logs: %s", logs)
			framework.Failf("Boom server pod did not sent any bad packet to the client")
		}
		framework.Logf("boom-server pod logs: %s", logs)
		framework.Logf("boom-server OK: did not receive any RST packet")
	})
})
