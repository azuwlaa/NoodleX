package help

import (
	"fmt"
	"html"
	"log"
	"regexp"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/error_handling"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/PaulSonOfLars/gotgbot/parsemode"
)

var markup ext.InlineKeyboardMarkup
var markdownHelpText string

func initMarkdownHelp() {
	markdownHelpText = "You can use markdown to make your messages more expressive. This is the markdown currently " +
		"supported:\n\n" +
		"<code>`code words`</code>: backticks allow you to wrap your words in monospace fonts.\n" +
		"<code>*bold*</code>: wrapping text with '*' will produce bold text\n" +
		"<code>_italics_</code>: wrapping text with '_' will produce italic text\n" +
		"<code>[hyperlink](example.com)</code>: this will create a link - the message will just show " +
		"<code>hyperlink</code>, and tapping on it will open the page at <code>example.com</code>\n\n" +
		"<code>[buttontext](buttonurl:example.com)</code>: this is a special enhancement to allow users to have " +
		"telegram buttons in their markdown. <code>buttontext</code> will be what is displayed on the button, and " +
		"<code>example.com</code> will be the url which is opened.\n\n" +
		"If you want multiple buttons on the same line, use :same, as such:\n" +
		"<code>[one](buttonurl://github.com)</code>\n" +
		"<code>[two](buttonurl://google.com:same)</code>\n" +
		"This will create two buttons on a single line, instead of one button per line.\n\n" +
		"Keep in mind that your message MUST contain some text other than just a button!"

}

func initHelpButtons() {
	helpButtons := [][]ext.InlineKeyboardButton{make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2),
		make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 2), make([]ext.InlineKeyboardButton, 1)}

	// First column
	helpButtons[0][0] = ext.InlineKeyboardButton{
		Text:         "Admin",
		CallbackData: fmt.Sprintf("help(%v)", "admin"),
	}
	helpButtons[1][0] = ext.InlineKeyboardButton{
		Text:         "Bans",
		CallbackData: fmt.Sprintf("help(%v)", "bans"),
	}
	helpButtons[2][0] = ext.InlineKeyboardButton{
		Text:         "Blacklists",
		CallbackData: fmt.Sprintf("help(%v)", "blacklist"),
	}
	helpButtons[3][0] = ext.InlineKeyboardButton{
		Text:         "Purge/Delete",
		CallbackData: fmt.Sprintf("help(%v)", "deleting"),
	}
	helpButtons[4][0] = ext.InlineKeyboardButton{
		Text:         "Alliance",
		CallbackData: fmt.Sprintf("help(%v)", "feds"),
	}
	helpButtons[5][0] = ext.InlineKeyboardButton{
		Text:         "Stickers",
		CallbackData: fmt.Sprintf("help(%v)", "stickers"),
	}

	// Second column
	helpButtons[0][1] = ext.InlineKeyboardButton{
		Text:         "Misc",
		CallbackData: fmt.Sprintf("help(%v)", "misc"),
	}
	helpButtons[1][1] = ext.InlineKeyboardButton{
		Text:         "Mutes",
		CallbackData: fmt.Sprintf("help(%v)", "muting"),
	}
	helpButtons[2][1] = ext.InlineKeyboardButton{
		Text:         "Notes",
		CallbackData: fmt.Sprintf("help(%v)", "notes"),
	}
	helpButtons[3][1] = ext.InlineKeyboardButton{
		Text:         "Users",
		CallbackData: fmt.Sprintf("help(%v)", "users"),
	}
	helpButtons[4][1] = ext.InlineKeyboardButton{
		Text:         "Warns",
		CallbackData: fmt.Sprintf("help(%v)", "warns"),
	}

	markup = ext.InlineKeyboardMarkup{InlineKeyboard: &helpButtons}
}

func help(b ext.Bot, u *gotgbot.Update) error {
	msg := b.NewSendableMessage(u.EffectiveChat.Id, fmt.Sprintf("Hey there! I'm <b>%v</b>, A fully feature packed group management bot. "+
		"My creator designed me to handle groups with features such as notes, filters and even a warn system"+
		"\n\n<b>Here are the list of helpful commands:</b>\n"+
		"- /start: starts the bot\n"+
		"- /help: sends you list of commands\n"+
		"\n\nIf you have any bugs or questions on how to use me, head to @NatalieSupport."+
		"\nAll commands can either be used with (/) slash, (.) dot or (!) exclamation", b.FirstName))
	msg.ParseMode = parsemode.Html
	msg.ReplyToMessageId = u.EffectiveMessage.MessageId
	msg.ReplyMarkup = &markup
	_, err := msg.Send()
	if err != nil {
		msg.ReplyToMessageId = 0
		_, err = msg.Send()
	}
	return err
}

func markdownHelp(_ ext.Bot, u *gotgbot.Update) error {
	chat := u.EffectiveChat
	if chat.Type != "private" {
		_, err := u.EffectiveMessage.ReplyText("This command is meant to be used in PM!")
		return err
	}

	_, err := u.EffectiveMessage.ReplyHTML(markdownHelpText)
	return err
}

func buttonHandler(b ext.Bot, u *gotgbot.Update) error {
	query := u.CallbackQuery
	pattern, _ := regexp.Compile(`help\((.+?)\)`)

	if pattern.MatchString(query.Data) {
		module := pattern.FindStringSubmatch(query.Data)[1]
		chat := u.EffectiveChat
		msg := b.NewSendableEditMessageText(chat.Id, u.EffectiveMessage.MessageId, "Module information unavailable")
		msg.ParseMode = parsemode.Html
		backButton := [][]ext.InlineKeyboardButton{{ext.InlineKeyboardButton{
			Text:         "Back",
			CallbackData: "help(back)",
		}}}
		backKeyboard := ext.InlineKeyboardMarkup{InlineKeyboard: &backButton}
		msg.ReplyMarkup = &backKeyboard

		switch module {
		case "admin":
			msg.Text = "Here is the help for the <b>Admin</b> module:\n\n" +
				"- /adminlist: list of admins in the chat\n\n" +
				"<b>Admin only:</b>" +
				html.EscapeString("- /pin: silently pins the message replied to - add 'loud' or 'notify' to give notifs to users.\n"+
					"- /unpin: unpins the currently pinned message\n"+
					"- /invitelink: gets invitelink\n"+
					"- /promote: promotes the user replied to\n"+
					"- /demote: demotes the user replied to\n")
			break
		case "bans":
			msg.Text = "Some people need to be publicly banned, mute or kicked; spammers, annoyances, or just trolls" +
				"This module allows you to do that easily, by exposing some common actions, so everyone will see!\n\n" +
				" - /kickme: kicks the user who issued the comman\n" +
				" - /banme: bans the user who issued the command\n\n" +
				"<b>Admin only</b>:\n" +
				html.EscapeString(
					" - /ban <userhandle>: bans a user. (via handle, or reply)\n"+
						" - /tban <userhandle> x(m/h/d): bans a user for x time. (via handle, or reply). m = minutes, h = hours,"+
						" d = days.\n"+
						"- /unban <userhandle>: unbans a user. (via handle, or reply)"+
						" - /kick <userhandle>: kicks a user, (via handle, or reply)"+
						" - /sban <userhandle>: bans a user, quietly, deletes your message (via handle, or reply)")

			break
		case "blacklist":
			msg.Text = "Here is the help for the <b>Word Blacklists</b> module:\n\n" +
				"Blacklists are used to stop certain triggers from being said in a group. Any time the trigger is " +
				"mentioned, the message will immediately be deleted. A good combo is sometimes to pair this up with " +
				"warn filters!\n\n" +
				"<b>NOTE:</b> blacklists do not affect group admins.\n\n" +
				" - /blacklist: View the current blacklisted words.\n\n" +
				"<b>Admin only:</b>\n" +
				html.EscapeString("- /addblacklist <triggers>: Add a trigger to the blacklist. Each line is "+
					"considered one trigger, so using different lines will allow you to add multiple triggers.\n"+
					"- /unblacklist <triggers>: Remove triggers from the blacklist. Same newline logic applies here, "+
					"so you can remove multiple triggers at once.\n"+
					" - /rmblacklist <triggers>: Same as above.")
			break
		case "deleting":
			msg.Text = "Here is the help for the <b>Purges</b> module:\n\n" +
				"<b>Admin only:</b>\n" +
				" - /del: deletes the message you replied to\n" +
				" - /purge: deletes all messages between this and the replied to message.\n" +
				html.EscapeString(" - /purge <amount>: number of messages to delete from the replied to message.\n")
			break
		case "feds":
			break
		case "misc":
			msg.Text = "An \"odds and ends\" module for small, simple commands which don't really fit anywhere:\n\n" +
				"- /id: get the current group id. If used by replying to a message, gets that user's id.\n" +
				"- /ping: get server respond time to google.com.\n" +
				"- /runs: reply a random string from an array of replies.\n" +
				"- <strike>/tl: use as reply or translate using text from any language to English.<b>(WIP)</b></strike>\n" +
				"- <strike>/rmkeyboard: removes nasty keyboard buttons from chat.<b>(WIP)</b></strike>\n" +
				"- <strike>/first: scrolls to first message of chat.<b>(WIP)</b></strike>\n" +
				"- <strike>/s: saves the message you reply to your chat with the bot.<b>(WIP)</b></strike>\n" +
				"- <strike>/slap: slap a user, or get slapped if not a reply.<b>(WIP)</b></strike>\n" +
				"- <strike>/hug: hug a user, or get hugged if not a reply.<b>(WIP)</b></strike>\n" +
				"- <strike>/kiss: kiss a user, or get kissed if not a reply.<b>(WIP)</b></strike>\n" +
				"- <strike>/punch: punch a user, or get punched if not a reply.<b>(WIP)</b></strike>\n" +
				"- /info: get information about a user.\n" +
				"- Hi {}: responds to user (to disable greet `/disable botgreet`; to enable greet `/enable botgreet`)\n" +
				"- <strike>/markdownhelp: quick summary of how markdown works in telegram - can only be called in private chats.<b>(WIP)</b></strike>\n" +
				"- <strike>y/n: get randomised answers to a question.<b>(WIP)</b></strike>\n" +
				"- <strike>/tts: convert texts to speech.<b>(WIP)</b></strike>\n" +
				"- <strike>/gifid: get GIF file id.<b>(WIP)</b></strike>"
			break
		case "muting":
			msg.Text = "Here is the help for the <b>Muting</b> module:\n\n" +
				"<b>Admin only:</b>\n" +
				html.EscapeString("- /mute <userhandle>: silences a user. Can also be used as a reply, muting the "+
					"replied to user.\n"+
					"- /tmute <userhandle> x(m/h/d): mutes a user for x time. (via handle, or reply). m = minutes, h = "+
					"hours, d = days.\n"+
					"- /unmute <userhandle>: unmutes a user. Can also be used as a reply, muting the replied to user.")
			break
		case "notes":
			msg.Text = "Here is the help for the <b>Notes</b> module:\n\n" +
				html.EscapeString("- /get <notename>: get the note with this notename\n"+
					"- #<notename>: same as /get\n"+
					"- /notes or /saved: list all saved notes in this chat\n\n"+
					"If you would like to retrieve the contents of a note without any formatting, use /get"+
					" <notename> noformat. This can be useful when updating a current note.\n\n") +
				"<b>Admin only:</b>\n" +
				html.EscapeString(" - /save <notename> <notedata>: saves notedata as a note with name notename\n"+
					"A button can be added to a note by using standard markdown link syntax - the link should just "+
					"be prepended with a buttonurl: section, as such: [somelink](buttonurl:example.com). Check "+
					"/markdownhelp for more info.\n"+
					" - /save <notename>: save the replied-to message as a note with name notename\n"+
					" - /clear <notename>: clear note with this name")
			break
		case "users":
			break
		case "warns":
			msg.Text = "Here is the help for the <b>Warnings</b> module:\n\n" +
				html.EscapeString(" - /warns <userhandle>: get a user's number, and reason, of warnings.\n"+
					" - /warnlist: list of all current warning filters\n\n") +
				"<b>Admin only:</b>\n" +
				html.EscapeString("- /warn <userhandle>: warn a user. After the warn limit, the user will be banned from the group. "+
					"Can also be used as a reply.\n"+
					" - /resetwarn <userhandle>: reset the warnings for a user. Can also be used as a reply.\n"+
					" - /addwarn <keyword> <reply message>: set a warning filter on a certain keyword. If you want your "+
					"keyword to be a sentence, encompass it with quotes, as such: /addwarn \"very angry\" "+
					"This is an angry user.\n"+
					"- /nowarn <keyword>: stop a warning filter\n"+
					"- /warnlimit <num>: set the warning limit\n"+
					" - /strongwarn <on/yes/off/no>: If set to on, exceeding the warn limit will result in a ban. "+
					"Else, will just kick.\n")
			break
		case "back":
			msg.Text = fmt.Sprintf("Hey there! I'm <b>%v</b>, A fully feature packed group management bot. "+
				"My creator designed me to handle groups with features such as notes, filters and even a warn system"+
				"\n\n<b>Here are the list of helpful commands:</b>\n"+
				"- /start: starts the bot\n"+
				"- /help: sends you list of commands\n"+
				"\n\nIf you have any bugs or questions on how to use me, head to @NatalieSupport."+
				"\nAll commands can either be used with (/) slash, (.) dot or (!) exclamation", b.FirstName)
			msg.ReplyMarkup = &markup
			break
		case "stickers":
			msg.Text = "Kanging or fetching ID of stickers are made easy! " +
				"With this stickers command you can easily get the png file " +
				"or tgs file for creating new sticker packs or fetch the ID of stickers.\n\n" +
				"<b>NOTE:</b>\n" +
				"Animated stickers are also supported!\n\n" +
				"- /stickerid: reply to a sticker to me to tell you its file ID.\n" +
				"- /getsticker: reply to a sticker to me to upload its raw PNG file.\n" +
				"- /kang: reply to a sticker to add it to your pack. Animated stickers are also supported!"
			msg.ParseMode = parsemode.Html
		}

		_, err := msg.Send()
		error_handling.HandleErr(err)
		_, err = b.AnswerCallbackQuery(query.Id)
		return err
	}
	return nil
}

func LoadHelp(u *gotgbot.Updater) {
	defer log.Println("Loaded module: help")
	initHelpButtons()
	initMarkdownHelp()
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("help", noodlex.BotConfig.Prefix, help))
	u.Dispatcher.AddHandler(handlers.NewCallback("help", buttonHandler))
	u.Dispatcher.AddHandler(handlers.NewPrefixCommand("markdownhelp", noodlex.BotConfig.Prefix, markdownHelp))
}
