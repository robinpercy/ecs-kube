
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


package ecsdeployment

import (
	"log"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	"github.com/robinpercy/ecs-kube/pkg/apis/ecskube/v1alpha1"
	"github.com/robinpercy/ecs-kube/pkg/controller/sharedinformers"
	listers "github.com/robinpercy/ecs-kube/pkg/client/listers_generated/ecskube/v1alpha1"
)

// +controller:group=ecskube,version=v1alpha1,kind=ECSDeployment,resource=ecsdeployments
type ECSDeploymentControllerImpl struct {
	builders.DefaultControllerFns

	// lister indexes properties about ECSDeployment
	lister listers.ECSDeploymentLister
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *ECSDeploymentControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {
	// Use the lister for indexing ecsdeployments labels
	c.lister = arguments.GetSharedInformers().Factory.Ecskube().V1alpha1().ECSDeployments().Lister()
}

// Reconcile handles enqueued messages
func (c *ECSDeploymentControllerImpl) Reconcile(u *v1alpha1.ECSDeployment) error {
	// Implement controller logic here
	log.Printf("Running reconcile ECSDeployment for %s\n", u.Name)
	return nil
}

func (c *ECSDeploymentControllerImpl) Get(namespace, name string) (*v1alpha1.ECSDeployment, error) {
	return c.lister.ECSDeployments(namespace).Get(name)
}
