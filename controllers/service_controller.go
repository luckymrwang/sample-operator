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
	"samp/enqueues"

	"github.com/go-logr/logr"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"samp/eventfilters"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch

func (r *ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	// name of our custom finalizer
	//myFinalizerName := "kubernetes"
	//// examine DeletionTimestamp to determine if object is under deletion
	//if svc.ObjectMeta.DeletionTimestamp.IsZero() {
	//	if !containsString(svc.ObjectMeta.Finalizers, myFinalizerName) {
	//		svc.ObjectMeta.Finalizers = append(svc.ObjectMeta.Finalizers, myFinalizerName)
	//		if err := r.Update(ctx, svc); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//	}
	//	svc.ObjectMeta.Finalizers = append(svc.ObjectMeta.Finalizers)
	//} else {
	//	// The object is being deleted
	//	if containsString(svc.ObjectMeta.Finalizers, myFinalizerName) {
	//		// our finalizer is present, so lets handle any external dependency
	//		// The object is being deleted
	//		if err := r.deleteDestinationRule(ctx, svc); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//
	//		if err := r.deleteVirtualService(ctx, svc); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//
	//		// remove our finalizer from the list and update it.
	//		svc.ObjectMeta.Finalizers = removeString(svc.ObjectMeta.Finalizers, myFinalizerName)
	//		if err := r.Update(ctx, svc); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//	}
	//
	//	// Stop reconciliation as the item is being deleted
	//	return ctrl.Result{}, nil
	//}

	//log.Info("---service list " + req.Name)
	//var svcs corev1.ServiceList
	//if err := r.List(ctx, &svcs, client.InNamespace(req.Namespace), client.MatchingLabels{"annotation": "icks.io_istio_workload_version"}); err != nil {
	//	return ctrl.Result{}, nil
	//}
	//for _, svc := range svcs.Items {
	//	fmt.Print("labels:", svc.Labels)
	//	fmt.Print("annotations:", svc.Annotations)
	//	fmt.Print("selector:", svc.Spec.Selector)
	//	fmt.Print("clusterIP:", svc.Spec.ClusterIP)
	//}
	_, ok := svc.Annotations[eventfilters.IstioAnnotation]
	if ok {
		if err := r.createDestinationRule(ctx, svc); err != nil {
			log.Error(err, "unable to create DestinationRule")
			// we'll ignore not-found errors, since they can't be fixed by an immediate
			// requeue (we'll need to wait for a new notification), and we can get them
			// on deleted requests.
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		if err := r.createVirtualService(ctx, svc); err != nil {
			log.Error(err, "unable to create VirtualService")
			// we'll ignore not-found errors, since they can't be fixed by an immediate
			// requeue (we'll need to wait for a new notification), and we can get them
			// on deleted requests.
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	} else {
		if err := r.deleteDestinationRule(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}

		if err := r.deleteVirtualService(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ServiceReconciler) createDestinationRule(ctx context.Context, svc *corev1.Service) error {
	log := r.Log.WithValues("DestinationRule", svc.Namespace)

	dr := new(v1beta1.DestinationRule)
	err := r.Get(ctx, client.ObjectKey{Namespace: svc.Namespace, Name: svc.Name}, dr)
	if apierrors.IsNotFound(err) {
		log.Info("could not find existing DestinationRule, creating one...")

		dr = r.buildDestinationRule(svc)
		if err := r.Create(ctx, dr); err != nil {
			log.Error(err, "failed to create DestinationRule resource")
			return err
		}

		r.Recorder.Eventf(dr, corev1.EventTypeNormal, "Created", "Created DestinationRule %s", dr.Name)
		log.Info("created DestinationRule resource")
		return nil
	}

	return err
}

func (r *ServiceReconciler) createVirtualService(ctx context.Context, svc *corev1.Service) error {
	log := r.Log.WithValues("VirtualService", svc.Namespace)

	vs := new(v1beta1.VirtualService)
	err := r.Get(ctx, client.ObjectKey{Namespace: svc.Namespace, Name: svc.Name}, vs)
	if apierrors.IsNotFound(err) {
		log.Info("could not find existing VirtualService, creating one...")

		vs = r.buildVirtualService(svc)
		if err := r.Create(ctx, vs); err != nil {
			log.Error(err, "failed to create VirtualService resource")
			return err
		}

		r.Recorder.Eventf(vs, corev1.EventTypeNormal, "Created", "Created VirtualService %s", vs.Name)
		log.Info("created VirtualService resource")
		return nil
	}

	return err
}

func (r *ServiceReconciler) buildDestinationRule(svc *corev1.Service) *v1beta1.DestinationRule {
	return &v1beta1.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{APIVersion: "v1", Controller: &[]bool{true}[0], BlockOwnerDeletion: &[]bool{true}[0], Kind: "Service", Name: svc.Name, UID: svc.UID},
			},
		},
		Spec: networkingv1beta1.DestinationRule{
			Host: svc.Name,
			Subsets: []*networkingv1beta1.Subset{
				{
					Name:   svc.Annotations[eventfilters.IstioAnnotation],
					Labels: map[string]string{"version": svc.Annotations[eventfilters.IstioAnnotation]},
				},
			},
		},
	}
}

func (r *ServiceReconciler) buildVirtualService(svc *corev1.Service) *v1beta1.VirtualService {
	return &v1beta1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: svc.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{APIVersion: "v1", Controller: &[]bool{true}[0], BlockOwnerDeletion: &[]bool{true}[0], Kind: "Service", Name: svc.Name, UID: svc.UID},
			},
		},
		Spec: networkingv1beta1.VirtualService{
			Hosts: []string{svc.Name},
			Http: []*networkingv1beta1.HTTPRoute{
				{
					Route: []*networkingv1beta1.HTTPRouteDestination{
						{
							Destination: &networkingv1beta1.Destination{Host: svc.Name, Subset: svc.Annotations[eventfilters.IstioAnnotation]},
						},
					},
				},
			},
		},
	}
}

func (r *ServiceReconciler) deleteDestinationRule(ctx context.Context, svc *corev1.Service) error {
	log := r.Log.WithValues("DestinationRule", svc.Namespace)

	dr := new(v1beta1.DestinationRule)
	err := r.Get(ctx, client.ObjectKey{Namespace: svc.Namespace, Name: svc.Name}, dr)
	if err == nil {
		if err = r.Delete(ctx, dr); err != nil {
			log.Error(err, "could not find DestinationRule "+svc.Name)
			r.Recorder.Eventf(dr, corev1.EventTypeWarning, "Finalize", "Deleted DestinationRule %q error", dr.Name)
			return err
		}

		r.Recorder.Eventf(dr, corev1.EventTypeNormal, "Deleted", "Deleted DestinationRule %q", dr.Name)
	}

	return nil
}

func (r *ServiceReconciler) deleteVirtualService(ctx context.Context, svc *corev1.Service) error {
	log := r.Log.WithValues("VirtualService", svc.Namespace)

	vs := new(v1beta1.VirtualService)
	err := r.Get(ctx, client.ObjectKey{Namespace: svc.Namespace, Name: svc.Name}, vs)
	if err == nil {
		if err = r.Delete(ctx, vs); err != nil {
			log.Error(err, "could not find VirtualService "+svc.Name)
			r.Recorder.Eventf(vs, corev1.EventTypeWarning, "Finalize", "Deleted VirtualService %q error", vs.Name)
			return err
		}

		r.Recorder.Eventf(vs, corev1.EventTypeNormal, "Deleted", "Deleted VirtualService %q", vs.Name)
	}

	return nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Service{}).
		Watches(
			&source.Kind{Type: &appv1.Deployment{}},
			&enqueues.EnqueueRequestForObject{}).
		WithEventFilter(&eventfilters.ServicePredicate{}).
		Complete(r)
}
