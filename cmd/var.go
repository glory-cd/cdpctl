/**
* @Author: xhzhang
* @Date: 2019/10/9 13:20
 */
package cmd

import "github.com/glory-cd/server/client"

//version
var versionString = "0.0.13"

const (
	ServerCertFileKey string = "server.certFile"
	ServerHostUrlKey  string = "server.hostUrl"

	ServerCertFileDefault string = "server.crt"
	ServerHostUrlKeyDefault string = "localhost:50051"
)

var MyConn *client.CDPClient
var certFile string
var hostUrl string

var (
	// add flag
	AddFlagOrgName    string
	AddFlagEnvName    string
	AddFlagProName    string
	FlagGroName    	  string
	FlagRelName       string
	FlagServiceIds    []string

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
	queryModuleNames  []string
	queryOnLine       bool
	queryFlagCronIDs  []int
	queryTaskName     string
	queryExecutionId  int

	//set config
	serverHost string
	serverCert string

	// deploy
	deployFile string
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

type showService struct {
	ID     string
	Name   string
	Dir    string
	HostIp string
}

