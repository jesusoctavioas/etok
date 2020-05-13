// Code generated by go generate; DO NOT EDIT.
package output

import (
	"context"

	v1alpha1 "github.com/leg100/stok/pkg/apis/stok/v1alpha1"
	"github.com/leg100/stok/pkg/controller/command"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_output")

// Add creates a new Output Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOutput{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("output-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Output
	err = c.Watch(&source.Kind{Type: &v1alpha1.Output{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Output
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Output{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to resource Workspace and requeue the associated Outputs
	err = c.Watch(&source.Kind{Type: &v1alpha1.Workspace{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
			rc := r.(*ReconcileOutput)
			cmdList := &v1alpha1.OutputList{}
			err = rc.client.List(context.TODO(), cmdList, client.InNamespace(a.Meta.GetNamespace()), client.MatchingLabels{
				"workspace": a.Meta.GetName(),
			})
			if err != nil {
				return []reconcile.Request{}
			}

			rr := []reconcile.Request{}
			for _, cmd := range cmdList.Items {
				rr = append(rr, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      cmd.GetName(),
						Namespace: cmd.GetNamespace(),
					},
				})
			}
			return rr
		}),
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileOutput implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileOutput{}

// ReconcileOutput reconciles a Output object
type ReconcileOutput struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Output object and makes changes based on the state read
// and what is in the Output.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOutput) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.V(1).Info("Reconciling Output")

	return command.Reconcile(r.client, r.scheme, request, &v1alpha1.Output{}, "output")
}