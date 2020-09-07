package bans

import (
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/chat_status"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/error_handling"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/extraction"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/helpers"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/string_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
)

func ban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userID, reason := extraction.ExtractUserAndText(message, args)
	if userID == 0 {
		_, err := message.ReplyText("You don't seem to be referring to a user.")
		return err
	}

	member, err := chat.GetMember(userID)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("I can't seem to find this user.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userID, member) {
		_, err := message.ReplyText("Why would I ban an admin? That sounds like a pretty dumb idea.")
		return err
	}

	if userID == bot.Id {
		_, err := message.ReplyText("I'm not gonna BAN myself, are you crazy?")
		return err
	}

	_, err = chat.KickMember(userID)
	if err != nil {
		return err
	}

	bannedUser, _ := bot.GetChat(userID)

	text := fmt.Sprintf("%v has been banned!", helpers.MentionHtml(userID, fmt.Sprintf("%v", bannedUser.FirstName)))
	text += fmt.Sprintf("\n<b>ID:</b> <code>%v</code>", userID)

	if reason != "" {
		text += fmt.Sprintf("\n<b>Reason:</b> %v", html.EscapeString(reason))
	}

	_, err = message.ReplyHTML(text)
	return err
}

func tempBan(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userID, reason := extraction.ExtractUserAndText(message, args)
	if userID == 0 {
		_, err := message.ReplyText("You don't seem to be referring to a user.")
		return err
	}

	member, err := chat.GetMember(userID)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("I can't seem to find this user.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userID, member) {
		_, err := message.ReplyText("Why would I ban an admin? That sounds like a pretty dumb idea.")
		return err
	}

	if userID == bot.Id {
		_, err := message.ReplyText("I'm not gonna BAN myself, are you crazy?")
		return err
	}

	if reason == "" {
		_, err := message.ReplyText("You haven't specified a time to ban this user for!")
		return err
	}

	splitReason := strings.SplitN(reason, " ", 2)
	timeVal := splitReason[0]
	banTime := string_handling.ExtractTime(message, timeVal)
	if banTime == -1 {
		return nil
	}
	newMsg := bot.NewSendableKickChatMember(chat.Id, userID)
	string_handling.ExtractTime(message, timeVal)
	newMsg.UntilDate = banTime
	_, err = newMsg.Send()
	if err != nil {
		_, err := message.ReplyText("Well damn, I can't ban that user.")
		error_handling.HandleErr(err)
	}

	bannedUser, _ := bot.GetChat(userID)

	text := fmt.Sprintf("%v has been temporarily banned for %s!",
		helpers.MentionHtml(userID, fmt.Sprintf("%v", bannedUser.FirstName)), timeVal)
	text += fmt.Sprintf("\n<b>ID:</b> <code>%v</code>", userID)

	_, err = message.ReplyHTML(text)
	return err
}

func kick(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userID, reason := extraction.ExtractUserAndText(message, args)
	if userID == 0 {
		_, err := message.ReplyText("You don't seem to be referring to a user.")
		return err
	}

	var member, err = chat.GetMember(userID)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("I can't seem to find this user.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userID, member) {
		_, err := message.ReplyText("One day I'll find out how to work around the bot API. Today is not that day.")
		return err
	}

	if userID == bot.Id {
		_, err := message.ReplyText("Yeahhh I'm not gonna do that.")
		return err
	}

	_, err = chat.UnbanMember(userID) // Apparently unban on current user = kick
	if err != nil {
		_, err = message.ReplyText("Hec, I can't seem to kick this user.")
		return err
	}

	kickedUser, _ := bot.GetChat(userID)

	text := fmt.Sprintf("%v has been kicked!", helpers.MentionHtml(userID, fmt.Sprintf("%v", kickedUser.FirstName)))
	text += fmt.Sprintf("\n<b>ID:</b> <code>%v</code>", userID)

	if reason != "" {
		text += fmt.Sprintf("\n<b>Reason:</b> %v", html.EscapeString(reason))
	}

	_, err = message.ReplyHTML(text)
	return err
}

func kickme(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		_, err := message.ReplyText("Yeahhh I'm not gonna do that.")
		error_handling.HandleErr(err)
		return gotgbot.EndGroups{}
	}

	Kickme, _ := chat.UnbanMember(user.Id) // kick the user
	if Kickme {
		_, err := message.ReplyText("Get out!")
		return err
	} else {
		_, err := message.ReplyText("Huh? I can't :/")
		return err
	}
}

func banme(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}

	if chat_status.IsUserAdmin(chat, user.Id) {
		_, err := message.ReplyText("Yeahhh I'm not gonna do that.")
		error_handling.HandleErr(err)
		return gotgbot.EndGroups{}
	}

	Banme, _ := chat.KickMember(user.Id)
	if Banme {
		_, err := message.ReplyText("Yeah right, get lost.")
		return err
	} else {
		_, err := message.ReplyText("Huh? I can't :/")
		return err
	}
}

func unban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) && chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userID, reason := extraction.ExtractUserAndText(message, args)

	if userID == 0 {
		_, err := message.ReplyText("You don't seem to be referring to a user.")
		return err
	}

	_, err := chat.GetMember(userID)
	if err != nil {
		_, err := message.ReplyText("I can't seem to find this user.")
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to unban users!")
		return err
	}

	if userID == bot.Id {
		_, err := message.ReplyText("How would I unban myself if I wasn't here...?")
		return err
	}

	if chat_status.IsUserInChat(chat, userID) {
		_, err := message.ReplyText("Why are you trying to unban someone that's already in the chat?")
		return err
	}

	_, err = chat.UnbanMember(userID)
	unbanUser, _ := bot.GetChat(userID)
	error_handling.HandleErr(err)

	text := fmt.Sprintf("Yep, %v can join again!", helpers.MentionHtml(userID, fmt.Sprintf("%v", unbanUser.FirstName)))
	text += fmt.Sprintf("\n<b>ID:</b> <code>%v</code>", userID)

	if reason != "" {
		text += fmt.Sprintf("\n<b>Reason:</b> %v", html.EscapeString(reason))
	}

	_, err = message.ReplyHTML(text)
	return err
}

func sban(bot ext.Bot, u *gotgbot.Update, args []string) error {
	chat := u.EffectiveChat
	user := u.EffectiveUser
	message := u.EffectiveMessage

	message.Delete()

	// Permission checks
	if u.EffectiveChat.Type == "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in a group!")
		return err
	}

	if !chat_status.RequireBotAdmin(chat, message) {
		return gotgbot.EndGroups{}
	}
	if !chat_status.RequireUserAdmin(chat, message, user.Id) {
		return gotgbot.EndGroups{}
	}

	userID, _ := extraction.ExtractUserAndText(message, args)
	if userID == 0 {
		_, err := message.ReplyText("You don't seem to be referring to a user.")
		return err
	}

	member, err := chat.GetMember(userID)
	if err != nil {
		if err.Error() == "User not found" {
			_, err = message.ReplyText("I can't seem to find this user.")
		}
		return err
	}

	userMember, _ := chat.GetMember(user.Id)
	if !userMember.CanRestrictMembers && userMember.Status != "creator" {
		_, err = message.ReplyText("You don't have permissions to ban users!")
		return err
	}

	if chat_status.IsUserBanProtected(chat, userID, member) {
		_, err := message.ReplyText("Why would I ban an admin? That sounds like a pretty dumb idea.")
		return err
	}

	if userID == bot.Id {
		_, err := message.ReplyText("I'm not gonna BAN myself, are you crazy?")
		return err
	}

	_, err = chat.KickMember(userID)
	if err != nil {
		return err
	}

	_, err = message.Delete()
	return err
}

func LoadBans(u *gotgbot.Updater) {
	defer log.Println("Loaded module: bans")
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("tban", noodlex.BotConfig.Prefix, tempBan))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("ban", noodlex.BotConfig.Prefix, ban))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("sban", noodlex.BotConfig.Prefix, sban))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("kick", noodlex.BotConfig.Prefix, kick))
	u.Dispatcher.AddHandler(handlers.NewPrefixArgsCommand("unban", noodlex.BotConfig.Prefix, unban))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("kickme", noodlex.BotConfig.Prefix, kickme))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("banme", noodlex.BotConfig.Prefix, banme))
}
