package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	var tmsCmd = &cobra.Command{
		Use:   "tms",
		Short: "The toDus message service",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	tmsCmd.AddCommand(sendCmd())

	err := tmsCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
