package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Temp struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

type Crypto struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

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
	reg            = "register"
	show           = "show"
	startCMD       = "/start"
	rate           = "rate"
	weather        = "weather"
	requestURL     = "http://api.openweathermap.org/data/2.5/weather?appid=%s&q=%s"
	apiKey         = "592ddec5b996fbe95cdd8a649ee3a92b"
	kelvinConstant = 273
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
			case rate:

				crypto := b.getCrypto()

				msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, crypto)

				_, err := bot.Send(msg)
				if err != nil {
					b.l.Error(err)
				}

			case weather:

				temp := b.getTemp()

				msg := tgbotapi.NewMessage(update.CallbackQuery.From.ID, temp)

				_, err := bot.Send(msg)
				if err != nil {
					b.l.Error(err)
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
				tgbotapi.NewInlineKeyboardButtonData("Показать ваши картинки", show),
			},

			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Актуальный курс криптовалют", rate),
			},
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("Погода в Беларуси", weather),
			},
		)

		msg = tgbotapi.NewMessage(message.Chat.ID, "Вот ваша клавиатура")
		msg.ReplyToMessageID = message.MessageID
		msg.ReplyMarkup = keyboard

	default:
		s := strings.Split(message.Text, " ")
		if len(s) == 2 {

			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вы зарегистрированы %s", startCMD))

			err := b.controller.AuthorizeTG(context.Background(), message.From.ID, s[0], s[1])
			if err != nil {
				msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Не получилось авторизоваться, нажмите %s", startCMD))
			}
		} else {
			msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Неправильная команда, нажмите %s", startCMD))

		}
	}

	_, err := b.botAPI.Send(msg)
	if err != nil {
		b.l.Error(err)
	}
}

func (b *Bot) getCrypto() string {
	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price?")
	if err != nil {
		b.l.Error(err)
	}

	var p []Crypto

	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		b.l.Error(err)
	}

	var s string

	for _, i := range p {
		for _, j := range cryptoConst() {
			if j == i.Symbol {
				s += fmt.Sprintf("Price: %s Crypto %s\n", i.Price, i.Symbol)
			}
		}
	}
	return s
}

func cryptoConst() []string {
	return []string{"BTCUSDT", "ETHUSDT", "ZECUSDT", "LTCUSDT", "XLMUSDT"}
}

func (b *Bot) getTemp() string {
	var s string

	for _, city := range cityConst() {
		resp, err := http.Get(fmt.Sprintf(requestURL, apiKey, city))
		if err != nil {
			b.l.Error(err)
		}

		var t Temp
		err = json.NewDecoder(resp.Body).Decode(&t)
		if err != nil {
			b.l.Error(err)
		}
		s += fmt.Sprintf("City: %s, Temp: %f\n", city, kelvinToCelsius(t.Main.Temp))
	}

	return s
}

func cityConst() []string {
	return []string{"Minsk", "Gomel", "Mogilev", "Brest", "Grodno", "Vitebsk"}
}

func kelvinToCelsius(temp float64) float64 {
	return temp - kelvinConstant
}
