package twchatbot

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

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

func TestAll(t *testing.T) {
	c, err := loadConfig("./config.yml")
	if err != nil {
		panic(err)
	}

	chatbot := New(c)
	recipientID := "1245969416694587393" // @R8472

	t.Run("sendMessage", func(t *testing.T) {
		scenarioTag := "s1"
		err = chatbot.SendMessage(recipientID, scenarioTag)
		if err != nil {
			t.Fatal(err)
		}
	})
}
