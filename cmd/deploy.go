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
	"os"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy services",
	Long:  `First, args need one ,it's release name,then can specify one or more service(-s),one group(-g) or define a yaml file that configures the necessary information for deployment(-f).`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("deploy must specify release name")
		}
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		MyConn, err = ConnServer(certFile, hostUrl)

		if MyConn == nil || err != nil {
			cmd.PrintErrf("conn server failed. %s\n", err)
			os.Exit(1)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		// 基本验证
		if len(FlagServiceIds) == 0 && FlagGroName == "" && deployFile == "" {
			cmd.PrintErrf("[Deploy]: must specify services or group or deploy file")
			return
		}


		if (len(FlagServiceIds) > 0 || FlagGroName != "") && deployFile != "" {
			cmd.PrintErrf("[Deploy]: -s(-g) and -f can't coexist")
			return
		}


		// 获取发布ID
		rid, modules, err := CheckReleaseNameIsLegal(args[0])
		if rid == 0 {
			cmd.PrintErrf("[Deploy]: release name is empty.")
			return
		}

		var deployInfos []client.DeployServiceDetail
		var hasAddService []string
		if deployFile == "" {
			var GroupCondition []string
			if FlagGroName != ""{
				GroupCondition = append(GroupCondition,FlagGroName)
			}

			// -s & -g
			services, err := MyConn.GetServices(client.WithServiceIds(FlagServiceIds), client.WithGroupNames(GroupCondition))
			if err != nil {
				cmd.PrintErrf("[Deploy]: get services failed. %s", err)
				return
			}

			// 校验
			for _, s := range services {
				if _, ok := modules[s.ModuleName]; !ok {
					cmd.PrintErrf("[Deploy]: [%s] not in this release")
					return
				}
			}
			for _, service := range services {
				deployInfos = append(deployInfos, client.DeployServiceDetail{ServiceID: service.ID})
			}

		} else {
			dLine := DeployLine{DFile: deployFile}
			d, err := dLine.CheckDeployFileIsLegal(deployFile)
			if err != nil {
				cmd.PrintErrf("[Deploy]: check deploy file failed. %v", err)
				return
			}

			err = dLine.CheckAgent()
			if err != nil {
				cmd.PrintErrf("[Deploy]: check agent failed. %v", err)
				return
			}

			gid, err := dLine.CheckGroupName()
			if err != nil {
				cmd.PrintErrf("[Deploy]: check group failed. %v\n", err)
				return
			}

			if gid == 0 {
				cmd.PrintErrf("[Deploy]: check group[%s] failed. non-existent.\n", dLine.GroupName)
				return
			}

			if err = dLine.CheckModuleName(modules); err != nil {
				cmd.PrintErrf("[Deploy]: check module failed. %s.\n", err)
				return
			}

			for _, s := range d.Services {
				serviceId, err := MyConn.AddService(s.Name, s.Dir, s.OsUser, s.OsPass, s.AgentID, s.ModuleName, client.WithGroupId(gid))
				if err != nil {
					cmd.PrintErrf("[Deploy]: add service[%s] failed. %v\n", s.Name, err)
					cmd.Println("[Deploy]: Rollback immediately...")
					rollBack(cmd, hasAddService, 0, d.TaskName)
					return
				}
				cmd.Printf("[Deploy] add service[%s]->[%s] successful.\n", s.Name, serviceId)
				hasAddService = append(hasAddService, serviceId)
				deployInfos = append(deployInfos, client.DeployServiceDetail{ServiceID: serviceId})
			}
		}

		taskName := "deploy_" + GetRandomString()
		taskId, err := MyConn.AddTask(taskName, client.WithReleaseId(rid), client.WithTaskDeploy(deployInfos), client.WithTaskShow(true))
		if err != nil {
			cmd.PrintErrf("[Deploy]: add task[%s] failed. %v\n", taskName, err)
			rollBack(cmd, hasAddService, int(taskId), taskName)
			return
		}
		cmd.Printf("[Deploy] task is ready, its ID is [%d].\n", taskId)
		result, err := MyConn.ExecuteTask(int32(taskId))
		if err != nil {
			cmd.PrintErrf("[Deploy]: execute task failed. %v", err)
			return
		}
		PrintGetResult(result)

	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVarP(&FlagGroName, "group", "g", "", "deploy all services under given group-name.")
	deployCmd.Flags().StringSliceVarP(&FlagServiceIds, "services", "s", []string{}, "deploy services under given service ids.")
	deployCmd.Flags().StringVarP(&deployFile, "file", "f", "", "deploy according to given yaml file.")
}

// 部署任务回滚
func rollBack(cmd *cobra.Command, hasAddService []string, taskId int, taskName string) {
	if taskId != 0 {
		if err := MyConn.DeleteTask(taskName); err != nil {
			cmd.PrintErrf("[RollBack]: delete task [%d] failed. %v\n", taskId, err)
		} else {
			cmd.Printf("[RollBack]: delete task [%d] successful.\n", taskId)
		}
	}

	for _, s := range hasAddService {
		err := MyConn.DeleteService(s)
		if err != nil {
			cmd.PrintErrf("[RollBack]: delete service [%s] failed. %v\n", s, err)
		} else {
			cmd.Printf("[RollBack]: delete service [%s] successful.\n", s)
		}
	}
}
