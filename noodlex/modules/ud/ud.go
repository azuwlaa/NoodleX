
package ud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NoodleSoup/NoodleX/noodlex"
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
	defer log.Println("Loaded module: ud")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ud", noodlex.BotConfig.Prefix, ud))
}
