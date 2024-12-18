package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type serviceMessage struct {
	waCli         *whatsmeow.Client
	fakeMessageId string
}

func initWaCli() *whatsmeow.Client {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	return client
}

func NewSendService() serviceMessage {
	return serviceMessage{
		waCli:         initWaCli(),
		fakeMessageId: "3EB01FF2E76C149A921B7C",
	}
}

func (service serviceMessage) findGroupByName(groupName string) (groupInfo types.GroupInfo, err error) {

	joinedGroups, err := service.waCli.GetJoinedGroups()
	if err != nil {
		return types.GroupInfo{}, err
	}

	var foundedGroup types.GroupInfo
	for _, element := range joinedGroups {
		if strings.Contains(element.GroupName.Name, groupName) {
			foundedGroup = *element
		}
	}

	if foundedGroup.GroupName.Name == "" {
		return types.GroupInfo{}, errors.New(fmt.Sprintf("Group with name %s was not found", groupName))
	}

	fmt.Println("Group with id", foundedGroup.JID, "and name", foundedGroup.GroupName.Name, "was found")
	return foundedGroup, nil
}

func (service serviceMessage) SendMessageToGroup(groupName string, participantToReply string, messageToReply string, messageText string) (response whatsmeow.SendResponse, err error) {
	groupInfo, err := service.findGroupByName(groupName)
	if err != nil {
		return whatsmeow.SendResponse{}, err
	}

	participant := participantToReply + types.DefaultUserServer
	message := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(messageText),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    proto.String(service.fakeMessageId),
				Participant: proto.String(participant),
				QuotedMessage: &waE2E.Message{
					Conversation: proto.String(messageToReply),
				},
			},
		},
	}

	return service.sendMessage(groupInfo.JID, message)
}

func (service serviceMessage) SendMessageToPerson(participantToReply string, messageToReply string, messageText string) (response whatsmeow.SendResponse, err error) {
	targetJID := types.NewJID(participantToReply, types.DefaultUserServer)

	participant := participantToReply + types.DefaultUserServer

	message := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: proto.String(messageText),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:    proto.String(service.fakeMessageId),
				Participant: proto.String(participant),
				QuotedMessage: &waE2E.Message{
					Conversation: proto.String(messageToReply),
				},
			},
		},
	}

	return service.sendMessage(targetJID, message)
}

func (service serviceMessage) sendMessage(targetJID types.JID, message *waE2E.Message) (response whatsmeow.SendResponse, err error) {
	return service.waCli.SendMessage(context.Background(), targetJID, message)
}

func (service serviceMessage) Disconnect() {
	service.waCli.Disconnect()
}
