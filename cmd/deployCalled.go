/**
* @Author: xhzhang
* @Date: 2019/9/16 16:53
 */
package cmd

import (
	"errors"
	"github.com/glory-cd/server/client"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Deploy struct {
	TaskName    string `yaml:"taskname"`
	GroupName   string `yaml:"groupname"`
	ReleaseName string `yaml:"releasename"`
	Services    []struct {
		Name       string `yaml:"name"`
		Dir        string `yaml:"dir"`
		OsUser     string `yaml:"osuser"`
		OsPass     string `yaml:"ospass"`
		ModuleName string `yaml:"moudlename"`
		AgentID    string `yaml:"nodeid"`
	}
}

type DeployLine struct {
	DFile       string
	TaskName    string
	GroupName   string
	ReleaseName string
	ModuleNames []string
	AgentIds    []string
}

/*
	校验部署文件yaml的合法性
    1. 文件打开是否正常
    2. 文件内容转换成struct是否正常
    3. 各项指标是否合法
    返回值: Deploy结构体，错误信息
*/
func (dl *DeployLine) CheckDeployFileIsLegal(file string) (*Deploy, error) {
	var d Deploy
	dFile, err := ioutil.ReadFile(file)
	if err != nil {
		return &d, err
	}

	err = yaml.Unmarshal(dFile, &d)
	if err != nil {
		return &d, err
	}
	if d.TaskName == "" {
		return &d, errors.New("[taskname] field cannot not be empty.")
	}

	if d.GroupName == "" {
		return &d, errors.New("[groupname] field cannot not be empty.")
	}

	if d.ReleaseName == "" {
		return &d, errors.New("[releasename] field cannot not be empty.")
	}

	if len(d.Services) == 0 {
		return &d, errors.New("[services] cannot not be empty.")
	} else {
		for _, s := range d.Services {
			if s.Name == "" {
				return &d, errors.New("[service->name] cannot not be empty.")
			}

			if s.Dir == "" {
				return &d, errors.New("[service->dir] cannot not be empty.")
			}

			if s.OsUser == "" {
				return &d, errors.New("[service->osuser] cannot not be empty.")
			}

			if s.OsPass == "" {
				return &d, errors.New("[service->ospass] cannot not be empty.")
			}

			if s.ModuleName == "" {
				return &d, errors.New("[service->moudlename] cannot not be empty.")
			}

			if s.AgentID == "" {
				return &d, errors.New("[service->nodeid] cannot not be empty.")
			}
		}
	}

	dl.TaskName = d.TaskName
	dl.GroupName = d.GroupName
	dl.ReleaseName = d.ReleaseName

	for _, s := range d.Services {
		dl.ModuleNames = append(dl.ModuleNames, s.ModuleName)
		dl.AgentIds = append(dl.AgentIds, s.AgentID)
	}
	return &d, nil
}

/*
	校验分组名称是否存在，如果存在返回组织ID
*/
func (dl *DeployLine) CheckGroupName() (int32, error) {
	groups, err := MyConn.GetGroups(client.WithNames([]string{dl.GroupName}))
	if err != nil {
		return 0, err
	}
	return groups.GetID(), nil
}


/*
	校验发布名称是否存在，如果存在返回组织ID
*/
func (dl *DeployLine) CheckReleaseName() (int32, error) {
	releases, err := MyConn.GetReleases(client.WithNames([]string{dl.ReleaseName}))
	if err != nil {
		return 0, err
	}
	return releases.GetID(), nil
}

/*
	校验服务的模块名是否在本次发布中，如果存在返回name和ID对应的Map.
*/
func (dl *DeployLine) CheckModuleName(releaseId int32) (map[string]int32, error) {
	mapMoudleNameId := map[string]int32{}
	mapReleaseCode, err := MyConn.GetReleaseCodeMap(releaseId)
	if err != nil {
		return mapMoudleNameId, err
	}

	for _, mn := range dl.ModuleNames {
		if _, ok := mapReleaseCode[mn]; !ok {
			return mapReleaseCode, errors.New("[" + mn + "]" + "not-exist.")
		}
		mapMoudleNameId[mn] = mapReleaseCode[mn]
	}
	return mapMoudleNameId, nil
}

/*
	校验服务所在agent是否在线.
*/
func (dl *DeployLine) CheckAgent() error {
	agentOnline, err := MyConn.GetAgents(client.WithAgentStatus(true))
	if err != nil {
		return err
	}
	mapaoid := map[string]string{}
	for _, ao := range agentOnline {
		mapaoid[ao.ID] = ao.Alias
	}

	for _, ai := range dl.AgentIds {
		if _, ok := mapaoid[ai]; !ok {
			return errors.New("node [" + ai + "]" + " not-exist or offline.")
		}
	}
	return nil
}
