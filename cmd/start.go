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
	"os"
)

// startupCmd represents the startup command
var startupCmd = &cobra.Command{
	Use:   "start",
	Short: "start service",
	Long:  `start service'`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		MyConn, err = ConnServer(certFile, hostUrl)

		if MyConn == nil  || err != nil{
			cmd.PrintErrf("conn server failed. %s\n", err)
			os.Exit(1)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(FlagServiceIds) == 0 && FlagGroName == "" {
			cmd.PrintErrf("[StartUp] backup must specify services or group.\n")
			return
		}

		var GroupCondition []string
		if FlagGroName != ""{
			GroupCondition = append(GroupCondition,FlagGroName)
		}

		services, err := MyConn.GetServices(client.WithGroupNames(GroupCondition),client.WithServiceIds(FlagServiceIds))
		if err != nil {
			cmd.PrintErrf("[StartUp] get service failed. %v\n", err)
			return
		}

		var startInfos []client.StaticServiceDetail

		for _, service := range services{
			startInfos = append(startInfos, client.StaticServiceDetail{ServiceID: service.ID, Op: client.OperateStart})
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
	RootCmd.AddCommand(startupCmd)
	startupCmd.Flags().StringVarP(&FlagGroName, "group", "g", "","startup all services under given group-name.")
	startupCmd.Flags().StringSliceVarP(&FlagServiceIds, "services", "s", []string{}, "backup services under given service ids.")
}
