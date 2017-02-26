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

package kubefed

import (
	"fmt"
	"io"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubernetes/federation/pkg/kubefed/util"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/kubectl"
	kubectlcmd "k8s.io/kubernetes/pkg/kubectl/cmd"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	// defaultClusterCIDR is the default CIDR range accepted by the
	// joining API server. See `apis/federation.ClusterSpec` for
	// details.
	// TODO(madhusudancs): Make this value customizable.
	defaultClientCIDR = "0.0.0.0/0"
)

var (
	join_long = templates.LongDesc(`
		Join a cluster to a federation.

        Current context is assumed to be a federation API
        server. Please use the --context flag otherwise.`)
	join_example = templates.Examples(`
		# Join a cluster to a federation by specifying the
		# cluster name and the context name of the federation
		# control plane's host cluster. Cluster name must be
		# a valid RFC 1123 subdomain name. Cluster context
		# must be specified if the cluster name is different
		# than the cluster's context in the local kubeconfig.
		kubefed join foo --host-cluster-context=bar`)
)

type joinFederation struct {
	commonOptions util.SubcommandOptions
	options       joinFederationOptions
}

type joinFederationOptions struct {
	clusterContext string
	secretName     string
	dryRun         bool
}

func (o *joinFederationOptions) Bind(flags *pflag.FlagSet) {
	flags.StringVar(&o.clusterContext, "cluster-context", "", "Name of the cluster's context in the local kubeconfig. Defaults to cluster name if unspecified.")
	flags.StringVar(&o.secretName, "secret-name", "", "Name of the secret where the cluster's credentials will be stored in the host cluster. This name should be a valid RFC 1035 label. Defaults to cluster name if unspecified.")
}

// NewCmdJoin defines the `join` command that joins a cluster to a
// federation.
func NewCmdJoin(f cmdutil.Factory, cmdOut io.Writer, config util.AdminConfig) *cobra.Command {
	opts := &joinFederation{}

	cmd := &cobra.Command{
		Use:     "join CLUSTER_NAME --host-cluster-context=HOST_CONTEXT",
		Short:   "Join a cluster to a federation",
		Long:    join_long,
		Example: join_example,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(opts.Complete(cmd, args))
			cmdutil.CheckErr(opts.Run(f, cmdOut, config, cmd))
		},
	}

	cmdutil.AddApplyAnnotationFlags(cmd)
	cmdutil.AddValidateFlags(cmd)
	cmdutil.AddPrinterFlags(cmd)
	cmdutil.AddGeneratorFlags(cmd, cmdutil.ClusterV1Beta1GeneratorName)

	flags := cmd.Flags()
	opts.commonOptions.Bind(flags)
	opts.options.Bind(flags)

	return cmd
}

// Complete ensures that options are valid and marshals them if necessary.
func (j *joinFederation) Complete(cmd *cobra.Command, args []string) error {
	err := j.commonOptions.SetName(cmd, args)
	if err != nil {
		return err
	}

	j.options.dryRun = cmdutil.GetDryRunFlag(cmd)

	if j.options.clusterContext == "" {
		j.options.clusterContext = j.commonOptions.Name
	}
	if j.options.secretName == "" {
		j.options.secretName = j.commonOptions.Name
	}

	glog.V(2).Infof("Args and flags: name %s, host: %s, host-system-namespace: %s, kubeconfig: %s, cluster-context: %s, secret-name: %s, dry-run: %s", j.commonOptions.Name, j.commonOptions.Host, j.commonOptions.FederationSystemNamespace, j.commonOptions.Kubeconfig, j.options.clusterContext, j.options.secretName, j.options.dryRun)

	return nil
}

// Run is the implementation of the `join federation` command.
func (j *joinFederation) Run(f cmdutil.Factory, cmdOut io.Writer, config util.AdminConfig, cmd *cobra.Command) error {
	po := config.PathOptions()
	po.LoadingRules.ExplicitPath = j.commonOptions.Kubeconfig
	clientConfig, err := po.GetStartingConfig()
	if err != nil {
		return err
	}
	generator, err := clusterGenerator(clientConfig, j.commonOptions.Name, j.options.clusterContext, j.options.secretName)
	if err != nil {
		glog.V(2).Infof("Failed creating cluster generator: %v", err)
		return err
	}
	glog.V(2).Infof("Created cluster generator: %#v", generator)

	hostFactory := config.HostFactory(j.commonOptions.Host, j.commonOptions.Kubeconfig)
	clientset, err := hostFactory.ClientSet()
	if err != nil {
		glog.V(2).Infof("Failed to serialize the kubeconfig for the given context %q: %v", j.options.clusterContext, err)
		return err
	}
	// We are not using the `kubectl create secret` machinery through
	// `RunCreateSubcommand` as we do to the cluster resource below
	// because we have a bunch of requirements that the machinery does
	// not satisfy.
	// 1. We want to create the secret in a specific namespace, which
	//    is neither the "default" namespace nor the one specified
	//    via the `--namespace` flag.
	// 2. `SecretGeneratorV1` requires LiteralSources in a string-ified
	//    form that it parses to generate the secret data key-value
	//    pairs. We, however, have the key-value pairs ready without a
	//    need for parsing.
	// 3. The result printing mechanism needs to be mostly quiet. We
	//    don't have to print the created secret in the default case.
	// Having said that, secret generation machinery could be altered to
	// suit our needs, but it is far less invasive and readable this way.
	_, err = createSecret(clientset, clientConfig, j.commonOptions.FederationSystemNamespace, j.options.clusterContext, j.options.secretName, j.options.dryRun)
	if err != nil {
		glog.V(2).Infof("Failed creating the cluster credentials secret: %v", err)
		return err
	}
	glog.V(2).Infof("Cluster credentials secret created")

	err = kubectlcmd.RunCreateSubcommand(f, cmd, cmdOut, &kubectlcmd.CreateSubcommandOptions{
		Name:                j.commonOptions.Name,
		StructuredGenerator: generator,
		DryRun:              j.options.dryRun,
		OutputFormat:        cmdutil.GetFlagString(cmd, "output"),
	})
	if !errors.IsAlreadyExists(err) && !j.options.dryRun {
		// Ensure to delete the secret from base cluster before exiting
		delErr := clientset.Core().Secrets(j.commonOptions.FederationSystemNamespace).Delete(j.options.secretName, &metav1.DeleteOptions{})
		if delErr != nil {
			return fmt.Errorf("%v, Cleanup failed after failed join:%v", err, delErr)
		}
	}

	return err
}

// minifyConfig is a wrapper around `clientcmdapi.MinifyConfig()` that
// sets the current context to the given context before calling
// `clientcmdapi.MinifyConfig()`.
func minifyConfig(clientConfig *clientcmdapi.Config, context string) (*clientcmdapi.Config, error) {
	// MinifyConfig inline-modifies the passed clientConfig. So we make a
	// copy of it before passing the config to it. A shallow copy is
	// sufficient because the underlying fields will be reconstructed by
	// MinifyConfig anyway.
	newClientConfig := *clientConfig
	newClientConfig.CurrentContext = context
	err := clientcmdapi.MinifyConfig(&newClientConfig)
	if err != nil {
		return nil, err
	}
	return &newClientConfig, nil
}

// createSecret extracts the kubeconfig for a given cluster and populates
// a secret with that kubeconfig.
func createSecret(clientset *internalclientset.Clientset, clientConfig *clientcmdapi.Config, namespace, contextName, secretName string, dryRun bool) (runtime.Object, error) {
	// Minify the kubeconfig to ensure that there is only information
	// relevant to the cluster we are registering.
	newClientConfig, err := minifyConfig(clientConfig, contextName)
	if err != nil {
		glog.V(2).Infof("Failed to minify the kubeconfig for the given context %q: %v", contextName, err)
		return nil, err
	}

	// Flatten the kubeconfig to ensure that all the referenced file
	// contents are inlined.
	err = clientcmdapi.FlattenConfig(newClientConfig)
	if err != nil {
		glog.V(2).Infof("Failed to flatten the kubeconfig for the given context %q: %v", contextName, err)
		return nil, err
	}

	return util.CreateKubeconfigSecret(clientset, newClientConfig, namespace, secretName, dryRun)
}

// clusterGenerator extracts the cluster information from the supplied
// kubeconfig and builds a StructuredGenerator for the
// `federation/cluster` API resource.
func clusterGenerator(clientConfig *clientcmdapi.Config, name, contextName, secretName string) (kubectl.StructuredGenerator, error) {
	// Get the context from the config.
	ctx, found := clientConfig.Contexts[contextName]
	if !found {
		return nil, fmt.Errorf("cluster context %q not found", contextName)
	}

	// Get the cluster object corresponding to the supplied context.
	cluster, found := clientConfig.Clusters[ctx.Cluster]
	if !found {
		return nil, fmt.Errorf("cluster endpoint not found for %q", name)
	}

	// Extract the scheme portion of the cluster APIServer endpoint and
	// default it to `https` if it isn't specified.
	scheme := extractScheme(cluster.Server)
	serverAddress := cluster.Server
	if scheme == "" {
		// Use "https" as the default scheme.
		scheme := "https"
		serverAddress = strings.Join([]string{scheme, serverAddress}, "://")
	}

	generator := &kubectl.ClusterGeneratorV1Beta1{
		Name:          name,
		ClientCIDR:    defaultClientCIDR,
		ServerAddress: serverAddress,
		SecretName:    secretName,
	}
	return generator, nil
}

// extractScheme parses the given URL to extract the scheme portion
// out of it.
func extractScheme(url string) string {
	scheme := ""
	segs := strings.SplitN(url, "://", 2)
	if len(segs) == 2 {
		scheme = segs[0]
	}
	return scheme
}
