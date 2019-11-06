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
	"os"
	"github.com/glory-cd/server/client"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart command",
	Long: `restart command`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		MyConn, err = ConnServer(certFile, hostUrl)

		if MyConn == nil  || err != nil{
			cmd.PrintErrf("conn server failed. %s\n", err)
			os.Exit(1)
		}

	},
}

var restartServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "restart service",
	Long:  `restart service'`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && FlagGroName == ""{
			return errors.New("restart service must specify services or group")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var startInfos []client.StaticServiceDetail
		if len(args) > 0 {
			for _, sId := range args {
				startInfos = append(startInfos, client.StaticServiceDetail{ServiceID: sId, Op: client.OperateRestart})
			}
		} else {
			services, err := MyConn.GetServices(client.WithGroupNames([]string{FlagGroName}))
			if err != nil {
				cmd.PrintErrf("[Restart] get group's service failed. %v\n", err)
				return
			}
			for _, service := range services {
				startInfos = append(startInfos, client.StaticServiceDetail{ServiceID: service.ID, Op: client.OperateRestart})
			}

		}

		taskName := "restart_" + GetRandomString()
		taskID, err := MyConn.AddTask(taskName, client.WithTaskStatic(startInfos),client.WithTaskShow(false))
		if err != nil {
			cmd.PrintErrf("[Restart] add restart task failed. %v\n", err)
		}
		result, err := MyConn.ExecuteTask(taskID)
		if err != nil {
			cmd.PrintErrf("[Restart] restart service failed. %v\n", err)
		}
		if len(result) >0{
			PrintGetResult(result)
		}


	},
}

var restartAgentCmd = &cobra.Command{
	Use:   "node",
	Short: "restart node",
	Long: `restart node`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && FlagGroName == ""{
			return errors.New("restart node must specify node or group")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		operateAgentIds := []string{}
		if len(args) > 0 {
			operateAgentIds = append(operateAgentIds, args...)
		} else {
			agentIds, err := MyConn.GetAgentIdFromGroup([]string{FlagGroName})
			if err != nil {
				cmd.PrintErrf("[Restart] restart node failed. %v\n", err)
				return
			}
			operateAgentIds = append(operateAgentIds, agentIds...)
		}

		err := MyConn.OperateAgent("SIGHUP", operateAgentIds...)
		if err != nil {
			cmd.PrintErrf("[Restart] restart node failed. %v\n", err)
			return
		}
		cmd.Println("[Restart] restart node successful.")
	},
}


func init() {
	RootCmd.AddCommand(restartCmd)
	restartCmd.AddCommand(restartServiceCmd)
	restartCmd.AddCommand(restartAgentCmd)

	restartServiceCmd.Flags().StringVarP(&FlagGroName, "group", "g", "","restart all services under given group-name.")
	restartAgentCmd.Flags().StringVarP(&FlagGroName, "group", "g", "","restart all nodes under given group-name.")
}
