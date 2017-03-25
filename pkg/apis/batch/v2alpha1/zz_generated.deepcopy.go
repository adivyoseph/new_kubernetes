// +build !ignore_autogenerated

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v2alpha1

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	v1 "k8s.io/kubernetes/pkg/api/v1"
)

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *CronJob) DeepCopyInto(out *CronJob) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new CronJob.
func (x *CronJob) DeepCopy() *CronJob {
	if x == nil {
		return nil
	}
	out := new(CronJob)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyObject will perform a deep copy of the receiver, creating a new object.
func (x *CronJob) DeepCopyObject() runtime.Object {
	if c := x.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *CronJobList) DeepCopyInto(out *CronJobList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CronJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new CronJobList.
func (x *CronJobList) DeepCopy() *CronJobList {
	if x == nil {
		return nil
	}
	out := new(CronJobList)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyObject will perform a deep copy of the receiver, creating a new object.
func (x *CronJobList) DeepCopyObject() runtime.Object {
	if c := x.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *CronJobSpec) DeepCopyInto(out *CronJobSpec) {
	*out = *in
	if in.StartingDeadlineSeconds != nil {
		in, out := &in.StartingDeadlineSeconds, &out.StartingDeadlineSeconds
		if *in == nil {
			*out = nil
		} else {
			*out = new(int64)
			**out = **in
		}
	}
	if in.Suspend != nil {
		in, out := &in.Suspend, &out.Suspend
		if *in == nil {
			*out = nil
		} else {
			*out = new(bool)
			**out = **in
		}
	}
	in.JobTemplate.DeepCopyInto(&out.JobTemplate)
	if in.SuccessfulJobsHistoryLimit != nil {
		in, out := &in.SuccessfulJobsHistoryLimit, &out.SuccessfulJobsHistoryLimit
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	if in.FailedJobsHistoryLimit != nil {
		in, out := &in.FailedJobsHistoryLimit, &out.FailedJobsHistoryLimit
		if *in == nil {
			*out = nil
		} else {
			*out = new(int32)
			**out = **in
		}
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new CronJobSpec.
func (x *CronJobSpec) DeepCopy() *CronJobSpec {
	if x == nil {
		return nil
	}
	out := new(CronJobSpec)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *CronJobStatus) DeepCopyInto(out *CronJobStatus) {
	*out = *in
	if in.Active != nil {
		in, out := &in.Active, &out.Active
		*out = make([]v1.ObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.LastScheduleTime != nil {
		in, out := &in.LastScheduleTime, &out.LastScheduleTime
		if *in == nil {
			*out = nil
		} else {
			*out = new(meta_v1.Time)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new CronJobStatus.
func (x *CronJobStatus) DeepCopy() *CronJobStatus {
	if x == nil {
		return nil
	}
	out := new(CronJobStatus)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *JobTemplate) DeepCopyInto(out *JobTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Template.DeepCopyInto(&out.Template)
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new JobTemplate.
func (x *JobTemplate) DeepCopy() *JobTemplate {
	if x == nil {
		return nil
	}
	out := new(JobTemplate)
	x.DeepCopyInto(out)
	return out
}

// DeepCopyObject will perform a deep copy of the receiver, creating a new object.
func (x *JobTemplate) DeepCopyObject() runtime.Object {
	if c := x.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto will perform a deep copy of the receiver, writing to out. in must be non-nil.
func (in *JobTemplateSpec) DeepCopyInto(out *JobTemplateSpec) {
	*out = *in
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy will perform a deep copy of the receiver, creating a new JobTemplateSpec.
func (x *JobTemplateSpec) DeepCopy() *JobTemplateSpec {
	if x == nil {
		return nil
	}
	out := new(JobTemplateSpec)
	x.DeepCopyInto(out)
	return out
}
