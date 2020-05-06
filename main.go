package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
)

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

func loadConfig(path string) (*ChatbotConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var c ChatbotConfig
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func newTwitterClient(ck, cs, at, as string) *twitter.Client {
	config := oauth1.NewConfig(ck, cs)
	token := oauth1.NewToken(at, as)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func sendMessage(client *twitter.Client, recipientID string, scenario *Scenario) error {
	options := make([]twitter.DirectMessageQuickReplyOption, len(scenario.QuickReply.Options))
	for i, v := range scenario.QuickReply.Options {
		options[i] = twitter.DirectMessageQuickReplyOption{
			Label:       v.Label,
			Description: v.Description,
			Metadata:    v.Next,
		}
	}

	_, _, err := client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
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

func run() error {
	c, err := loadConfig("./config.yml")
	if err != nil {
		return err
	}

	client := newTwitterClient(
		c.Account.ConsumerKey,
		c.Account.ConsumerSecret,
		c.Account.AccessToken,
		c.Account.AccessTokenSecret,
	)

	recipientID := "1245969416694587393" // @R8472
	next := "s1"

	s, ok := c.Scenario[next]
	if !ok {
		err = fmt.Errorf("unknown scenario: %s", next)
		return err
	}

	err = sendMessage(client, recipientID, &s)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
