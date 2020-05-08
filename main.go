package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/yaml.v2"
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

func New(c *ChatbotConfig) *Chatbot {
	config := oauth1.NewConfig(c.Account.ConsumerKey, c.Account.ConsumerSecret)
	token := oauth1.NewToken(c.Account.AccessToken, c.Account.AccessTokenSecret)
	client := config.Client(oauth1.NoContext, token)
	t := twitter.NewClient(client)

	return &Chatbot{c, t}
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

func (c *Chatbot) CRC(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello, World!")
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

func run() error {
	c, err := loadConfig("./config.yml")
	if err != nil {
		return err
	}

	chatbot := New(c)
	recipientID := "1245969416694587393" // @R8472
	scenarioTag := "s1"

	err = chatbot.SendMessage(recipientID, scenarioTag)
	if err != nil {
		return err
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/crc", chatbot.CRC)

	e.Logger.Fatal(e.Start(":3000"))

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
