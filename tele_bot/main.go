package main

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	// Replace with your actual bot token
	botToken := ""

	// Initialize the bot
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle /start command
	bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Send /open to open Instagram inside Telegram.")
	})

	// Handle /open command with inline button
	bot.Handle("/open", func(c tele.Context) error {
		markup := &tele.ReplyMarkup{}
		instagramButton := markup.URL("Open Instagram", "https://www.instagram.com")

		markup.Inline(
			markup.Row(instagramButton),
		)

		return c.Send("Click the button below to open Instagram inside Telegram:", markup)
	})

	// Start the bot
	bot.Start()
}
