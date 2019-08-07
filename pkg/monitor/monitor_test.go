package monitor

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var s = NewServer(AccessKey, AccessKeySecret)

/*
SDK.ServerError
        ErrorCode: Forbidden
        Recommend:
        RequestId: 3BADF51A-63CD-4923-9DE6-D97C8808A06E
		Message: You are not authorized to operate the specified resource.

    /home/wen/git/site-monitor-operator/pkg/monitor/monitor_test.go:21: err [SDK.JsonUnmarshalError] Failed to unmarshal response, but you can get the data via response.GetHttpStatusCode() and response.GetHttpContentString()
        caused by:
		cms.DescribeSiteMonitorListResponse.SiteMonitors: cms.SiteMonitorsInDescribeSiteMonitorList.SiteMonitor: []cms.SiteMonitor: cms.SiteMonitor.OptionsJson: fuzzyStringDecoder: not number or string or bool, error found in #10 byte of ...|ionsJson":{"http_met|..., bigger context ...|s://loan.rongba.com/h5tuiguang/aff","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime|...

    /home/wen/git/site-monitor-operator/pkg/monitor/monitor_test.go:27: err request err: SDK.ServerError
        ErrorCode: InvalidPageSize
        Recommend: https://error-center.aliyun.com/status/search?Keyword=InvalidPageSize&source=PopGw
        RequestId: 4F68A43C-D9DD-4576-868C-7E2FD9F4393E
		Message: Specified parameter PageSize is not valid

    /home/wen/git/site-monitor-operator/pkg/monitor/monitor_test.go:33: err request err: SDK.ServerError
        ErrorCode: InvalidPage
        Recommend: https://error-center.aliyun.com/status/search?Keyword=InvalidPage&source=PopGw
        RequestId: 1BB40233-823F-4CB7-A2C6-6F167B76AADD
        Message: Specified parameter Page is not valid.
*/

/*
monitor.SiteMonitor{Interval:1, CreateTime:"2018-12-20 17:26:29", Address:"https://loan.rongba.com/h5tuiguang/aff", OptionsJSON:struct { HTTPMethod string "json:\"http_method\""; TimeOut int "json:\"time_out\"" }{HTTPMethod:"get", TimeOut:30000}, UpdateTime:"2019-07-12 09:41:46", TaskID:"f5003de3-b957-48d3-8e0c-69271e04559e", TaskName:"flow_center-loan.rongba.com", TaskState:1, TaskType:"HTTP"},
monitor.SiteMonitor{Interval:1, CreateTime:"2018-12-20 17:30:56", Address:"http://oc.haodai.com/Home/OrderApi/orderAdd", OptionsJSON:struct { HTTPMethod string "json:\"http_method\""; TimeOut int "json:\"time_out\"" }{HTTPMethod:"get", TimeOut:30000}, UpdateTime:"2019-07-12 09:41:46", TaskID:"a518d7a3-eefb-43e1-9a80-66a9b2994661", TaskName:"flow_center-order_center_later", TaskState:1, TaskType:"HTTP"},
monitor.SiteMonitor{Interval:1, CreateTime:"2018-12-20 17:31:58", Address:"https://openapi.haodai.com/SdTuiguang/checkstatus", OptionsJSON:struct { HTTPMethod string "json:\"http_method\""; TimeOut int "json:\"time_out\"" }{HTTPMethod:"get", TimeOut:30000}, UpdateTime:"2019-07-12 09:41:46", TaskID:"4d77518b-4be1-4c05-b51d-2b9f04e60019", TaskName:"flow_center-openapi-1", TaskState:1, TaskType:"HTTP"},
*/

/*
(monitor.SiteMonitor) {
 Interval: (int) 1,
 CreateTime: (string) (len=19) "2018-12-20 17:26:29",
 Address: (string) (len=38) "https://loan.rongba.com/h5tuiguang/aff",
 OptionsJSON: (monitor.OptionsJSON) {
  HTTPMethod: (string) (len=3) "get",
  TimeOut: (int) 30000
 },
 UpdateTime: (string) (len=19) "2019-07-12 09:41:46",
 TaskID: (string) (len=36) "f5003de3-b957-48d3-8e0c-69271e04559e",
 TaskName: (string) (len=27) "flow_center-loan.rongba.com",
 TaskState: (int) 1,
 TaskType: (string) (len=4) "HTTP"
}

(monitor.SiteMonitor) {
 Interval: (int) 1,
 CreateTime: (string) (len=19) "2019-06-11 10:50:09",
 Address: (string) (len=23) "http://a.clwen.com:9000",
 OptionsJSON: (monitor.OptionsJSON) {
  HTTPMethod: (string) (len=3) "get",
  TimeOut: (int) 30000
 },
 UpdateTime: (string) (len=19) "2019-07-12 10:50:33",
 TaskID: (string) (len=36) "b47607e8-9826-4ba0-8f04-d27c8ca0670b",
 TaskName: (string) (len=4) "test",
 TaskState: (int) 2,
 TaskType: (string) (len=4) "HTTP"
}
*/
func TestList(t *testing.T) {
	// fmt.Println("s", s)
	r, err := s.ListMonitor("")
	if err != nil {
		t.Error("err", err)
		return
	}
	fmt.Printf("got %v monitors\n", len(r.Data.SiteMonitors))
	// return
	for _, v := range r.Data.SiteMonitors {
		// if i == 0 || v.TaskName == "test" || v.TaskName == "baidu.com" {
		if v.TaskName == "baidu.com" {
			spew.Dump(v)
		}
		fmt.Printf("%v->%v\n", v.TaskName, v.Address)
		// spew.Dump("got",v)
	}
	// fmt.Printf("r: %v\n", r)
}

func TestGetSiteMonitor(t *testing.T) {
	sm, err := s.getSiteMonitor("baidu.com")
	if err != nil {
		t.Error("err", err)
		return
	}
	spew.Dump("sm", sm)
}

// yaml tag doesn't work, the following example is ok
// func TestUnmarshalSiteMonitor(t *testing.T) {
// 	a :=
// 		`
// TaskName: baidu.com
// TaskType: HTTP
// HttpMethod: get
// Timeout: 1000
// Interval: 1
// Address: http://www.baidu.com
// TaskState: true  # default to true
// ContactGroups: defaultGroup
// `

// 	m := SiteMonitor{}
// 	err := yaml.Unmarshal([]byte(a), &m)
// 	if err != nil {
// 		t.Error("err", err)
// 		return
// 	}
// 	spew.Dump(m)
// }

/*
{"Interval":1,"CreateTime":"2018-12-20 17:26:29","Address":"https://loan.rongba.com/h5tuiguang/aff","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"f5003de3-b957-48d3-8e0c-69271e04559e","TaskName":"flow_center-loan.rongba.com","TaskState":1,"TaskType":"HTTP"}

    /home/wen/git/site-monitor-operator/pkg/monitor/monitor_test.go:109: err CreateSiteMonitor err: SDK.ServerError
        ErrorCode: NameRepeat
        Recommend:
        RequestId: 7F392789-803B-4B97-A5C2-A39DCAE4D98E
		Message: Task name repeat

because SiteMonitor exist, we now delete it first, if it exist

    /home/wen/git/site-monitor-operator/pkg/monitor/monitor_test.go:117: err CreateDefaultMetric err: create metricrule ResponseTime for baidu.com err: SDK.ServerError
        ErrorCode: 400
        Recommend:
        RequestId: 370A508B-414A-419D-8EB5-2A61F6EF60A7
		Message: Bad Request

because missing req.RuleName = m.TaskID, fixed
*/

func TestCreateMonitor(t *testing.T) {
	sm := SiteMonitor{
		Interval:   1,
		Address:    "http://www.baidu.com",
		HttpMethod: "get",
		Timeout:    30000,
		TaskName:   "baidu.com",
		TaskType:   "HTTP",
		TaskState:  false,
	}
	err := s.CreateMonitor(sm)
	if err != nil {
		t.Error("err", err)
		return
	}
}

func TestDeleteMonitor(t *testing.T) {
	err := s.DeleteMonitor("baidu.com")
	if err != nil {
		t.Error("err", err)
		return
	}
}

func TestCreateDefaultMetric(t *testing.T) {
	err := s.CreateDefaultMetric("baidu.com", "wen1")
	if err != nil {
		t.Errorf("CreateDefaultMetric err: %v", err)
		return
	}
}

func TestDisableMonitor(t *testing.T) {
	err := s.DisableMonitor("baidu.com")
	if err != nil {
		t.Error("err", err)
		return
	}
}

func TestEnableMonitor(t *testing.T) {
	err := s.EnableMonitor("baidu.com")
	if err != nil {
		t.Error("err", err)
		return
	}
}

var examplebody = `
{"PageNumber":1,"TotalCount":28,"Message":"请求成功","PageSize":10,"RequestId":"","SiteMonitors":{"SiteMonitor":[{"Interval":1,"CreateTime":"2018-12-20 17:26:29","Address":"https://loan.rongba.com/h5tuiguang/aff","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"f5003de3-b957-48d3-8e0c-69271e04559e","TaskName":"flow_center-loan.rongba.com","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 17:30:56","Address":"http://oc.haodai.com/Home/OrderApi/orderAdd","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"a518d7a3-eefb-43e1-9a80-66a9b2994661","TaskName":"flow_center-order_center_later","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 17:31:58","Address":"https://openapi.haodai.com/SdTuiguang/checkstatus","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"4d77518b-4be1-4c05-b51d-2b9f04e60019","TaskName":"flow_center-openapi-1","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 17:31:58","Address":"https://openapi.haodai.com/Capi/Index/index_newlist","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"9df8744c-a950-4da2-b677-d5fe9eb6c026","TaskName":"flow_center-openapi-2","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 17:32:32","Address":"https://loanapi.haodai.com/api/about/help","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"39366442-8c56-4f4c-9be2-40c1ed264988","TaskName":"flow_center-loanapi","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 18:17:34","Address":"http://cloud.haodai.com/HostCheck/main","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"06ce82c0-8d16-49fa-b42a-5ec3fec31c7a","TaskName":"flow_center-8.yun","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 18:18:09","Address":"https://click.haodai.com/tool/hostcheck.php","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"d8616be2-5afb-47b0-9ae7-3564ec09434f","TaskName":"flow_center-click.haodai.com","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 18:18:42","Address":"http://yun.haodai.com/HostCheck/main","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"a137fd83-290a-4fa2-b943-400d4b89294f","TaskName":"flow_center-yun.haodai.com","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 18:19:39","Address":"http://www.hanbaojinrong.com/Index/Tuiguang/test","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"46febe57-527e-4cdb-9dce-fd54de2db871","TaskName":"flow_center-hamburg","TaskState":1,"TaskType":"HTTP"},{"Interval":1,"CreateTime":"2018-12-20 18:22:14","Address":"http://guanfang.hanbaofamily.com/Index/Tuiguang/test","OptionsJson":{"http_method":"get","time_out":30000},"UpdateTime":"2019-07-12 09:41:46","TaskId":"630b4f0f-a83c-4e56-9435-c583906dc66f","TaskName":"flow_center-agent_system","TaskState":1,"TaskType":"HTTP"}]},"Success":true,"Code":"200"}
`

func TestUnquote(t *testing.T) {
	// fmt.Println("s", s)
	// s, err := unquote(examplebody)
	// if err != nil {
	// 	t.Error("unquote", err)
	// }
	// fmt.Println(s)

	// r := &SiteMonitorResp{}
	// err = json.Unmarshal([]byte(s), r)
	r, err := unmrshalMonitorsBody(examplebody)
	if err != nil {
		t.Error("Unmarshal err", err)
	}
	fmt.Println(r)
}
