# twchatbot

Make Twitter chatbot

## Installation

[WIP]

## Usage

[WIP]

config.yml:

```yml
account:
  consumer_key: '...'
  consumer_secret: '...'
  access_token: '...'
  access_token_secret: '...'
scenario:
  s1:
    text: Nice to meet you!
    quick_reply:
      options:
        -
          label: Pardon?
          description: Repeat again.
          next: s1
        -
          label: Nice meeting you too!
          description: Go to the next step
          next: s2
      default:
        text: Please choose from the options.
        next: s1
  s2:
    text: Good bye!
    quick_reply:
      options:
        -
          label: Good bye!
```

```go
package main

import (
	"github.com/dqn/twchatbot"
)

func main() {
	var c twchatbot.ChatbotConfig

	// unmarshal YAML

	chatbot := twchatbot.New(&c)

	recipientID := "..."
	scenarioID := "s1"

	err := chatbot.SendMessage(recipientID, scenarioID)
	if err != nil {
		// handle error
	}
}
```
