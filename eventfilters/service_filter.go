package eventfilters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var IstioAnnotation = "icks.io_istio_workload_version"

type ServicePredicate struct {
}

func (r *ServicePredicate) Create(e event.CreateEvent) bool {
	return r.predict(e.Object)
}

func (r *ServicePredicate) Update(e event.UpdateEvent) bool {
	return r.predict(e.ObjectNew) || r.predict(e.ObjectOld)
}

func (r *ServicePredicate) Delete(e event.DeleteEvent) bool {
	return r.predict(e.Object)
}

func (r *ServicePredicate) Generic(e event.GenericEvent) bool {
	return r.predict(e.Object)
}

func (r *ServicePredicate) predict(obj metav1.Object) bool {
	//return containsAnnotation(obj, IstioAnnotation)
	if obj.GetNamespace() == "testing" {
		return true
	}
	return false
}

func containsAnnotation(obj metav1.Object, annotation string) bool {
	annotations := obj.GetAnnotations()
	_, ok := annotations[annotation]

	return ok
}
