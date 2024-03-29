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
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"my.domain/guestbook/eventfilters"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	if req.Name != "foo" {
		return ctrl.Result{}, nil
	}
	log := r.Log.WithValues("namespace", req.NamespacedName)

	instance := new(corev1.Namespace)
	if err := r.Get(ctx, client.ObjectKey{Name: req.Name}, instance); err != nil {
		log.Info("unable to fetch namespace")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Println("namespace sleep 2s")
	time.Sleep(2 * time.Second)

	fmt.Println(instance.Annotations)
	annotations := instance.Annotations
	switch annotations["event"] {
	case "CREATE":
		fmt.Println("CREATING >>>>>>>>")
	case "UPDATE":
		fmt.Println("UPDATING >>>>>>>>")
	default:
	}

	return ctrl.Result{}, nil
}

func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		WithEventFilter(&eventfilters.NamespacePredicate{}).
		Complete(r)
}
