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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverHost string
	serverCert string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configure for tool ",
	Long: `configure the server connection address and certificate path`,
	Run: func(cmd *cobra.Command, args []string) {
		if serverHost != ""{
			viper.Set(ServerHostUrlKey, serverHost)
		}

		if serverCert != ""{
			viper.Set(ServerCertFileKey, serverCert)
		}
		err := viper.WriteConfig()
		if err != nil{
			cmd.PrintErrf("[SetConfig] set config failed. %s\n",err)
			return
		}
		cmd.Println("[SetConfig] set config success.")
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&serverHost,"server","","","")
	configCmd.Flags().StringVarP(&serverCert,"cert","","","")

}
