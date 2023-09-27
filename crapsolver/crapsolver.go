package crapsolver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	headerContentTypeJson = []byte("application/json")

	client = &fasthttp.Client{
		ReadTimeout:                   10 * time.Second,
		WriteTimeout:                  10 * time.Second,
		MaxIdleConnDuration:           time.Minute,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
)

func NewSolver() *Solver {
	return &Solver{
		ServerAddr: SERVER_ADDR,
		Client:     client,
		WaitTime:   3 * time.Second,
	}
}

func (S *Solver) SetWaitTime(t time.Duration) error {
	if t.Seconds() > 30 {
		return fmt.Errorf("waiting is too long, 30s max")
	}

	if t.Seconds() < 1 {
		return fmt.Errorf("waiting is too small, 1s min")
	}

	S.WaitTime = t

	return nil
}

func (S *Solver) NewTask(config *TaskConfig) (resp *TaskResponse, err error) {
	if config.Domain == "" {
		config.Domain = "accounts.hcaptcha.com"
	}
	if config.SiteKey == "" {
		config.SiteKey = "2eaf963b-eeab-4516-9599-9daa18cd5138"
	}
	if config.UserAgent == "" {
		config.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36"
	}
	if config.Turbo && (config.TurboSt == 0 || config.TurboSt > 10000) {
		config.TurboSt = 3000
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf(`%s/api/task/new`, S.ServerAddr))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(payload)

	response := fasthttp.AcquireResponse()
	err = client.Do(req, response)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(response)

	if err := json.Unmarshal(response.Body(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (S *Solver) GetResult(T *TaskResponse) (resp *CheckResponse, err error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/task/%s", S.ServerAddr, T.Data[0].ID))
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.SetContentTypeBytes(headerContentTypeJson)

	response := fasthttp.AcquireResponse()
	err = client.Do(req, response)
	fasthttp.ReleaseRequest(req)

	if err := json.Unmarshal(response.Body(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (S *Solver) Solve(config *TaskConfig) (string, error) {
	task, err := S.NewTask(config)
	if err != nil {
		return "", err
	}

	if config.Turbo {
		time.Sleep(time.Duration(config.TurboSt) * time.Millisecond)
	} else {
		time.Sleep(5 * time.Second)
	}

	for {
		resp, err := S.GetResult(task)
		if err != nil {
			return "", err
		}

		switch resp.Data.Status {
		case STATUS_SOLVING:
			time.Sleep(S.WaitTime)
		case STATUS_SOLVED:
			return resp.Data.Token, nil
		case STATUS_ERROR:
			return "", fmt.Errorf(resp.Data.Error)
		}
	}
}
