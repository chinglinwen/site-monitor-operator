// to interact with ali monitor
package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	//	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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

func (s *Server) ListMonitor(key string) (rs SiteMonitorResp, err error) {
	client, err := s.getclient()
	req := cms.CreateDescribeSiteMonitorListRequest()
	req.Scheme = "https"
	// req.PageSize = requests.Integer(100)
	// req.PageSize = requests.Integer(10000) // we fetch all at once // InvalidPageSize
	// req.Keyword = key

	total, size := 1, 1
	for i := 1; ; i++ {
		if i > 1 {
			if i > total/size+1 {
				log.Println("stop request due to no more pages after ", i)
				break
			}
		}
		log.Printf("requesting page: %v\n", i)
		// req.Page = requests.Integer(i)  // must increase the page number
		r, e := request(client, req)
		if e != nil {
			err = fmt.Errorf("request err: %v", e)
			return
		}
		if i == 1 {
			rs = r
			continue
		}
		total, size = r.TotalCount, r.PageSize

		rs.Data.SiteMonitors = append(rs.Data.SiteMonitors, r.Data.SiteMonitors...)
		// break
	}
	if len(rs.Data.SiteMonitors) != rs.TotalCount {
		fmt.Printf("got wrong number monitors, expect %v, got %v\n", rs.TotalCount, len(rs.Data.SiteMonitors))
	}
	return
}

func request(client *cms.Client, request *cms.DescribeSiteMonitorListRequest) (r SiteMonitorResp, err error) {
	response, err := client.DescribeSiteMonitorList(request)
	if err != nil {
		// original unmarshal is error, we unmarshal ourselves
		if strings.Contains(err.Error(), "JsonUnmarshalError") {
			// log.Println("got sdk unmarshal err", err)
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

type SiteMonitor struct {
	Interval    int    `json:"Interval"`
	CreateTime  string `json:"CreateTime"`
	Address     string `json:"Address"`
	OptionsJSON struct {
		HTTPMethod string `json:"http_method"`
		TimeOut    int    `json:"time_out"`
	} `json:"OptionsJson"`
	UpdateTime string `json:"UpdateTime"`
	TaskID     string `json:"TaskId"`
	TaskName   string `json:"TaskName"`
	TaskState  int    `json:"TaskState"`
	TaskType   string `json:"TaskType"`
}

type SiteMonitorResp struct {
	PageNumber int    `json:"PageNumber"`
	TotalCount int    `json:"TotalCount"`
	Message    string `json:"Message"`
	PageSize   int    `json:"PageSize"`
	RequestID  string `json:"RequestId"`
	Data       struct {
		SiteMonitors []SiteMonitor `json:"SiteMonitor"`
	} `json:"SiteMonitors"`
	Success bool   `json:"Success"`
	Code    string `json:"Code"`
}

func (s *Server) CreateMonitor() (err error) {
	client, err := s.getclient()

	request := cms.CreateModifySiteMonitorRequest()
	request.Scheme = "https"

	response, err := client.ModifySiteMonitor(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
	return
}
