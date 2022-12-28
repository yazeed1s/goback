package cmd

import (
	tool "goback/app"

	"github.com/spf13/cobra"
)

func Execute() error {
	var cmd = &cobra.Command{
		Use:   "goback",
		Short: "goback is a handy commands history browser",
		Long:  `0.0.1`,
		Run: func(cmd *cobra.Command, args []string) {
			tool.InitTool()
		},
	}
	return cmd.Execute()
}
