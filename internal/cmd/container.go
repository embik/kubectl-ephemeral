package cmd

import (
	"os"

	corev1 "k8s.io/api/core/v1"
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

	return &container, nil
}
