package eventfilters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type Deployment2Predicate struct {
}

func (r *Deployment2Predicate) Create(e event.CreateEvent) bool {
	return r.predict(e.Object)
}

func (r *Deployment2Predicate) Update(e event.UpdateEvent) bool {
	return r.predict(e.ObjectNew) || r.predict(e.ObjectOld)
}

func (r *Deployment2Predicate) Delete(e event.DeleteEvent) bool {
	return r.predict(e.Object)
}

func (r *Deployment2Predicate) Generic(e event.GenericEvent) bool {
	return r.predict(e.Object)
}

func (r *Deployment2Predicate) predict(obj metav1.Object) bool {
	return true
}
