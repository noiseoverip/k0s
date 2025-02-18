/*
Copyright 2021 k0s Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package install

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"

	config "github.com/k0sproject/k0s/pkg/apis/v1beta1"
	"github.com/k0sproject/k0s/pkg/constant"
)

type K0sStatus struct {
	Version       string
	Pid           int
	PPid          int
	Role          string
	SysInit       string
	StubFile      string
	Output        string
	Workloads     bool
	Args          []string
	ClusterConfig *config.ClusterConfig
	K0sVars       constant.CfgVars
}

func GetStatusInfo(socketPath string) (status *K0sStatus, err error) {
	status = &K0sStatus{}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	response, err := httpc.Get("http://localhost")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseData, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}
