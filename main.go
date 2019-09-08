package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

// environment variables used are placed under a truct
type envConfig struct {
	Port              string `envconfig:"PORT" default:"5000"`
	BotToken          string `envconfig:"BOT_TOKEN" required:"true"`
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`
	BotID             string `envconfig:"BOT_ID" required:"true"`
	ChannelID         string `envconfig:"CHANNEL_ID" required:"false"`
}

// added the struct below to verify this app from slack
type slackResponse struct {
	Token     string
	Challenge string
	Type      string
}

func main() {
	//used github.com/kelseyhightower/envconfig library to get all the list of env variables
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		fmt.Printf("[ERROR] Failed to process env var: %s", err)
	}
	// github.com/nlopes/slack library supports most if not all of the api.slack.com REST calls,
	// as well as the Real-Time Messaging protocol over websocket, in a fully managed way.
	api := slack.New(env.BotToken)

	slackListener := &SlackListener{
		api:       api,
		botID:     env.BotID,
		channelID: env.ChannelID,
	}

	go slackListener.ListenAndResponse()

	mux := http.NewServeMux()
	mux.HandleFunc("/learnAmharic", learnAmharic)

	http.ListenAndServe(":5000", mux)
}

func learnAmharic(writer http.ResponseWriter, request *http.Request) {

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
