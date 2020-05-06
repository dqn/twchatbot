package twchatbot

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Chatbot struct {
	ChatbotConfig *ChatbotConfig
	client        *twitter.Client
}

type ChatbotConfig struct {
	Account  Account             `yaml:"account"`
	Scenario map[string]Scenario `yaml:"scenario"`
}

type Account struct {
	ConsumerKey       string `yaml:"consumer_key"`
	ConsumerSecret    string `yaml:"consumer_secret"`
	AccessToken       string `yaml:"access_token"`
	AccessTokenSecret string `yaml:"access_token_secret"`
}

type Scenario struct {
	Text       string     `yaml:"text"`
	QuickReply QuickReply `yaml:"quick_reply"`
}

type QuickReply struct {
	Options []QuickReplyOption `yaml:"options"`
	Default QuickReplyDefault  `yaml:"default"`
}

type QuickReplyOption struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Next        string `yaml:"next"`
}

type QuickReplyDefault struct {
	Text string `yaml:"text"`
	Next string `yaml:"next"`
}

func New(config *ChatbotConfig) *Chatbot {
	httpClient := oauth1.NewConfig(
		config.Account.ConsumerKey,
		config.Account.ConsumerSecret,
	).Client(
		oauth1.NoContext,
		oauth1.NewToken(config.Account.AccessToken, config.Account.AccessTokenSecret),
	)

	return &Chatbot{
		config,
		twitter.NewClient(httpClient),
	}
}

func (c *Chatbot) SendMessage(recipientID, scenarioID string) error {
	scenario, ok := c.ChatbotConfig.Scenario[scenarioID]
	if !ok {
		err := fmt.Errorf("unknown scenario ID: %s", scenarioID)
		return err
	}

	options := make([]twitter.DirectMessageQuickReplyOption, len(scenario.QuickReply.Options))
	for i, v := range scenario.QuickReply.Options {
		options[i] = twitter.DirectMessageQuickReplyOption{
			Label:       v.Label,
			Description: v.Description,
			Metadata:    v.Next,
		}
	}

	_, _, err := c.client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type: "message_create",
			Message: &twitter.DirectMessageEventMessage{
				Target: &twitter.DirectMessageTarget{
					RecipientID: recipientID,
				},
				Data: &twitter.DirectMessageData{
					Text: scenario.Text,
					QuickReply: &twitter.DirectMessageQuickReply{
						Type:    "options",
						Options: options,
					},
				},
			},
		},
	})

	return err
}
