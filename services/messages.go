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

func initWaCli() (c *whatsmeow.Client, error error) {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func NewSendService() (service serviceMessage, error error) {
	client, err := initWaCli()
	if err != nil {
		return serviceMessage{}, err
	}

	return serviceMessage{
		waCli:         client,
		fakeMessageId: "3EB01FF2E76C149A921B7C",
	}, nil
}

func (service serviceMessage) FindGroupByName(groupName string) (groupInfo types.GroupInfo, err error) {

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

// SendMessageToGroup send a message to a group with the provided group name replying a participant with the number provided in participantToReply
// It will reply a message adding the text provided in `messageToReply` with the text provided in `messageText`
// It return the send response and any error encountered
func (service serviceMessage) SendMessageToGroup(groupName string, participantToReply string, messageToReply string, messageText string) (response whatsmeow.SendResponse, err error) {
	groupInfo, err := service.FindGroupByName(groupName)
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

// SendMessageToPerson send a private message to a person with the provided number
// It will reply a message adding the text provided in messageToReply with the text provided in messageText
// It return the send response and any error encountered
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
