/**
* @Author: xhzhang
* @Date: 2019/9/19 15:33
 */
package cmd

import (
	"github.com/modood/table"
	"github.com/spf13/cobra"
)

func PrintAddResult(cmd *cobra.Command, id interface{}, err error, category string) {
	if err == nil {
		cmd.Printf("[Success]: add %s successful，ID is [%d]", category, id)
	} else {
		cmd.PrintErrf("[Error]: add %s failed. => %v", category, err)
	}
}

func PrintDelResult(cmd *cobra.Command, err error, category string) {
	if err == nil {
		cmd.Printf("delete %s successful", category)
	} else {
		cmd.PrintErrf("[Error]: delete %s failed. => %v", category, err)
	}
}


func PrintGetResult(result interface{})  {
	table.Output(result)
}


