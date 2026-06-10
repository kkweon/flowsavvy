package cmd

import "github.com/spf13/cobra"

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Read the rendered schedule",
}

var scheduleGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get schedule entries within a date range (max 31 days)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		f := cmd.Flags()
		req := api.ScheduleAPI.GetSchedule(ctx)
		if v, _ := f.GetString("start-date"); v != "" {
			req = req.StartDate(v)
		}
		if v, _ := f.GetString("end-date"); v != "" {
			req = req.EndDate(v)
		}
		resp, _, err := req.Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp)
	},
}

var recalculateCmd = &cobra.Command{
	Use:   "recalculate",
	Short: "Run FlowSavvy's auto-scheduling engine and reschedule auto-scheduled tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		req := api.ScheduleAPI.Recalculate(ctx)
		if cmd.Flags().Changed("reschedule-past-tasks") {
			v, _ := cmd.Flags().GetBool("reschedule-past-tasks")
			req = req.ReschedulePastTasks(v)
		}
		if _, err := req.Execute(); err != nil {
			return apiError(err)
		}
		cmd.Println("recalculated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd, recalculateCmd)
	scheduleCmd.AddCommand(scheduleGetCmd)

	scheduleGetCmd.Flags().String("start-date", "", "Inclusive range start date (yyyy-MM-dd)")
	scheduleGetCmd.Flags().String("end-date", "", "Inclusive range end date (yyyy-MM-dd)")

	recalculateCmd.Flags().Bool("reschedule-past-tasks", false, "Also reschedule past and in-progress tasks")
}
