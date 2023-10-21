package crapsolver

import (
	"time"

	"github.com/valyala/fasthttp"
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

	ApiKey string
}

type TaskConfig struct {
	// @domain (str, required): The domain where the captcha is presented.
	Domain string `json:"domain"`

	// @sitekey (str, required): The sitekey associated with the captcha.
	SiteKey string `json:"site_key"`

	// @useragent (str, optional): The user agent to use when making requests.
	// -> Defaults to a common user agent string.
	UserAgent string `json:"user_agent"`

	// @proxy (str, required): The proxy to use for making requests.
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

	// @href (string, optional): href of the actual page where the captcha spawn, get it via motionData.
	// Defaults to https://domain.
	Href string `json:"href"`
}

type TaskDataResponse struct {
	ID         string `json:"id"`
	Status     int    `json:"status"`
	Token      string `json:"token"`
	Error      string `json:"error"`
	Success    bool   `json:"success"`
	Expiration int    `json:"expiration"`
}

type TaskResponse struct {
	Data    []TaskDataResponse `json:"data,omitempty"`
	Message string             `json:"message"`
	Success bool               `json:"success"`
}

type CheckResponse struct {
	Data struct {
		Error      string `json:"error"`
		Expiration int    `json:"expiration"`
		ID         string `json:"id"`
		Status     int    `json:"status"`
		Success    bool   `json:"success"`
		Token      string `json:"token"`
		UserAgent  string `json:"user_agent"`
		Req        string `json:"req"`
	} `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type GetUserResp struct {
	Data struct {
		Error                 string `json:"error"`
		Balance               string `json:"balance"`
		ApiKey                string `json:"id"`
		SolvedHcaptcha        int    `json:"solved_hcaptcha"`
		ThreadUsedHcaptcha    int    `json:"thread_used_hcaptcha"`
		ThreadMaxHcaptcha     int    `json:"thread_max_hcaptcha"`
		BypassRestrictedSites bool   `json:"bypass_restricted_sites"`
	} `json:"data"`
	Success bool `json:"success"`
}

type Restrictions struct {
	// Minimum turbo submit time
	MinSubmitTime int

	// Maximum turbo submit time
	MaxSubmitTime int

	// Domain associated to the site-key
	Domain string

	// If 'a11y_tfe' must be enabled
	AlwaysText bool

	// If 'oneclick_only' must be enabled
	OneclickOnly bool

	// If the site-key is enabled
	Enabled bool

	// Price / 1000 captchas
	Rate float64
}

type CheckRestrictionResp struct {
	Success bool         `json:"success"`
	Data    Restrictions `json:"data"`
}
