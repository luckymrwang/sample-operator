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

package controllers

import (
	"context"
	"fmt"
	"samp/eventfilters"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/rest"

	"github.com/go-logr/logr"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceReconciler reconciles a Service object
type DestinationRuleReconciler struct {
	*rest.Config
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch

func (r *DestinationRuleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("service", req.NamespacedName)

	svc := new(corev1.Service)
	if err := r.Get(ctx, req.NamespacedName, svc); err != nil {
		log.Info("unable to fetch Service")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *DestinationRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	apiextensionsClient, _ := apiextensionsclientset.NewForConfig(r.Config)
	_, err := apiextensionsClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), "destinationrules.networking.istio.io", metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			fmt.Println("not found")
			return err
		}
		fmt.Println(err)
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.DestinationRule{}).
		WithEventFilter(&eventfilters.DestinationRulePredicate{}).
		Complete(r)
}
