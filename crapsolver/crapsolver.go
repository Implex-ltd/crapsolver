package crapsolver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewSolver() *Solver {
	return &Solver{
		ServerAddr: SERVER_ADDR,
		Client:     &http.Client{},
	}
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

	response, err := S.Client.Post(fmt.Sprintf(`%s/api/task/new`, S.ServerAddr), "application/json", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (S *Solver) GetResult(T *TaskResponse) (resp *CheckResponse, err error) {
	response, err := S.Client.Get(fmt.Sprintf("%s/api/task/%s", S.ServerAddr, T.Data[0].ID))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &resp); err != nil {
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
		time.Sleep(time.Duration(config.TurboSt))
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
			time.Sleep(1 * time.Second)
		case STATUS_SOLVED:
			return resp.Data.Token, nil
		case STATUS_ERROR:
			return "", fmt.Errorf(resp.Data.Error)
		}
	}
}
