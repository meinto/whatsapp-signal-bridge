package bridge

import (
	"strings"
)

type MESSAGE_TYPE string

const (
	TEXT_MESSAGE_TYPE     MESSAGE_TYPE = "text"
	IMAGE_MESSAGE_TYPE                 = "image"
	AUDIO_MESSAGE_TYPE                 = "audio"
	VIDEO_MESSAGE_TYPE                 = "video"
	DOCUMENT_MESSAGE_TYPE              = "document"
)

type Attachment struct {
	Bytes []byte
	Type  string
}

type Quote struct {
	MessageType MESSAGE_TYPE
	Body        *string
}

type Message interface {
	ID() string
	SetID(string) Message

	ChatID() string
	SetChatID(string) Message

	ChatName() string
	SetChatName(string) Message

	Sender() string
	SetSender(string) Message

	Body() string
	SetBody(string) Message

	Quote() *Quote
	SetQuote(*Quote) Message

	Attachment() *Attachment
	SetAttachment(*Attachment) Message
}

type message struct {
	id         string
	chatID     string
	chatName   string
	sender     string
	body       string
	quote      *Quote
	attachment *Attachment
}

func PlainMessage() Message {
	return &message{}
}

func ErrorMessage(err error, description ...string) Message {
	msg := ""
	if len(description) > 0 {
		msg += strings.Join(description, "\n")
	}
	msg += "\n" + err.Error()
	return &message{body: msg}
}

func PlainTextMessage(txt string) Message {
	return &message{body: txt}
}

func (m *message) ID() string {
	return m.id
}
func (m *message) SetID(id string) Message {
	m.id = id
	return m
}

func (m *message) ChatID() string {
	return m.chatID
}
func (m *message) SetChatID(id string) Message {
	m.chatID = id
	return m
}

func (m *message) ChatName() string {
	return m.chatName
}
func (m *message) SetChatName(name string) Message {
	m.chatName = name
	return m
}

func (m *message) Sender() string {
	return m.sender
}
func (m *message) SetSender(sender string) Message {
	m.sender = sender
	return m
}

func (m *message) Body() string {
	return m.body
}
func (m *message) SetBody(body string) Message {
	m.body = body
	return m
}

func (m *message) Quote() *Quote {
	return m.quote
}
func (m *message) SetQuote(quote *Quote) Message {
	m.quote = quote
	return m
}

func (m *message) Attachment() *Attachment {
	return m.attachment
}
func (m *message) SetAttachment(a *Attachment) Message {
	m.attachment = a
	return m
}
