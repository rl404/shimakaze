package main

import (
	"github.com/spf13/cobra"

	_ "github.com/rl404/shimakaze/docs"
	"github.com/rl404/shimakaze/internal/utils"
)

// @title Shimakaze API
// @description Shimakaze API.
// @BasePath /
// @schemes http https
func main() {
	cmd := cobra.Command{
		Use:   "shimakaze",
		Short: "Shimakaze API",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Run API server",
		RunE: func(*cobra.Command, []string) error {
			return server()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "consumer",
		Short: "Run message consumer",
		RunE: func(*cobra.Command, []string) error {
			return consumer()
		},
	})

	cronCmd := cobra.Command{
		Use:   "cron",
		Short: "Cron",
	}

	cronCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update old data",
		RunE: func(*cobra.Command, []string) error {
			return cronUpdate()
		},
	})

	cronCmd.AddCommand(&cobra.Command{
		Use:   "fill",
		Short: "Fill missing data",
		RunE: func(*cobra.Command, []string) error {
			return cronFill()
		},
	})

	cmd.AddCommand(&cronCmd)

	if err := cmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}
