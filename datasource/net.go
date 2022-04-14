package datasource

import (
	"encoding/json"
	"github.com/imroc/req"
	"strings"
)

type ReqType string

const (
	POST ReqType = "post"
	GET  ReqType = "get"
)

type Net struct {
	url     string
	apiKey  string
	header  req.Header
	param   req.Param
	reqType ReqType
	isJson  bool
}

func NewNet(url, apiKey string, header req.Header, param req.Param, reqType ReqType) *Net {
	return &Net{url: url, apiKey: apiKey, header: header, param: param, reqType: reqType}
}

func InitHeader(header map[string]string) (req.Header, bool) {
	authHeader := req.Header{}
	hasJson := false

	for k, v := range header {
		authHeader[k] = v
		if hasJsonInHeader(k, v) {
			hasJson = true
		}
	}
	return authHeader, hasJson
}

func hasJsonInHeader(key, value string) bool {
	return strings.ToLower(key) == "accept" && strings.Contains(strings.ToLower(value), "json")
}

func InitParam(params map[string]string) req.Param {
	reqParams := req.Param{}
	for k, v := range params {
		reqParams[k] = v
	}
	return reqParams
}

func (n *Net) Request() (string, error) {
	switch n.reqType {
	case POST:
		return n.post()
	case GET:
		return n.get()
	default:
		return n.get()
	}
}

func (n *Net) post() (string, error) {
	var reqResp = &req.Resp{}
	var err error
	if n.isJson {
		jsonParam, _ := json.Marshal(n.param)
		reqResp, err = req.Post(n.url, jsonParam, n.header)
	} else {
		reqResp, err = req.Post(n.url, n.param, n.header)
	}
	return reqResp.String(), err
}

func (n *Net) get() (string, error) {
	resp, err := req.Get(n.url, n.header)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}
