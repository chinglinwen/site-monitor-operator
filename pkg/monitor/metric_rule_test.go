package monitor

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMetricRuleList(t *testing.T) {
	// fmt.Println("s", s)
	r, err := s.MetricRuleList()
	if err != nil {
		t.Error("err", err)
		return
	}
	fmt.Printf("got %v alarms\n", len(r.Alarms.Alarm))
	for i, v := range r.Alarms.Alarm {
		if i == 0 || v.RuleName == "test" || v.RuleName == "baidu.com" {
			// spew.Dump(v)
			pretty("alarm", v)
		}
	}
	// 	fmt.Printf("%v->%v\n", v.TaskName, v.Address)
	// 	// spew.Dump("got",v)
	// }
	// fmt.Printf("r: %v\n", r)
}

func TestGetMetricRule(t *testing.T) {
	// fmt.Println("s", s)
	r, err := s.GetMetricRule("baidu.com")
	if err != nil {
		t.Error("err", err)
		return
	}
	spew.Dump(r)
}

/*
   /home/wen/git/site-monitor-operator/pkg/monitor/metric_rule_test.go:36: err SDK.ServerError
       ErrorCode: MissingGroupId
       Recommend: https://error-center.aliyun.com/status/search?Keyword=MissingGroupId&source=PopGw
       RequestId: 3BC703B0-6514-44B5-821E-C49ECC22F6C5
	   Message: GroupId is mandatory for this action.

    /home/wen/git/site-monitor-operator/pkg/monitor/metric_rule_test.go:36: err SDK.ServerError
        ErrorCode: 404
        Recommend:
        RequestId: 1C87B9A9-B198-4E43-AAFB-0BAF88D2ED69
		Message: The resource 1 not found.

https://help.aliyun.com/document_detail/115071.html
changed to use PutResourceMetricRule

    /home/wen/git/site-monitor-operator/pkg/monitor/metric_rule_test.go:42: err SDK.ServerError
        ErrorCode: MissingRuleId
        Recommend: https://error-center.aliyun.com/status/search?Keyword=MissingRuleId&source=PopGw
        RequestId: AAC2D4AC-3070-4887-984A-66652C696E17
		Message: RuleId is mandatory for this action.

req.RuleId = "b47607e8-9826-4ba0-8f04-d27c8ca0670b_StatusCode"

    /home/wen/git/site-monitor-operator/pkg/monitor/metric_rule_test.go:51: err SDK.ServerError
        ErrorCode: MissingResources
        Recommend: https://error-center.aliyun.com/status/search?Keyword=MissingResources&source=PopGw
        RequestId: EEF8C6D7-AC4A-4003-A1BB-2332492F7B09
		Message: Resources is mandatory for this action.

req.Resources = "[{\"taskId\":\"31bc7843-afd1-44ab-b9ed-67e5314d6095\"}]"

    /home/wen/git/site-monitor-operator/pkg/monitor/metric_rule_test.go:61: err SDK.ServerError
        ErrorCode: MissingContactGroups
        Recommend: https://error-center.aliyun.com/status/search?Keyword=MissingContactGroups&source=PopGw
        RequestId: 49F4673B-B232-4401-A04B-D3E779BEDBBD
		Message: ContactGroups is mandatory for this action.

req.ContactGroups = "wen"

    /home/wen/git/site-monitor-operator/pkg/monitor/metric_rule_test.go:67: err SDK.ServerError
        ErrorCode: 400
        Recommend:
        RequestId: F353EF37-9DD8-47FF-8CF4-070CC205934E
        Message: alert master internal error:Failed to invoke controller method
        HandlerMethod details:
        Controller [com.aliyun.tianji.alert.restful.v2.AlarmController]
        Method [public com.aliyun.tianji.alert.domain.result.DataResult&lt;java.lang.String&gt; com.aliyun.tianji.alert.restful.v2.AlarmController.putMetricAlarm(java.lang.String

req.MetricName = "StatusCode"

response is &cms.PutResourceMetricRuleResponse{BaseResponse:(*responses.BaseResponse)(0xc000034e80), Success:true, Code:"200", Message:"", RequestId:"d8542f89-d314-404e-ac99-794686ac26e7"}
--- PASS: TestMetricRuleCreate (0.72s)
*/

// func TestMetricRuleCreate(t *testing.T) {
// 	// fmt.Println("s", s)
// 	err := s.CreateMetricRule()
// 	if err != nil {
// 		t.Error("err", err)
// 		return
// 	}
// }
