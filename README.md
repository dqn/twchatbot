# twchatbot

Twitter chatbot using Account Activity API

## Usage

The account you use for bot must be able to receive Account Activity.
 - [Getting started with webhooks — Twitter Developers](https://developer.twitter.com/en/docs/accounts-and-users/subscribe-account-activity/guides/getting-started-with-webhooks)

```bash
$ cp config.yml.sample config.yml
```

edit `config.yml`.

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

```bash
$ ./twchatbot
```

If you send a message to the account set up for the bot, the message will be returned!
