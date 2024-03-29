package eventfilters

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type DeploymentPredicate struct {
}

func (r *DeploymentPredicate) Create(e event.CreateEvent) bool {
	return r.predict(e.Object, "CREATE")
}

func (r *DeploymentPredicate) Update(e event.UpdateEvent) bool {
	return r.predict(e.ObjectNew, "UPDATE") || r.predict(e.ObjectOld, "UPDATE")
}

func (r *DeploymentPredicate) Delete(e event.DeleteEvent) bool {
	return r.predict(e.Object, "DELETE")
}

func (r *DeploymentPredicate) Generic(e event.GenericEvent) bool {
	return r.predict(e.Object, "GENERIC")
}

func (r *DeploymentPredicate) predict(obj metav1.Object, event string) bool {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations["event"] = event
	obj.SetAnnotations(annotations)
	return true
}
