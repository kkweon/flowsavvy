package cmd

import "github.com/spf13/cobra"

var calendarsCmd = &cobra.Command{
	Use:   "calendars",
	Short: "Calendar reference data",
}

var calendarsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List calendars",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		resp, _, err := api.CalendarsAPI.ListCalendars(ctx).Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp)
	},
}

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Task-list reference data",
}

var listsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List task lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		resp, _, err := api.ListsAPI.ListLists(ctx).Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp)
	},
}

var schedulingHoursCmd = &cobra.Command{
	Use:   "scheduling-hours",
	Short: "Scheduling-hours reference data",
}

var schedulingHoursListCmd = &cobra.Command{
	Use:   "list",
	Short: "List scheduling hours",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		resp, _, err := api.SchedulingHoursAPI.ListSchedulingHours(ctx).Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp)
	},
}

func init() {
	rootCmd.AddCommand(calendarsCmd, listsCmd, schedulingHoursCmd)
	calendarsCmd.AddCommand(calendarsListCmd)
	listsCmd.AddCommand(listsListCmd)
	schedulingHoursCmd.AddCommand(schedulingHoursListCmd)
}
