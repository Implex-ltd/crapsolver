package crapsolver

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

// Get restriction for a specific site-keys
func GetRestrictions(sitekey string) (*Restrictions, error) {
	Server, err := Nodes.Next()
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/misc/check/%s", Server, sitekey))
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("crapsolver-go")
	req.Header.SetContentTypeBytes(headerContentTypeJson)

	response := fasthttp.AcquireResponse()
	if err = client.Do(req, response); err != nil {
		return nil, err
	}

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

func GetUser(uid string) (*GetUserResp, error) {
	Server, err := Nodes.Next()
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/user/%s", Server, uid))
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetUserAgent("crapsolver-go")
	req.Header.SetContentTypeBytes(headerContentTypeJson)

	response := fasthttp.AcquireResponse()
	if err = client.Do(req, response); err != nil {
		return nil, err
	}

	fasthttp.ReleaseRequest(req)

	var resp GetUserResp
	if err := json.Unmarshal(response.Body(), &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("error: %v", string(response.Body()))
	}

	return &resp, nil
}
