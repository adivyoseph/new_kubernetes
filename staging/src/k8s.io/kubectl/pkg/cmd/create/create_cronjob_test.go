/*
Copyright 2018 The Kubernetes Authors.

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

package create

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateCronJob(t *testing.T) {
	cronjobName := "test-job"
	tests := map[string]struct {
		image     string
		command   []string
		schedule  string
		restart   string
		labels    string
		expected  *batchv1.CronJob
		expectErr bool
	}{
		"just image and OnFailure restart policy": {
			image:    "busybox",
			schedule: "0/5 * * * ?",
			restart:  "OnFailure",
			expected: &batchv1.CronJob{
				TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "CronJob"},
				ObjectMeta: metav1.ObjectMeta{
					Name: cronjobName,
				},
				Spec: batchv1.CronJobSpec{
					Schedule: "0/5 * * * ?",
					JobTemplate: batchv1.JobTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: cronjobName,
						},
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  cronjobName,
											Image: "busybox",
										},
									},
									RestartPolicy: corev1.RestartPolicyOnFailure,
								},
							},
						},
					},
				},
			},
		},
		"image, command , schedule and Never restart policy": {
			image:    "busybox",
			command:  []string{"date"},
			schedule: "0/5 * * * ?",
			restart:  "Never",
			expected: &batchv1.CronJob{
				TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "CronJob"},
				ObjectMeta: metav1.ObjectMeta{
					Name: cronjobName,
				},
				Spec: batchv1.CronJobSpec{
					Schedule: "0/5 * * * ?",
					JobTemplate: batchv1.JobTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: cronjobName,
						},
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:    cronjobName,
											Image:   "busybox",
											Command: []string{"date"},
										},
									},
									RestartPolicy: corev1.RestartPolicyNever,
								},
							},
						},
					},
				},
			},
		},
		"create configmap with label": {
			image:    "busybox",
			schedule: "0/5 * * * ?",
			restart:  "OnFailure",
			labels:   "key1=val1",
			expected: &batchv1.CronJob{
				TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "CronJob"},
				ObjectMeta: metav1.ObjectMeta{
					Name: cronjobName,
					Labels: map[string]string{
						"key1": "val1",
					},
				},
				Spec: batchv1.CronJobSpec{
					Schedule: "0/5 * * * ?",
					JobTemplate: batchv1.JobTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: cronjobName,
						},
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  cronjobName,
											Image: "busybox",
										},
									},
									RestartPolicy: corev1.RestartPolicyOnFailure,
								},
							},
						},
					},
				},
			},
		},
		"create configmap with multiple labels": {
			image:    "busybox",
			schedule: "0/5 * * * ?",
			restart:  "OnFailure",
			labels:   "key1=val1,key2=val2,key3=val3",
			expected: &batchv1.CronJob{
				TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "CronJob"},
				ObjectMeta: metav1.ObjectMeta{
					Name: cronjobName,
					Labels: map[string]string{
						"key1": "val1",
						"key2": "val2",
						"key3": "val3",
					},
				},
				Spec: batchv1.CronJobSpec{
					Schedule: "0/5 * * * ?",
					JobTemplate: batchv1.JobTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name: cronjobName,
						},
						Spec: batchv1.JobSpec{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{
										{
											Name:  cronjobName,
											Image: "busybox",
										},
									},
									RestartPolicy: corev1.RestartPolicyOnFailure,
								},
							},
						},
					},
				},
			},
		},
		"create configmap with invalid labels": {
			image:     "busybox",
			schedule:  "0/5 * * * ?",
			restart:   "OnFailure",
			labels:    "key1=val1,key2-val2,key3=val3",
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			o := &CreateCronJobOptions{
				Name:     cronjobName,
				Image:    tc.image,
				Command:  tc.command,
				Schedule: tc.schedule,
				Restart:  tc.restart,
				Labels:   tc.labels,
			}
			cronjob, err := o.createCronJob()
			if !tc.expectErr && err != nil {
				t.Errorf("test %s, unexpected error: %v", name, err)
			}
			if tc.expectErr && err == nil {
				t.Errorf("test %s was expecting an error but no error occurred", name)
			}
			if !apiequality.Semantic.DeepEqual(cronjob, tc.expected) {
				t.Errorf("expected:\n%#v\ngot:\n%#v", tc.expected, cronjob)
			}
		})
	}
}
