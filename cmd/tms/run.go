package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yvcruz/tms"
)

func sendCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send message to User",
		Long:  "send command allow to send a todus message to an user",
		Run: func(cmd *cobra.Command, args []string) {
			to, _ := cmd.Flags().GetString("to")
			message, _ := cmd.Flags().GetString("message")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			t := tms.NewTodusMessageService(tms.TodusMessageServiceConfig{
				Url:      "https://broadcast.mprc.cu/api/v1",
				Username: username,
				Password: password,
				Uid:      "",
			})

			if t.SendMessageToUser(to, message) {
				fmt.Println("Send message body to defined user")
			}
		},
	}

	sendCmd.Flags().String("to", "", "the user to send notification")
	sendCmd.Flags().String("message", "", "the notification message body")
	sendCmd.Flags().String("username", "", "the service username")
	sendCmd.Flags().String("password", "", "the service password")
	return sendCmd
}
