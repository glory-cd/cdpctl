/**
* @Author: xhzhang
* @Date: 2019/10/18 11:05
 */
package main

import (
	"cdpctl/cmd"
	"github.com/spf13/cobra/doc"
	"log"
)

func main() {
	//err := doc.GenMarkdownTree(cmd.RootCmd, "./")
	err := doc.GenReSTTree(cmd.RootCmd, "./")
	if err != nil {
		log.Fatal(err)
	}
}
