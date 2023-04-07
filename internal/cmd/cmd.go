package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/utils/pointer"

	"github.com/embik/kubectl-ephemeral/internal/options"
)

func NewEphemeralContainerCmd(streams genericclioptions.IOStreams) *cobra.Command {
	o := options.NewEphemeralContainerOptions(streams)

	cmd := &cobra.Command{
		Use:          "ephemeral [pod name] [flags]",
		Short:        "Creates an ephemeral container in a target Pod from a YAML specification.",
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
	cmd.Flags().StringVarP(&o.TargetContainerName, "container", "c", "", "container within the target pod that the ephemeral container will attach itself to")
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

	containerValid := false
	for _, container := range pod.Spec.Containers {
		if container.Name == opts.TargetContainerName {
			containerValid = true
			break
		}
	}

	if !containerValid {
		return fmt.Errorf("container '%s' not found in pod '%s'", opts.TargetContainerName, opts.TargetPodName)
	}

	container.TargetContainerName = opts.TargetContainerName

	pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, *container)

	_, err = client.Pods(*namespace).UpdateEphemeralContainers(ctx, pod.Name, pod, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
