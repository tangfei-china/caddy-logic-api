// Copyright 2021 Tang Fei
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logic

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

func init() {
	caddy.RegisterModule(RequestDebugger{})
}

// RequestDebugger is a middleware which displays the content of the request it
// handles. It helps troubleshooting web requests by exposing headers
// (e.g. cookies), URL parameters, etc.
type RequestDebugger struct {
	// Enables or disables the plugin.
	Disabled bool `json:"disabled,omitempty"`
	// Adds a tag to a log message
	Tag string `json:"tag,omitempty"`
	// Adds response buffering and debugging
	ResponseDebugEnabled bool `json:"response_debug_enabled,omitempty"`
	// Request api address returns true or false
	Url string `json:"url,omitempty"`
	// Redirect URL
	Redirect string `json:"redirect,omitempty"`

	logger *zap.Logger
}

type RequestPara struct {
	Address string        `json:"address,omitempty"`
	Host    string        `json:"host,omitempty"`
	Agent   string        `json:"agent,omitempty"`
	Tag     string        `json:"tag,omitempty"`
	Proto   string        `json:"proto,omitempty"`
	Cookies []CookiesPara `json:"cookies,omitempty"`
}

type CookiesPara struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (RequestDebugger) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.logic_api",
		New: func() caddy.Module { return new(RequestDebugger) },
	}
}

// Provision sets up RequestDebugger.
func (dbg *RequestDebugger) Provision(ctx caddy.Context) error {
	// dbg.logger = ctx.Logger(dbg)
	if dbg.logger == nil {
		dbg.logger = initLogger()
	}
	return nil
}

func (dbg RequestDebugger) ServeHTTP(resp http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	if dbg.Disabled {
		return next.ServeHTTP(resp, req)
	}

	para := dbg.setRequest(req)

	client := resty.New()
	var check, _ = client.R().SetBody(para).Post(dbg.Url)
	var result map[string]interface{}
	json.Unmarshal([]byte(check.String()), &result)
	//如果判断假，就重定向指定的地址
	if result["success"] != nil && result["success"] == false {

		// Redirect to external provider
		resp.Header().Set("Cache-Control", "no-store")
		resp.Header().Set("Pragma", "no-cache")
		http.Redirect(resp, req, dbg.Redirect, http.StatusFound)
		return nil

	}


	return next.ServeHTTP(resp, req)
}

func (dbg *RequestDebugger) setRequest(r *http.Request) RequestPara {

	var para RequestPara

	cookies := r.Cookies()

	para.Address = r.RemoteAddr
	para.Agent = r.UserAgent()
	para.Host = r.Host
	para.Proto = r.Proto
	para.Tag = dbg.Tag

	for i := range cookies {
		para.Cookies = append(para.Cookies, CookiesPara{cookies[i].Name, cookies[i].Value})
	}

	//dbg.logger.Debug(
	//	"请求的参数集",
	//	zap.Any("参数集", para),
	//)

	return para
}

func initLogger() *zap.Logger {
	logAtom := zap.NewAtomicLevel()
	logAtom.SetLevel(zapcore.DebugLevel)
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logEncoderConfig.TimeKey = "time"
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEncoderConfig),
		zapcore.Lock(os.Stdout),
		logAtom,
	))
	return logger

}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*RequestDebugger)(nil)
