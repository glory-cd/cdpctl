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

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "config info",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if serverHost == "" && serverCert == ""{
			cmd.Printf("[GetConfig]server host: %s\n",hostUrl)
			cmd.Printf("[GetConfig]cert file: %s\n",certFile)
		}else{
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
		}
	},
}

func init() {
	RootCmd.AddCommand(confCmd)
	confCmd.Flags().StringVarP(&serverHost,"server","","","")
	confCmd.Flags().StringVarP(&serverCert,"cert","","","")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// confCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// confCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
