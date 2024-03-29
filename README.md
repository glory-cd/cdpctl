### 概述
cdpctl是一个持续部署检查的命令行工具，依附于glory/server和glory/agent服务，共同完成一次部署操作，一次检查操作等操作。
### 要求
* server已部署，并且知道server暴露的rpc服务的地址，例如：host:ip
* server端提供的证书
### 命令行
通过命令行可以完成所有功能的操作。

可以通过基本命令来完成操作指令，就是会较繁琐；也可以通过操作命令，直接完成操作，但是需要基本命令的辅助。

#### 基本命令
##### add命令
* cdpctl add org
```text
添加组织，必须提供组织名称
example: 
        cdpctl add org 无锡不锈钢
```
* cdpctl add env
```text
添加环境, 必须提供环境名称，例如：测试环境，准生产环境等
example: 
      cdpctl add env 测试环境
```
* cdpctl add project
```text
添加项目，必须提供项目名称，例如：交易系统
example: 
       cdpctl add project 交易系统
```
* cdpctl add group
```text
添加分组，必须提供分组名称，分组所属于的组织(-o)，环境(-e)，项目(-p)可选，若不指定，则为系统默认
example:
        cdpctl add group group1
        cdpctl add group group1 -o 无锡不锈钢
        cdpctl add group group1 -o 无锡不锈钢 -e 测试环境
        cdpctl add group group1 -o 无锡不锈钢 -e 测试环境 -p 交易系统      
```
* cdpctl add release
```text
添加发布，必须提供发布名称，版本号，组织名称，项目名称，这几个参数是位置参数，顺序不能改变
发布代码(-c)为可选参数，可以在添加发布的时候通过-c指定，也可以通过set releasecode 命令指定.
发布代码参数的格式为: ModuleName1:RelativePath1;ModuleName2:RelativePath2
example:
       cdpctl add release 发布1 1.0.0 无锡不锈钢 交易系统
       cdpctl add release 发布1 1.0.0 无锡不锈钢 交易系统 -c Gateway:/hfp/1.0.0/Gateway.zip;Member:/hfp/1.0.0/Member.zip
```
* cdpctl add service
```text
 添加服务,必须提供服务名，服务部署路径，服务所在用户和用户密码，服务模块名，以及所属的节点ID，这几个参数是位置参数，顺序不能改变
 example:
        cdpctl add service Gateway1 /home/gateway/Gateway gateway gateway Gateway 10ddf1f2-80de-45c6-a17e-03b6b53dd7f5
```
* cdpctl add task
```text
添加任务其实是本工具的核心操作,几乎可以说，基本命令的功能都是为了这一步的顺利执行而服务的
首先，我们将操作类型分为两大类,即动态的和静态的，动态的包含deploy,upgrade;剩余的皆为静态的
其次，任务是可以针对一个group,也可以是针对一个或者多个服务的.
     * 若是针对一个group,需(-g)指定分组,操作(-o)指定操作类型,此处不支持deploy操作
       注意: (1) 若操作类型为静态，则不需要指定发布版本(-r);若操作类型为动态中的upgrade,则必须指定发布版本(-r)
     * 若任务是真多一个或者多个服务的，可以通过不同的flag指定不同的操作,-d,-u,-s可混合使用.详见下面示例 
       (1) -d: 部署任务,需要指定服务ID和服务模块名称,显而易见，此处服务必须已经添加
       (2) -u: 升级任务
       (3) -s: 其他静态任务       
example:
       cdpctl add task TaskName -o start -g GroupName
       cdpctl add task TaskName -d "ServiceID:ModuleName;ServiceID:ModuleName" -r ReleaseName
       cdpctl add task TaskName -u "ServiceID;ServiceID:CustomUpgradeDir1,CustomUpgradeDir1" -r ReleaseName
       cdpctl add task TaskName -s "ServiceID:OpMode;ServiceID:OpMode"`           
```
##### del命令
* cdpctl del org
```text
   删除组织
```
* cdpctl del env
```text
   删除环境
```
* cdpctl del project
```text
   删除项目
```
* cdpctl del group
```text
   删除分组
```
* cdpctl del release
```text
   删除发布
```
* cdpctl del service
```text
   删除服务
```
* cdpctl del task
```text
   删除任务
```
* cdpctl del cron
```text
   删除定时任务
```
##### get命令
* cdpctl get org
```text
   查询组织信息
```
* cdpctl get env
```text
   查询环境信息
```
* cdpctl get project
```text
   查询项目信息
```
* cdpctl get group
```text
   查询分组信息
```
* cdpctl get node
```text
   查询节点信息
```
* cdpctl get release
```text
   查询发布信息
```
* cdpctl get releasecode
```text
   查询发布代码信息
```
* cdpctl get service
```text
   查询服务信息
```
* cdpctl get task
```text
   查询任务信息
```
* cdpctl get cron
```text
   查询定时任务
```
* cdpctl get work
```text
   查询任务切片
```
* cdpctl get step
```text
   查询切片详情
```
##### exec命令
* cdpctl exec task
```text
   执行任务
```
##### set命令
* cdpctl set node
```text
设置节点别名,必须提供nodeID和别名
example:
     cdpctl set node 10ddf1f2-80de-45c6-a17e-03b6b53dd7f5 node111
```
* cdpctl set releasecode
```text
设置发布代码详情
example:
        cdpctl set releasecode 发布1 "Gateway:/hfp/1.0.0/Gateway.zip;Member:/hfp/1.0.0/Member.zip"
```
* cdpctl set task
```text
设置任务为定时任务.定时任务时间格式参照Linux crontab中的即可
example:
        cdpctl set task task_id task_time
        
```

#### 操作命令
#### deploy命令
deploy命令完成的是一次部署任务,前提条件是需要配置一个yaml文件.格式如下:
```yaml
    taskname: gateway部署6
    groupname: 交易系统3              
    releasename: 发布1
    services:
      - name: gateway24
        dir: /home/gateway/Gateway24
        osuser: gateway
        ospass: gateway
        moudlename: Gateway
        nodeid: 7a6bec18-7b62-40f2-a941-c0e051cf1d86
      - name: gateway25
        dir: /home/gateway/Gateway25
        osuser: gateway
        ospass: gateway
        moudlename: Gateway
        nodeid: e0217e6f-63b2-470c-8968-b18ffa612cd2
```
#### upgrade命令
upgrade命令完成一次升级任务. 可以是针对
#### backup命令
#### rollback命令
#### check命令
#### start命令
#### stop命令
#### restart命令



