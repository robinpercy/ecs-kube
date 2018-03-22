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

// This file was automatically generated by lister-gen

package internalversion

import (
	ecskube "github.com/robinpercy/ecs-kube/pkg/apis/ecskube"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ECSDeploymentLister helps list ECSDeployments.
type ECSDeploymentLister interface {
	// List lists all ECSDeployments in the indexer.
	List(selector labels.Selector) (ret []*ecskube.ECSDeployment, err error)
	// ECSDeployments returns an object that can list and get ECSDeployments.
	ECSDeployments(namespace string) ECSDeploymentNamespaceLister
	ECSDeploymentListerExpansion
}

// eCSDeploymentLister implements the ECSDeploymentLister interface.
type eCSDeploymentLister struct {
	indexer cache.Indexer
}

// NewECSDeploymentLister returns a new ECSDeploymentLister.
func NewECSDeploymentLister(indexer cache.Indexer) ECSDeploymentLister {
	return &eCSDeploymentLister{indexer: indexer}
}

// List lists all ECSDeployments in the indexer.
func (s *eCSDeploymentLister) List(selector labels.Selector) (ret []*ecskube.ECSDeployment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*ecskube.ECSDeployment))
	})
	return ret, err
}

// ECSDeployments returns an object that can list and get ECSDeployments.
func (s *eCSDeploymentLister) ECSDeployments(namespace string) ECSDeploymentNamespaceLister {
	return eCSDeploymentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ECSDeploymentNamespaceLister helps list and get ECSDeployments.
type ECSDeploymentNamespaceLister interface {
	// List lists all ECSDeployments in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*ecskube.ECSDeployment, err error)
	// Get retrieves the ECSDeployment from the indexer for a given namespace and name.
	Get(name string) (*ecskube.ECSDeployment, error)
	ECSDeploymentNamespaceListerExpansion
}

// eCSDeploymentNamespaceLister implements the ECSDeploymentNamespaceLister
// interface.
type eCSDeploymentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ECSDeployments in the indexer for a given namespace.
func (s eCSDeploymentNamespaceLister) List(selector labels.Selector) (ret []*ecskube.ECSDeployment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*ecskube.ECSDeployment))
	})
	return ret, err
}

// Get retrieves the ECSDeployment from the indexer for a given namespace and name.
func (s eCSDeploymentNamespaceLister) Get(name string) (*ecskube.ECSDeployment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(ecskube.Resource("ecsdeployment"), name)
	}
	return obj.(*ecskube.ECSDeployment), nil
}