/*
Copyright 2018 Pressinfra SRL

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

package v1alpha1

import (
	v1alpha1 "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
	scheme "github.com/presslabs/wordpress-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// WordpressesGetter has a method to return a WordpressInterface.
// A group's client should implement this interface.
type WordpressesGetter interface {
	Wordpresses(namespace string) WordpressInterface
}

// WordpressInterface has methods to work with Wordpress resources.
type WordpressInterface interface {
	Create(*v1alpha1.Wordpress) (*v1alpha1.Wordpress, error)
	Update(*v1alpha1.Wordpress) (*v1alpha1.Wordpress, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Wordpress, error)
	List(opts v1.ListOptions) (*v1alpha1.WordpressList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Wordpress, err error)
	WordpressExpansion
}

// wordpresses implements WordpressInterface
type wordpresses struct {
	client rest.Interface
	ns     string
}

// newWordpresses returns a Wordpresses
func newWordpresses(c *WordpressV1alpha1Client, namespace string) *wordpresses {
	return &wordpresses{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the wordpress, and returns the corresponding wordpress object, and an error if there is any.
func (c *wordpresses) Get(name string, options v1.GetOptions) (result *v1alpha1.Wordpress, err error) {
	result = &v1alpha1.Wordpress{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("wordpresses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Wordpresses that match those selectors.
func (c *wordpresses) List(opts v1.ListOptions) (result *v1alpha1.WordpressList, err error) {
	result = &v1alpha1.WordpressList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("wordpresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested wordpresses.
func (c *wordpresses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("wordpresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a wordpress and creates it.  Returns the server's representation of the wordpress, and an error, if there is any.
func (c *wordpresses) Create(wordpress *v1alpha1.Wordpress) (result *v1alpha1.Wordpress, err error) {
	result = &v1alpha1.Wordpress{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("wordpresses").
		Body(wordpress).
		Do().
		Into(result)
	return
}

// Update takes the representation of a wordpress and updates it. Returns the server's representation of the wordpress, and an error, if there is any.
func (c *wordpresses) Update(wordpress *v1alpha1.Wordpress) (result *v1alpha1.Wordpress, err error) {
	result = &v1alpha1.Wordpress{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("wordpresses").
		Name(wordpress.Name).
		Body(wordpress).
		Do().
		Into(result)
	return
}

// Delete takes name of the wordpress and deletes it. Returns an error if one occurs.
func (c *wordpresses) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("wordpresses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *wordpresses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("wordpresses").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched wordpress.
func (c *wordpresses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Wordpress, err error) {
	result = &v1alpha1.Wordpress{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("wordpresses").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
