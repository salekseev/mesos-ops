// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/gogo/protobuf/proto"
	"github.com/ondrej-smola/mesos-go-http/lib/client/leader"
	"github.com/ondrej-smola/mesos-go-http/lib/operator/master"
	"github.com/spf13/cobra"
)

type config struct {
	endpoints     []string
	printResponse func(proto.Message)
	ctx           context.Context
	credentials   string
}

type Framework struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var Output []Framework

// frameworksCmd represents the frameworks command
var frameworksCmd = &cobra.Command{
	Use:   "frameworks",
	Short: "Will output frameworks actively registered with the Mesos leader in YAML",
	Run: func(cmd *cobra.Command, args []string) {
		e, err := RootCmd.Flags().GetStringArray("endpoint")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		getFrameworks(e)
	},
}

func getFrameworks(e []string) {

	w := log.NewSyncWriter(os.Stderr)
	logger := log.NewLogfmtLogger(w)

	cfg := &config{
		ctx:       context.Background(),
		endpoints: e,
	}

	getLeader := func(opts ...leader.Opt) *master.Client {
		return master.New(
			leader.New(
				cfg.endpoints,
				leader.WithLogger(logger),
			),
		)
	}

	f, err := getLeader().GetFrameworks(cfg.ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, p := range f.Frameworks {
		Output = append(
			Output,
			Framework{
				ID:   p.GetFrameworkInfo().GetId().GetValue(),
				Name: p.GetFrameworkInfo().GetName()})
	}

	j, err := json.Marshal(&Output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(j))

}

func init() {
	RootCmd.AddCommand(frameworksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// frameworksCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// frameworksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
