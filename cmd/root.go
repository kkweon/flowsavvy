package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/kkweon/flowsavvy/client"
)

var rootCmd = &cobra.Command{
	Use:   "flowsavvy",
	Short: "Command-line client for the FlowSavvy API",
	Long: "flowsavvy is a command-line client for the FlowSavvy API (https://my.flowsavvy.app).\n\n" +
		"Authentication: set FLOWSAVVY_API_KEY to your API key\n" +
		"(Settings → Integrations → API in the FlowSavvy app, requires Pro).",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command and exits non-zero on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// setup returns an authenticated context and an API client, or an error if the
// API key is missing.
func setup() (context.Context, *client.APIClient, error) {
	token := os.Getenv("FLOWSAVVY_API_KEY")
	if token == "" {
		token = os.Getenv("FLOWSAVVY_TOKEN")
	}
	if token == "" {
		return nil, nil, errors.New("FLOWSAVVY_API_KEY is not set (Settings → Integrations → API; requires Pro)")
	}
	cfg := client.NewConfiguration()
	cfg.UserAgent = "flowsavvy-cli"
	ctx := context.WithValue(context.Background(), client.ContextAccessToken, token)
	return ctx, client.NewAPIClient(cfg), nil
}

// printJSON pretty-prints any value as JSON to stdout.
func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// apiError enriches a FlowSavvy API error with the response body when present.
func apiError(err error) error {
	var apiErr *client.GenericOpenAPIError
	if errors.As(err, &apiErr) && len(apiErr.Body()) > 0 {
		return fmt.Errorf("%w: %s", err, string(apiErr.Body()))
	}
	return err
}

// readFileOrStdin reads the named file, or stdin when path is "-".
func readFileOrStdin(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}
