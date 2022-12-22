package cmd

import (
   // "fmt"
    "github.com/spf13/cobra"
	tool "goback/app"
)

func Execute() error {

	var cmd = &cobra.Command{
		Use:   "goback",
		Short: " ", // for now
		Long:  ` `, // same
		Run: func(cmd *cobra.Command, args []string) {
			tool.InitTool()
		},
	}
	return cmd.Execute()
}