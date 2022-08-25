package cmd

import (
	"github.com/spf13/cobra"
)

func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "letscrum",
		Short: "Open Source Lightweight Scrum/Agile Project Management System.",
	}

	cmd.AddCommand(GetServerCommand())

	return cmd
}

var cfgFile string
