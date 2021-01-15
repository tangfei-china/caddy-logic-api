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
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"strings"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("logic_api", parseCaddyfileRequestDebugger)
}


func parseCaddyfileRequestDebugger(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var dbg RequestDebugger
	for h.Next() {
		args := h.RemainingArgs()
		if len(args) == 0 {
			dbg.Disabled = false
			return dbg, nil
		}
		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("unsupported argument: %s", arg)
			}
			k := parts[0]
			v := strings.Trim(parts[1], "\"")
			switch k {
			case "tag":
				dbg.Tag = v
			case "disabled":
				if !isSwitchArg(v) {
					return nil, fmt.Errorf("%s argument value of %s is unsupported", k, v)
				}
				if isEnabledArg(v) {
					dbg.Disabled = true
				}
			case "response_debug":
				if !isSwitchArg(v) {
					return nil, fmt.Errorf("%s argument value of %s is unsupported", k, v)
				}
				if isEnabledArg(v) {
					dbg.ResponseDebugEnabled = true
				}
			case "url":
				dbg.Url = v
			case "redirect":
				dbg.Redirect = v
			default:
				return nil, fmt.Errorf("unsupported argument: %s", arg)
			}
		}
	}
	return dbg, nil
}

func isEnabledArg(s string) bool {
	if s == "yes" || s == "true" || s == "on" {
		return true
	}
	return false
}

func isSwitchArg(s string) bool {
	if s == "yes" || s == "true" || s == "on" {
		return true
	}
	if s == "no" || s == "false" || s == "off" {
		return true
	}
	return false
}
