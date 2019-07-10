package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scnewma/go-mcsrvstat"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mcsrvstat",
	Short: "mcsrvstat is a command line interface for the mcsrvstat API.",
}

var statusCmd = &cobra.Command{
	Use:   "status SERVER",
	Short: "fetch the status of the server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Printf("Must provide exactly one SERVER to get status for\n")
			os.Exit(1)
		}

		addr := args[0]

		client := mcsrvstat.NewClient(nil)
		status, _, err := client.Status(context.Background(), addr)
		if err != nil {
			fmt.Printf("Failed to get status for server: %s\n", addr)
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		fmt.Println("Status:")
		json.NewEncoder(os.Stdout).Encode(status)
	},
}

func init() {
	addCommands()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func addCommands() {
	rootCmd.AddCommand(statusCmd)
}
