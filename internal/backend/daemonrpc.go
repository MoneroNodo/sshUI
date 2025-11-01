package backend

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/moneronodo/sshui/internal/base"
	"github.com/moneronodo/sshui/internal/model/daemonrpc"
)

type DaemonRequestBody struct {
	body         []byte
	responseType daemonrpc.DaemonRPCResponseWrapper
}

func DaemonPost(url string, body DaemonRequestBody) *daemonrpc.DaemonRPCResponse {
	c := &http.Client{Timeout: 3 * time.Second}
	resp, err := c.Post(url, "application/json", bytes.NewReader(body.body))
	if err != nil {
		spew.Fprintf(base.Dump, "http: %v\n", err)
		return &daemonrpc.DaemonRPCResponse{}
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	j := daemonrpc.MakeDaemonRPCResponse(body.responseType)
	if err := dec.Decode(j); err != nil {
		spew.Fprintf(base.Dump, "Decode: %v\n", err)
		return &daemonrpc.DaemonRPCResponse{}
	}
	return j
}

func daemonMakeRequestBody(method string, params any) []byte {
	req := &daemonrpc.DaemonRPCRequest{
		Jsonrpc: "2.0",
		Id:      0,
		Method:  method,
		Params:  params,
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil
	}
	return b
}

func DaemonRequestBodyGetInfo() DaemonRequestBody {
	method := "get_info"
	body := DaemonRequestBody{
		daemonMakeRequestBody(method, nil),
		daemonrpc.DaemonRPCResponseWrapper(&daemonrpc.DaemonResponseBodyGetInfo{}),
	}
	return body
}

func DaemonRequestBodyGetVersion() DaemonRequestBody {
	method := "get_version"
	body := DaemonRequestBody{
		daemonMakeRequestBody(method, nil),
		daemonrpc.DaemonRPCResponseWrapper(&daemonrpc.DaemonResponseBodyGetVersion{}),
	}
	return body
}
