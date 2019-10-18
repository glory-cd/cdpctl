/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"errors"
	"github.com/glory-cd/server/client"
	"github.com/spf13/cobra"
)

// startupCmd represents the startup command
var startupCmd = &cobra.Command{
	Use:   "start",
	Short: "start service",
	Long:  `start service'`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && FlagGroName == ""{
			return errors.New("startup must specify services or group")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var startInfos []client.StaticServiceDetail
		if len(args) > 0 {
			for _, sId := range args {
				startInfos = append(startInfos, client.StaticServiceDetail{ServiceID: sId, Op: client.OperateStart})
			}
		} else {
			services, err := MyConn.GetServices(client.WithGroupNames([]string{FlagGroName}))
			if err != nil {
				cmd.PrintErrf("[StartUp] get group's service failed. %v\n", err)
				return
			}
			for _, service := range services {
				startInfos = append(startInfos, client.StaticServiceDetail{ServiceID: service.ID, Op: client.OperateStart})
			}

		}

		taskName := "startup_" + GetRandomString()
		taskID, err := MyConn.AddTask(taskName, client.WithTaskStatic(startInfos),client.WithTaskShow(false))
		if err != nil {
			cmd.PrintErrf("[StartUp] add startup task failed. %v\n", err)
		}
		result, err := MyConn.ExecuteTask(taskID)
		if err != nil {
			cmd.PrintErrf("[StartUp] startup service failed. %v\n", err)
		}
		if len(result) >0{
			PrintGetResult(result)
		}


	},
}

func init() {
	rootCmd.AddCommand(startupCmd)
	startupCmd.Flags().StringVarP(&FlagGroName, "group", "g", "","startup all services under given group-name.")

}
