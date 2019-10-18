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
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "delete command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("del called")
	},
}

var delOrgCmd = &cobra.Command{
	Use:   "org",
	Short: "organization command",
	Long:  `delete Organizations.`,
	Example:`  del org OrgName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("del org command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "cdporg"{
			cmd.PrintErrf("[Error]: cdporg is default organization,does not delete.")
			return
		}
		err := MyConn.DeleteOrganization(args[0])
		PrintDelResult(cmd, err, "organization")
	},
}

var delEnvCmd = &
	cobra.Command{
	Use:   "env",
	Short: "environment command",
	Long:  `del environment.`,
		Example:`  del env EnvName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New(" del env command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "cdpenv"{
			cmd.PrintErrf("[Error]: cdpenv is default environment,does not delete.")
			return
		}
		err := MyConn.DeleteEnvironment(args[0])
		PrintDelResult(cmd, err, "environment")
	},
}

var delProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "project command",
	Long:  `delete project`,
	Example:`  del project ProName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("delete project command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "cdppro"{
			cmd.PrintErrf("[Error]: cdppro is default project,does not delete.")
			return
		}
		err := MyConn.DeleteProject(args[0])
		PrintDelResult(cmd, err, "project")
	},
}

var delGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "group command",
	Long:  `delete group`,
	Example:`  del group GroupName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("delete group command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MyConn.DeleteGroup(args[0])
		PrintDelResult(cmd, err, "group")
	},
}

var delReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "release command",
	Long:  `delete release`,
	Example:`  del release ReleaseName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("delete release command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MyConn.DeleteRelease(args[0])
		PrintDelResult(cmd, err, "release")
	},
}

var delServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "service command",
	Long:  `delete service`,
	Example:`  del service ServiceName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("delete service command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MyConn.DeleteService(args[0])
		PrintDelResult(cmd, err, "service")
	},
}

var delTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "task command",
	Long:  `delete task`,
	Example: `del task TaskName`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("delete task command requires one arg. its name you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MyConn.DeleteTask(args[0])
		PrintDelResult(cmd, err, "release")
	},
}

var delCronCmd = &cobra.Command{
	Use:   "cron",
	Short: "del cron-task",
	Long:  `delete cron-task`,
	Example:`  del cron CronID`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("[Delete CronTask]remove timedtask command requires one arg. its cron id you need provide.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		delEntryId, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.PrintErrf("[Delete CronTask]: %v", err)
			return
		}
		err = MyConn.RemoveTimedTask(int32(delEntryId))
		if err != nil {
			cmd.PrintErrf("[Delete CronTask]: %v", err)
			return
		}
		PrintDelResult(cmd, err, "cron-task")
	},
}

func init() {
	RootCmd.AddCommand(delCmd)

	delCmd.AddCommand(delOrgCmd)
	delCmd.AddCommand(delEnvCmd)
	delCmd.AddCommand(delProjectCmd)
	delCmd.AddCommand(delGroupCmd)
	delCmd.AddCommand(delReleaseCmd)
	delCmd.AddCommand(delServiceCmd)
	delCmd.AddCommand(delTaskCmd)
	delCmd.AddCommand(delCronCmd)

}
