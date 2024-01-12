/*
Copyright 2024 The cert-manager Authors.

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

package heuristics

import (
	"reflect"
	"testing"
)

func TestRecutNewLines(t *testing.T) {
	type args struct {
		lines []string
	}
	type want struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"MultiLineWithIndents",
			args{
				[]string{
					" Labels to apply to all resources",
					" Please note that this does not add labels to the resources created dynamically by the controllers.",
					" For these resources, you have to add the labels in the template in the cert-manager custom resource:",
					"",
					" eg. podTemplate/ ingressTemplate in ACMEChallengeSolverHTTP01Ingress",
					"    ref: https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEChallengeSolverHTTP01Ingress",
					" eg. secretTemplate in CertificateSpec",
					"    ref: https://cert-manager.io/docs/reference/api-docs/#cert-manager.io/v1.CertificateSpec",
				},
			},
			want{
				[]string{
					"Labels to apply to all resources",
					"Please note that this does not add labels to the resources created dynamically by the controllers. For these resources, you have to add the labels in the template in the cert-manager custom resource:",
					"",
					"eg. podTemplate/ ingressTemplate in ACMEChallengeSolverHTTP01Ingress",
					"   ref: https://cert-manager.io/docs/reference/api-docs/#acme.cert-manager.io/v1.ACMEChallengeSolverHTTP01Ingress",
					"eg. secretTemplate in CertificateSpec",
					"   ref: https://cert-manager.io/docs/reference/api-docs/#cert-manager.io/v1.CertificateSpec",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RecutNewLines(tt.args.lines); !reflect.DeepEqual(got, tt.want.lines) {
				t.Errorf("RecutNewLines() = %v, want %v", got, tt.want.lines)
			}
		})
	}
}
