package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dqn/twchatbot/chatbot"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/yaml.v2"
)

func loadConfig(path string) (*chatbot.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var c chatbot.Config
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

	chatbot := chatbot.New(c)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/crc", chatbot.CRC)
	e.GET("/webhook", chatbot.Webhook)

	e.Logger.Fatal(e.Start(":3000"))

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
