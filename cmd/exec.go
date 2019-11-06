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
	"github.com/spf13/cobra"
	"strconv"
	"os"
)



// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "execute one task",
	Long: `For example: exec taskid`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		MyConn, err = ConnServer(certFile, hostUrl)

		if MyConn == nil  || err != nil{
			cmd.PrintErrf("conn server failed. %s\n", err)
			os.Exit(1)
		}

	},
}

var execTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "execute one task",
	Long:  `execute one task`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("exec task command requires one arg, its task id you want to execute.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		execTaskId, err := strconv.Atoi(args[0])
		if err != nil {
			cmd.PrintErrf("[Execute]: %v", err)
			return
		}

		result, err := MyConn.ExecuteTask(int32(execTaskId))
		if err != nil {
			cmd.PrintErrf("[Execute]: %v", err)
			return
		}
		PrintGetResult(result)
	},
}





func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.AddCommand(execTaskCmd)
}
