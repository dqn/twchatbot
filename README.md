# twchatbot

Twitter chatbot using Account Activity API

## Installation

[WIP]

## Usage

[WIP]

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
