package service

import (
	"fmt"
	"os"
	"time"
	"x-ui/logger"
	"x-ui/util/common"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var tgBotid int
var isRunning bool

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Get Usage", "get_usage"),
	),
)

type LoginStatus byte

const (
	LoginSuccess LoginStatus = 1
	LoginFail    LoginStatus = 0
)

type Tgbot struct {
	inboundService InboundService
	settingService SettingService
}

func (t *Tgbot) NewTgbot() *Tgbot {
	return new(Tgbot)
}

func (t *Tgbot) Start() error {
	tgBottoken, err := t.settingService.GetTgBotToken()
	if err != nil || tgBottoken == "" {
		logger.Warning("Get TgBotToken failed:", err)
		return err
	}

	tgBotid, err = t.settingService.GetTgBotChatId()
	if err != nil {
		logger.Warning("Get GetTgBotChatId failed:", err)
		return err
	}

	bot, err = tgbotapi.NewBotAPI(tgBottoken)
	if err != nil {
		fmt.Println("Get tgbot's api error:", err)
		return err
	}
	bot.Debug = false

	// listen for TG bot income messages
	if !isRunning {
		logger.Info("Telegram receiver starting")
		go t.OnReceive()
		isRunning = true
	}

	return nil
}

func (t *Tgbot) IsRunnging() bool {
	return isRunning
}

func (t *Tgbot) Stop() {
	bot.StopReceivingUpdates()
	logger.Info("Send Kill to Telegram listener ...")
	isRunning = false
}

func (t *Tgbot) OnReceive() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {

			if update.CallbackQuery != nil {
				// Respond to the callback query, telling Telegram to show the user
				// a message with the data received.
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := bot.Request(callback); err != nil {
					logger.Warning(err)
				}

				// And finally, send a message containing the data received.
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")

				switch update.CallbackQuery.Data {
				case "get_usage":
					msg.Text = "for get your usage send command like this : \n <code>/usage uuid | id</code> \n example : <code>/usage fc3239ed-8f3b-4151-ff51-b183d5182142</code>"
					msg.ParseMode = "HTML"
				}
				if _, err := bot.Send(msg); err != nil {
					logger.Warning(err)
				}
			}

			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "What you need?"
			msg.ReplyMarkup = numericKeyboard
		case "start":
			msg.Text = "Hi :) \n What you need?"
			msg.ReplyMarkup = numericKeyboard

		case "status":
			msg.Text = "bot is ok."

		case "usage":
			msg.Text = t.getClientUsage(update.Message.CommandArguments())
		default:
			msg.Text = "I don't know that command, /help"
			msg.ReplyMarkup = numericKeyboard
		}

		if _, err := bot.Send(msg); err != nil {
			logger.Warning(err)
		}
	}
}

func (t *Tgbot) SendMsgToTgbot(msg string) {
	info := tgbotapi.NewMessage(int64(tgBotid), msg)
	//msg.ReplyToMessageID = int(tgBotid)
	bot.Send(info)
}

func (t *Tgbot) UserLoginNotify(username string, ip string, time string, status LoginStatus) {
	if username == "" || ip == "" || time == "" {
		logger.Warning("UserLoginNotify failed,invalid info")
		return
	}
	var msg string
	// Get hostname
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("get hostname error:", err)
		return
	}
	if status == LoginSuccess {
		msg = fmt.Sprintf("Successfully logged-in to the panel\r\nHostname:%s\r\n", name)
	} else if status == LoginFail {
		msg = fmt.Sprintf("Login to the panel was unsuccessful\r\nHostname:%s\r\n", name)
	}
	msg += fmt.Sprintf("Time:%s\r\n", time)
	msg += fmt.Sprintf("Username:%s\r\n", username)
	msg += fmt.Sprintf("IP:%s\r\n", ip)
	t.SendMsgToTgbot(msg)
}

func (t *Tgbot) getClientUsage(id string) string {
	traffic, err := t.inboundService.GetClientTrafficById(id)
	if err != nil {
		logger.Warning(err)
		return "something wrong!"
	}
	expiryTime := ""
	if traffic.ExpiryTime == 0 {
		expiryTime = "unlimited"
	} else {
		expiryTime = time.Unix((traffic.ExpiryTime / 1000), 0).Format("2006-01-02 15:04:05")
	}
	total := ""
	if traffic.Total == 0 {
		total = "unlimited"
	} else {
		total = common.FormatTraffic((traffic.Total))
	}
	output := fmt.Sprintf("💡 Active: %t\r\n📧 Email: %s\r\n🔼 Upload↑: %s\r\n🔽 Download↓: %s\r\n🔄 Total: %s / %s\r\n📅 Expire in: %s\r\n",
		traffic.Enable, traffic.Email, common.FormatTraffic(traffic.Up), common.FormatTraffic(traffic.Down), common.FormatTraffic((traffic.Up + traffic.Down)),
		total, expiryTime)

	return output
}
