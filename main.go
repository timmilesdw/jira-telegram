package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

var (
	jiraUrl  = flag.String("jiraUrl", "", "Jira instance URL")
	chatId   = flag.String("chatId", "", "Telegram chat to send messages")
	botToken = os.Getenv("BOT_TOKEN")
)

func main() {
	flag.Parse()

	chatId, err := strconv.Atoi(*chatId)
	if err != nil {
		log.Fatalf("Couldn't convert chatId to int")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:     botToken,
		ParseMode: telebot.ModeHTML,
	})
	if err != nil {
		log.Fatalf("Couldn't initialize telegram bot, err: %s", err)
	}

	app := fiber.New()

	app.Post("/webhook", func(c *fiber.Ctx) error {

		event, err := Parse(c.Request())
		if err != nil {
			log.Error("Couldn't parse event")
		}
		str := fmt.Sprintf(
			`<a href="%s">[%s]</a><b>'%s'</b>`,
			fmt.Sprintf("%s/browse/%s", *jiraUrl, event.Issue.Key),
			event.Issue.Key,
			event.Issue.Fields.Summary,
		)

		switch event.WebhookEvent {
		case "jira:issue_updated":
			str = str + "\n\nUpdated "
		case "jira:issue_created":
			str = str + "\n\nCreated "
		}

		str = str + " " + fmt.Sprintf(
			"by <b>%s</b>\nStatus: <b>%s</b>\nAsignee: <b>%s</b>\nPriority: <b>%s</b>",
			event.User.DisplayName,
			event.Issue.Fields.Status.Name,
			event.Issue.Fields.Assignee.DisplayName,
			event.Issue.Fields.Priority.Name,
		)

		bot.Send(telebot.ChatID(chatId), str)

		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":3000")
}
