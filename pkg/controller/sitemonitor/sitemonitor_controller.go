package sitemonitor

import (
	"context"
	"reflect"

	sitemonitorv1alpha1 "wen/site-monitor-operator/pkg/apis/sitemonitor/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_sitemonitor")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SiteMonitor Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSiteMonitor{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sitemonitor-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SiteMonitor
	err = c.Watch(&source.Kind{Type: &sitemonitorv1alpha1.SiteMonitor{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SiteMonitor
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sitemonitorv1alpha1.SiteMonitor{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSiteMonitor implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSiteMonitor{}

// ReconcileSiteMonitor reconciles a SiteMonitor object
type ReconcileSiteMonitor struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SiteMonitor object and makes changes based on the state read
// and what is in the SiteMonitor.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSiteMonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SiteMonitor")

	// Fetch the SiteMonitor instance
	instance := &sitemonitorv1alpha1.SiteMonitor{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// spew.Dump("instance", instance)

	// Define a new site monitor, it will fetch,compare,update, or create if not exist
	err = updateSiteMonitorForCR(instance)
	if err != nil {
		reqLogger.Error(err, "newSiteMonitorForCR", instance.GetName())
		return reconcile.Result{}, err
	}

	// Check monitor status
	state, err := S.GetTaskState(instance.Spec.TaskName)
	if err != nil {
		reqLogger.Error(err, "GetTaskState", instance.GetName())
		return reconcile.Result{}, err
	}

	alertstate, err := S.GetAlertState(instance.Spec.TaskName)
	if err != nil {
		reqLogger.Error(err, "GetAlertState", instance.GetName())
		return reconcile.Result{}, err
	}

	// Ensure the monitor status is the same as the spec
	if instance.Spec.TaskState != state {
		// we don't update state spec, so no need this?

		// Spec updated - return and requeue
		return reconcile.Result{Requeue: true}, nil
	}

	var updatestatus bool
	// Update status.Nodes if needed
	if !reflect.DeepEqual(state, instance.Status.TaskState) {
		instance.Status.TaskState = state
		updatestatus = true
	}
	if !reflect.DeepEqual(alertstate, instance.Status.AlertState) {
		instance.Status.AlertState = alertstate
		updatestatus = true
	}
	if updatestatus {
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update sitemonitor status.")
			return reconcile.Result{}, err
		}
	}

	// Check if the SiteMonitor instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isSiteMonitorMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isSiteMonitorMarkedToBeDeleted {
		if contains(instance.GetFinalizers(), sitemonitorFinalizer) {
			// Run finalization logic for instanceFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeSiteMonitor(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}

			// Remove instanceFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			instance.SetFinalizers(remove(instance.GetFinalizers(), sitemonitorFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !contains(instance.GetFinalizers(), sitemonitorFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// // Set SiteMonitor instance as the owner and controller
	// if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Check if this Pod already exists
	// found := &corev1.Pod{}
	// err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	// if err != nil && errors.IsNotFound(err) {
	// 	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	// 	err = r.client.Create(context.TODO(), pod)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}

	// 	// Pod created successfully - don't requeue
	// 	return reconcile.Result{}, nil
	// } else if err != nil {
	// 	return reconcile.Result{}, err
	// }

	// // Pod already exists - don't requeue
	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

	return reconcile.Result{}, nil
}

// // newPodForCR returns a busybox pod with the same name/namespace as the cr
// func newPodForCR(cr *sitemonitorv1alpha1.SiteMonitor) *corev1.Pod {
// 	labels := map[string]string{
// 		"app": cr.Name,
// 	}
// 	return &corev1.Pod{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      cr.Name + "-pod",
// 			Namespace: cr.Namespace,
// 			Labels:    labels,
// 		},
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{
// 				{
// 					Name:    "busybox",
// 					Image:   "busybox",
// 					Command: []string{"sleep", "3600"},
// 				},
// 			},
// 		},
// 	}
// }
