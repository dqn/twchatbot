package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Account struct {
		ConsumerKey       string `yaml:"consumer_key"`
		ConsumerSecret    string `yaml:"consumer_secret"`
		AccessToken       string `yaml:"access_token"`
		AccessTokenSecret string `yaml:"access_token_secret"`
	} `yaml:"account"`
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var c Config
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

	_, _, err = client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type: "message_create",
			Message: &twitter.DirectMessageEventMessage{
				Target: &twitter.DirectMessageTarget{
					RecipientID: "1245969416694587393", // @R8472
				},
				Data: &twitter.DirectMessageData{
					Text: "test",
					QuickReply: &twitter.DirectMessageQuickReply{
						Type: "options",
						Options: []twitter.DirectMessageQuickReplyOption{
							{Label: "hoge", Description: "abc", Metadata: "1-1"},
							{Label: "fuga", Description: "abc", Metadata: "1-2"},
							{Label: "foo", Description: "abc", Metadata: "1-3"},
							{Label: "bar", Description: "abc", Metadata: "1-4"},
						},
					},
				},
			},
		},
	})
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
