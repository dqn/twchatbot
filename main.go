package main

import (
	"fmt"
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
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Fprintf(os.Stderr, "Request: %v\n", string(reqBody))
	}))

	e.GET("/webhook", chatbot.GetWebhook)
	e.POST("/webhook", chatbot.PostWebhook)

	e.Logger.Fatal(e.Start(":3000"))

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}
