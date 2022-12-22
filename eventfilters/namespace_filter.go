package eventfilters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type NamespacePredicate struct {
}

func (r *NamespacePredicate) Create(e event.CreateEvent) bool {
	return r.predict(e.Object, "CREATE")
}

func (r *NamespacePredicate) Update(e event.UpdateEvent) bool {
	return r.predict(e.ObjectNew, "UPDATE") || r.predict(e.ObjectOld, "UPDATE")
}

func (r *NamespacePredicate) Delete(e event.DeleteEvent) bool {
	return r.predict(e.Object, "DELETE")
}

func (r *NamespacePredicate) Generic(e event.GenericEvent) bool {
	return r.predict(e.Object, "GENERIC")
}

func (r *NamespacePredicate) predict(obj metav1.Object, event string) bool {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["event"] = event
	obj.SetAnnotations(annotations)
	return true
}
