package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var itemsCmd = &cobra.Command{
	Use:   "items",
	Short: "List, get, and delete items (events and tasks)",
}

var itemsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List items, ordered by id",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		f := cmd.Flags()
		req := api.ItemsEventsAndTasksAPI.ListItems(ctx)
		if v, _ := f.GetString("item-type"); v != "" {
			req = req.ItemType(v)
		}
		if f.Changed("completed") {
			v, _ := f.GetBool("completed")
			req = req.Completed(v)
		}
		if v, _ := f.GetString("list-id"); v != "" {
			req = req.ListId(v)
		}
		if v, _ := f.GetString("calendar-id"); v != "" {
			req = req.CalendarId(v)
		}
		if v, _ := f.GetString("query"); v != "" {
			req = req.Query(v)
		}
		if v, _ := f.GetString("modified-after"); v != "" {
			req = req.ModifiedAfter(v)
		}
		if v, _ := f.GetString("page-token"); v != "" {
			req = req.PageToken(v)
		}
		if v, _ := f.GetInt32("limit"); v != 0 {
			req = req.Limit(v)
		}
		resp, _, err := req.Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp)
	},
}

var itemsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a single event or task by id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		req := api.ItemsEventsAndTasksAPI.GetItem(ctx, args[0])
		if v, _ := cmd.Flags().GetString("occurrence-date"); v != "" {
			req = req.OccurrenceDate(v)
		}
		resp, _, err := req.Execute()
		if err != nil {
			return apiError(err)
		}
		return printJSON(resp.GetActualInstance())
	},
}

var itemsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Permanently delete an item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, api, err := setup()
		if err != nil {
			return err
		}
		req := api.ItemsEventsAndTasksAPI.DeleteItem(ctx, args[0])
		if v, _ := cmd.Flags().GetString("occurrence-date"); v != "" {
			req = req.OccurrenceDate(v)
		}
		if v, _ := cmd.Flags().GetString("scope"); v != "" {
			req = req.Scope(v)
		}
		if _, err := req.Execute(); err != nil {
			return apiError(err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "deleted %s\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(itemsCmd)
	itemsCmd.AddCommand(itemsListCmd, itemsGetCmd, itemsDeleteCmd)

	lf := itemsListCmd.Flags()
	lf.String("item-type", "", "Filter by type: task or event")
	lf.Bool("completed", false, "Only items with this completion status")
	lf.String("list-id", "", "Only tasks in this list ('inbox' for inbox); requires --item-type task")
	lf.String("calendar-id", "", "Only events in this calendar; requires --item-type event")
	lf.String("query", "", "Case-insensitive search over title, notes, and location")
	lf.String("modified-after", "", "Only items modified strictly after this UTC instant (YYYY-MM-DDThh:mm:ssZ)")
	lf.String("page-token", "", "Page token for pagination")
	lf.Int32("limit", 0, "Max items per response (1-200; default 50)")

	itemsGetCmd.Flags().String("occurrence-date", "", "Select one occurrence of a repeating series (yyyy-MM-dd)")

	itemsDeleteCmd.Flags().String("occurrence-date", "", "Occurrence to delete (yyyy-MM-dd); for repeating items")
	itemsDeleteCmd.Flags().String("scope", "", "thisOccurrence | thisAndFutureOccurrences (required for repeating items)")
}
