/**
* @Author: xhzhang
* @Date: 2019/10/9 13:20
 */
package cmd

import "github.com/glory-cd/server/client"

var MyConn *client.CDPClient

var (
	// add flag
	AddFlagOrgName    string
	AddFlagEnvName    string
	AddFlagProName    string
	FlagGroName    	  string
	FlagRelName       string

	addServiceCodePattern string
	addServiceStopCmd     string
	addServicePidFile     string


	addReleaseCodes   string
	addTaskOpMode     string
	addTaskDeploy     string
	addTaskUpgrade    string
	addTaskStatic     string
	addTimedSpec      string

	// query flag
	QueryFlagNames    []string
	QueryFlagIDs      []string
	QueryFlagOrgNames []string
	QueryFlagEnvNames []string
	QueryFlagProNames []string
	QueryFlagGroNames []string
	QueryFlagRelNames []string
	QueryFlagTasNames []string
	queryAgentIds     []string
	queryMoudleNames  []string
	queryOnLine       bool
	queryFlagCronIDs  []int
	queryTaskName     string
	queryExecutionId  int
)



var OpMap = map[string]client.OpMode{"": client.OperateDefault,
	"deploy":   client.OperateDeploy,
	"upgrade":  client.OperateUpgrade,
	"startup":    client.OperateStart,
	"stop":     client.OperateStop,
	"restart":  client.OperateRestart,
	"check":    client.OperateCheck,
	"backup":   client.OperateBackUp,
	"rollback": client.OperateRollBack}


const (
	UsageFormat = "Usage:\n  %s\n\nFlags:\n%s"
	AddUsageFormat = "cdpctl add command [args] [flags]\n\n"
)
