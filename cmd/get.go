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
	"github.com/glory-cd/server/client"
	"github.com/spf13/cobra"
	"os"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get command",
	Long:  `get command can obtain information you want.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		MyConn, err = ConnServer(certFile, hostUrl)

		if MyConn == nil || err != nil {
			cmd.PrintErrf("conn server failed. %s\n", err)
			os.Exit(1)
		}

	},
}

var getOrgCmd = &cobra.Command{
	Use:   "org",
	Short: "get organization",
	Long:  `query Organizations.`,
	Run: func(cmd *cobra.Command, args []string) {
		intIDs, err := String2IntSlice(QueryFlagIDs)
		queryOrgs, err := MyConn.GetOrganizations(client.WithInt32Ids(intIDs), client.WithNames(QueryFlagNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get organization failed. %v", err)
			return
		}
		PrintGetResult(queryOrgs)
	},
}

var getEnvCmd = &cobra.Command{
	Use:   "env",
	Short: "get environment",
	Long:  `query environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		intIDs, err := String2IntSlice(QueryFlagIDs)
		queryEnvs, err := MyConn.GetEnvironments(client.WithInt32Ids(intIDs), client.WithNames(QueryFlagNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get environment failed. %v", err)
			return
		}
		PrintGetResult(queryEnvs)
	},
}

var getProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "get project",
	Long:  `query project`,
	Run: func(cmd *cobra.Command, args []string) {
		intIDs, err := String2IntSlice(QueryFlagIDs)
		queryProjects, err := MyConn.GetProjects(client.WithInt32Ids(intIDs), client.WithNames(QueryFlagNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get project failed. %v", err)
			return
		}
		PrintGetResult(queryProjects)
	},
}

var getGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "get group",
	Long:  `query group`,
	Run: func(cmd *cobra.Command, args []string) {
		intIDs, err := String2IntSlice(QueryFlagIDs)
		if err != nil {
			cmd.PrintErrf("[Get]: get group failed. %v", err)
			return
		}
		queryGroups, err := MyConn.GetGroups(client.WithInt32Ids(intIDs), client.WithNames(QueryFlagNames), client.WithOrgNames(QueryFlagOrgNames), client.WithEnvNames(QueryFlagEnvNames), client.WithProNames(QueryFlagProNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get group failed. %v", err)
			return
		}
		PrintGetResult(queryGroups)
	},
}

var getReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "get release",
	Long:  `query release`,
	Run: func(cmd *cobra.Command, args []string) {
		intIDs, err := String2IntSlice(QueryFlagIDs)
		if err != nil {
			cmd.PrintErrf("[Get]: get release failed. %v", err)
			return
		}
		queryReleases, err := MyConn.GetReleases(client.WithInt32Ids(intIDs), client.WithNames(QueryFlagNames), client.WithOrgNames(QueryFlagOrgNames), client.WithProNames(QueryFlagProNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get release failed. %v", err)
			return
		}
		PrintGetResult(queryReleases)
	},
}

var getReleaseCodeCmd = &cobra.Command{
	Use:   "releasecode",
	Short: "get releasecode",
	Long:  `query releasecode`,
	Run: func(cmd *cobra.Command, args []string) {
		releases, err := MyConn.GetReleases(client.WithNames(QueryFlagRelNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get release code failed. %v", err)
			return
		}
		releaseIDs := []int32{}
		for _, r := range releases {
			releaseIDs = append(releaseIDs, r.ID)
		}

		queryReleaseCode, err := MyConn.GetReleaseCodes(releaseIDs)
		if err != nil {
			cmd.PrintErrf("[Get]: get release code failed. %v", err)
			return
		}

		PrintGetResult(queryReleaseCode)
	},
}

var getServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "get service",
	Long:  `get service`,
	Run: func(cmd *cobra.Command, args []string) {
		queryServices, err := MyConn.GetServices(client.WithServiceIds(QueryFlagIDs), client.WithNames(QueryFlagNames), client.WithGroupNames(QueryFlagGroNames), client.WithAgentIds(queryAgentIds), client.WithModuleNames(queryModuleNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get service failed. %v", err)
			return
		}
		if queryFlagShow {
			for _, s := range queryServices {
				cmd.Println("Service:")
				cmd.Printf("\tID=%s\n", s.ID)
				cmd.Printf("\tName=%s\n", s.Name)
				cmd.Printf("\tDir=%s\n", s.Dir)
				cmd.Printf("\tModule=%s\n", s.ModuleName)
				cmd.Printf("\tOsUser=%s\n", s.OsUser)
				cmd.Printf("\tCodes=%s\n", s.CodePattern)
				cmd.Printf("\tPid=%s\n", s.PidFile)
				cmd.Printf("\tStart=%s\n", s.StartCmd)
				cmd.Printf("\tStop=%s\n", s.StopCmd)
				cmd.Printf("\tHost=%s\n", s.HostIp)
				cmd.Printf("\tAgentID=%s\n", s.AgentID)
				cmd.Printf("\tAgentName=%s\n", s.AgentName)
				cmd.Printf("\tGroup=%s\n", s.GroupName)
				cmd.Printf("\tCTime=%s\n", s.CreatTime)
				cmd.Printf("\tUTime=%s\n", s.LatestTime)
			}
		} else {
			var showServices []showService
			for _, s := range queryServices {
				showServices = append(showServices, showService{ID: s.ID, Name: s.Name, Dir: s.Dir, HostIp: s.HostIp, Group: s.GroupName, CreateTime: s.CreatTime, LatestTime: s.LatestTime})
			}
			PrintGetResult(showServices)
		}

	},
}

var getTaskCmd = &cobra.Command{
	Use:   "task",
	Short: "get task",
	Long:  `get task`,
	Run: func(cmd *cobra.Command, args []string) {
		intIDs, err := String2IntSlice(QueryFlagIDs)
		if err != nil {
			cmd.PrintErrf("[Get]: get task failed. %v", err)
			return
		}

		queryTasks, err := MyConn.GetTasks(client.WithInt32Ids(intIDs), client.WithNames(QueryFlagNames), client.WithGroupNames(QueryFlagGroNames), client.WithReleaseNames(QueryFlagRelNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get task failed. %v", err)
			return
		}

		var showTasks client.TaskSlice
		for _, t := range queryTasks {
			if t.IsShow {
				showTasks = append(showTasks, t)
			}
		}

		PrintGetResult(showTasks)
	},
}

var getAgentCmd = &cobra.Command{
	Use:   "node",
	Short: "get node",
	Long:  `query node`,
	Run: func(cmd *cobra.Command, args []string) {
		queryAgents, err := MyConn.GetAgents(client.WithAgentStatus(queryOnLine), client.WithAgentIds(QueryFlagIDs), client.WithNames(QueryFlagNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get node failed. %v", err)
			return
		}
		PrintGetResult(queryAgents)
	},
}

var getCronCmd = &cobra.Command{
	Use:              "cron",
	Short:            "get timed task",
	Long:             `get cron-task`,
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		newCronIDs := []int32{}
		for _, id := range queryFlagCronIDs {
			newCronIDs = append(newCronIDs, int32(id))
		}
		queryCrons, err := MyConn.GetTimedTask(client.WithCronEntryIds(newCronIDs), client.WithTaskNames(QueryFlagTasNames))
		if err != nil {
			cmd.PrintErrf("[Get]: get timed failed. %v", err)
			return
		}
		PrintGetResult(queryCrons)
	},
}

var getExecutionCmd = &cobra.Command{
	Use:   "work",
	Short: "get task's work slice",
	Long:  `get task's work slice`,
	Run: func(cmd *cobra.Command, args []string) {
		ts, err := MyConn.GetTasks(client.WithNames([]string{queryTaskName}))
		if err != nil {
			cmd.PrintErrf("[Get]: get executions failed. %v", err)
			return
		}
		es, err := MyConn.GetTaskExecutions(ts.GetID())
		if err != nil {
			cmd.PrintErrf("[Get]: get executions failed. %v", err)
			return
		}
		PrintGetResult(es)
	},
}

var getExecutionDetailCmd = &cobra.Command{
	Use:   "step",
	Short: "get work's step slice",
	Long:  `get work's step slice`,
	Run: func(cmd *cobra.Command, args []string) {
		eds, err := MyConn.GetTaskExecutionDetails(int32(queryExecutionId))
		if err != nil {
			cmd.PrintErrf("[Get]: get executions detail failed. %v", err)
			return
		}
		if len(eds) > 0 {
			PrintGetResult(eds)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.PersistentFlags().StringSliceVarP(&QueryFlagNames, "name", "n", []string{}, "Obtain category records based on gaven names.")
	getCmd.PersistentFlags().StringSliceVarP(&QueryFlagIDs, "id", "i", []string{}, "Obtain category records based on gaven ids.")

	getCmd.AddCommand(getOrgCmd)
	getCmd.AddCommand(getEnvCmd)
	getCmd.AddCommand(getProjectCmd)
	getCmd.AddCommand(getGroupCmd)
	getCmd.AddCommand(getReleaseCmd)
	getCmd.AddCommand(getReleaseCodeCmd)
	getCmd.AddCommand(getServiceCmd)
	getCmd.AddCommand(getTaskCmd)
	getCmd.AddCommand(getAgentCmd)
	getCmd.AddCommand(getCronCmd)
	getCmd.AddCommand(getExecutionCmd)
	getCmd.AddCommand(getExecutionDetailCmd)

	getGroupCmd.Flags().StringSliceVarP(&QueryFlagOrgNames, "orgs", "o", []string{}, "organization  names")
	getGroupCmd.Flags().StringSliceVarP(&QueryFlagEnvNames, "envs", "e", []string{}, "environment  names")
	getGroupCmd.Flags().StringSliceVarP(&QueryFlagProNames, "pros", "p", []string{}, "project names")

	getReleaseCmd.Flags().StringSliceVarP(&QueryFlagOrgNames, "orgs", "o", []string{}, "organization names")
	getReleaseCmd.Flags().StringSliceVarP(&QueryFlagProNames, "pros", "p", []string{}, "project names")

	getReleaseCodeCmd.Flags().StringSliceVarP(&QueryFlagRelNames, "releases", "r", []string{}, "release names")

	getServiceCmd.Flags().StringSliceVarP(&QueryFlagGroNames, "groups", "g", []string{}, "group names")
	getServiceCmd.Flags().StringSliceVarP(&queryAgentIds, "agents", "a", []string{}, "agent ids")
	getServiceCmd.Flags().StringSliceVarP(&queryModuleNames, "modules", "m", []string{}, "module names")
	getServiceCmd.Flags().BoolVarP(&queryFlagShow, "describe", "d", false, "show service details or not")

	getTaskCmd.Flags().StringSliceVarP(&QueryFlagGroNames, "groups", "g", []string{}, "group name")
	getTaskCmd.Flags().StringSliceVarP(&QueryFlagRelNames, "releases", "r", []string{}, "release name")

	getAgentCmd.Flags().BoolVarP(&queryOnLine, "online", "l", false, "")

	getCronCmd.Flags().IntSliceVarP(&queryFlagCronIDs, "crons", "c", []int{}, "cron ids")
	getCronCmd.Flags().StringSliceVarP(&QueryFlagTasNames, "tasks", "", []string{}, "task name")

	getExecutionCmd.Flags().StringVarP(&queryTaskName, "tasks", "t", "", "task name")

	getExecutionDetailCmd.Flags().IntVarP(&queryExecutionId, "work", "e", 0, "work id")
}
