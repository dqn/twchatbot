package chatbot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo"
)

func (c *Chatbot) CRC(ctx echo.Context) error {
	crcToken := ctx.QueryParam("crc_token")
	mac := hmac.New(sha256.New, []byte(c.config.Account.ConsumerSecret))
	mac.Write([]byte(crcToken))
	responseToken := "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return ctx.JSON(http.StatusOK, map[string]string{
		"response_token": responseToken,
	})
}

func (c *Chatbot) Webhook(ctx echo.Context) error {
	recipientID := "1245969416694587393" // @R8472
	scenarioTag := "s1"

	err := c.SendMessage(recipientID, scenarioTag)
	if err != nil {
		return err
	}

	return ctx.String(http.StatusOK, "ok")
}
