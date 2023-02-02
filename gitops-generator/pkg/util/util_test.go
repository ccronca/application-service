//
// Copyright 2023 Red Hat, Inc.
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

package util

import (
	"reflect"
	"testing"

	"github.com/devfile/library/v2/pkg/devfile/parser"
	appstudiov1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"
	gitopsgenv1alpha1 "github.com/redhat-developer/gitops-generator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetMappedComponent(t *testing.T) {

	other := make([]interface{}, 1)
	other[0] = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment1",
		},
	}

	tests := []struct {
		name                string
		component           appstudiov1alpha1.Component
		kubernetesResources parser.KubernetesResources
		want                gitopsgenv1alpha1.GeneratorOptions
	}{
		{
			name: "Test01ComponentSpecFilledIn",
			component: appstudiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testcomponent",
					Namespace: "testnamespace",
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "frontEnd",
					Application:   "AppTest001",
					Secret:        "Secret",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceLimitsCPU: resource.MustParse("1"),
							corev1.ResourceMemory:    resource.MustParse("1Gi"),
						},
					},
					TargetPort: 8080,
					Route:      "https://testroute",
					Env: []corev1.EnvVar{
						{
							Name:  "env1",
							Value: "env1Value",
						},
					},
					ContainerImage:               "myimage:image",
					SkipGitOpsResourceGeneration: false,
					Source: appstudiov1alpha1.ComponentSource{
						ComponentSourceUnion: appstudiov1alpha1.ComponentSourceUnion{
							GitSource: &appstudiov1alpha1.GitSource{
								URL:           "https://host/git-repo.git",
								Revision:      "1.0",
								Context:       "/context",
								DevfileURL:    "https://mydevfileurl",
								DockerfileURL: "https://mydockerfileurl",
							},
						},
					},
				},
			},
			want: gitopsgenv1alpha1.GeneratorOptions{
				Name:        "testcomponent",
				Namespace:   "testnamespace",
				Application: "AppTest001",
				Secret:      "Secret",
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceLimitsCPU: resource.MustParse("1"),
						corev1.ResourceMemory:    resource.MustParse("1Gi"),
					},
				},
				TargetPort: 8080,
				Route:      "https://testroute",
				BaseEnvVar: []corev1.EnvVar{
					{
						Name:  "env1",
						Value: "env1Value",
					},
				},
				ContainerImage: "myimage:image",
				GitSource: &gitopsgenv1alpha1.GitSource{
					URL: "https://host/git-repo.git",
				},
			},
		},
		{
			name: "Test02EmptyComponentSource",
			component: appstudiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testcomponent",
					Namespace: "testnamespace",
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "frontEnd",
					Application:   "AppTest002",
					Source:        appstudiov1alpha1.ComponentSource{},
				},
			},
			want: gitopsgenv1alpha1.GeneratorOptions{
				Name:        "testcomponent",
				Namespace:   "testnamespace",
				Application: "AppTest002",
				GitSource:   &gitopsgenv1alpha1.GitSource{},
			},
		},
		{
			name: "Test03NoSource",
			component: appstudiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testcomponent",
					Namespace: "testnamespace",
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "frontEnd",
					Application:   "AppTest003",
				},
			},
			want: gitopsgenv1alpha1.GeneratorOptions{
				Name:        "testcomponent",
				Namespace:   "testnamespace",
				Application: "AppTest003",
			},
		},
		{
			name: "Test04EmptyComponentSourceUnion",
			component: appstudiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testcomponent",
					Namespace: "testnamespace",
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "frontEnd",
					Application:   "AppTest004",
					Source: appstudiov1alpha1.ComponentSource{
						ComponentSourceUnion: appstudiov1alpha1.ComponentSourceUnion{},
					},
				},
			},
			want: gitopsgenv1alpha1.GeneratorOptions{
				Name:        "testcomponent",
				Namespace:   "testnamespace",
				Application: "AppTest004",
				GitSource:   &gitopsgenv1alpha1.GitSource{},
			},
		},
		{
			name: "Test05EmptyGitSource",
			component: appstudiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testcomponent",
					Namespace: "testnamespace",
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "frontEnd",
					Application:   "AppTest005",
					Source: appstudiov1alpha1.ComponentSource{
						ComponentSourceUnion: appstudiov1alpha1.ComponentSourceUnion{
							GitSource: &appstudiov1alpha1.GitSource{},
						},
					},
				},
			},
			want: gitopsgenv1alpha1.GeneratorOptions{
				Name:        "testcomponent",
				Namespace:   "testnamespace",
				Application: "AppTest005",
				GitSource:   &gitopsgenv1alpha1.GitSource{},
			},
		},
		{
			name: "Test06KubernetesResources",
			component: appstudiov1alpha1.Component{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testcomponent",
					Namespace: "testnamespace",
				},
				Spec: appstudiov1alpha1.ComponentSpec{
					ComponentName: "frontEnd",
					Application:   "AppTest005",
					Source: appstudiov1alpha1.ComponentSource{
						ComponentSourceUnion: appstudiov1alpha1.ComponentSourceUnion{
							GitSource: &appstudiov1alpha1.GitSource{
								URL: "url",
							},
						},
					},
				},
			},
			kubernetesResources: parser.KubernetesResources{
				Deployments: []appsv1.Deployment{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "deployment1",
						},
					},
				},
				Services: []corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name: "service1",
						},
					},
				},
				Others: other,
			},
			want: gitopsgenv1alpha1.GeneratorOptions{
				Name:        "testcomponent",
				Namespace:   "testnamespace",
				Application: "AppTest005",
				GitSource: &gitopsgenv1alpha1.GitSource{
					URL: "url",
				},
				KubernetesResources: gitopsgenv1alpha1.KubernetesResources{
					Deployments: []appsv1.Deployment{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "deployment1",
							},
						},
					},
					Services: []corev1.Service{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name: "service1",
							},
						},
					},
					Others: other,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mappedComponent := GetMappedGitOpsComponent(tt.component, tt.kubernetesResources)
			assert.True(t, tt.want.Name == mappedComponent.Name, "Expected ObjectMeta.Name: %s, is different than actual: %s", tt.want.Name, mappedComponent.Name)
			assert.True(t, tt.want.Namespace == mappedComponent.Namespace, "Expected ObjectMeta.Namespace: %s, is different than actual: %s", tt.want.Namespace, mappedComponent.Namespace)
			assert.True(t, tt.want.Application == mappedComponent.Application, "Expected Spec.Application: %s, is different than actual: %s", tt.want.Application, mappedComponent.Application)
			assert.True(t, tt.want.Secret == mappedComponent.Secret, "Expected Spec.Secret: %s, is different than actual: %s", tt.want.Secret, mappedComponent.Secret)
			assert.True(t, reflect.DeepEqual(tt.want.Resources, mappedComponent.Resources), "Expected Spec.Resources: %s, is different than actual: %s", tt.want.Resources, mappedComponent.Resources)
			assert.True(t, tt.want.Route == mappedComponent.Route, "Expected Spec.Route: %s, is different than actual: %s", tt.want.Route, mappedComponent.Route)
			assert.True(t, reflect.DeepEqual(tt.want.BaseEnvVar, mappedComponent.BaseEnvVar), "Expected Spec.Env: %s, is different than actual: %s", tt.want.BaseEnvVar, mappedComponent.BaseEnvVar)
			assert.True(t, tt.want.ContainerImage == mappedComponent.ContainerImage, "Expected Spec.ContainerImage: %s, is different than actual: %s", tt.want.ContainerImage, mappedComponent.ContainerImage)

			if tt.want.GitSource != nil {
				assert.True(t, tt.want.GitSource.URL == mappedComponent.GitSource.URL, "Expected GitSource URL: %s, is different than actual: %s", tt.want.GitSource.URL, mappedComponent.GitSource.URL)
			}

			if !reflect.DeepEqual(tt.want.KubernetesResources, gitopsgenv1alpha1.KubernetesResources{}) {
				for _, wantDeployment := range tt.want.KubernetesResources.Deployments {
					matched := false
					for _, gotDeployment := range mappedComponent.KubernetesResources.Deployments {
						if wantDeployment.Name == gotDeployment.Name {
							matched = true
							break
						}
					}
					assert.True(t, matched, "Expected Deployment: %s, but didnt find in actual", wantDeployment.Name)
				}

				for _, wantService := range tt.want.KubernetesResources.Services {
					matched := false
					for _, gotService := range mappedComponent.KubernetesResources.Services {
						if wantService.Name == gotService.Name {
							matched = true
							break
						}
					}
					assert.True(t, matched, "Expected Service: %s, but didnt find in actual", wantService.Name)
				}

				for _, wantRoute := range tt.want.KubernetesResources.Routes {
					matched := false
					for _, gotRoute := range mappedComponent.KubernetesResources.Routes {
						if wantRoute.Name == gotRoute.Name {
							matched = true
							break
						}
					}
					assert.True(t, matched, "Expected Route: %s, but didnt find in actual", wantRoute.Name)
				}

				for _, wantIngress := range tt.want.KubernetesResources.Ingresses {
					matched := false
					for _, gotIngress := range mappedComponent.KubernetesResources.Ingresses {
						if wantIngress.Name == gotIngress.Name {
							matched = true
							break
						}
					}
					assert.True(t, matched, "Expected Ingress: %s, but didnt find in actual", wantIngress.Name)
				}

				for _, wantOther := range tt.want.KubernetesResources.Others {
					matched := false
					wantDeployment := wantOther.(appsv1.Deployment)

					for _, gotOther := range mappedComponent.KubernetesResources.Others {
						gotDeployment := gotOther.(appsv1.Deployment)
						if wantDeployment.Name == gotDeployment.Name {
							matched = true
							break
						}
					}
					assert.True(t, matched, "Expected Other: %s, but didnt find in actual", wantDeployment.Name)
				}
			}
		})
	}
}

func TestGetRemoteURL(t *testing.T) {
	tests := []struct {
		name      string
		gitopsUrl string
		token     string
		want      string
		wantErr   bool
	}{
		{
			name:      "Basic URL",
			gitopsUrl: "https://github.com/redhat-appstudio-appdata/test",
			token:     "my-token",
			want:      "https://my-token@github.com/redhat-appstudio-appdata/test",
		},
		{
			name:      "Invalid URL",
			gitopsUrl: "http://github.com/?org\nrepo",
			token:     "my-token",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remote, err := GetRemoteURL(tt.gitopsUrl, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestGetRemoteURL: unexpected error %v", err)
			}
			if !tt.wantErr && (remote != tt.want) {
				t.Errorf("TestGetRemoteURL: want %v, got %v", tt.want, remote)
			}
		})
	}
}

func TestProcessGitOpsStatus(t *testing.T) {
	tests := []struct {
		name         string
		gitopsStatus appstudiov1alpha1.GitOpsStatus
		gitToken     string
		wantURL      string
		wantBranch   string
		wantContext  string
		wantErr      bool
	}{
		{
			name: "gitops status processed as expected",
			gitopsStatus: appstudiov1alpha1.GitOpsStatus{
				RepositoryURL: "https://github.com/myrepo",
				Branch:        "notmain",
				Context:       "context",
			},
			gitToken:    "token",
			wantURL:     "https://token@github.com/myrepo",
			wantBranch:  "notmain",
			wantContext: "context",
		},
		{
			name: "gitops url is empty",
			gitopsStatus: appstudiov1alpha1.GitOpsStatus{
				RepositoryURL: "",
			},
			wantErr: true,
		},
		{
			name: "gitops branch and context not set",
			gitopsStatus: appstudiov1alpha1.GitOpsStatus{
				RepositoryURL: "https://github.com/myrepo",
			},
			gitToken:    "token",
			wantURL:     "https://token@github.com/myrepo",
			wantBranch:  "main",
			wantContext: "/",
		},
		{
			name: "gitops url parse err",
			gitopsStatus: appstudiov1alpha1.GitOpsStatus{
				RepositoryURL: "http://foo.com/?foo\nbar",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitopsURL, gitopsBranch, gitopsContext, err := ProcessGitOpsStatus(tt.gitopsStatus, tt.gitToken)
			if tt.wantErr && (err == nil) {
				t.Error("wanted error but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("got unexpected error %v", err)
			} else {
				assert.Equal(t, tt.wantURL, gitopsURL, "should be equal")
				assert.Equal(t, tt.wantBranch, gitopsBranch, "should be equal")
				assert.Equal(t, tt.wantContext, gitopsContext, "should be equal")
			}
		})
	}
}
