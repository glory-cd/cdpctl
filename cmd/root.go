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
  "fmt"
  "github.com/glory-cd/server/client"
  "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "os"
)


var cfgFile string

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
  Use:   "cdpctl",
  Short: "A Continuous Deployment Process Tool",
  Long:  ``,
  // Uncomment the following line if your bare application
  // has an action associated with it:
  //	Run: func(cmd *cobra.Command, args []string) { },
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
    var err error
    cdpattr := client.CDPCClientAttr{CertFile: "cert/server.crt", Address: "localhost:50051"}
    //cdpattr := client.CDPCClientAttr{CertFile: "cert_75/server.crt", Address: "192.168.1.75:30051"}
    MyConn, err = client.NewClient(cdpattr)
    if err != nil {
      cmd.PrintErrf("Conn Server failed. [%s]", err)
      return
    }
  },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  //rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cdpctl.yaml)")


  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  //rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".cdpctl" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".cdpctl")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}

