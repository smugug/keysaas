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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	keysaascontrollerv1 "github.com/smugug/keysaas/pkg/apis/keysaascontroller/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// KeysaasLister helps list Keysaases.
// All objects returned here must be treated as read-only.
type KeysaasLister interface {
	// List lists all Keysaases in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*keysaascontrollerv1.Keysaas, err error)
	// Keysaases returns an object that can list and get Keysaases.
	Keysaases(namespace string) KeysaasNamespaceLister
	KeysaasListerExpansion
}

// keysaasLister implements the KeysaasLister interface.
type keysaasLister struct {
	listers.ResourceIndexer[*keysaascontrollerv1.Keysaas]
}

// NewKeysaasLister returns a new KeysaasLister.
func NewKeysaasLister(indexer cache.Indexer) KeysaasLister {
	return &keysaasLister{listers.New[*keysaascontrollerv1.Keysaas](indexer, keysaascontrollerv1.Resource("keysaas"))}
}

// Keysaases returns an object that can list and get Keysaases.
func (s *keysaasLister) Keysaases(namespace string) KeysaasNamespaceLister {
	return keysaasNamespaceLister{listers.NewNamespaced[*keysaascontrollerv1.Keysaas](s.ResourceIndexer, namespace)}
}

// KeysaasNamespaceLister helps list and get Keysaases.
// All objects returned here must be treated as read-only.
type KeysaasNamespaceLister interface {
	// List lists all Keysaases in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*keysaascontrollerv1.Keysaas, err error)
	// Get retrieves the Keysaas from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*keysaascontrollerv1.Keysaas, error)
	KeysaasNamespaceListerExpansion
}

// keysaasNamespaceLister implements the KeysaasNamespaceLister
// interface.
type keysaasNamespaceLister struct {
	listers.ResourceIndexer[*keysaascontrollerv1.Keysaas]
}
