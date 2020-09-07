package deleting

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/chat_status"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func purge(bot ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	user := u.EffectiveUser
	chat := u.EffectiveChat

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}

	if msg.ReplyToMessage != nil {
		if chat_status.CanDelete(chat, bot.Id) {
			msgId := msg.ReplyToMessage.MessageId
			deleteTo := msg.MessageId - 1

			if args != nil && len(args) > 0 {
				if _, err := strconv.Atoi(args[0]); err == nil {
					convInt, _ := strconv.Atoi(args[0])
					newDel := msgId + (convInt - 1)
					if newDel < deleteTo {
						deleteTo = newDel
					}
				}
			}

			for mId := deleteTo; mId > msgId-1; mId-- {
				_, err := bot.DeleteMessage(chat.Id, mId)
				if err != nil {
					if err.Error() == "Bad Request: message can't be deleted" {
						_, err := bot.SendMessage(chat.Id, "Cannot delete all messages. The messages may be too old, I might "+
							"not have delete rights, or this might not be a supergroup.")
						error_handling.HandleErr(err)
					} else if err.Error() != "Bad Request: message to delete not found" {
						error_handling.HandleErr(err)
					}
				}
			}
			_, err := msg.Delete()
			if err != nil {
				if err.Error() == "Bad Request: message can't be deleted" {
					_, err := bot.SendMessage(chat.Id, "Cannot delete all messages. The messages may be too old, I might "+
						"not have delete rights, or this might not be a supergroup.")
					error_handling.HandleErr(err)
				} else if err.Error() == "Bad Request: message to delete not found" {
					error_handling.HandleErr(err)
				}
			}
			delMsg, err := bot.SendMessage(chat.Id, fmt.Sprintf("Purged %v messages.", strconv.Itoa((deleteTo-msgId)+1)))
			error_handling.HandleErr(err)
			time.Sleep(2 * time.Second)
			_, err = bot.DeleteMessage(chat.Id, delMsg.MessageId)
			return err
		}
		return nil
	} else {
		_, err := msg.ReplyText("Reply to a message to select where to start purging from.")
		return err
	}
}

func delMessage(bot ext.Bot, u *gotgbot.Update) error {
	msg := u.EffectiveMessage
	user := u.EffectiveUser
	chat := u.EffectiveChat

	// Permission checks
	if !chat_status.RequireBotAdmin(chat, msg) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, msg, user.Id) {
		return gotgbot.EndGroups{}
	}

	if msg.ReplyToMessage != nil {
		if chat_status.CanDelete(chat, bot.Id) {
			_, err := msg.ReplyToMessage.Delete()
			error_handling.HandleErr(err)
			_, err = msg.Delete()
			return err
		}
	} else {
		_, err := msg.ReplyText("Whaddya want to delete?")
		return err
	}
	return nil
}

func LoadDelete(u *gotgbot.Updater) {
	defer log.Println("Loaded module: message_deleting")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("purge", []rune{'/', '!', '?'}, purge))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("del", noodlex.BotConfig.Prefix, delMessage))
}
