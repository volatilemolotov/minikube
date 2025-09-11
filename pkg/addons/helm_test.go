/*
Copyright 2023 The Kubernetes Authors All rights reserved.

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

package addons

import (
	"context"
	"strings"
	"testing"

	"k8s.io/minikube/pkg/minikube/assets"
)

func TestHelmCommand(t *testing.T) {
	tests := []struct {
		description string
		chart       *assets.HelmChart
		enable      bool
		expected    string
	}{
		{
			description: "enable an addon",
			chart: &assets.HelmChart{
				Name:      "addon-name",
				Repo:      "addon-repo/addon-chart",
				Namespace: "addon-namespace",
				Values:    []string{"--set", "key=value"},
			},
			enable:   true,
			expected: "sudo KUBECONFIG=/var/lib/minikube/kubeconfig helm install addon-name addon-repo/addon-chart --create-namespace --namespace addon-namespace --set key=value",
		},
		{
			description: "enable an addon without namespace",
			chart: &assets.HelmChart{
				Name: "addon-name",
				Repo: "addon-repo/addon-chart",
			},
			enable:   true,
			expected: "sudo KUBECONFIG=/var/lib/minikube/kubeconfig helm install addon-name addon-repo/addon-chart --create-namespace",
		},
		{
			description: "disable an addon",
			chart: &assets.HelmChart{
				Name:      "addon-name",
				Namespace: "addon-namespace",
			},
			enable:   false,
			expected: "sudo KUBECONFIG=/var/lib/minikube/kubeconfig helm uninstall addon-name --namespace addon-namespace",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			command := helmCommand(context.Background(), test.chart, test.enable)
			actual := strings.Join(command.Args, " ")
			if actual != test.expected {
				t.Errorf("helm command mismatch:\nexpected: %s\nactual:   %s", test.expected, actual)
			}
		})
	}
}