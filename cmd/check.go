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
	"github.com/glory-cd/server/client"
	"github.com/spf13/cobra"
	"errors"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check service",
	Long:  `check service'`,
	Example:`  check service ServiceID ServiceID
  check service -g GroupName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && FlagGroName == ""{
			return errors.New("check must specify services or group")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var checkInfos []client.StaticServiceDetail
		if len(args) > 0 {
			for _, sId := range args {
				checkInfos = append(checkInfos, client.StaticServiceDetail{ServiceID: sId, Op: client.OperateCheck})
			}
		} else {
			services, err := MyConn.GetServices(client.WithGroupNames([]string{FlagGroName}))
			if err != nil {
				cmd.PrintErrf("[Check] get group's service failed. %v\n", err)
				return
			}
			for _, service := range services {
				checkInfos = append(checkInfos, client.StaticServiceDetail{ServiceID: service.ID, Op: client.OperateCheck})
			}

		}

		taskName := "check_" + GetRandomString()
		taskID, err := MyConn.AddTask(taskName, client.WithTaskStatic(checkInfos),client.WithTaskShow(false))
		if err != nil {
			cmd.PrintErrf("[Check] add check task failed. %v\n", err)
		}
		result, err := MyConn.ExecuteTask(taskID)
		if err != nil {
			cmd.PrintErrf("[Check] check service failed. %v\n", err)
		}
		if len(result) >0{
			PrintGetResult(result)
		}


	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&FlagGroName, "group", "g", "","check all services under given group-name.")
}
