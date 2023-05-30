/*
Copyright 2016 The Kubernetes Authors.

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

package images

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
)

type pullResult struct {
	imageRef     string
	err          error
	pullDuration time.Duration
}

type imagePuller interface {
	pullImage(context.Context, kubecontainer.ImageSpec, []v1.Secret, chan<- pullResult, *runtimeapi.PodSandboxConfig)
	pullImageWithProgress(context.Context, kubecontainer.ImageSpec, []v1.Secret, chan<- pullResult, *runtimeapi.PodSandboxConfig, *v1.ObjectReference)
}

var _, _ imagePuller = &parallelImagePuller{}, &serialImagePuller{}

type parallelImagePuller struct {
	imageService kubecontainer.ImageService
	tokens       chan struct{}
}

func newParallelImagePuller(imageService kubecontainer.ImageService, maxParallelImagePulls *int32) imagePuller {
	if maxParallelImagePulls == nil || *maxParallelImagePulls < 1 {
		return &parallelImagePuller{imageService, nil}
	}
	return &parallelImagePuller{imageService, make(chan struct{}, *maxParallelImagePulls)}
}

func (pip *parallelImagePuller) pullImage(ctx context.Context, spec kubecontainer.ImageSpec, pullSecrets []v1.Secret, pullChan chan<- pullResult, podSandboxConfig *runtimeapi.PodSandboxConfig) {
	go func() {
		if pip.tokens != nil {
			pip.tokens <- struct{}{}
			defer func() { <-pip.tokens }()
		}
		startTime := time.Now()
		imageRef, err := pip.imageService.PullImage(ctx, spec, pullSecrets, podSandboxConfig)
		pullChan <- pullResult{
			imageRef:     imageRef,
			err:          err,
			pullDuration: time.Since(startTime),
		}
	}()
}

func (pip *parallelImagePuller) pullImageWithProgress(ctx context.Context, spec kubecontainer.ImageSpec, pullSecrets []v1.Secret, pullChan chan<- pullResult, podSandboxConfig *runtimeapi.PodSandboxConfig, ref *v1.ObjectReference) {
	go func() {
		if pip.tokens != nil {
			pip.tokens <- struct{}{}
			defer func() { <-pip.tokens }()
		}
		startTime := time.Now()
		imageRef, err := pip.imageService.PullImageWithProgress(ctx, spec, pullSecrets, podSandboxConfig, ref)
		pullChan <- pullResult{
			imageRef:     imageRef,
			err:          err,
			pullDuration: time.Since(startTime),
		}
	}()
}

// Maximum number of image pull requests than can be queued.
const maxImagePullRequests = 10

type serialImagePuller struct {
	imageService kubecontainer.ImageService
	pullRequests chan *imagePullRequest
}

func newSerialImagePuller(imageService kubecontainer.ImageService) imagePuller {
	imagePuller := &serialImagePuller{imageService, make(chan *imagePullRequest, maxImagePullRequests)}
	go wait.Until(imagePuller.processImagePullRequests, time.Second, wait.NeverStop)
	return imagePuller
}

type imagePullRequest struct {
	ctx              context.Context
	spec             kubecontainer.ImageSpec
	pullSecrets      []v1.Secret
	pullChan         chan<- pullResult
	podSandboxConfig *runtimeapi.PodSandboxConfig
	objRef           *v1.ObjectReference
}

func (sip *serialImagePuller) pullImage(ctx context.Context, spec kubecontainer.ImageSpec, pullSecrets []v1.Secret, pullChan chan<- pullResult, podSandboxConfig *runtimeapi.PodSandboxConfig) {
	sip.pullRequests <- &imagePullRequest{
		ctx:              ctx,
		spec:             spec,
		pullSecrets:      pullSecrets,
		pullChan:         pullChan,
		podSandboxConfig: podSandboxConfig,
		objRef:           nil,
	}
}

func (sip *serialImagePuller) pullImageWithProgress(ctx context.Context, spec kubecontainer.ImageSpec, pullSecrets []v1.Secret, pullChan chan<- pullResult, podSandboxConfig *runtimeapi.PodSandboxConfig, ref *v1.ObjectReference) {
	sip.pullRequests <- &imagePullRequest{
		ctx:              ctx,
		spec:             spec,
		pullSecrets:      pullSecrets,
		pullChan:         pullChan,
		podSandboxConfig: podSandboxConfig,
		objRef:           ref,
	}
}

func (sip *serialImagePuller) processImagePullRequests() {
	for pullRequest := range sip.pullRequests {
		startTime := time.Now()
		var (
			imageRef string
			err      error
		)
		if pullRequest.objRef == nil {
			imageRef, err = sip.imageService.PullImage(pullRequest.ctx, pullRequest.spec, pullRequest.pullSecrets, pullRequest.podSandboxConfig)
		} else {
			imageRef, err = sip.imageService.PullImageWithProgress(pullRequest.ctx, pullRequest.spec, pullRequest.pullSecrets, pullRequest.podSandboxConfig, pullRequest.objRef)
		}
		pullRequest.pullChan <- pullResult{
			imageRef:     imageRef,
			err:          err,
			pullDuration: time.Since(startTime),
		}
	}
}
