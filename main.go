package main

import (
	"fmt"
	"os"

	"github.com/felipesena/whatsapp_go/services"
)

func main() {

	argsLen := len(os.Args)
	if argsLen > 1 {
		toTelephone := os.Args[1]
		replyMessage := os.Args[2]
		textMessage := os.Args[3]

		fmt.Println(toTelephone)
		fmt.Println(replyMessage)
		fmt.Println(textMessage)
	}

	serviceMessage := services.NewSendService()

	serviceMessage.Disconnect()
}
