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

	consumerCmd := cobra.Command{
		Use:   "consumer",
		Short: "Run message consumer",
	}

	consumerCmd.AddCommand(&cobra.Command{
		Use:   "parse-vtuber",
		Short: "Run message consumer parse vtuber",
		RunE: func(*cobra.Command, []string) error {
			return consumerParseVtuber()
		},
	})

	consumerCmd.AddCommand(&cobra.Command{
		Use:   "parse-agency",
		Short: "Run message consumer parse agency",
		RunE: func(*cobra.Command, []string) error {
			return consumerParseAgency()
		},
	})

	cmd.AddCommand(&consumerCmd)

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
