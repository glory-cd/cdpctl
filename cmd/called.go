/**
* @Author: xhzhang
* @Date: 2019/10/9 13:33
 */
package cmd

import (
	"errors"
	"github.com/glory-cd/server/client"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ConnServer(certFile, url string) (conn *client.CDPClient, err error) {
	if certFile != "" && url != ""{
		cdpAttr := client.CDPCClientAttr{CertFile: certFile, Address: url}
		conn, err = client.NewClient(cdpAttr)
	}
	return
}

/*
	添加发布时，检查发布的代码，将其格式化.
    要求: "name:path;name:path"
*/
func CheckReleaseCodes(codes string) (client.ReleaseCodeSlice, error) {
	rcObj := client.ReleaseCodeSlice{}
	if codes == "" {
		return rcObj, errors.New("release codes string is empty.")
	}
	for _, rCode := range strings.Split(codes, ";") {
		rCodeSlice := strings.Split(rCode, ":")
		if len(rCodeSlice) != 2 {
			return rcObj, errors.New("release code format-error")
		}
		rcObj = append(rcObj, client.ReleaseCode{CodeName: rCodeSlice[0], CodePath: rCodeSlice[1]})
	}
	return rcObj, nil
}

func CheckReleaseVersion(version string) error {
	if version == "" {
		return errors.New("release version is empty.")
	}
	return nil
}

/*
	添加任务时，检查操作模式是否合法，若合法，返回指定格式
*/
func CheckTaskOpModeIsLegal(opMode string) (client.OpMode, error) {
	if _, ok := OpMap[opMode]; !ok {
		return client.OperateDefault, errors.New("opmode-illegal")
	}
	return OpMap[opMode], nil
}

/*
	校验发布名称是否存在，如果存在返回组织ID，和module的slice
*/
func CheckReleaseNameIsLegal(rName string) (int32, map[string]int32, error) {
	ms := map[string]int32{}
	if rName == "" {
		return 0, ms, nil
	}

	releases, err := MyConn.GetReleases(client.WithNames([]string{rName}))
	if err != nil {
		return 0, ms, err
	}
	relseaseID := releases.GetID()
	releaseCodes, err := MyConn.GetReleaseCodeMap(relseaseID)
	if err != nil {
		return 0, ms, err
	}
	ms = releaseCodes
	return relseaseID, ms, nil
}

/*
	添加任务时，检查操作部署详情字符串是否合法，若合法，返回指定格式
	要求: "serviceid;serviceid"
*/
func CheckTaskDeploysIsLegal(releaseId int32, deploy string) ([]client.DeployServiceDetail, error) {
	if deploy == "" {
		return nil, nil
	}
	var dd []client.DeployServiceDetail
	for _, dmeta := range strings.Split(deploy, ";") {
		dd = append(dd, client.DeployServiceDetail{ServiceID: dmeta})
	}
	return dd, nil
}

/*
	添加任务时，检查操作升级详情字符串是否合法，若合法，返回指定格式
	要求: "serviceid;serviceid:lib,config/aaa.xml"
*/
func CheckTaskUpgradeIsLegal(releaseId int32, upgrade string) ([]client.UpgradeServiceDetail, error) {
	if upgrade == "" {
		return nil, nil
	}

	var ud []client.UpgradeServiceDetail
	for _, umeta := range strings.Split(upgrade, ";") {
		dslice := strings.Split(umeta, ":")
		if len(dslice) != 2 {
			return ud, errors.New("[" + umeta + "]" + "format-error")
		}
		serviceId := dslice[0]
		customPattern := strings.Split(dslice[1], ",")
		ud = append(ud, client.UpgradeServiceDetail{ServiceID: serviceId, CustomUpgradePattern: customPattern})
	}
	return ud, nil
}

/*
	添加任务时，检查操作升级详情字符串是否合法，若合法，返回指定格式
	要求: "serviceid:op;serviceid:op"
*/
func CheckTaskStaticIsLegal(releaseId int32, static string) ([]client.StaticServiceDetail, error) {
	if static == "" {
		return nil, nil
	}
	var sd []client.StaticServiceDetail
	for _, smeta := range strings.Split(static, ";") {
		ssliece := strings.Split(smeta, ";")
		if len(ssliece) != 2 {
			return sd, errors.New("[" + smeta + "]" + "format-error")
		}
		serviceid := ssliece[0]
		serviceop := ssliece[1]
		if _, ok := OpMap[serviceop]; !ok {
			return sd, errors.New("[" + serviceid + "]" + "serviceOp-error")
		}
		sd = append(sd, client.StaticServiceDetail{ServiceID: serviceid, Op: OpMap[serviceop]})
	}
	return sd, nil
}

/*
	change string slice to int slice. so string slice value must be int string.
*/
func String2IntSlice(s []string) ([]int32, error) {
	i := []int32{}
	for _, si := range s {
		ii, err := strconv.Atoi(si)
		if err != nil {
			return i, err
		}
		i = append(i, int32(ii))
	}
	return i, nil
}

/*
	Determine whether strings are all numbers. if is,conert it from string to int32.
*/
func ParseStringIsDigit(ids []string) ([]int32, error) {
	var idNumber []int32
	if len(ids) == 0 {
		return idNumber, errors.New("Not specified release id.")
	}
	for _, id := range ids {
		result, _ := regexp.MatchString("\\d+", id)
		if result {
			reals, err := strconv.Atoi(id)
			if err != nil {
				return idNumber, err
			}
			idNumber = append(idNumber, int32(reals))
		} else {
			return idNumber, errors.New("not numbers")
		}
	}
	return idNumber, nil
}

func GetRandomString() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIGKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
