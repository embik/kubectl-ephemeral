package options

import (
	"fmt"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type EphemeralContainerOptions struct {
	genericclioptions.IOStreams

	ConfigFlags *genericclioptions.ConfigFlags

	args                []string
	ContainerFilePath   string
	TargetPodName       string
	TargetContainerName string
}

func NewEphemeralContainerOptions(streams genericclioptions.IOStreams) *EphemeralContainerOptions {
	return &EphemeralContainerOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// Complete sets all information required for updating the current context
func (o *EphemeralContainerOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	if len(o.args) > 0 {
		o.TargetPodName = args[0]
	}

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *EphemeralContainerOptions) Validate() error {
	if len(o.args) != 1 {
		return fmt.Errorf("missing a target pod")
	}

	if o.ContainerFilePath == "" {
		return fmt.Errorf("must pass a container file path")
	}

	if o.TargetContainerName == "" {
		return fmt.Errorf("must pass a container name within the target pod")
	}

	return nil
}
