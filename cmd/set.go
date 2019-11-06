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
	"github.com/spf13/cobra"
	"strconv"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set command",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		MyConn, err = ConnServer(certFile, hostUrl)

		if MyConn == nil  || err != nil{
			cmd.PrintErrf("conn server failed. %s\n", err)
			os.Exit(1)
		}

	},
}

var setAgentCmd = &cobra.Command{
	Use:   "node",
	Short: "set node alias",
	Long:  `set node alias`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("need two agrs. set agent alais, must provide node id and its alais name.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := MyConn.SetAgentAlias(args[0], args[1])
		if err != nil {
			cmd.PrintErrf("[Set node Alias] %v", err)
			return
		}
		cmd.Println("set node alias successful.")
	},
}

var setReleaseCodeCmd = &cobra.Command{
	Use:   "release",
	Short: "set release",
	Long:  `set release`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("set release code only for two agrs. they are release name and release code.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		releaseName := args[0]
		releaseCodeString := args[1]

		rid, err := CheckReleaseNameIsLegal(releaseName)
		if err != nil {
			cmd.PrintErrf("[SetRelease] set release code failed. %v", rid)
			return
		}

		rc, err := CheckReleaseCodes(releaseCodeString)
		if err != nil {
			cmd.PrintErrf("[SetRelease] set release code failed. %v", rid)
		}

		err = MyConn.SetReleaseCode(rid, rc)
		if err != nil {
			cmd.PrintErrf("[SetRelease] set release code failed.", err)
			return
		}
		cmd.Println("set release code successful.")
	},
}

var setTimedTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "set task to timed",
	Long: `just set existing task to timed,and specify execution time.
For example:
  set task task_id 
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("set task command requires two arg, they are task_id and execution time,witch Refer to the crontab format under Linux.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		setTaskId, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.PrintErrf("[Set]: set timed-task failed. %v", err)
			return
		}

		//当timeSpec不为空时，认为执行定时任务转换
		if addTimedSpec != "" {
			id, err := MyConn.SetTaskToTimed(int32(setTaskId), args[1])
			if err != nil {
				cmd.PrintErrf("[Set]: set timed-task failed.%v", err)
				return
			}
			cmd.Println("set timed-task successful. CronID: [%d]", id)

		}

	},
}

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.AddCommand(setAgentCmd)
	setCmd.AddCommand(setReleaseCodeCmd)
	setCmd.AddCommand(setTimedTaskCmd)
}
