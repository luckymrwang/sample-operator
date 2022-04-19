/*


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

package v1

import (
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	ExtraSchemeBuilder = &scheme.Builder{GroupVersion: schema.GroupVersion{Group: "networking.istio.io", Version: "v1beta1"}}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddExtraToScheme = ExtraSchemeBuilder.AddToScheme
)

func init() {
	ExtraSchemeBuilder.Register(
		&v1beta1.DestinationRule{},
		&v1beta1.DestinationRuleList{},
		&v1beta1.VirtualService{},
		&v1beta1.VirtualServiceList{},
	)
}
