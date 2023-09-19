//
// Copyright 2022 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhooks

import (
	"testing"

	appstudiov1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"

	"github.com/stretchr/testify/assert"
)

func TestApplicationValidatingWebhook(t *testing.T) {

	originalApplication := appstudiov1alpha1.Application{
		Spec: appstudiov1alpha1.ApplicationSpec{
			DisplayName: "My App",
			AppModelRepository: appstudiov1alpha1.ApplicationGitRepository{
				URL: "http://appmodelrepo",
			},
			GitOpsRepository: appstudiov1alpha1.ApplicationGitRepository{
				URL: "http://gitopsrepo",
			},
		},
	}

	tests := []struct {
		name      string
		updateApp appstudiov1alpha1.Application
		err       string
	}{
		{
			name: "app model repo cannot be changed",
			err:  "app model repository cannot be updated",
			updateApp: appstudiov1alpha1.Application{
				Spec: appstudiov1alpha1.ApplicationSpec{
					DisplayName: "My App",
					AppModelRepository: appstudiov1alpha1.ApplicationGitRepository{
						URL: "http://appmodelrepo1",
					},
					GitOpsRepository: appstudiov1alpha1.ApplicationGitRepository{
						URL: "http://gitopsrepo",
					},
				},
			},
		},
		{
			name: "gitops repo cannot be changed",
			err:  "gitops repository cannot be updated",
			updateApp: appstudiov1alpha1.Application{
				Spec: appstudiov1alpha1.ApplicationSpec{
					DisplayName: "My App",
					AppModelRepository: appstudiov1alpha1.ApplicationGitRepository{
						URL: "http://appmodelrepo",
					},
					GitOpsRepository: appstudiov1alpha1.ApplicationGitRepository{
						URL: "http://gitopsrepo1",
					},
				},
			},
		},
		{
			name: "display name can be changed",
			updateApp: appstudiov1alpha1.Application{
				Spec: appstudiov1alpha1.ApplicationSpec{
					DisplayName: "My App 2",
					AppModelRepository: appstudiov1alpha1.ApplicationGitRepository{
						URL: "http://appmodelrepo",
					},
					GitOpsRepository: appstudiov1alpha1.ApplicationGitRepository{
						URL: "http://gitopsrepo",
					},
				},
			},
		},
		{
			name: "not application",
			err:  "runtime object is not of type Application",
			updateApp: appstudiov1alpha1.Application{
				Spec: appstudiov1alpha1.ApplicationSpec{
					DisplayName: "My App",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error
			if test.name == "not application" {
				originalComponent := appstudiov1alpha1.Component{
					Spec: appstudiov1alpha1.ComponentSpec{
						ComponentName: "component",
						Application:   "application",
					},
				}
				err = test.updateApp.ValidateUpdate(&originalComponent)
			} else {
				err = test.updateApp.ValidateUpdate(&originalApplication)
			}

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}
		})
	}
}

func TestApplicationDeleteValidatingWebhook(t *testing.T) {

	tests := []struct {
		name string
		app  appstudiov1alpha1.Application
		err  string
	}{
		{
			name: "ValidateDelete should return nil, it's unimplimented",
			err:  "",
			app:  appstudiov1alpha1.Application{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.app.ValidateDelete()

			if test.err == "" {
				assert.Nil(t, err)
			} else {
				assert.Contains(t, err.Error(), test.err)
			}
		})
	}
}
