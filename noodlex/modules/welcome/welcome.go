/*
 *    Copyright © 2020 Haruka Network Development
 *    This file is part of Haruka X.
 *
 *    Haruka X is free software: you can redistribute it and/or modify
 *    it under the terms of the Raphielscape Public License as published by
 *    the Devscapes Open Source Holding GmbH., version 1.d
 *
 *    Haruka X is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    Devscapes Raphielscape Public License for more details.
 *
 *    You should have received a copy of the Devscapes Raphielscape Public License
 */

package welcome

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/sql"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/handlers/Filters"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

//var VALID_WELCOME_FORMATTERS = []string{"first", "last", "fullname", "username", "id", "count", "chatname", "mention"}

// EnumFuncMap map of welcome type to function to execute
var EnumFuncMap = map[int]func(ext.Bot, int, string) (*ext.Message, error){
	sql.TEXT:        ext.Bot.SendMessage,
	sql.BUTTON_TEXT: ext.Bot.SendMessage,
}

// AltEnumFuncMap map of alternative welcome types to function to execute
var AltEnumFuncMap = map[int]func(ext.Bot, int, ext.InputFile) (*ext.Message, error){
	sql.STICKER:  ext.Bot.SendSticker,
	sql.DOCUMENT: ext.Bot.SendDocument,
	sql.PHOTO:    ext.Bot.SendPhoto,
	sql.AUDIO:    ext.Bot.SendAudio,
	sql.VOICE:    ext.Bot.SendVoice,
	sql.VIDEO:    ext.Bot.SendVideo,
}

func send(bot ext.Bot, u *gotgbot.Update, message string, keyboard *ext.InlineKeyboardMarkup, backupMessage string, reply bool) *ext.Message {
	msg := bot.NewSendableMessage(u.EffectiveChat.Id, message)
	msg.ParseMode = parsemode.Html
	if reply {
		msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	}
	msg.ReplyMarkup = keyboard
	m, err := msg.Send()
	if err != nil {
		m, _ = u.EffectiveMessage.ReplyText(backupMessage + "Note: The current message was invalid due to some issues.")
	}
	return m
}

func newMember(bot ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	newMembers := u.EffectiveMessage.NewChatMembers
	welcPrefs := sql.GetWelcomePrefs(strconv.Itoa(chat.Id))
	var firstName = ""
	var fullName = ""
	var username = ""
	var res = ""
	var keyb = make([][]ext.InlineKeyboardButton, 0)

	if welcPrefs.ShouldWelcome {
		for _, mem := range newMembers {
			if mem.Id == bot.Id {
				continue
			}

			if welcPrefs.WelcomeType != sql.TEXT && welcPrefs.WelcomeType != sql.BUTTON_TEXT {
				if welcPrefs.WelcomeType > 1 {
					switch welcomeType := welcPrefs.WelcomeType; welcomeType {
					case 4:
						newPhoto := bot.NewSendablePhoto(chat.Id, welcPrefs.CustomWelcome)
						newPhoto.Photo = bot.NewFileId(welcPrefs.Content)
						_, err := newPhoto.Send()
						if err != nil {
							return nil
						}
						break
					case 7:
						newVideo := bot.NewSendableVideo(chat.Id, welcPrefs.CustomWelcome)
						newVideo.Video = bot.NewFileId(welcPrefs.Content)
						_, err := newVideo.Send()
						if err != nil {
							return nil
						}
						break
					default:
						inputFile := bot.NewFileId(welcPrefs.CustomWelcome)
						_, err := AltEnumFuncMap[welcPrefs.WelcomeType](bot, chat.Id, inputFile)
						if err != nil {
							return err
						}
						break
					}
				} else {
					_, err := EnumFuncMap[welcPrefs.WelcomeType](bot, chat.Id, welcPrefs.CustomWelcome)
					if err != nil {
						return err
					}
				}
			}
			firstName = mem.FirstName
			if len(mem.FirstName) <= 0 {
				firstName = "PersonWithNoName"
			}

			if welcPrefs.CustomWelcome != "" {
				if mem.LastName != "" {
					fullName = firstName + " " + mem.LastName
				} else {
					fullName = firstName
				}
				count, _ := chat.GetMembersCount()
				mention := helpers.MentionHtml(mem.Id, firstName)

				if mem.Username != "" {
					username = "@" + html.EscapeString(mem.Username)
				} else {
					username = mention
				}

				r := strings.NewReplacer(
					"{first}", html.EscapeString(firstName),
					"{last}", html.EscapeString(mem.LastName),
					"{fullname}", html.EscapeString(fullName),
					"{username}", username,
					"{mention}", mention,
					"{count}", strconv.Itoa(count),
					"{chatname}", html.EscapeString(chat.Title),
					"{id}", strconv.Itoa(mem.Id),
					"{rules}", "",
				)
				res = r.Replace(welcPrefs.CustomWelcome)
				buttons := sql.GetWelcomeButtons(strconv.Itoa(chat.Id))
				if strings.Contains(welcPrefs.CustomWelcome, "{rules}") {
					rulesButton := sql.WelcomeButton{
						Id:       0,
						ChatId:   strconv.Itoa(u.EffectiveChat.Id),
						Name:     "Rules",
						Url:      fmt.Sprintf("t.me/%v?start=%v", bot.UserName, u.EffectiveChat.Id),
						SameLine: false,
					}
					buttons = append(buttons, rulesButton)
				}
				keyb = helpers.BuildWelcomeKeyboard(buttons)
			} else {
				r := strings.NewReplacer("{first}", firstName)
				res = r.Replace(sql.DefaultWelcome)
			}

			if welcPrefs.ShouldMute {
				if !sql.IsUserHuman(strconv.Itoa(mem.Id), strconv.Itoa(chat.Id)) {
					if !sql.HasUserClickedButton(strconv.Itoa(mem.Id), strconv.Itoa(chat.Id)) {
						_, _ = bot.RestrictChatMember(chat.Id, mem.Id)
					}
				}
				kb := make([]ext.InlineKeyboardButton, 1)
				kb[0] = ext.InlineKeyboardButton{Text: "Click here to prove you're human", CallbackData: "unmute"}
				keyb = append(keyb, kb)
			}

			keyboard := &ext.InlineKeyboardMarkup{InlineKeyboard: &keyb}
			r := strings.NewReplacer("{first}", firstName)
			sent := send(bot, u, res, keyboard, r.Replace(sql.DefaultWelcome), !welcPrefs.DelJoined)

			if welcPrefs.CleanWelcome != 0 {
				_, _ = bot.DeleteMessage(chat.Id, welcPrefs.CleanWelcome)
				sql.SetCleanWelcome(strconv.Itoa(chat.Id), sent.MessageId)
			}

			if welcPrefs.DelJoined {
				_, _ = u.EffectiveMessage.Delete()
			}
		}
	}
	return nil
}

func unmuteCallback(bot ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	user := u.EffectiveUser
	chat := u.EffectiveChat

	if !sql.IsUserHuman(strconv.Itoa(user.Id), strconv.Itoa(chat.Id)) {
		if !sql.HasUserClickedButton(strconv.Itoa(user.Id), strconv.Itoa(chat.Id)) {
			_, err := bot.UnRestrictChatMember(chat.Id, user.Id)
			if err != nil {
				return err
			}
			go sql.UserClickedButton(strconv.Itoa(user.Id), strconv.Itoa(chat.Id))
			_, _ = bot.AnswerCallbackQueryText(query.Id, "You've proved that you are a human! "+
				"You can now talk in this group.", false)
			return nil
		}
	}

	_, _ = bot.AnswerCallbackQueryText(query.Id, "This action is invalid for you.", false)
	return gotgbot.EndGroups{}
}

// LoadWelcome load welcome module
func LoadWelcome(u *gotgbot.Updater) {
	defer log.Println("Loading module welcome")
	u.Dispatcher.AddHandler(handlers.NewMessage(Filters.NewChatMembers(), newMember))
	u.Dispatcher.AddHandler(handlers.NewCallback("unmute", unmuteCallback))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("welcome", noodlex.BotConfig.Prefix, welcome))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("setwelcome", noodlex.BotConfig.Prefix, setWelcome))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("resetwelcome", noodlex.BotConfig.Prefix, resetWelcome))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("cleanwelcome", noodlex.BotConfig.Prefix, cleanWelcome))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("deljoined", noodlex.BotConfig.Prefix, delJoined))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("welcomemute", noodlex.BotConfig.Prefix, welcomeMute))
}
