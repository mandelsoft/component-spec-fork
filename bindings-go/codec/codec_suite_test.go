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

package codec_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v2 "github.com/gardener/component-spec/bindings-go/apis/v2"
	"github.com/gardener/component-spec/bindings-go/codec"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Test Suite")
}

var _ = Describe("serializer", func() {

	It("should decode a simple component", func() {
		data, err := ioutil.ReadFile("../../language-independent/test-resources/component_descriptor_v2.yaml")
		Expect(err).ToNot(HaveOccurred())

		var comp v2.ComponentDescriptor
		err = codec.Decode(data, &comp)
		Expect(err).ToNot(HaveOccurred())

		Expect(comp.Name).To(Equal("github.com/gardener/gardener"))
		Expect(comp.Version).To(Equal("v1.7.2"))
		Expect(comp.ExternalResources).To(HaveLen(1))

		extDep := comp.ExternalResources[0]
		Expect(extDep.Name).To(Equal("grafana"))
		Expect(extDep.Version).To(Equal("7.0.3"))
		Expect(extDep.GetType()).To(Equal(v2.OCIImageType))
		Expect(extDep.Access).To(BeAssignableToTypeOf(&v2.OCIRegistryAccess{}))

		ociAccess := extDep.Access.(*v2.OCIRegistryAccess)
		Expect(ociAccess.ImageReference).To(Equal("registry-1.docker.io/grafana/grafana/7.0.3"))
	})

})
