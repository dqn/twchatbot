package chatbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Chatbot struct {
	config *Config
	client *twitter.Client
}

type Config struct {
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

func New(c *Config) *Chatbot {
	config := oauth1.NewConfig(c.Account.ConsumerKey, c.Account.ConsumerSecret)
	token := oauth1.NewToken(c.Account.AccessToken, c.Account.AccessTokenSecret)
	client := config.Client(oauth1.NoContext, token)
	t := twitter.NewClient(client)

	return &Chatbot{c, t}
}

func (c *Chatbot) makeResponseToken(crcToken string) string {
	mac := hmac.New(sha256.New, []byte(c.config.Account.ConsumerSecret))
	mac.Write([]byte(crcToken))
	return "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *Chatbot) SendMessage(recipientID, scenarioID string) error {
	scenario, ok := c.config.Scenario[scenarioID]
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

	params := &twitter.DirectMessageEventsNewParams{
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
	}

	_, _, err := c.client.DirectMessages.EventsNew(params)

	return err
}
