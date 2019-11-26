/**
* @Author: xhzhang
* @Date: 2019/11/20 16:06
 */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ReadLine(filePth string) ([]string,error) {
	lines := []string{}
	f, err := os.Open(filePth)
	if err != nil {
		return lines,err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		lines = append(lines,string(line)) //放在错误处理前面，即使发生错误，也会处理已经读取到的数据。
		if err != nil { 				   //遇到任何错误立即返回，并忽略 EOF 错误信息
			if err == io.EOF {
				return lines,nil
			}
			return lines,err
		}
	}
}

var file string

func init() {
	flag.StringVar(&file, "file", "", "test file")
}

func main() {
	flag.Parse()
	fmt.Println(file)
	testLines, err := ReadLine(file)
	if err != nil {
		_ = fmt.Errorf("[Test]: read test file failed. %s", err)
		return
	}
	for _, tl := range testLines {
		ntl := strings.Trim(tl, "\r\n")
		args := strings.Split(ntl, " ")
		fmt.Print(ntl + " => ")
		cmd := exec.Command("cdpctl.exe", args...)
		cmd.Stderr = os.Stderr
		_, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}
