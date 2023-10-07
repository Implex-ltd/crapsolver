package crapsolver

import (
	"github.com/valyala/fasthttp"
	"time"
)

var (
	TASKTYPE_ENTERPRISE = 0
	TASKTYPE_NORMAL     = 1
)

var (
	STATUS_SOLVING = 0
	STATUS_SOLVED  = 1
	STATUS_ERROR   = 2
)

type Solver struct {
	Client   *fasthttp.Client
	WaitTime time.Duration
}

type TaskConfig struct {
	// @domain (str, optional): The domain where the captcha is presented.
	// -> Defaults to "accounts.hcaptcha.com".
	Domain string `json:"domain"`

	// @sitekey (str, optional): The sitekey associated with the captcha.
	// -> Defaults to "2eaf963b-eeab-4516-9599-9daa18cd5138".
	SiteKey string `json:"site_key"`

	// @useragent (str, optional): The user agent to use when making requests.
	// -> Defaults to a common user agent string.
	UserAgent string `json:"user_agent"`

	// @proxy (str, optional): The proxy to use for making requests.
	// -> Defaults to an empty string.
	Proxy string `json:"proxy"`

	// @task_type (TASK_TYPE, optional): The type of captcha-solving task.
	// -> Defaults to TASK_TYPE.TYPE_NORMAL.
	TaskType int `json:"task_type"`

	// @invisible (bool, optional): Whether the captcha is invisible.
	// -> Defaults to False.
	Invisible bool `json:"invisible"`

	// @rqdata (str, optional): Additional request data.
	// -> Defaults to an empty string.
	Rqdata string `json:"rqdata"`

	// @text_free_entry (bool, optional): Whether free text entry is allowed.
	// -> Defaults to False.
	A11YTfe bool `json:"a11y_tfe"`

	// @turbo (bool, optional): Whether turbo mode is enabled.
	// -> Defaults to False.
	Turbo bool `json:"turbo"`

	// @turbo_st (int, optional): The turbo mode submit time in milliseconds.
	// -> Defaults to 3000 (3s).
	TurboSt int `json:"turbo_st"`

	// @hc_accessibility (string, optional): hc_accessibility cookie, instant pass normal website.
	HcAccessibility string `json:"hc_accessibility"`

	// @oneclick_only (bool, optional): If captcha images spawn, task will be stopped and error returned.
	OneclickOnly bool `json:"oneclick_only"`
}

type TaskResponse struct {
	Data []struct {
		CreatedAt  time.Time `json:"CreatedAt"`
		UpdatedAt  time.Time `json:"UpdatedAt"`
		DeletedAt  time.Time `json:"DeletedAt"`
		ID         string    `json:"id"`
		Status     int       `json:"status"`
		Token      string    `json:"token"`
		Error      string    `json:"error"`
		Success    bool      `json:"success"`
		Expiration int       `json:"expiration"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type CheckResponse struct {
	Data struct {
		Error      string `json:"error"`
		Expiration int    `json:"expiration"`
		ID         string `json:"id"`
		Status     int    `json:"status"`
		Success    bool   `json:"success"`
		Token      string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
