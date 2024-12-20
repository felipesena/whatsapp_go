package main

import (
	"flag"
	"fmt"

	"github.com/felipesena/whatsapp_go/services"
)

var (
	telephoneToVar  string
	replyMessageVar string
	textVar         string
	groupNameVar    string
)

func init() {
	flag.StringVar(&telephoneToVar, "telephone", "", "phone number")
	flag.StringVar(&replyMessageVar, "reply", "", "reply message")
	flag.StringVar(&textVar, "text", "", "text")
	flag.StringVar(&groupNameVar, "group", "", "group name")

	flag.Parse()

	if telephoneToVar == "" {
		fmt.Println("Telephone should be informed")
		return
	}

	if replyMessageVar == "" {
		fmt.Println("Reply message should be informed")
		return
	}

	if textVar == "" {
		fmt.Println("Text message should be informed")
		return
	}
}

func main() {
	fmt.Println(telephoneToVar)
	fmt.Println(replyMessageVar)
	fmt.Println(textVar)
	fmt.Println(groupNameVar)

	serviceMessage, err := services.NewSendService()
	if err != nil {
		fmt.Println("Error connecting to WhatsApp client", err)
		return
	}

	if groupNameVar == "" {
		//SEND MESSAGE TO PERSON
		_, err := serviceMessage.SendMessageToPerson(telephoneToVar, replyMessageVar, textVar)
		if err != nil {
			errorMessage := fmt.Sprintf("Error sending message to %s", telephoneToVar)
			fmt.Println(errorMessage, err)
			return
		}

	} else {
		//SEND MESSAGE TO GROUP
		_, err := serviceMessage.SendMessageToGroup(groupNameVar, telephoneToVar, replyMessageVar, textVar)
		if err != nil {
			errorMessage := fmt.Sprintf("Error sending message in group %s to %s", groupNameVar, telephoneToVar)
			fmt.Println(errorMessage, err)
			return
		}
	}

	serviceMessage.Disconnect()
}
