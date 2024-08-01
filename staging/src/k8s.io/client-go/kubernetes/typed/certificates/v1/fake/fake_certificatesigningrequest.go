/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/api/certificates/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certificatesv1 "k8s.io/client-go/applyconfigurations/certificates/v1"
	gentype "k8s.io/client-go/gentype"
	clientset "k8s.io/client-go/kubernetes/typed/certificates/v1"
	testing "k8s.io/client-go/testing"
)

// fakeCertificateSigningRequests implements CertificateSigningRequestInterface
type fakeCertificateSigningRequests struct {
	*gentype.FakeClientWithListAndApply[*v1.CertificateSigningRequest, *v1.CertificateSigningRequestList, *certificatesv1.CertificateSigningRequestApplyConfiguration]
	Fake *FakeCertificatesV1
}

func newFakeCertificateSigningRequests(fake *FakeCertificatesV1) clientset.CertificateSigningRequestInterface {
	return &fakeCertificateSigningRequests{
		gentype.NewFakeClientWithListAndApply[*v1.CertificateSigningRequest, *v1.CertificateSigningRequestList, *certificatesv1.CertificateSigningRequestApplyConfiguration](
			fake.Fake,
			"",
			v1.SchemeGroupVersion.WithResource("certificatesigningrequests"),
			v1.SchemeGroupVersion.WithKind("CertificateSigningRequest"),
			func() *v1.CertificateSigningRequest { return &v1.CertificateSigningRequest{} },
			func() *v1.CertificateSigningRequestList { return &v1.CertificateSigningRequestList{} },
			func(dst, src *v1.CertificateSigningRequestList) { dst.ListMeta = src.ListMeta },
			func(list *v1.CertificateSigningRequestList) []*v1.CertificateSigningRequest {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1.CertificateSigningRequestList, items []*v1.CertificateSigningRequest) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}

// UpdateApproval takes the representation of a certificateSigningRequest and updates it. Returns the server's representation of the certificateSigningRequest, and an error, if there is any.
func (c *fakeCertificateSigningRequests) UpdateApproval(ctx context.Context, certificateSigningRequestName string, certificateSigningRequest *v1.CertificateSigningRequest, opts metav1.UpdateOptions) (result *v1.CertificateSigningRequest, err error) {
	emptyResult := &v1.CertificateSigningRequest{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(c.Resource(), "approval", certificateSigningRequest, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1.CertificateSigningRequest), err
}
