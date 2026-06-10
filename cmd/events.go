package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kkweon/flowsavvy/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Create and update events",
}

// eventFlagSet registers the event fields shared by create and update.
func eventFlagSet(fs *pflag.FlagSet) {
	fs.String("title", "", "Event title")
	fs.String("notes", "", "Notes (rich-text HTML)")
	fs.String("start", "", "Start date/time (YYYY-MM-DDThh:mm:ss)")
	fs.String("end", "", "End date/time (YYYY-MM-DDThh:mm:ss, exclusive)")
	fs.String("calendar-id", "", "Calendar id (omit for default calendar)")
	fs.String("location", "", "Location")
	fs.Bool("all-day", false, "All-day event")
	fs.Bool("busy", false, "Busy/free availability")
	fs.String("recurrence", "", "RRULE without the RRULE: prefix (e.g. FREQ=WEEKLY)")
	fs.String("time-zone", "", "IANA time zone id or 'Floating'")
}

// applyEventFlags copies any explicitly-set event flags onto e.
func applyEventFlags(cmd *cobra.Command, e *client.Event) {
	f := cmd.Flags()
	if f.Changed("title") {
		v, _ := f.GetString("title")
		e.SetTitle(v)
	}
	if f.Changed("notes") {
		v, _ := f.GetString("notes")
		e.SetNotes(v)
	}
	if f.Changed("start") {
		v, _ := f.GetString("start")
		e.SetStartDateTime(v)
	}
	if f.Changed("end") {
		v, _ := f.GetString("end")
		e.SetEndDateTime(v)
	}
	if f.Changed("calendar-id") {
		v, _ := f.GetString("calendar-id")
		e.SetCalendarId(v)
	}
	if f.Changed("location") {
		v, _ := f.GetString("location")
		e.SetLocation(v)
	}
	if f.Changed("all-day") {
		v, _ := f.GetBool("all-day")
		e.SetAllDay(v)
	}
	if f.Changed("busy") {
		v, _ := f.GetBool("busy")
		e.SetBusy(v)
	}
	if f.Changed("recurrence") {
		v, _ := f.GetString("recurrence")
		e.SetRecurrenceRule(v)
	}
	if f.Changed("time-zone") {
		v, _ := f.GetString("time-zone")
		e.SetTimeZone(v)
	}
}

func readEventJSON(path string) (*client.Event, error) {
	data, err := readFileOrStdin(path)
	if err != nil {
		return nil, err
	}
	var e client.Event
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("parsing event JSON: %w", err)
	}
	return &e, nil
}

var eventsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an event",
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Flags()
		start, _ := f.GetString("start")
		end, _ := f.GetString("end")
		if start == "" || end == "" {
			return errors.New("--start and --end are required")
		}
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		ev := client.NewEvent(end, start, "event")
		applyEventFlags(cmd, ev)
		created, _, err := api.ItemsEventsAndTasksAPI.CreateEvent(ctx).Event(*ev).Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(created)
	},
}

var eventsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an event (fetches the current event, applies flags, then replaces it)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		f := cmd.Flags()
		occ, _ := f.GetString("occurrence-date")
		scope, _ := f.GetString("scope")

		var ev *client.Event
		if path, _ := f.GetString("from-json"); path != "" {
			ev, err = readEventJSON(path)
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
			if cur.Event == nil {
				return fmt.Errorf("item %s is not an event", args[0])
			}
			ev = cur.Event
			applyEventFlags(cmd, ev)
		}

		putReq := api.ItemsEventsAndTasksAPI.UpdateEvent(ctx, args[0]).Event(*ev)
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

func init() {
	rootCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(eventsCreateCmd, eventsUpdateCmd)

	eventFlagSet(eventsCreateCmd.Flags())
	eventFlagSet(eventsUpdateCmd.Flags())
	eventsUpdateCmd.Flags().String("from-json", "", "Read a full Event JSON from file ('-' for stdin) instead of using flags")
	eventsUpdateCmd.Flags().String("occurrence-date", "", "Occurrence to edit (yyyy-MM-dd) for repeating events")
	eventsUpdateCmd.Flags().String("scope", "", "thisOccurrence | thisAndFutureOccurrences (required for repeating events)")
}
