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

package node

import (
	"fmt"

	"k8s.io/kubernetes/cmd/kubeadm/app/api"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
	clientcmdapi "k8s.io/kubernetes/pkg/client/unversioned/clientcmd/api"
	"k8s.io/kubernetes/pkg/kubelet/util/csr"
	certutil "k8s.io/kubernetes/pkg/util/cert"
)

// PerformTLSBootstrap creates a RESTful client in order to execute certificate signing request.
func PerformTLSBootstrap(connection *api.ConnectionDetails) (*clientcmdapi.Config, error) {
	csrClient := connection.CertClient.CertificateSigningRequests()

	fmt.Println("<node/csr> created API client to obtain unique certificate for this node, generating keys and certificate signing request")

	key, err := certutil.MakeEllipticPrivateKeyPEM()
	if err != nil {
		return nil, fmt.Errorf("<node/csr> failed to generating private key [%v]", err)
	}

	cert, err := csr.RequestNodeCertificate(csrClient, key, connection.NodeName)
	if err != nil {
		return nil, fmt.Errorf("<node/csr> failed to request signed certificate from the API server [%v]", err)
	}

	// TODO(phase1+) print some basic info about the cert
	fmt.Println("<node/csr> received signed certificate from the API server, generating kubelet configuration")

	bareClientConfig := kubeadmutil.CreateBasicClientConfig("kubernetes", connection.Endpoint, connection.CACert)
	finalConfig := kubeadmutil.MakeClientConfigWithCerts(
		bareClientConfig, "kubernetes", fmt.Sprintf("kubelet-%s", connection.NodeName),
		key, cert,
	)

	return finalConfig, nil
}
