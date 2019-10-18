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

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "upgrade service",
	Long:  `upgrade service'`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && FlagGroName == ""{
			return errors.New("upgrade must specify services or group")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var upgradeInfos []client.UpgradeServiceDetail
		if len(args) > 0 {
			for _, sId := range args {
				upgradeInfos = append(upgradeInfos, client.UpgradeServiceDetail{ServiceID: sId})
			}
		} else {
			services, err := MyConn.GetServices(client.WithGroupNames([]string{FlagGroName}))
			if err != nil {
				cmd.PrintErrf("[StartUp] get group's service failed. %v\n", err)
				return
			}
			for _, service := range services {
				upgradeInfos = append(upgradeInfos, client.UpgradeServiceDetail{ServiceID: service.ID})
			}

		}

		taskName := "upgrade_" + GetRandomString()
		// 获取发布ID
		rid, err := CheckReleaseNameIsLegal(FlagRelName)
		if rid == 0 {
			cmd.PrintErrf("[Upgrade]: release name is empty.")
			return
		}
		if err != nil {
			cmd.PrintErrf("[Upgrade]: release name is illegal. %v", err)
			return
		}
		taskID, err := MyConn.AddTask(taskName, client.WithTaskUpgrade(upgradeInfos), client.WithReleaseId(rid),client.WithTaskShow(true))
		if err != nil {
			cmd.PrintErrf("[Upgrade] add upgrade task failed. %v\n", err)
		}
		result, err := MyConn.ExecuteTask(taskID)
		if err != nil {
			cmd.PrintErrf("[Upgrade] upgrade service failed. %v\n", err)
		}
		if len(result) > 0 {
			PrintGetResult(result)
		}

	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().StringVarP(&FlagGroName, "group", "g", "","upgrade all services under given group-name.")
	upgradeCmd.Flags().StringVarP(&FlagRelName, "release", "r", "","release-name involved in upgrade task.")

}
