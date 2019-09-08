package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type envConfig struct {
	// Port is server port to be listened.
	Port string `envconfig:"PORT" default:"5000"`

	// BotToken is bot user token to access to slack API.
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`

	// VerificationToken is used to validate interactive messages from slack.
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`

	// BotID is bot user ID.
	BotID string `envconfig:"BOT_ID" required:"true"`

	// ChannelID is slack channel ID where bot is working.
	// Bot responses to the mention in this channel.
	ChannelID string `envconfig:"CHANNEL_ID" required:"false"`
}

// SlackListener data type
type SlackListener struct {
	api       *slack.Client
	botID     string
	channelID string
}

type slackResponse struct {
	Token     string
	Challenge string
	Type      string
}

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		fmt.Printf("[ERROR] Failed to process env var: %s", err)
	}

	api := slack.New(env.BotToken)

	slackListener := &SlackListener{
		api:       api,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}

	go slackListener.ListenAndResponse()

	mux := http.NewServeMux()
	mux.HandleFunc("/echo", authenticateSlack)
	mux.HandleFunc("/learnAmharic", learnAmharic)

	http.ListenAndServe(":5000", mux)
}

// "ListenAndResponse method respond to event changes on slack channel"
func (s *SlackListener) ListenAndResponse() {
	rtm := s.api.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				fmt.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// value is passed to message handler when request is approved.
	if _, _, err := s.api.PostMessage(ev.Channel, slack.MsgOptionText("Some text", true)); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}

func authenticateSlack(writer http.ResponseWriter, request *http.Request) {

	slackMessage := slackResponse{}

	err := json.NewDecoder(request.Body).Decode(&slackMessage)
	if err != nil {
		fmt.Println("The marshallong has some issue")
	}
	buf, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("This is the message from the slack ", string(buf))
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain")
	writer.Write([]byte(slackMessage.Challenge))
}

func learnAmharic(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while getting request to the bot via /learningAmharic")
	}
	fmt.Println(string(buf))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("success"))
}
