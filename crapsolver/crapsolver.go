package crapsolver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xF7A4C6/GoCycle"
	"github.com/valyala/fasthttp"
)

const (
	DEFAULT_UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"
)

var (
	Nodes = GoCycle.New(&[]string{
		//"http://127.0.0.1:80",
		"https://node01.nikolahellatrigger.solutions",
		"https://node02.nikolahellatrigger.solutions",
		"https://node03.nikolahellatrigger.solutions",
	})
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

func NewSolver(apiKey string, nodeList ...string) (*Solver, error) {
	if len(nodeList) > 0 {
		Nodes = GoCycle.New(&nodeList)
	}

	if apiKey == "" {
		return nil, errors.New("please provide api-key")
	}

	if !strings.HasPrefix(apiKey, "user:") {
		return nil, errors.New("apikey format is invalid, be sure to provide 'user:xxxxxxxxxx' format")
	}

	return &Solver{
		ApiKey:   apiKey,
		Client:   client,
		WaitTime: 3 * time.Second,
	}, nil
}

// delete node address
func RemoveNode(url string) error {
	if !Nodes.IsInList(url) {
		return fmt.Errorf("node not in list")
	}

	Nodes.Remove(url)
	return nil
}

// Edit the delay between check task result
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

// Create captcha task
func (S *Solver) NewTask(config *TaskConfig, Server string) (resp *TaskResponse, err error) {
	if config.Domain == "" || config.SiteKey == "" || config.Proxy == "" {
		return nil, errors.New("please provide site-key, domain & proxy")
	}

	if config.TaskType == TASKTYPE_NORMAL {
		return nil, errors.New("normal task is disabled")
	}

	if config.Turbo && (config.TurboSt == 0 || config.TurboSt > 10000) {
		config.TurboSt = 3000
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf(`%s/api/task/new`, Server))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Add("Authorization", S.ApiKey)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.SetBodyRaw(payload)

	response := fasthttp.AcquireResponse()
	err = client.Do(req, response)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(response)

	if err := json.Unmarshal(response.Body(), &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf(resp.Message)
	}

	return resp, nil
}

// Get captcha task result
func (S *Solver) GetResult(T *TaskResponse, Server string) (resp *CheckResponse, err error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("%s/api/task/%s", Server, T.Data[0].ID))
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Add("Authorization", S.ApiKey)
	req.Header.SetContentTypeBytes(headerContentTypeJson)

	response := fasthttp.AcquireResponse()
	err = client.Do(req, response)
	fasthttp.ReleaseRequest(req)

	if err := json.Unmarshal(response.Body(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// Solve captcha and raise error if can't solve
func (S *Solver) Solve(config *TaskConfig) (string, string, error) {
	server, err := Nodes.Next()
	if err != nil {
		return "", "", err
	}

	task, err := S.NewTask(config, server)
	if err != nil {
		return "", "", err
	}

	if config.Turbo {
		time.Sleep(time.Duration(config.TurboSt) * time.Millisecond)
	} else {
		time.Sleep(7 * time.Second)
	}

	for {
		resp, err := S.GetResult(task, server)
		if err != nil {
			return "", "", err
		}

		switch resp.Data.Status {
		case STATUS_SOLVING:
			time.Sleep(S.WaitTime)
		case STATUS_SOLVED:
			return resp.Data.Token, resp.Data.UserAgent, nil
		case STATUS_ERROR:
			return "", "", fmt.Errorf(resp.Data.Error)
		}
	}
}

// Attempt to solve a captcha indefinitely until success, leave the "retry" parameter at 0 to solve indefinitely, modify it to stop if X errors occur and return a list of errors.
func (S *Solver) SolveUntil(config *TaskConfig, retry ...int) (string, string, error) {
	errs := []error{}

	maxRetry := 0
	if len(retry) > 0 {
		maxRetry = retry[0]
	}

	for {
		if len(errs) > maxRetry && maxRetry != 0 {
			return "", "", fmt.Errorf("max retry reached errors: %v", errs)
		}

		token, ua, err := S.Solve(config)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		return token, ua, nil
	}
}
