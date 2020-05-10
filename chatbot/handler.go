package chatbot

import (
	"net/http"

	"github.com/labstack/echo"
)

func (c *Chatbot) GetWebhook(ctx echo.Context) error {
	crcToken := ctx.QueryParam("crc_token")
	responseToken := c.makeResponseToken(crcToken)

	return ctx.JSON(http.StatusOK, map[string]string{
		"response_token": responseToken,
	})
}

func (c *Chatbot) PostWebhook(ctx echo.Context) error {
	scenarioTag := "s1"

	var event Event
	if err := ctx.Bind(&event); err != nil {
		return err
	}

	for _, v := range event.DirectMessageEvents {
		err := c.SendMessage(v.MessageCreate.SenderID, scenarioTag)
		if err != nil {
			return err
		}
	}

	crcToken := ctx.QueryParam("crc_token")
	responseToken := c.makeResponseToken(crcToken)

	return ctx.JSON(http.StatusOK, map[string]string{
		"response_token": responseToken,
	})
}
