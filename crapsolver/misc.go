package crapsolver

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func GetRestrictions(sitekey string) (*Restrictions, error) {
	Server, err := Nodes.Next()
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/misc/check/%s", Server, sitekey))
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentTypeBytes(headerContentTypeJson)

	response := fasthttp.AcquireResponse()
	err = client.Do(req, response)
	fasthttp.ReleaseRequest(req)

	var resp CheckRestrictionResp
	if err := json.Unmarshal(response.Body(), &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("error: %v", string(response.Body()))
	}

	return &resp.Data, nil
}
