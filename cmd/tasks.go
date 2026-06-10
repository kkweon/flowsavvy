package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kkweon/flowsavvy/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Create, update, and complete tasks",
}

// taskFlagSet registers the task fields shared by create and update.
func taskFlagSet(fs *pflag.FlagSet) {
	fs.String("title", "", "Task title")
	fs.String("notes", "", "Notes (rich-text HTML)")
	fs.String("list-id", "", "List id ('inbox' for inbox; omit for default list)")
	fs.String("due", "", "Due date/time (YYYY-MM-DDThh:mm:ss) — auto-scheduled tasks")
	fs.Int32("duration", 0, "Duration in minutes — auto-scheduled tasks")
	fs.String("start", "", "Start date/time (YYYY-MM-DDThh:mm:ss) — fixed-time tasks")
	fs.String("end", "", "End date/time (YYYY-MM-DDThh:mm:ss) — fixed-time tasks")
	fs.String("priority", "", "low | normal | high | asap")
	fs.Bool("auto-scheduled", false, "Whether FlowSavvy auto-schedules the task")
	fs.String("scheduling-hours-id", "", "Scheduling hours id — auto-scheduled tasks")
	fs.Int32("min-length", 0, "Minimum scheduled chunk length in minutes")
	fs.String("recurrence", "", "RRULE without the RRULE: prefix (e.g. FREQ=DAILY;INTERVAL=1)")
	fs.String("time-zone", "", "IANA time zone id or 'Floating'")
	fs.String("can-be-started-at", "", "Earliest start (YYYY-MM-DDThh:mm:ss) — auto-scheduled tasks")
}

// applyTaskFlags copies any explicitly-set task flags onto t.
func applyTaskFlags(cmd *cobra.Command, t *client.Task) {
	f := cmd.Flags()
	if f.Changed("title") {
		v, _ := f.GetString("title")
		t.SetTitle(v)
	}
	if f.Changed("notes") {
		v, _ := f.GetString("notes")
		t.SetNotes(v)
	}
	if f.Changed("list-id") {
		v, _ := f.GetString("list-id")
		t.SetListId(v)
	}
	if f.Changed("due") {
		v, _ := f.GetString("due")
		t.SetDueDateTime(v)
	}
	if f.Changed("duration") {
		v, _ := f.GetInt32("duration")
		t.SetDurationMinutes(v)
	}
	if f.Changed("start") {
		v, _ := f.GetString("start")
		t.SetStartDateTime(v)
	}
	if f.Changed("end") {
		v, _ := f.GetString("end")
		t.SetEndDateTime(v)
	}
	if f.Changed("priority") {
		v, _ := f.GetString("priority")
		t.SetPriority(v)
	}
	if f.Changed("auto-scheduled") {
		v, _ := f.GetBool("auto-scheduled")
		t.SetIsAutoScheduled(v)
	}
	if f.Changed("scheduling-hours-id") {
		v, _ := f.GetString("scheduling-hours-id")
		t.SetSchedulingHoursId(v)
	}
	if f.Changed("min-length") {
		v, _ := f.GetInt32("min-length")
		t.SetMinLengthMinutes(v)
	}
	if f.Changed("recurrence") {
		v, _ := f.GetString("recurrence")
		t.SetRecurrenceRule(v)
	}
	if f.Changed("time-zone") {
		v, _ := f.GetString("time-zone")
		t.SetTimeZone(v)
	}
	if f.Changed("can-be-started-at") {
		v, _ := f.GetString("can-be-started-at")
		t.SetCanBeStartedAt(v)
	}
}

func readTaskJSON(path string) (*client.Task, error) {
	data, err := readFileOrStdin(path)
	if err != nil {
		return nil, err
	}
	var t client.Task
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing task JSON: %w", err)
	}
	return &t, nil
}

var tasksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if title, _ := cmd.Flags().GetString("title"); title == "" {
			return errors.New("--title is required")
		}
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		task := client.NewTask("task")
		applyTaskFlags(cmd, task)
		created, _, err := api.ItemsEventsAndTasksAPI.CreateTask(ctx).Task(*task).Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(created)
	},
}

var tasksUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a task (fetches the current task, applies flags, then replaces it)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		f := cmd.Flags()
		occ, _ := f.GetString("occurrence-date")
		scope, _ := f.GetString("scope")

		var task *client.Task
		if path, _ := f.GetString("from-json"); path != "" {
			task, err = readTaskJSON(path)
			if err != nil {
				return err
			}
		} else {
			getReq := api.ItemsEventsAndTasksAPI.GetItem(ctx, args[0])
			if occ != "" {
				getReq = getReq.OccurrenceDate(occ)
			}
			cur, _, gerr := getReq.Execute()
			if gerr != nil {
				return apiError(gerr)
			}
			if cur.Task == nil {
				return fmt.Errorf("item %s is not a task", args[0])
			}
			task = cur.Task
			applyTaskFlags(cmd, task)
		}

		putReq := api.ItemsEventsAndTasksAPI.UpdateTask(ctx, args[0]).Task(*task)
		if occ != "" {
			putReq = putReq.OccurrenceDate(occ)
		}
		if scope != "" {
			putReq = putReq.Scope(scope)
		}
		updated, _, err := putReq.Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(updated)
	},
}

var tasksCompleteCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Mark a task complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		req := api.ItemsEventsAndTasksAPI.CompleteTask(ctx, args[0])
		if v, _ := cmd.Flags().GetString("occurrence-date"); v != "" {
			req = req.OccurrenceDate(v)
		}
		resp, _, err := req.Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp)
	},
}

var tasksUncompleteCmd = &cobra.Command{
	Use:   "uncomplete <id>",
	Short: "Mark a task incomplete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		if _, err := api.ItemsEventsAndTasksAPI.UncompleteTask(ctx, args[0]).Execute(); err != nil {
			return apiError(err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "uncompleted %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tasksCmd)
	tasksCmd.AddCommand(tasksCreateCmd, tasksUpdateCmd, tasksCompleteCmd, tasksUncompleteCmd)

	taskFlagSet(tasksCreateCmd.Flags())
	taskFlagSet(tasksUpdateCmd.Flags())
	tasksUpdateCmd.Flags().String("from-json", "", "Read a full Task JSON from file ('-' for stdin) instead of using flags")
	tasksUpdateCmd.Flags().String("occurrence-date", "", "Occurrence to edit (yyyy-MM-dd) for repeating tasks")
	tasksUpdateCmd.Flags().String("scope", "", "thisOccurrence | thisAndFutureOccurrences (required for repeating tasks)")
	tasksCompleteCmd.Flags().String("occurrence-date", "", "Occurrence to complete (yyyy-MM-dd) for repeating tasks")
}
