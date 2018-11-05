package main

import "fmt"

import (
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdExposeService(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewExposeServiceOptions(streams)

	validArgs := []string{}
	resources := regexp.MustCompile(`\s*,`).Split(exposeResources, -1)
	for _, r := range resources {
		validArgs = append(validArgs, strings.Fields(r)[0])
	}

	cmd := &cobra.Command{
		Use:                   "expose (-f FILENAME | TYPE NAME) [--port=port] [--protocol=TCP|UDP|SCTP] [--target-port=number-or-name] [--name=name] [--external-ip=external-ip-of-service] [--type=type]",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Take a replication controller, service, deployment or pod and expose it as a new Kubernetes Service"),
		Long:                  exposeLong,
		Example:               exposeExample,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.RunExpose(cmd, args))
		},
		ValidArgs: validArgs,
	}

	o.RecordFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	cmd.Flags().String("generator", "service/v2", i18n.T("The name of the API generator to use. There are 2 generators: 'service/v1' and 'service/v2'. The only difference between them is that service port in v1 is named 'default', while it is left unnamed in v2. Default is 'service/v2'."))
	cmd.Flags().String("protocol", "", i18n.T("The network protocol for the service to be created. Default is 'TCP'."))
	cmd.Flags().String("port", "", i18n.T("The port that the service should serve on. Copied from the resource being exposed, if unspecified"))
	cmd.Flags().String("type", "", i18n.T("Type for this service: ClusterIP, NodePort, LoadBalancer, or ExternalName. Default is 'ClusterIP'."))
	cmd.Flags().String("load-balancer-ip", "", i18n.T("IP to assign to the LoadBalancer. If empty, an ephemeral IP will be created and used (cloud-provider specific)."))
	cmd.Flags().String("selector", "", i18n.T("A label selector to use for this service. Only equality-based selector requirements are supported. If empty (the default) infer the selector from the replication controller or replica set.)"))
	cmd.Flags().StringP("labels", "l", "", "Labels to apply to the service created by this call.")
	cmd.Flags().String("container-port", "", i18n.T("Synonym for --target-port"))
	cmd.Flags().MarkDeprecated("container-port", "--container-port will be removed in the future, please use --target-port instead")
	cmd.Flags().String("target-port", "", i18n.T("Name or number for the port on the container that the service should direct traffic to. Optional."))
	cmd.Flags().String("external-ip", "", i18n.T("Additional external IP address (not managed by Kubernetes) to accept for the service. If this IP is routed to a node, the service can be accessed by this IP in addition to its generated service IP."))
	cmd.Flags().String("overrides", "", i18n.T("An inline JSON override for the generated object. If this is non-empty, it is used to override the generated object. Requires that the object supply a valid apiVersion field."))
	cmd.Flags().String("name", "", i18n.T("The name for the newly created object."))
	cmd.Flags().String("session-affinity", "", i18n.T("If non-empty, set the session affinity for the service to this; legal values: 'None', 'ClientIP'"))
	cmd.Flags().String("cluster-ip", "", i18n.T("ClusterIP to be assigned to the service. Leave empty to auto-allocate, or set to 'None' to create a headless service."))

	usage := "identifying the resource to expose a service"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	cmdutil.AddDryRunFlag(cmd)
	cmdutil.AddApplyAnnotationFlags(cmd)
	return cmd
}
