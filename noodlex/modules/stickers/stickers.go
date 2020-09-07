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

package stickers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"go.uber.org/zap"
)

func stickerID(_ ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	if msg.ReplyToMessage != nil && msg.ReplyToMessage.Sticker != nil {
		msg.ReplyHTMLf("Sticker  ID:\n<code>%v</code>", msg.ReplyToMessage.Sticker.FileId)
	} else {
		msg.ReplyText("Sticker ID not found.")
	}
	return nil
}

func getSticker(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat

	if msg.ReplyToMessage != nil && msg.ReplyToMessage.Sticker != nil && msg.ReplyToMessage.Sticker.IsAnimated == false {
		fileID := msg.ReplyToMessage.Sticker.FileId

		inputFile, r, err := getInputFile(b, fileID, "sticker.png")
		if r != nil {
			defer r.Body.Close()
		}
		if err != nil {
			return err
		}

		newDoc := b.NewSendableDocument(chat.Id, "Sticker")
		newDoc.Document = inputFile
		newDoc.Send()
	} else if msg.ReplyToMessage != nil && msg.ReplyToMessage.Sticker != nil && msg.ReplyToMessage.Sticker.IsAnimated == true {
		fileID := msg.ReplyToMessage.Sticker.FileId

		inputFile, r, err := getInputFile(b, fileID, "sticker.rename")
		if r != nil {
			defer r.Body.Close()
		}
		if err != nil {
			return err
		}

		newDoc := b.NewSendableDocument(chat.Id, "Go to @Stickers bot and rename this file to .tgs then use "+
			"/newanimated or /addsticker and send this file")
		newDoc.Document = inputFile
		newDoc.Send()
	} else {
		msg.ReplyText("Please reply to a sticker for me to upload its PNG.")
	}
	return nil
}

func kangSticker(b ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	user := u.EffectiveUser
	packnum := 0
	packname := fmt.Sprintf("a%v_by_%v", strconv.Itoa(user.Id), b.UserName)
	if msg.ReplyToMessage == nil {
		msg.ReplyText("What are you trying to kang?")
		var err error
		return err
	}
	if msg.ReplyToMessage.Sticker == nil {
		msg.ReplyText("Can't kang that mate.")
		var err error
		return err
	}
	if msg.ReplyToMessage.Sticker.IsAnimated == true {
		packname = fmt.Sprintf("b%v_by_%v", strconv.Itoa(user.Id), b.UserName)
	}
	packnameFound := 0
	maxStickers := 120
	for packnameFound == 0 {
		if msg.ReplyToMessage.Sticker.IsAnimated == true {
			stickerset, err := b.GetStickerSet(packname)

			if err != nil {
				packnameFound = 1
				break
			}

			if len(stickerset.Stickers) >= maxStickers {
				packnum++
				packname = fmt.Sprintf("b%v_%v_by_%v", strconv.Itoa(packnum), strconv.Itoa(user.Id), b.UserName)
			} else {
				packnameFound = 1
			}
		} else {
			stickerset, err := b.GetStickerSet(packname)

			if err != nil {
				packnameFound = 1
				break
			}

			if len(stickerset.Stickers) >= maxStickers {
				packnum++
				packname = fmt.Sprintf("a%v_%v_by_%v", strconv.Itoa(packnum), strconv.Itoa(user.Id), b.UserName)
			} else {
				packnameFound = 1
			}
		}
	}
	if msg.ReplyToMessage != nil {
		var fileID string
		var stickerEmoji string
		var success bool
		var err error
		animTitle := "nil"
		if msg.ReplyToMessage.Sticker != nil {
			fileID = msg.ReplyToMessage.Sticker.FileId
		} else {
			msg.ReplyText("Please reply to a sticker for me to kang.")
		}

		if msg.ReplyToMessage.Sticker != nil && msg.ReplyToMessage.Sticker.Emoji != "nil" {
			stickerEmoji = msg.ReplyToMessage.Sticker.Emoji
		} else {
			stickerEmoji = "🛡"
		}

		if msg.ReplyToMessage.Sticker.IsAnimated == true {
			inputFile, r, err := getInputFile(b, fileID, "sticker.tgs")
			if r != nil {
				defer r.Body.Close()
			}
			if err != nil {
				return err
			}
			success, err = b.AddTgsStickerToSet(user.Id, packname, inputFile, stickerEmoji)
			animTitle = "%v's animated pack %v"
		} else {
			inputFile, r, err := getInputFile(b, fileID, "sticker.png")
			if r != nil {
				defer r.Body.Close()
			}
			if err != nil {
				return err
			}
			success, err = b.AddPngStickerToSet(user.Id, packname, inputFile, stickerEmoji)
		}

		if err != nil {
			err := makeInternal(msg, user, fileID, stickerEmoji, b, packname, packnum, animTitle)
			if err != nil {
				msg.ReplyText("Something went wrong with kanging.")
				return err
			}
		}

		if success {
			msg.ReplyMarkdownf("Sticker successfully added to [pack](t.me/addstickers/%v)\nEmoji is: %v", packname, stickerEmoji)
		}
	} else {
		msg.ReplyText("What even fam.")
	}
	return nil
}

func makeInternal(msg *ext.Message, user *ext.User, fileID string, emoji string, bot ext.Bot, packname string, packnum int, animTitle string) error {
	name := user.FirstName
	extraVersion := ""
	title := "%v's pack %v"
	if packnum > 0 {
		extraVersion = " " + strconv.Itoa(packnum)
	}

	if animTitle != "nil" {
		title = animTitle
	}
	newStick := bot.NewSendableCreateNewStickerSet(user.Id, packname, fmt.Sprintf(title, name, extraVersion), emoji)
	if animTitle != "nil" {
		inputFile, r, err := getInputFile(bot, fileID, "sticker.tgs")
		if r != nil {
			defer r.Body.Close()
		}
		if err != nil {
			return err
		}
		newStick.TgsSticker = &inputFile
	} else {
		inputFile, r, err := getInputFile(bot, fileID, "sticker.png")
		if r != nil {
			defer r.Body.Close()
		}
		if err != nil {
			return err
		}
		newStick.PngSticker = &inputFile
	}

	success, err := newStick.Send()

	if err != nil {
		bot.Logger.Warnw("No sticker file.", zap.Error(err))
		return err
	}

	if success == true {
		msg.ReplyHTMLf(fmt.Sprintf("Successfully created pack with name: %v. Get it <a href=\"t.me/addstickers/%v\">here</a>", packname, packname))
	}

	return nil
}

func getInputFile(bot ext.Bot, fileID string, fileName string) (ext.InputFile, *http.Response, error) {
	file, err := bot.GetFile(fileID)
	var inputFile ext.InputFile
	var r *http.Response
	if err != nil {
		print("Cannot get the file!")
		return inputFile, r, err
	}

	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/file/bot%v/%v", noodlex.BotConfig.ApiKey, file.FilePath))

	if err != nil {
		return inputFile, r, err
	}

	inputFile = bot.NewFileReader(fileName, io.Reader(resp.Body))
	return inputFile, resp, nil
}

// LoadStickers - Add commands from module to the bot
func LoadStickers(u *gotgbot.Updater) {
	defer log.Println("Loading module stickers")
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("stickerid", noodlex.BotConfig.Prefix, stickerID))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("getsticker", noodlex.BotConfig.Prefix, getSticker))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("kang", noodlex.BotConfig.Prefix, kangSticker))
}
