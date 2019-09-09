package main

import (
	"fmt"
	"strings"

	"github.com/nlopes/slack"
)

// SlackListener data type
type SlackListener struct {
	api       *slack.Client
	botID     string
	channelID string
}

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

// "ListenAndResponse method respond to event changes on slack channel"
func (s *SlackListener) ListenAndResponse(dbConnection dbAccess) {
	rtm := s.api.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev, dbConnection); err != nil {
				fmt.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent, dbConnection dbAccess) error {
	dbConnection.getListofCommands()
	// make sure the BOT ID is mentioned
	if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)) {
		fmt.Println(ev.Msg.Text)
	}

	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || m[0] != "translate" {
		return fmt.Errorf("invalid message")
	}

	// get the actual work to translate
	text := m[1:]
	if len(text) == 0 {
		return fmt.Errorf("No word or text to translate")
	}

	// value is passed to message handler when request is approved.
	if _, _, err := s.api.PostMessage(ev.Channel, slack.MsgOptionText(translateToAmharic(text), true)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
