/*
Copyright 2024 The Kubernetes Authors.

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

package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"

	authnv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"
	"k8s.io/kubernetes/test/integration/framework"
	"k8s.io/kubernetes/test/utils/ktesting"
)

func TestAuthnToKAS(t *testing.T) {
	tCtx := ktesting.Init(t)

	frontProxyCA, frontProxyClient, frontProxyKey, err := newTestCAWithClient(
		pkix.Name{
			CommonName: "test-front-proxy-ca",
		},
		big.NewInt(43),
		pkix.Name{
			CommonName:   "test-aggregated-apiserver",
			Organization: []string{"system:masters"},
		},
		big.NewInt(86),
	)
	if err != nil {
		t.Error(err)
		return
	}

	modifyOpts := func(setUIDHeaders bool) func(opts *options.ServerRunOptions) {
		return func(opts *options.ServerRunOptions) {
			opts.GenericServerRunOptions.MaxRequestBodyBytes = 1024 * 1024

			// rewrite the client + request header CA certs with our own content
			frontProxyCAFilename := opts.Authentication.RequestHeader.ClientCAFile
			if err := os.WriteFile(frontProxyCAFilename, frontProxyCA, 0644); err != nil {
				t.Fatal(err)
			}

			opts.Authentication.RequestHeader.AllowedNames = append(opts.Authentication.RequestHeader.AllowedNames, "test-aggregated-apiserver")
			if setUIDHeaders {
				opts.Authentication.RequestHeader.UIDHeaders = []string{"X-Remote-Uid"}
			}
		}
	}

	for _, tt := range []struct {
		name   string
		setUID bool
	}{
		{
			name:   "KAS without UID config",
			setUID: false,
		},
		{
			name:   "KAS with UID config",
			setUID: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, kubeConfig, tearDownFn := framework.StartTestServer(tCtx, t, framework.TestServerSetup{
				ModifyServerRunOptions: modifyOpts(tt.setUID),
			})
			defer tearDownFn()

			// Test an aggregated apiserver client (signed by the new front proxy CA) is authorized
			extensionApiserverClient, err := kubernetes.NewForConfig(&rest.Config{
				Host: kubeConfig.Host,
				TLSClientConfig: rest.TLSClientConfig{
					CAData:     kubeConfig.TLSClientConfig.CAData,
					CAFile:     kubeConfig.TLSClientConfig.CAFile,
					ServerName: kubeConfig.TLSClientConfig.ServerName,
					KeyData:    frontProxyKey,
					CertData:   frontProxyClient,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			selfInfo := &authnv1.SelfSubjectReview{}
			err = extensionApiserverClient.AuthenticationV1().RESTClient().
				Post().
				Resource("selfsubjectreviews").
				VersionedParams(&metav1.CreateOptions{}, scheme.ParameterCodec).
				Body(&authnv1.SelfSubjectReview{}).
				SetHeader("X-Remote-Uid", "test-uid").
				SetHeader("X-Remote-User", "testuser").
				SetHeader("X-Remote-Groups", "group1", "group2").
				Do(tCtx).
				Into(selfInfo)
			if err != nil {
				t.Fatalf("failed to retrieve self-info: %v", err)
			}

			if selfUID := selfInfo.Status.UserInfo.UID; (len(selfUID) != 0) != tt.setUID {
				t.Errorf("UID should be set: %v, but we got %v", tt.setUID, selfUID)
			}
		})
	}

}

func newTestCAWithClient(caSubject pkix.Name, caSerial *big.Int, clientSubject pkix.Name, subjectSerial *big.Int) (caPEMBytes, clientCertPEMBytes, clientKeyPEMBytes []byte, err error) {
	ca := &x509.Certificate{
		SerialNumber:          caSerial,
		Subject:               caSubject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	clientCert := &x509.Certificate{
		SerialNumber: subjectSerial,
		Subject:      clientSubject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	clientCertPrivateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, nil, err
	}

	clientCertPrivateKeyPEM := new(bytes.Buffer)
	err = pem.Encode(clientCertPrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientCertPrivateKey),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	clientCertBytes, err := x509.CreateCertificate(rand.Reader, clientCert, ca, &clientCertPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	clientCertPEM := new(bytes.Buffer)
	err = pem.Encode(clientCertPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCertBytes,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return caPEM.Bytes(),
		clientCertPEM.Bytes(),
		clientCertPrivateKeyPEM.Bytes(),
		nil
}
