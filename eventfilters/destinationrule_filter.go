package eventfilters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type DestinationRulePredicate struct {
}

func (r *DestinationRulePredicate) Create(e event.CreateEvent) bool {
	return r.predict(e.Meta, "CREATE")
}

func (r *DestinationRulePredicate) Update(e event.UpdateEvent) bool {
	return r.predict(e.MetaNew, "UPDATE") || r.predict(e.MetaOld, "UPDATE")
}

func (r *DestinationRulePredicate) Delete(e event.DeleteEvent) bool {
	return r.predict(e.Meta, "DELETE")
}

func (r *DestinationRulePredicate) Generic(e event.GenericEvent) bool {
	return r.predict(e.Meta, "GENERIC")
}

func (r *DestinationRulePredicate) predict(obj metav1.Object, event string) bool {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["event"] = event
	obj.SetAnnotations(annotations)
	return true
}
