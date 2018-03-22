/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha1

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/robinpercy/ecs-kube/pkg/apis/ecskube"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ECSDeployment
// +k8s:openapi-gen=true
// +resource:path=ecsdeployments,strategy=ECSDeploymentStrategy
type ECSDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ECSDeploymentSpec   `json:"spec,omitempty"`
	Status ECSDeploymentStatus `json:"status,omitempty"`
}

// ECSDeploymentSpec defines the desired state of ECSDeployment
// TODO: rather than using a TaskSpec, should we watch configmaps references from it?
// See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ecs-taskdefinition.html
// Note: JSON names are capitalized so that ECS configs can be copy-pasted into the spec
type ECSDeploymentSpec struct {
	Task    ECSTask    `json:"Task"`
	Service ECSService `json:"Service"`
}

// ECSDeploymentStatus defines the observed state of ECSDeployment
type ECSDeploymentStatus struct {
	State string `json:"state,omitempty"`
}
type ECSTask struct {
	Type       string            `json:"Type"` // TODO: Handle
	Properties ECSTaskProperties `json:"Properties"`
}

type ECSService struct {
	Cluster string `json:"Cluster,omitempty"` // TODO: Handle
	//DeploymentConfiguration // TODO: Handle
	DesiredCount                  int32             `json:"DesiredCount"`
	HealthCheckGracePeriodSeconds int               `json:"HealthCheckGracePeriodSeconds,omitempty"` // TODO: Handle
	LaunchType                    string            `json:"LaunchType,omitempty"`                    // TODO: Handle
	LoadBalancers                 []ECSLoadBalancer `json:"LoadBalancers,omitempty"`
	// NetworkConfiguration // TODO: Handle
	PlacementConstraints ECSPlacement `json:"PlacementConstraints,omitempty"` // TODO: Handle
	// PlacementStrategies // TODO: Handle
	PlatformVersion string `json:"PlatformVersion,omitempty"` // TODO: Handle
	Role            string `json:"Role,omitempty"`            // TODO: Handle
	ServiceName     string `json:"ServiceName,omitempty"`     // TODO: Handle
	TaskDefinition  string `json:"TaskDefinition,omitempty"`  // TODO: Handle
}

type ECSLoadBalancer struct {
	ContainerName    string `json:"ContainerName"`
	ContainerPort    int32  `json:"ContainerPort"`
	LoadBalancerName string `json:"LoadBalancerName,omitempty"`
	TargetGroupArn   string `json:"TargetGroupArn,omitempty"`
}

type ECSTaskProperties struct {
	Volumes                 ECSVolumeHost  `json:"Volumes,omitempty"`                 // TODO: Handle
	CPU                     string         `json:"Cpu,omitempty"`                     // TODO: Handle
	ExecutionRoleARN        string         `json:"ExecutionRoleArn,omitempty"`        // TODO: Handle
	Family                  string         `json:"Family,omitempty"`                  // TODO: Handle
	Memory                  string         `json:"Memory,omitempty"`                  // TODO: Handle
	NetworkMode             string         `json:"NetworkMode,omitempty"`             // TODO: Handle
	PlacementConstraints    ECSPlacement   `json:"PlacementConstraints,omitempty"`    // TODO: Handle
	RequiresCompatibilities []string       `json:"RequiresCompatibilities,omitempty"` // TODO: Handle
	TaskRoleARN             string         `json:"TaskRoleArn,omitempty"`             // TODO: Handle
	ContainerDefinitions    []ECSContainer `json:"ContainerDefinitions"`
}
type ECSContainer struct {
	Command                []string            `json:"Command,omitempty"`               // TODO: Handle
	CPU                    int                 `json:"Cpu,omitempty"`                   // TODO: Handle
	DisableNetworking      bool                `json:"DisableNetworking,omitempty"`     // TODO: Handle
	DNSSearchDomains       []string            `json:"DnsSearchDomains,omitempty"`      // TODO: Handle
	DNSServers             []string            `json:"DnsServers,omitempty"`            // TODO: Handle
	DockerLabels           map[string]string   `json:"DockerLabels,omitempty"`          // TODO: Handle
	DockerSecurityOptions  string              `json:"DockerSecurityOptions,omitempty"` // TODO: Handle
	EntryPoint             string              `json:"EntryPoint,omitempty"`            // TODO: Handle
	Environment            []ECSKeyValuePair   `json:"Environment,omitempty"`           // TODO: Handle
	Essential              bool                `json:"Essential,omitempty"`             // TODO: Handle
	ExtraHosts             []ECSHostEntry      `json:"ExtraHosts,omitempty"`            // TODO: Handle
	Hostname               string              `json:"Hostname,omitempty"`              // TODO: Handle
	Image                  string              `json:"Image"`
	Links                  []string            `json:"Links,omitempty"`             // TODO: Handle
	LinuxParameters        ECSLinuxParameters  `json:"LinuxParameters,omitempty"`   // TODO: Handle
	LogConfiguration       ECSLogConfiguration `json:"LogConfiguration,omitempty"`  // TODO: Handle
	Memory                 int                 `json:"Memory,omitempty"`            // TODO: Handle
	MemoryReservation      int                 `json:"MemoryReservation,omitempty"` // TODO: Handle
	MountPoints            []ECSMountPoint     `json:"MountPoints,omitempty"`       // TODO: Handle
	Name                   string              `json:"Name"`
	PortMappings           []ECSPortMapping    `json:"PortMappings,omitempty"`           // TODO: Handle
	Privileged             bool                `json:"Privileged,omitempty"`             // TODO: Handle
	ReadonlyRootFilesystem bool                `json:"ReadonlyRootFilesystem,omitempty"` // TODO: Handle
	Ulimits                []ECSULimit         `json:"Ulimits,omitempty"`                // TODO: Handle
	User                   string              `json:"User,omitempty"`                   // TODO: Handle
	VolumesFrom            []ECSVolumeFrom     `json:"VolumesFrom,omitempty"`            // TODO: Handle
	WorkingDirectory       string              `json:"WorkingDirectory,omitempty"`       // TODO: Handle
}

type ECSKeyValuePair struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

type ECSHostEntry struct {
	Hostname  string `json:"Hostname"`
	IPAddress string `json:"IpAddress"`
}

type ECSLinuxParameters struct {
	Capabilities       []ECSKernelCapabilities `json:"Capabilities,omitempty"`
	Devices            []ECSDevice             `json:"Devices,omitempty"`
	InitProcessEnabled bool                    `json:"InitProcessEnabled,omitempty"`
}

type ECSKernelCapabilities struct {
	Add  []string `json:"Add,omitempty"`
	Drop []string `json:"Drop,omitempty"`
}

type ECSDevice struct {
	ContainerPath string   `json:"ContainerPath,omitempty"`
	HostPath      string   `json:"HostPath"`
	Permissions   []string `json:"Permissions,omitempty"`
}

type ECSLogConfiguration struct {
	LogDriver string            `json:"LogDriver"`
	Options   map[string]string `json:"Options,omitempty"`
}

type ECSMountPoint struct {
	ContainerPath string `json:"ContainerPath"`
	SourceVolume  string `json:"SourceVolume"`
	ReadOnly      bool   `json:"ReadOnly,omitempty"`
}

type ECSPortMapping struct {
	ContainerPort int32  `json:"ContainerPort"`
	HostPort      int32  `json:"HostPort,omitempty"`
	Protocol      string `json:"Protocol,omitempty"`
}

type ECSULimit struct {
	HardLimit int    `json:"HardLimit"`
	Name      string `json:"Name,omitempty"`
	SoftLimit int    `json:"SoftLimit"`
}

type ECSVolumeFrom struct {
	SourceContainer string `json:"SourceContainer"`
	ReadOnly        bool   `json:"ReadOnly,omitempty"`
}
type ECSVolume struct {
	Name string        `json:"Name"`
	Host ECSVolumeHost `json:"Host,omitempty"`
}

type ECSVolumeHost struct {
	SourcePath string `json:"SourcePath,omitempty"`
}

type ECSPlacement struct {
	Type       string `json:"Type"`
	Expression string `json:"Expression,omitempty"`
}

// Validate checks that an instance of ECSDeployment is well formed
func (ECSDeploymentStrategy) Validate(ctx request.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*ecskube.ECSDeployment)
	log.Printf("Validating fields for ECSDeployment %s\n", o.Name)
	errors := field.ErrorList{}
	// perform validation here and add to errors using field.Invalid
	return errors
}

// DefaultingFunction sets default ECSDeployment field values
func (ECSDeploymentSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*ECSDeployment)
	// set default field values here
	log.Printf("Defaulting fields for ECSDeployment %s\n", obj.Name)
}
