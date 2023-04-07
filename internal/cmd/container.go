package cmd

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/yaml"
)

func getContainerFromFile(filepath string) (*corev1.EphemeralContainer, error) {
	var container corev1.EphemeralContainer

	containerData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	if err := yaml.UnmarshalStrict(containerData, &container); err != nil {
		return nil, err
	}

	// add a random suffix to the container name since ephemeral containers stay around.
	container.Name = fmt.Sprintf("%s-%s", container.Name, utilrand.String(6))

	return &container, nil
}
