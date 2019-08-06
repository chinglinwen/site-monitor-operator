// to interact with ali monitor
package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	// "log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/tidwall/gjson"
	//	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hokaccha/go-prettyjson"
)

type Server struct {
	key    string
	secret string
	client *cms.Client
}

func NewServer(key, secret string) *Server {
	return &Server{
		key:    key,
		secret: secret,
	}
}

func (s *Server) getclient() (client *cms.Client, err error) {
	if s.client == nil {
		client, err = cms.NewClientWithAccessKey("cn-hangzhou", s.key, s.secret)
		if err != nil {
			return
		}
		s.client = client
	}
	return s.client, nil
}

// Keyword is the taskname,  empty key will list all
func (s *Server) ListMonitor(Keyword string) (r SiteMonitorResp, err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}
	req := cms.CreateDescribeSiteMonitorListRequest()
	req.Scheme = "https"
	req.PageSize = requests.NewInteger(10000) // we fetch all at once
	req.Keyword = Keyword

	r, err = request(client, req)
	if err != nil {
		err = fmt.Errorf("request err: %v", err)
		return
	}

	if len(r.Data.SiteMonitors) != r.TotalCount {
		fmt.Printf("got wrong number monitors, expect %v, got %v\n", r.TotalCount, len(r.Data.SiteMonitors))
	}
	return
}

func request(client *cms.Client, request *cms.DescribeSiteMonitorListRequest) (r SiteMonitorResp, err error) {
	response, err := client.DescribeSiteMonitorList(request)
	if err != nil {
		// original unmarshal is error, we unmarshal ourselves
		if strings.Contains(err.Error(), "JsonUnmarshalError") {
			log.Println("got sdk unmarshal err", err)
			err = nil
		} else {
			return
		}
	}
	// detect code first by response.GetHttpStatusCode()?
	return unmrshalMonitorsBody(response.GetHttpContentString())
}

// invalid syntax
func unquote(body string) (s string, err error) {
	// s = body
	// s := strings.TrimSuffix(strings.TrimPrefix(body, "\""), "\"")
	// s = "`" + s + "`" // not working
	// s = "\\\"" + s + "\\\""  // not working
	// return strconv.Unquote(s)
	s = strings.Replace(body, `\"`, `"`, -1)
	return
}

func unmrshalMonitorsBody(body string) (r SiteMonitorResp, err error) {
	// fmt.Printf("body: %v\n", body)
	s, err := unquote(body)
	if err != nil {
		err = fmt.Errorf("unquote err: %v, body: %v", err, body)
		return
	}
	r = SiteMonitorResp{}
	err = json.Unmarshal([]byte(s), &r)
	if err != nil {
		err = fmt.Errorf("unmarshal sitemonintor list response err: %v", err)
		return
	}
	return
}

type OptionsJSON struct {
	HTTPMethod string `json:"http_method"`
	TimeOut    int    `json:"time_out"`
}

type SiteMonitor struct {
	TaskName   string `json:"TaskName"`
	TaskType   string `json:"TaskType"`
	HTTPMethod string `json:"http_method"`
	TimeOut    int    `json:"time_out"`

	Interval  int    `json:"Interval"`
	Address   string `json:"Address"`
	TaskState bool   `json:"TaskState"` // default enabled
}

// for unmarshal purpose
type SiteMonitorRaw struct {
	TaskName string `json:"TaskName"`
	TaskType string `json:"TaskType"`
	Interval int    `json:"Interval"`
	Address  string `json:"Address"`

	TaskID    string `json:"TaskId"`    // unmarshal needed this
	TaskState int    `json:"TaskState"` // default enabled

	OptionsJSON OptionsJSON `json:"OptionsJson"`
	// OptionsJSON struct {
	// 	HTTPMethod string `json:"http_method"`
	// 	TimeOut    int    `json:"time_out"`
	// } `json:"OptionsJson"`
	CreateTime string `json:"CreateTime"`
	UpdateTime string `json:"UpdateTime"`
}

type SiteMonitorResp struct {
	PageNumber int    `json:"PageNumber"`
	TotalCount int    `json:"TotalCount"`
	Message    string `json:"Message"`
	PageSize   int    `json:"PageSize"`
	RequestID  string `json:"RequestId"`
	Data       struct {
		SiteMonitors []SiteMonitorRaw `json:"SiteMonitor"`
	} `json:"SiteMonitors"`
	Success bool   `json:"Success"`
	Code    string `json:"Code"`
}

/*
	OptionsJson string `position:"Query" name:"OptionsJson"`
	Address     string `position:"Query" name:"Address"`
	TaskType    string `position:"Query" name:"TaskType"`
	AlertIds    string `position:"Query" name:"AlertIds"`
	TaskName    string `position:"Query" name:"TaskName"`
	Interval    string `position:"Query" name:"Interval"`
	IspCities   string `position:"Query" name:"IspCities"`
*/
func SiteMonitorToCreateReq(s SiteMonitor) (r *cms.CreateSiteMonitorRequest) {
	r = cms.CreateCreateSiteMonitorRequest()

	r.TaskType = s.TaskType
	r.TaskName = s.TaskName
	r.Interval = strconv.Itoa(s.Interval)
	r.Address = s.Address
	// r.IspCities

	op, _ := json.Marshal(OptionsJSON{
		HTTPMethod: s.HTTPMethod,
		TimeOut:    s.TimeOut,
	})
	r.OptionsJson = string(op)
	return
}

func SiteMonitorToModReq(s SiteMonitor) (r *cms.ModifySiteMonitorRequest) {
	r = cms.CreateModifySiteMonitorRequest()

	// r.TaskType = s.TaskType
	r.TaskName = s.TaskName
	r.Interval = strconv.Itoa(s.Interval)
	r.Address = s.Address

	op, _ := json.Marshal(OptionsJSON{
		HTTPMethod: s.HTTPMethod,
		TimeOut:    s.TimeOut,
	})
	r.OptionsJson = string(op)
	return
}

/*
// https://godoc.org/github.com/aliyun/alibaba-cloud-sdk-go/services/cms#CreateHostAvailabilityRequest
type CreateHostAvailabilityRequest struct {
    *requests.RpcRequest
    InstanceList                       *[]string                                          `position:"Query" name:"InstanceList"  type:"Repeated"`
    TaskType                           string                                             `position:"Query" name:"TaskType"`
    TaskOptionHttpMethod               string                                             `position:"Query" name:"TaskOption.HttpMethod"`
    AlertConfigEscalationList          *[]CreateHostAvailabilityAlertConfigEscalationList `position:"Query" name:"AlertConfigEscalationList"  type:"Repeated"`
    GroupId                            requests.Integer                                   `position:"Query" name:"GroupId"`
    TaskName                           string                                             `position:"Query" name:"TaskName"`
    AlertConfigSilenceTime             requests.Integer                                   `position:"Query" name:"AlertConfig.SilenceTime"`
    TaskOptionHttpResponseCharset      string                                             `position:"Query" name:"TaskOption.HttpResponseCharset"`
    AlertConfigEndTime                 requests.Integer                                   `position:"Query" name:"AlertConfig.EndTime"`
    TaskOptionHttpURI                  string                                             `position:"Query" name:"TaskOption.HttpURI"`
    TaskOptionHttpNegative             requests.Boolean                                   `position:"Query" name:"TaskOption.HttpNegative"`
    TaskScope                          string                                             `position:"Query" name:"TaskScope"`
    AlertConfigNotifyType              requests.Integer                                   `position:"Query" name:"AlertConfig.NotifyType"`
    AlertConfigStartTime               requests.Integer                                   `position:"Query" name:"AlertConfig.StartTime"`
    TaskOptionTelnetOrPingHost         string                                             `position:"Query" name:"TaskOption.TelnetOrPingHost"`
    TaskOptionHttpResponseMatchContent string                                             `position:"Query" name:"TaskOption.HttpResponseMatchContent"`
    AlertConfigWebHook                 string                                             `position:"Query" name:"AlertConfig.WebHook"`
}
*/

func (s *Server) CreateMonitor(sm SiteMonitor) (err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}

	request := SiteMonitorToCreateReq(sm)

	_, err = client.CreateSiteMonitor(request)
	if err != nil {
		if strings.Contains(err.Error(), "NameRepeat") {
			log.Printf("delete exist monitor: %v first\n", sm.TaskName)
			err = s.DeleteMonitor(sm.TaskName)
			if err != nil {
				err = fmt.Errorf("create client with delete first for %v err: %v", sm.TaskName, err)
				return
			}
			_, err = client.CreateSiteMonitor(request)
			if err != nil {
				err = fmt.Errorf("CreateSiteMonitor err: %v", err)
				return
			}
		} else {
			err = fmt.Errorf("CreateSiteMonitor err: %v", err)
			return
		}

	}
	// fmt.Printf("response is %#v\n", response)

	err = s.CreateDefaultMetric(sm.TaskName)
	if err != nil {
		err = fmt.Errorf("CreateDefaultMetric err: %v", err)
		return
	}
	return
}

func (s *Server) CreateDefaultMetric(taskname string) (err error) {
	taskid, err := s.gettaskid(taskname)
	if err != nil {
		err = fmt.Errorf("gettaskid for %v err: %v", taskname, err)
		return
	}

	m1 := &MetricRule{
		// TaskName: taskname,
		TaskID:        taskid,
		MetricName:    "ResponseTime",
		ContactGroups: "wen",

		EffectiveInterval: "00:00-23:59",
		SilenceTime:       86400,
		Period:            "60",

		Namespace: "acs_networkmonitor",
		Webhook:   "http://alialert.haodai.net",

		EscalationsWarnStatistics:         "Average",
		EscalationsWarnComparisonOperator: "GreaterThanThreshold",
		EscalationsWarnThreshold:          "12000",
		EscalationsWarnTimes:              3,
	}
	err = s.CreateMetricRule(m1)
	if err != nil {
		err = fmt.Errorf("create metricRule ResponseTime for %v err: %v", taskname, err)
		return
	}

	m2 := &MetricRule{
		// TaskName: taskname,
		TaskID:        taskid,
		MetricName:    "Availability",
		ContactGroups: "wen",

		EffectiveInterval: "00:00-23:59",
		SilenceTime:       86400,
		Period:            "60",

		Namespace: "acs_networkmonitor",
		Webhook:   "http://alialert.haodai.net",

		EscalationsWarnStatistics:         "Availability",
		EscalationsWarnComparisonOperator: "LessThanThreshold",
		EscalationsWarnThreshold:          "90",
		EscalationsWarnTimes:              3,
	}
	err = s.CreateMetricRule(m2)
	if err != nil {
		err = fmt.Errorf("create metricRule Availability for %v err: %v", taskname, err)
		return
	}
	return
}

func pretty(prefix, a interface{}) {
	out, _ := prettyjson.Marshal(a)
	fmt.Printf("%v: %s\n", prefix, out)
}

/*
type ModifySiteMonitorRequest struct {
	*requests.RpcRequest
	OptionsJson string `position:"Query" name:"OptionsJson"`
	Address     string `position:"Query" name:"Address"`
	AlertIds    string `position:"Query" name:"AlertIds"`
	TaskName    string `position:"Query" name:"TaskName"`
	Interval    string `position:"Query" name:"Interval"`
	TaskId      string `position:"Query" name:"TaskId"`
	IspCities   string `position:"Query" name:"IspCities"`
}
*/

// https://help.aliyun.com/document_detail/115049.html?spm=a2c4g.11186623.6.807.39c9d8e0dcuSVE
func (s *Server) ModifyMonitor(sm SiteMonitor) (err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}
	request := SiteMonitorToModReq(sm)
	response, err := client.ModifySiteMonitor(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	return
}

func (s *Server) DeleteMonitor(taskname string) (err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client for deletemonitor err: %v", err)
		return
	}
	req := cms.CreateDeleteSiteMonitorsRequest()
	req.Scheme = "https"
	taskid, err := s.gettaskid(taskname)
	if err != nil {
		err = fmt.Errorf("gettaskid for %v err: %v", taskname, err)
		return
	}
	req.TaskIds = taskid
	_, err = client.DeleteSiteMonitors(req)
	if err != nil {
		if strings.Contains(err.Error(), "JsonUnmarshalError") {
			err = nil
			return
		}
		err = fmt.Errorf("deletemonitor for %v err: %v", taskname, err)
		return
	}
	// fmt.Printf("response is %#v\n", response)
	return
}

func (s *Server) gettaskid(taskname string) (taskid string, err error) {
	r, err := s.ListMonitor(taskname)
	if err != nil {
		err = fmt.Errorf("try find taskid for %v by listmonitor err: %v", taskname, err)
		return
	}
	n := len(r.Data.SiteMonitors)
	if n != 1 {
		err = fmt.Errorf("error find taskid for %v, expect 1 result, got %v results", taskname, n)
		return
	}
	taskid = r.Data.SiteMonitors[0].TaskID
	return
}

func (s *Server) DisableMonitor(taskname string) (err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}

	taskid, err := s.gettaskid(taskname)
	if err != nil {
		err = fmt.Errorf("gettaskid for %v err: %v", taskname, err)
		return
	}

	req := cms.CreateDisableSiteMonitorsRequest()
	req.Scheme = "https"
	req.TaskIds = taskid

	response, err := client.DisableSiteMonitors(req)
	if err != nil {
		resp := response.GetHttpContentString()
		if gjson.Get(resp, "Code").String() == "200" {
			err = nil
			return
		}
		err = fmt.Errorf("DisableMonitor for %v err: %v", taskname, err)
		return
	}
	// fmt.Printf("response is %#v\n", response)

	return
}

func (s *Server) EnableMonitor(taskname string) (err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}

	taskid, err := s.gettaskid(taskname)
	if err != nil {
		err = fmt.Errorf("gettaskid for %v err: %v", taskname, err)
		return
	}

	req := cms.CreateEnableSiteMonitorsRequest()
	req.Scheme = "https"
	req.TaskIds = taskid

	response, err := client.EnableSiteMonitors(req)
	if err != nil {
		resp := response.GetHttpContentString()
		if gjson.Get(resp, "Code").String() == "200" {
			err = nil
			return
		}
		err = fmt.Errorf("EnableMonitor for %v err: %v", taskname, err)
		return
	}
	// fmt.Printf("response is %#v\n", response)

	return
}
