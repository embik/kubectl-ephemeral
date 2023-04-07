package main

import (
	"os"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/embik/kubectl-ephemeral-container/internal/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-ephemeral-container", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewEphemeralContainerCmd(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
