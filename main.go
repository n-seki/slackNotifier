package main

import (
	"bytes"
	"log"

	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use: "slacknotifier",
	Run: func(cmd *cobra.Command, args []string) {
		token, err := cmd.Flags().GetString("token")
		if err != nil {
			log.Fatal(err)
		}
		channelID, err := cmd.Flags().GetString("channel")
		if err != nil {
			log.Fatal(err)
		}
		header, err := cmd.Flags().GetString("header")
		if err != nil {
			log.Fatal(err)
		}

		var message string
		if len(args) >= 1 {
			message = args[0]
		} else {
			buf := new(bytes.Buffer)
			buf.ReadFrom(cmd.InOrStdin())
			message = buf.String()
		}

		notify(token, channelID, header, message)
	},
}

func init() {
	cobra.OnInitialize()
	cmd.PersistentFlags().StringP("token", "t", "", "slack api token")
	cmd.MarkPersistentFlagRequired("token")
	cmd.PersistentFlags().StringP("channel", "c", "", "target slack channel id")
	cmd.MarkPersistentFlagRequired("channel")
	cmd.PersistentFlags().StringP("header", "", "", "message header")
}

func main() {
	cmd.Execute()
}

func notify(token string, channelID string, header string, message string) {
	api := slack.New(token)

	text := "```\n" + message

	if message[len(message)-1] != '\n' {
		text = text + "\n```"
	} else {
		text = text + "```"
	}
	if len(header) > 0 {
		text = header + "\n\n" + text
	}

	_, _, err := api.PostMessage(
		channelID,
		slack.MsgOptionText(text, false),
	)

	if err != nil {
		log.Fatal(err)
	}
}
