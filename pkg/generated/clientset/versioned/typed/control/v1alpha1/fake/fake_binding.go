/*
Copyright The KubeStellar Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha1 "github.com/kubestellar/kubestellar/api/control/v1alpha1"
)

// FakeBindings implements BindingInterface
type FakeBindings struct {
	Fake *FakeControlV1alpha1
}

var bindingsResource = v1alpha1.SchemeGroupVersion.WithResource("bindings")

var bindingsKind = v1alpha1.SchemeGroupVersion.WithKind("Binding")

// Get takes name of the binding, and returns the corresponding binding object, and an error if there is any.
func (c *FakeBindings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(bindingsResource, name), &v1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Binding), err
}

// List takes label and field selectors, and returns the list of Bindings that match those selectors.
func (c *FakeBindings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.BindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(bindingsResource, bindingsKind, opts), &v1alpha1.BindingList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.BindingList{ListMeta: obj.(*v1alpha1.BindingList).ListMeta}
	for _, item := range obj.(*v1alpha1.BindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested bindings.
func (c *FakeBindings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(bindingsResource, opts))
}

// Create takes the representation of a binding and creates it.  Returns the server's representation of the binding, and an error, if there is any.
func (c *FakeBindings) Create(ctx context.Context, binding *v1alpha1.Binding, opts v1.CreateOptions) (result *v1alpha1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(bindingsResource, binding), &v1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Binding), err
}

// Update takes the representation of a binding and updates it. Returns the server's representation of the binding, and an error, if there is any.
func (c *FakeBindings) Update(ctx context.Context, binding *v1alpha1.Binding, opts v1.UpdateOptions) (result *v1alpha1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(bindingsResource, binding), &v1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Binding), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeBindings) UpdateStatus(ctx context.Context, binding *v1alpha1.Binding, opts v1.UpdateOptions) (*v1alpha1.Binding, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(bindingsResource, "status", binding), &v1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Binding), err
}

// Delete takes name of the binding and deletes it. Returns an error if one occurs.
func (c *FakeBindings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(bindingsResource, name, opts), &v1alpha1.Binding{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBindings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(bindingsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.BindingList{})
	return err
}

// Patch applies the patch and returns the patched binding.
func (c *FakeBindings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Binding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(bindingsResource, name, pt, data, subresources...), &v1alpha1.Binding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Binding), err
}
