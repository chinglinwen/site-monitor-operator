package monitor

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
)

type MetricRule struct {
	// TaskName   string // to generate ruleid?
	TaskID     string
	MetricName string // "ResponseTime" //StatusCode
	// RuleName      string // "test-a"
	// RuleId        string // "31bc7843-afd1-44ab-b9ed-67e5314d6095_ResponseTime"
	// Resources     string // "[{\"taskId\":\"31bc7843-afd1-44ab-b9ed-67e5314d6095\"}]"
	ContactGroups string // "wen"

	EffectiveInterval string // "00:00-23:59"
	SilenceTime       int    // "86400"
	Period            string // "60"

	Namespace string // "acs_networkmonitor"
	Webhook   string // "http://alialert.haodai.net"

	EscalationsWarnStatistics         string // "Average"
	EscalationsWarnComparisonOperator string // "GreaterThanThreshold"
	EscalationsWarnThreshold          string //"12000"
	EscalationsWarnTimes              int    // "3"
}

func (s *Server) MetricRuleList() (response *cms.DescribeMetricRuleListResponse, err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}

	req := cms.CreateDescribeMetricRuleListRequest()
	req.Scheme = "https"
	req.PageSize = "10000"

	return client.DescribeMetricRuleList(req)
}

func (s *Server) GetMetricRule(taskname string) (alarms []cms.Alarm, err error) {
	taskid, err := s.gettaskid(taskname)
	if err != nil {
		err = fmt.Errorf("gettaskid for %v err: %v", taskname, err)
		return
	}
	r, err := s.MetricRuleList()
	if err != nil {
		err = fmt.Errorf("MetricRuleList err: %v", err)
		return
	}

	// fmt.Printf("got %v alarms\n", len(r.Alarms.Alarm))
	for _, v := range r.Alarms.Alarm {
		if v.RuleName == taskid {
			alarms = append(alarms, v)
		}
	}
	if len(alarms) == 0 {
		err = fmt.Errorf("no metricrules been found")
		return
	}
	return
}

func (s *Server) GetAlertState(taskname string) (alertstate string, err error) {
	alarms, err := s.GetMetricRule(taskname)
	if err != nil {
		return
	}
	for i, v := range alarms {
		if i == 0 {
			alertstate = fmt.Sprintf("%v: %v", v.MetricName, v.AlertState)
			continue
		}
		alertstate = fmt.Sprintf("%v, %v: %v", alertstate, v.MetricName, v.AlertState)
	}
	return
}

func (s *Server) compareAndUpdateMetricRule(taskname, contactgroups string) (err error) {
	alarms, err := s.GetMetricRule(taskname)
	if err != nil {
		return
	}
	var update bool
	for _, v := range alarms {
		if v.ContactGroups != contactgroups {
			update = true
		}
	}
	if update {
		err = s.CreateDefaultMetric(taskname, contactgroups)
		if err != nil {
			err = fmt.Errorf("CreateDefaultMetric err: %v", err)
			return
		}
	}
	return
}

// convert MetricRule to request struct
func metricRuleToCreateReq(m *MetricRule) (req *cms.PutResourceMetricRuleRequest) {
	req = cms.CreatePutResourceMetricRuleRequest()
	req.Scheme = "https"

	req.RuleName = m.TaskID                                   // TODO: append generation-id?
	req.MetricName = m.MetricName                             // ResponseTime, Availability, StatusCode
	req.RuleId = fmt.Sprintf("%v_%v", m.TaskID, m.MetricName) // "31bc7843-afd1-44ab-b9ed-67e5314d6095_ResponseTime"
	req.Resources = fmt.Sprintf("[{\"taskId\":\"%v\"}]", m.TaskID)
	req.ContactGroups = m.ContactGroups
	req.EffectiveInterval = m.EffectiveInterval          //"00:00-23:59"
	req.SilenceTime = requests.NewInteger(m.SilenceTime) //"86400"
	req.Period = m.Period                                //"60"

	req.Namespace = m.Namespace // "acs_networkmonitor"
	req.Webhook = m.Webhook     // "http://alialert.haodai.net"

	req.EscalationsWarnStatistics = m.EscalationsWarnStatistics                 // "Average"
	req.EscalationsWarnComparisonOperator = m.EscalationsWarnComparisonOperator // "GreaterThanThreshold"
	req.EscalationsWarnThreshold = m.EscalationsWarnThreshold                   //"12000"
	req.EscalationsWarnTimes = requests.NewInteger(m.EscalationsWarnTimes)      //"3"

	return
}

func (s *Server) CreateMetricRule(m *MetricRule) (err error) {
	client, err := s.getclient()
	if err != nil {
		err = fmt.Errorf("create client err: %v", err)
		return
	}
	if m == nil {
		err = fmt.Errorf("empty metric rule parameter")
		return
	}

	req := metricRuleToCreateReq(m)

	_, err = client.PutResourceMetricRule(req)
	if err != nil {
		err = fmt.Errorf("PutResourceMetricRule err: %v", err)
		return
	}
	return
}

// CreateGroupMetricRules is for app group, we don't use this

// func (s *Server) MetricRuleCreate() (err error) {
// 	client, err := s.getclient()
// 	if err != nil {
// 		err = fmt.Errorf("create client err: %v", err)
// 		return
// 	}
// 	req := cms.CreateCreateGroupMetricRulesRequest()
// 	req.Scheme = "https"
// 	req.GroupId = requests.NewInteger(1)

// 	req.GroupMetricRules = &[]cms.CreateGroupMetricRulesGroupMetricRules{
// 		{
// 			RuleName:                          "test-a",
// 			Namespace:                         "acs_networkmonitor",
// 			MetricName:                        "baidu.com",
// 			EffectiveInterval:                 "00:00-23:59",
// 			SilenceTime:                       "86400",
// 			Period:                            "60",
// 			Webhook:                           "http://alialert.haodai.net",
// 			EscalationsWarnStatistics:         "Average",
// 			EscalationsWarnComparisonOperator: "GreaterThanThreshold",
// 			EscalationsWarnThreshold:          "12000",
// 			EscalationsWarnTimes:              "3",
// 		},
// 	}

// 	response, err := client.CreateGroupMetricRules(req)
// 	if err != nil {
// 		fmt.Print(err.Error())
// 	}
// 	fmt.Printf("response is %#v\n", response)

// 	return
// }
