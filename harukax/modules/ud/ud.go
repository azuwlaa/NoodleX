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

package ud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HarukaNetwork/HarukaX/harukax"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

type Response struct {
	List []Info
}

type Info struct {
	Definition string
	Example    string
}

func ud(_ ext.Bot, u *gotgbot.Update, args []string) error {
	oText := strings.Join(args, " ")
	qText := oText
	if len(args) > 1 {
		qText = strings.Replace(oText, " ", "+", -1)
	}
	resp, err := http.Get(fmt.Sprintf("http://api.urbandictionary.com/v0/define?term=%v", qText))

	defer resp.Body.Close()

	var data Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	text := fmt.Sprintf("Word: %v\n\nDefinition: %v\n\n<i>%v</i>", oText, data.List[0].Definition, data.List[0].Example)

	_, err = u.EffectiveMessage.ReplyHTML(text)
	return err
}

func LoadUd(u *gotgbot.Updater) {
	defer log.Println("Loading module ud")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ud", harukax.BotConfig.Prefix, ud))
}
