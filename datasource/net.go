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
	header  req.Header
	param   req.Param
	reqType ReqType
	isJson  bool
}

func NewNet(url string, header req.Header, param req.Param, reqType ReqType) *Net {
	return &Net{url: url, header: header, param: param, reqType: reqType}
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

func InitParam(params map[string]interface{}) req.Param {
	reqParams := req.Param{}
	for k, v := range params {
		reqParams[k] = v
	}
	return reqParams
}

func (n *Net) SetJson(isJson bool) {
	n.isJson = isJson
}

func (n *Net) Request() ([]byte, error) {
	switch n.reqType {
	case POST:
		return n.post()
	case GET:
		return n.get()
	default:
		return n.get()
	}
}

func (n *Net) post() ([]byte, error) {
	var reqResp = &req.Resp{}
	var err error
	if n.isJson {
		jsonParam, _ := json.Marshal(n.param)
		reqResp, err = req.Post(n.url, jsonParam, n.header)
	} else {
		reqResp, err = req.Post(n.url, n.param, n.header)
	}
	return reqResp.Bytes(), err
}

func (n *Net) get() ([]byte, error) {
	resp, err := req.Get(n.url, n.header)
	if err != nil {
		return nil, err
	}
	return resp.Bytes(), nil
}
