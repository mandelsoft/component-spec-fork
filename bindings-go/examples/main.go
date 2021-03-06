// Copyright 2020 Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/codec"
)

func main() {
	data := []byte(`
meta:
  schemaVersion: 'v2'

component:
  name: 'github.com/gardener/gardener'
  version: 'v1.7.2'

  provider: internal

  sources: []
  references: []

  localResources:
  - name: 'apiserver'
    version: 'v1.7.2'
    type: 'ociImage'
    access:
      type: 'ociRegistry'
      imageReference: 'eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2'

  externalResources:
  - name: 'hyperkube'
    version: 'v1.16.4'
    type: 'ociImage'
    access:
      type: 'ociRegistry'
      imageReference: 'k8s.gcr.io/hyperkube:v1.16.4'
`)

	component := &v2.ComponentDescriptor{}
	err := codec.Decode(data, component)
	check(err)

	// get a specific local resource
	res, err := component.GetLocalResource(v2.OCIImageType, "apiserver", "v1.7.2")
	check(err)
	fmt.Printf("%v\n", res)

	// get a specific external resource
	res, err = component.GetExternalResource(v2.OCIImageType, "hyperkube", "v1.16.4")
	check(err)
	fmt.Printf("%v\n", res)

	// get the access for a resource
	// known types implement the TypedObjectAccessor interface and can be cast to the specific type.
	ociAccess := res.Access.(*v2.OCIRegistryAccess)
	fmt.Println(ociAccess.ImageReference) // prints: eu.gcr.io/gardener-project/gardener/apiserver:v1.7.2
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
