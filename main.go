package main

import (
	"fmt"

	"github.com/alecthomas/kingpin/v2"
	"github.com/felipesena/whatsapp_go/services"
)

var (
	telephoneToVar  = kingpin.Flag("telephone", "phone number from the person to reply to").Required().Short('t').String()
	replyMessageVar = kingpin.Flag("reply", "reply message").Required().Short('r').String()
	textVar         = kingpin.Flag("message", "text message").Required().Short('m').String()
	groupNameVar    = kingpin.Flag("group", "(optional) group name if you want to reply a message in a group").Short('g').String()
)

func main() {
	kingpin.Parse()

	serviceMessage, err := services.NewSendService()
	kingpin.FatalIfError(err, "Error connecting to WhatsApp client")

	if *groupNameVar == "" {
		// SEND MESSAGE TO PERSON
		fmt.Println("Sending message to person")
		_, err := serviceMessage.SendMessageToPerson(*telephoneToVar, *replyMessageVar, *textVar)

		errorMessage := fmt.Sprintf("Error sending message to %s", *telephoneToVar)
		kingpin.FatalIfError(err, errorMessage)
	} else {
		// SEND MESSAGE TO GROUP
		fmt.Println("Sending message to group", *groupNameVar)
		_, err := serviceMessage.SendMessageToGroup(*groupNameVar, *telephoneToVar, *replyMessageVar, *textVar)

		errorMessage := fmt.Sprintf("Error sending message in group %s to %s", *groupNameVar, *telephoneToVar)
		kingpin.FatalIfError(err, errorMessage)
	}

	serviceMessage.Disconnect()
}
