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

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy services",
	Long:  `First, define a yaml file that configures the necessary information for deployment, and then use -f to perform deployment.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("deploy command must specify a deploy file as args, which format is yaml.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		deployFile := args[0]
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

		releaseId, err := dLine.CheckReleaseName()
		if err != nil {
			cmd.PrintErrf("[Deploy]: check release failed. %v\n", err)
			return
		}

		moduleNameIDMap, err := dLine.CheckMoudleName(releaseId)
		if err != nil {
			cmd.PrintErrf("[Deploy]: check module failed. %v\n", err)
			return
		}

		var deployDetails []client.DeployServiceDetail
		var hasAddService []string
		for _, s := range d.Services {
			serviceId, err := MyConn.AddService(s.Name, s.Dir, s.OsUser, s.OsPass, s.StartCMD, s.AgentID, s.MoudleName, client.WithGroupId(gid))
			if err != nil {
				cmd.PrintErrf("[Deploy]: add service[%s] failed. %v\n", s.Name, err)
				cmd.Println("[Deploy]: Rollback immediately...")
				rollBack(cmd, hasAddService, 0, d.TaskName)
				return
			}
			cmd.Printf("[Deploy] add service[%s]->[%s] successful.\n", s.Name, serviceId)
			hasAddService = append(hasAddService, serviceId)
			deployDetails = append(deployDetails, client.DeployServiceDetail{ServiceID: serviceId, ReleaseCodeID: moduleNameIDMap[s.MoudleName]})
		}

		taskId, err := MyConn.AddTask(d.TaskName, client.WithReleaseId(releaseId), client.WithTaskDeploy(deployDetails),client.WithTaskShow(true))
		if err != nil {
			cmd.PrintErrf("[Deploy]: add task[%s] failed. %v\n", d.TaskName, err)
			rollBack(cmd, hasAddService, int(taskId), d.TaskName)
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
