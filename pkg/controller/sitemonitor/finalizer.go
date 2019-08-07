package sitemonitor

import (
	"context"

	sitemonitorv1alpha1 "wen/site-monitor-operator/pkg/apis/sitemonitor/v1alpha1"

	"github.com/go-logr/logr"
)

const sitemonitorFinalizer = "finalizer.sitemonitor.haodai.com"

func (r *ReconcileSiteMonitor) finalizeSiteMonitor(reqLogger logr.Logger, m *sitemonitorv1alpha1.SiteMonitor) error {
	// TODO(user): Add the cleanup steps that the operator
	// needs to do before the CR can be deleted. Examples
	// of finalizers include performing backups and deleting
	// resources that are not owned by this CR, like a PVC.

	err := S.DeleteMonitor(m.Spec.TaskName)
	if err != nil {
		reqLogger.Error(err, "Failed to delete sitemonitor with finalizer")
		return err
	}
	reqLogger.Info("Successfully finalized sitemonitor")
	return nil
}

func (r *ReconcileSiteMonitor) addFinalizer(reqLogger logr.Logger, m *sitemonitorv1alpha1.SiteMonitor) error {
	reqLogger.Info("Adding Finalizer for the Memcached")
	m.SetFinalizers(append(m.GetFinalizers(), sitemonitorFinalizer))

	// Update CR
	err := r.client.Update(context.TODO(), m)
	if err != nil {
		reqLogger.Error(err, "Failed to update sitemonitor with finalizer")
		return err
	}
	return nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func remove(list []string, s string) []string {
	for i, v := range list {
		if v == s {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}
