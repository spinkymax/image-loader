package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

type controller interface {
	AuthorizeTG(ctx context.Context, tgID int64, login, password string) error
	GetImageObjects(ctx context.Context, tgID int64) ([]io.Reader, error)
}

type Bot struct {
	controller controller
	botAPI     *tgbotapi.BotAPI
	l          *logrus.Logger
}

const (
	reg      = "register"
	show     = "🍿 Показать картинки 🍿"
	startCMD = "/start"
)

func NewBot(token string, l *logrus.Logger, c controller) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		botAPI:     bot,
		l:          l,
		controller: c,
	}, nil
}

func (b *Bot) StartBot() {
	bot := b.botAPI

	b.l.Infof("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			b.ProcessMessage(update.Message)

		} else if update.CallbackQuery != nil {
			b.l.Info(update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {
			case show:
				images, err := b.controller.GetImageObjects(context.Background(), update.CallbackQuery.From.ID)
				if err != nil {
					b.l.Error(err)
				}

				for i, image := range images {
					byt, err := io.ReadAll(image)
					if err != nil {
						b.l.Error(err)
					}

					msg := tgbotapi.NewPhoto(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileBytes{
						Name:  strconv.Itoa(i) + ".jpg",
						Bytes: byt,
					})

					_, err = b.botAPI.Send(msg)
					if err != nil {
						b.l.Error(err)
					}
				}
			case reg:
				b.l.Info("register")
			}
		}
	}
}

func (b *Bot) ProcessMessage(message *tgbotapi.Message) {
	var msg tgbotapi.MessageConfig

	switch message.Text {
	case startCMD:
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Показать картинки", show),
			},
		)

		msg = tgbotapi.NewMessage(message.Chat.ID, "Вот ваша клавиатура")
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = keyboard

	default:
		s := strings.Split(message.Text, " ")

		msg = tgbotapi.NewMessage(message.Chat.ID, "Вы зарегистрированы")

		err := b.controller.AuthorizeTG(context.Background(), message.From.ID, s[0], s[1])
		if err != nil {
			msg = tgbotapi.NewMessage(message.Chat.ID, "Не получилось авторизоваться")
		}
	}

	_, err := b.botAPI.Send(msg)
	if err != nil {
		b.l.Error(err)
	}
}
