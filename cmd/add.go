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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add command",
	Long:  `add one of categories`,
	Run: func(cmd *cobra.Command, args []string) {
		//_ = cmd.Help()
		return
	},
}

var addOrgCmd = &cobra.Command{
	Use:   "org",
	Short: "add organization",
	Long:  `add Organization`,
    Example: "  add org OrganizationName",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("add org command requires one args. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := MyConn.AddOrganization(args[0])
		PrintAddResult(cmd, id, err, "organization")
	},
}

var addEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "add environment",
	Long:  `add environment with name`,
	Example:"  add env EnvironmentName",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("add env command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := MyConn.AddEnvironment(args[0])
		PrintAddResult(cmd, id, err, "environment")
	},
}

var addProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "add project",
	Long:  `add project with name`,
	Example: "  add project ProjectName",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("add project command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := MyConn.AddProject(args[0])
		PrintAddResult(cmd, id, err, "project")
	},
}

var addGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "add group",
	Long:  `add group with name, you can design the organization,the environment,the project to which it belongs`,
    Example:`  add group group_name
  add group GroupName -o OrgName
  add group GroupName -o OrgName -e EnvName
  add group GroupName -o OrgName -e EnvName -p ProName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("add group command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := MyConn.AddGroup(args[0], client.WithOrgName(AddFlagOrgName), client.WithEnvName(AddFlagEnvName), client.WithProName(AddFlagProName))
		PrintAddResult(cmd, id, err, "group")
	},
}

var addReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "add release",
	Long:  `add release with name, version,organization-name,project-name,
you can specify code now or you can specify it with the set command.`,
	Example:`  add release ReleaseName ReleaseVersion OrganizationName ProjectName
  add release ReleaseName ReleaseVersion OrganizationName ProjectName -c name:path;name:path`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 4 {
			return errors.New("add release command requires only 4 args, they are release_name,release_version,organization_name,project_name.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		releaseName := args[0]
		releaseVersion := args[1]
		releaseOrgName := args[2]
		releaseProName := args[3]

		var err error
		rc := client.ReleaseCodeSlice{}
		if addReleaseCodes != "" {
			rc, err = CheckReleaseCodes(addReleaseCodes)
			if err != nil {
				cmd.PrintErrf("[Add]: add release failed. %v", err)
				return
			}
		}

		id, err := MyConn.AddRelease(releaseName, releaseVersion, releaseOrgName, releaseProName, rc)
		PrintAddResult(cmd, id, err, "release")
	},
}

var addServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "add service",
	Long:  `add service with name,directory,os-user,os-password,module-name,agent-id`,
	Example:`  add service ServiceName ServicePath ServiceOsUser ServiceOsPass ServiceModuleName serviceAgentId`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 6 {
			return errors.New("add service command requires 6 args.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		serviceDir := args[1]
		serviceOsUser := args[2]
		serviceOsPass := args[3]
		serviceModuleName := args[4]
		serviceAgentId := args[5]

		id,err := MyConn.AddService(serviceName, serviceDir, serviceOsUser, serviceOsPass, serviceAgentId, serviceModuleName,
			client.WithGroupName(FlagGroName),
			client.WithStopCmd(addServiceStopCmd),
			client.WithPidFile(addServicePidFile),
			client.WithCodePattern(addServiceCodePattern))
		PrintAddResult(cmd, id, err, "service")

	},
}

var addTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "add task",
	Long:  `add task with name and operate information`,
	Example:`  add task TaskName -o start -g GroupName
  add task TaskName -d ServiceID:ModuleName;ServiceID:ModuleName -r ReleaseName
  add task TaskName -u ServiceID;ServiceID:CustomUpgradeDir1,CustomUpgradeDir1 -r ReleaseName
  add task TaskName -s ServiceID:OpMode;ServiceID:OpMode`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("add task command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		op, err := CheckTaskOpModeIsLegal(addTaskOpMode)
		if err != nil {
			cmd.PrintErrf("[Add]: add task failed. %v", err)
			return
		}

		// 获取发布ID
		rid, err := CheckReleaseNameIsLegal(FlagRelName)
		if err != nil {
			cmd.PrintErrf("[Add]: add task failed. %v", err)
			return
		}

		// 检查部署详情
		deployInfos, err := CheckTaskDeploysIsLegal(rid, addTaskDeploy)
		if err != nil {
			cmd.PrintErrf("[Add]: add task failed. %v", err)
			return
		}
		// 检查升级详情
		upgradeInfos, err := CheckTaskUpgradeIsLegal(rid, addTaskUpgrade)
		if err != nil {
			cmd.PrintErrf("[Add]: add task failed. %v", err)
			return
		}

		// 检查静态详情
		staticInfos, err := CheckTaskStaticIsLegal(rid, addTaskStatic)
		if err != nil {
			cmd.PrintErrf("[Add]: add task failed. %v", err)
			return
		}
		id, err := MyConn.AddTask(args[0], client.WithTaskShow(true),client.WithTaskOp(op), client.WithReleaseId(rid), client.WithGroupName(FlagGroName),client.WithTaskDeploy(deployInfos), client.WithTaskUpgrade(upgradeInfos), client.WithTaskStatic(staticInfos))
		PrintAddResult(cmd, id, err, "task")

	},
}



func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addOrgCmd)
	addCmd.AddCommand(addEnvCmd)
	addCmd.AddCommand(addProjectCmd)
	addCmd.AddCommand(addGroupCmd)
	addCmd.AddCommand(addReleaseCmd)
	addCmd.AddCommand(addServiceCmd)
	addCmd.AddCommand(addTaskCmd)

	addGroupCmd.Flags().StringVarP(&AddFlagOrgName, "org", "o", "cdporg", "organization name")
	addGroupCmd.Flags().StringVarP(&AddFlagEnvName, "env", "e", "cdpenv", "environment name")
	addGroupCmd.Flags().StringVarP(&AddFlagProName, "pro", "p", "cdppro", "project name")

	addReleaseCmd.Flags().StringVarP(&addReleaseCodes, "code", "c", "", `release codes.[format]: name:path;name:path`)

	addServiceCmd.Flags().StringVarP(&addServiceCodePattern, "code", "c", "", "service code directories.[format]: code1;code2")
	addServiceCmd.Flags().StringVarP(&addServiceStopCmd, "stopcmd", "s", "", "service stop cmd.")
	addServiceCmd.Flags().StringVarP(&addServicePidFile, "pidfile", "p", "", "service pid-file path.")
	addServiceCmd.Flags().StringVarP(&FlagGroName, "group", "g", "cdpgro", "group name. ")

	addTaskCmd.Flags().StringVarP(&FlagGroName, "group", "g", "cdpgro", "group name. ")
	addTaskCmd.Flags().StringVarP(&FlagRelName, "release", "r", "", "release name.")
	addTaskCmd.Flags().StringVarP(&addTaskOpMode, "op", "o", "", "operate mode.")
	addTaskCmd.Flags().StringVarP(&addTaskDeploy, "deploy", "d", "", `deploy info. [format]: serviceid:moudlename;serviceid:moudlename`)
	addTaskCmd.Flags().StringVarP(&addTaskUpgrade, "upgrade", "u", "", `upgrade info. [format]: serviceid;serviceid:lib,config/aaa.xml`)
	addTaskCmd.Flags().StringVarP(&addTaskStatic, "static", "s", "", `static info. [format]:serviceid:op;serviceid:op`)

}
