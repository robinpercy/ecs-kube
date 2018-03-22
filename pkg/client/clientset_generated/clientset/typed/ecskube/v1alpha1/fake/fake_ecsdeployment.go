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
package fake

import (
	v1alpha1 "github.com/robinpercy/ecs-kube/pkg/apis/ecskube/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeECSDeployments implements ECSDeploymentInterface
type FakeECSDeployments struct {
	Fake *FakeEcskubeV1alpha1
	ns   string
}

var ecsdeploymentsResource = schema.GroupVersionResource{Group: "ecskube.ecskube", Version: "v1alpha1", Resource: "ecsdeployments"}

var ecsdeploymentsKind = schema.GroupVersionKind{Group: "ecskube.ecskube", Version: "v1alpha1", Kind: "ECSDeployment"}

// Get takes name of the eCSDeployment, and returns the corresponding eCSDeployment object, and an error if there is any.
func (c *FakeECSDeployments) Get(name string, options v1.GetOptions) (result *v1alpha1.ECSDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(ecsdeploymentsResource, c.ns, name), &v1alpha1.ECSDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ECSDeployment), err
}

// List takes label and field selectors, and returns the list of ECSDeployments that match those selectors.
func (c *FakeECSDeployments) List(opts v1.ListOptions) (result *v1alpha1.ECSDeploymentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(ecsdeploymentsResource, ecsdeploymentsKind, c.ns, opts), &v1alpha1.ECSDeploymentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ECSDeploymentList{}
	for _, item := range obj.(*v1alpha1.ECSDeploymentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested eCSDeployments.
func (c *FakeECSDeployments) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(ecsdeploymentsResource, c.ns, opts))

}

// Create takes the representation of a eCSDeployment and creates it.  Returns the server's representation of the eCSDeployment, and an error, if there is any.
func (c *FakeECSDeployments) Create(eCSDeployment *v1alpha1.ECSDeployment) (result *v1alpha1.ECSDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(ecsdeploymentsResource, c.ns, eCSDeployment), &v1alpha1.ECSDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ECSDeployment), err
}

// Update takes the representation of a eCSDeployment and updates it. Returns the server's representation of the eCSDeployment, and an error, if there is any.
func (c *FakeECSDeployments) Update(eCSDeployment *v1alpha1.ECSDeployment) (result *v1alpha1.ECSDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(ecsdeploymentsResource, c.ns, eCSDeployment), &v1alpha1.ECSDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ECSDeployment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeECSDeployments) UpdateStatus(eCSDeployment *v1alpha1.ECSDeployment) (*v1alpha1.ECSDeployment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(ecsdeploymentsResource, "status", c.ns, eCSDeployment), &v1alpha1.ECSDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ECSDeployment), err
}

// Delete takes name of the eCSDeployment and deletes it. Returns an error if one occurs.
func (c *FakeECSDeployments) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(ecsdeploymentsResource, c.ns, name), &v1alpha1.ECSDeployment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeECSDeployments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(ecsdeploymentsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ECSDeploymentList{})
	return err
}

// Patch applies the patch and returns the patched eCSDeployment.
func (c *FakeECSDeployments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ECSDeployment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(ecsdeploymentsResource, c.ns, name, data, subresources...), &v1alpha1.ECSDeployment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ECSDeployment), err
}
