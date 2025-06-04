package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/felipesena/whatsapp_go/services"
)

var (
	app = kingpin.New("whatsapp", "A command-line interface to exploit a WhatsApp Vulnerability")

	searchGroup     = app.Command("search", "Search group from joined groups")
	searchGroupName = searchGroup.Arg("search", "Group name to search").Required().String()

	send             = app.Command("send", "Send a message in WhatsApp")
	sendTelephoneTo  = send.Flag("telephone", "phone number from the person to reply to").Required().Short('t').String()
	sendReplyMessage = send.Flag("reply", "reply message").Required().Short('r').String()
	sendText         = send.Flag("message", "text message").Required().Short('m').String()
	sendGroupName    = send.Flag("group", "(optional) group name if you want to reply a message in a group").Short('g').String()
)

func main() {

	serviceMessage, err := services.NewSendService()
	kingpin.FatalIfError(err, "Error connecting to WhatsApp client")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case searchGroup.FullCommand():
		_, err := serviceMessage.FindGroupByName(*searchGroupName)
		errorMessage := fmt.Sprintf("Error finding group with name %s", *searchGroupName)
		kingpin.FatalIfError(err, errorMessage)
		return
	case send.FullCommand():

		if *sendGroupName == "" {
			// SEND MESSAGE TO PERSON
			fmt.Println("Sending message to person")
			_, err := serviceMessage.SendMessageToPerson(*sendTelephoneTo, *sendReplyMessage, *sendText)

			errorMessage := fmt.Sprintf("Error sending message to %s", *sendTelephoneTo)
			kingpin.FatalIfError(err, errorMessage)
		} else {
			// SEND MESSAGE TO GROUP
			fmt.Println("Sending message to group", *sendGroupName)
			_, err := serviceMessage.SendMessageToGroup(*sendGroupName, *sendTelephoneTo, *sendReplyMessage, *sendText)

			errorMessage := fmt.Sprintf("Error sending message in group %s to %s", *sendGroupName, *sendTelephoneTo)
			kingpin.FatalIfError(err, errorMessage)
		}
	}

	serviceMessage.Disconnect()
}
