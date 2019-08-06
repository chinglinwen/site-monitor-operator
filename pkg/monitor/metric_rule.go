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
	// req.PageSize = requests.Integer("10000")
	// req.TaskName = "test"

	return client.DescribeMetricRuleList(req)
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
