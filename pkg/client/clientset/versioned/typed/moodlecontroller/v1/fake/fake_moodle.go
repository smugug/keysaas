/*
Copyright The Kubernetes Authors.

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
	moodlecontrollerv1 "github.com/cloud-ark/kubeplus-operators/moodle/pkg/apis/moodlecontroller/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMoodles implements MoodleInterface
type FakeMoodles struct {
	Fake *FakeMoodlecontrollerV1
	ns   string
}

var moodlesResource = schema.GroupVersionResource{Group: "moodlecontroller.kubeplus", Version: "v1", Resource: "moodles"}

var moodlesKind = schema.GroupVersionKind{Group: "moodlecontroller.kubeplus", Version: "v1", Kind: "Moodle"}

// Get takes name of the moodle, and returns the corresponding moodle object, and an error if there is any.
func (c *FakeMoodles) Get(name string, options v1.GetOptions) (result *moodlecontrollerv1.Moodle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(moodlesResource, c.ns, name), &moodlecontrollerv1.Moodle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*moodlecontrollerv1.Moodle), err
}

// List takes label and field selectors, and returns the list of Moodles that match those selectors.
func (c *FakeMoodles) List(opts v1.ListOptions) (result *moodlecontrollerv1.MoodleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(moodlesResource, moodlesKind, c.ns, opts), &moodlecontrollerv1.MoodleList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &moodlecontrollerv1.MoodleList{ListMeta: obj.(*moodlecontrollerv1.MoodleList).ListMeta}
	for _, item := range obj.(*moodlecontrollerv1.MoodleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested moodles.
func (c *FakeMoodles) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(moodlesResource, c.ns, opts))

}

// Create takes the representation of a moodle and creates it.  Returns the server's representation of the moodle, and an error, if there is any.
func (c *FakeMoodles) Create(moodle *moodlecontrollerv1.Moodle) (result *moodlecontrollerv1.Moodle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(moodlesResource, c.ns, moodle), &moodlecontrollerv1.Moodle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*moodlecontrollerv1.Moodle), err
}

// Update takes the representation of a moodle and updates it. Returns the server's representation of the moodle, and an error, if there is any.
func (c *FakeMoodles) Update(moodle *moodlecontrollerv1.Moodle) (result *moodlecontrollerv1.Moodle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(moodlesResource, c.ns, moodle), &moodlecontrollerv1.Moodle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*moodlecontrollerv1.Moodle), err
}

// Delete takes name of the moodle and deletes it. Returns an error if one occurs.
func (c *FakeMoodles) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(moodlesResource, c.ns, name), &moodlecontrollerv1.Moodle{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMoodles) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(moodlesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &moodlecontrollerv1.MoodleList{})
	return err
}

// Patch applies the patch and returns the patched moodle.
func (c *FakeMoodles) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *moodlecontrollerv1.Moodle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(moodlesResource, c.ns, name, pt, data, subresources...), &moodlecontrollerv1.Moodle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*moodlecontrollerv1.Moodle), err
}
