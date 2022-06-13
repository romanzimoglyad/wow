package model

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	None = iota
	GetChallenge
	GetMessage
	Error
)

type Message struct {
	Type int
	Body string
}

func (m *Message) String() string {
	return fmt.Sprintf("%d|%s", m.Type, m.Body)
}

func FromString(str string) (*Message, error) {
	str = strings.TrimSpace(str)
	parts := strings.Split(str, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("wrong request")
	}
	msgType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cannot parse type")
	}
	return &Message{
		Type: msgType,
		Body: parts[1],
	}, nil
}
