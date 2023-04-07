package cmd

import (
	"context"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"

	"github.com/embik/kubectl-ephemeral-container/internal/options"
)

func NewEphemeralContainerCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := options.NewEphemeralContainerOptions(streams)

	cmd := &cobra.Command{
		Use:          "ephemeral-container [pod name] [flags]",
		Short:        "Creates an ephemeral container in a target Pod from a YAML specification and execs into it.",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			if err := Run(o, c.Context()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.ContainerFilePath, "file", "f", "", "file containing a YAML container spec to create an ephemeral container from")
	o.ConfigFlags.AddFlags(cmd.Flags())

	return cmd
}

func Run(opts *options.EphemeralContainerOptions, ctx context.Context) error {
	container, err := getContainerFromFile(opts.ContainerFilePath)
	if err != nil {
		return err
	}

	clientConfig, err := opts.ConfigFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return err
	}

	client, err := corev1.NewForConfig(clientConfig)
	if err != nil {
		return err
	}

	namespace := opts.ConfigFlags.Namespace
	if namespace == nil || *namespace == "" {
		namespace = pointer.String(metav1.NamespaceDefault)
	}

	pod, err := client.Pods(*namespace).Get(ctx, opts.TargetPodName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, *container)

	_, err = client.Pods(*namespace).UpdateEphemeralContainers(ctx, pod.Name, pod, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
