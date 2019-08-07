package sitemonitor

import (
	"fmt"
	sitemonitorv1alpha1 "wen/site-monitor-operator/pkg/apis/sitemonitor/v1alpha1"
	"wen/site-monitor-operator/pkg/monitor"

	prettyjson "github.com/hokaccha/go-prettyjson"
)

var S *monitor.Server

func updateSiteMonitorForCR(cr *sitemonitorv1alpha1.SiteMonitor) (err error) {
	log.Info("creating monitor:", "name", cr.GetName())
	pretty("sitemonitor:", cr)
	err = S.UpdateMonitor(cr.Spec.SiteMonitor)
	if err != nil {
		err = fmt.Errorf("UpdateMonitor err:%v", err)
		return
	}
	return
}

func pretty(prefix, a interface{}) {
	out, _ := prettyjson.Marshal(a)
	fmt.Printf("%v: %s\n", prefix, out)
}
