package misc

import (
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/sql"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/error_handling"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/extraction"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/helpers"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/PaulSonOfLars/gotgbot/ext"
	"github.com/PaulSonOfLars/gotgbot/handlers"
	"github.com/sirupsen/logrus"
	"github.com/tcnksm/go-httpstat"
)

var runStrings = [59]string{"Where do you think you're going?",
	"Huh? what? did they get away?",
	"ZZzzZZzz... Huh? what? oh, just them again, nevermind.",
	"Get back here!",
	"Not so fast...",
	"Look out for the wall!",
	"Don't leave me alone with them!!",
	"You run, you die.",
	"Energy drinks makes you run faster!",
	"Stop walking and start to run",
	"Jokes on you, I'm everywhere",
	"You're gonna regret that...",
	"You could also try /kickme, I hear that's fun.",
	"Go bother someone else, no-one here cares.",
	"You can run, but you can't hide.",
	"Is that all you've got?",
	"I'm behind you...",
	"You've got company!",
	"We can do this the easy way, or the hard way.",
	"You just don't get it, do you?",
	"Yeah, you better run!",
	"Please, remind me how much I care?",
	"I'd run faster if I were you.",
	"That's definitely the droid we're looking for.",
	"May the odds be ever in your favour.",
	"Famous last words.",
	"If you disappear, don't call for help...",
	"Run for your life!",
	"And they disappeared forever, never to be seen again.",
	"\"Oh, look at me! I'm so cool, I can run from a bot!\" - this person",
	"Yeah yeah, just tap /kickme already.",
	"Here, take this ring and head to Mordor while you're at it.",
	"Legend has it, they're still running...",
	"Unlike Harry Potter, your parents can't protect you from me.",
	"Fear leads to anger. Anger leads to hate. Hate leads to suffering. If you keep running in fear, you might " +
		"be the next Vader.",
	"Multiple calculations later, I have decided my interest in your shenanigans is exactly 0.",
	"Legend has it, they're still running.",
	"Keep it up, not sure we want you here anyway.",
	"You're a wiza- Oh. Wait. You're not Harry, keep moving.",
	"NO RUNNING IN THE HALLWAYS!",
	"Hasta la vista, baby.",
	"Run carelessly you might get tripped.",
	"You have done a wonderful job, Keep it up...",
	"I see an evil spirits here, Let's expel them!\n\n" +
		"Exorcizamus te, omnis immunde spiritus, omni satanica potestas, omnis incursio infernalis adversarii," +
		" omnis legio, omnis congregatio et secta diabolica, in nomini et virtute Domini nostri Jesu Christi, eradicare " +
		"et effugare a Dei Ecclesia, ab animabus ad imaginem Dei conditis ac pretioso divini Agni sanguini redemptis.",
	"Who let the dogs out?",
	"It's funny, because no one cares.",
	"That's cool, just hit on seppuku /banme already.",
	"Ah, what a waste. I liked that one.",
	"Frankly, my dear, I don't give a damn.",
	"My flowers brings all the girls to yard... So run faster!",
	"You can't HANDLE the truth!",
	"A long time ago, in a galaxy far far away... Someone would've cared about that. Not anymore though.",
	"Hey, look at them! They're running from the inevitable banhammer... Cute.",
	"Han shot first. So will I.",
	"What are you running after, a white rabbit?",
	"As The Doctor would say... RUN!"}

var yornStrings = [9]string{"Yes.",
	"No.",
	"Get Rekt.",
	"Absolutely.",
	"Nah.",
	"No one cares.",
	"I dont care!",
	"Go for it!",
	"In your dreams."}

func getID(bot ext.Bot, u *gotgbot.Update, args []string) error {
	userId := extraction.ExtractUser(u.EffectiveMessage, args)
	if userId != 0 {
		if u.EffectiveMessage.ReplyToMessage != nil && u.EffectiveMessage.ReplyToMessage.ForwardFrom != nil {
			user1 := u.EffectiveMessage.ReplyToMessage.From
			user2 := u.EffectiveMessage.ReplyToMessage.ForwardFrom
			_, err := u.EffectiveMessage.ReplyHTMLf("The original sender, %v, has an ID of <code>%v</code>.\n"+
				"The forwarder, %v, has an ID of <code>%v</code>.", html.EscapeString(user2.FirstName),
				user2.Id,
				html.EscapeString(user1.FirstName),
				user1.Id)
			return err
		} else {
			user, err := bot.GetChat(userId)
			error_handling.HandleErr(err)
			_, err = u.EffectiveMessage.ReplyHTMLf("%v's ID is <code>%v</code>", html.EscapeString(user.FirstName), user.Id)
		}
	} else {
		chat := u.EffectiveChat
		if chat.Type == "private" {
			_, err := u.EffectiveMessage.ReplyHTMLf("Your ID is <code>%v</code>", chat.Id)
			return err
		} else {
			_, err := u.EffectiveMessage.ReplyHTMLf("This group's ID is <code>%v</code>", chat.Id)
			return err
		}
	}
	return nil
}

func info(bot ext.Bot, u *gotgbot.Update, args []string) error {
	msg := u.EffectiveMessage
	chat := u.EffectiveChat
	userId := extraction.ExtractUser(msg, args)
	var user *ext.User

	if userId != 0 {
		userChat, _ := bot.GetChat(userId)
		user = &ext.User{
			Id:        userChat.Id,
			FirstName: userChat.FirstName,
			LastName:  userChat.LastName,
		}

	} else if msg.ReplyToMessage == nil && len(args) <= 0 {
		user = msg.From
		userId = msg.From.Id

	} else if _, err := strconv.Atoi(args[0]); msg.ReplyToMessage == nil && (len(args) <= 0 || (len(args) >= 1 &&
		strings.HasPrefix(args[0], "@") && err != nil && msg.ParseEntities()[0].Type != "TEXT_MENTION")) {
		_, err := msg.ReplyText("You don't seem to be referring to a user.")
		return err
	} else {
		return nil
	}

	text := fmt.Sprintf("<b>User info:</b>"+
		"\nID: <code>%v</code>"+
		"\nFirst Name: %v", userId, html.EscapeString(user.FirstName))

	if user.LastName != "" {
		text += fmt.Sprintf("\nLast Name: %v", user.LastName)
	}

	if user.Username != "" {
		text += fmt.Sprintf("\nUsername: @%v", user.Username)
	}

	text += fmt.Sprintf("\nPermanent user link: %v", helpers.MentionHtml(user.Id, fmt.Sprintf("link")))

	fed := sql.GetChatFed(strconv.Itoa(chat.Id))
	if fed != nil {
		fban := sql.GetFbanUser(fed.Id, strconv.Itoa(userId))
		if fban != nil {
			text += fmt.Sprintf("\n\nThis user is fedbanned in the current federation - "+
				"<code>%v</code>", fed.FedName)
		} else {
			text += "\n\nThis user is not fedbanned in the current federation."
		}
	}

	if user.Id == noodlex.BotConfig.OwnerId {
		text += "\n\n<i>This person is my creator! They have total power over me.</i>"
	} else {
		for _, id := range noodlex.BotConfig.SudoUsers {
			if strconv.Itoa(user.Id) == id {
				text += "\n\n<i>This person is one of my sudo users! Nearly as powerful as my owner - so watch it.</i>"
			}
		}
	}

	text += fmt.Sprintf("\n\nI've seen them in <code>%v</code> chats in total.", sql.GetUserChats(userId))

	_, err := u.EffectiveMessage.ReplyHTML(text)
	return err
}

func ping(_ ext.Bot, u *gotgbot.Update) error {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	error_handling.HandleErr(err)

	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	error_handling.HandleErr(err)

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		logrus.Println(err)
	}

	_ = res.Body.Close()

	text := fmt.Sprintf("Ping: <b>%d</b> ms", result.ServerProcessing/time.Millisecond)

	_, err = u.EffectiveMessage.ReplyHTML(text)
	return err
}

func runs(bot ext.Bot, u *gotgbot.Update) error {
	rand.Seed(time.Now().Unix())
	u.EffectiveMessage.ReplyText(runStrings[rand.Intn(len(runStrings))])
	return nil
}

func yorn(bot ext.Bot, u *gotgbot.Update) error {
	rand.Seed(time.Now().Unix())
	msg := u.EffectiveMessage
	msg.ReplyText(yornStrings[rand.Intn(len(yornStrings))])
	return nil
}

func LoadMisc(u *gotgbot.Updater) {
	defer log.Println("Loaded module: misc")
	CHandler := u.Dispatcher.AddHandler
	CHandler(handlers.NewPrefixArgsCommand("id", noodlex.BotConfig.Prefix, getID))
	CHandler(handlers.NewPrefixArgsCommand("info", noodlex.BotConfig.Prefix, info))
	CHandler(handlers.NewPrefixCommand("ping", noodlex.BotConfig.Prefix, ping))
	CHandler(handlers.NewPrefixCommand("runs", noodlex.BotConfig.Prefix, runs))
	CHandler(handlers.NewRegex("(?i)y/n", yorn))
}
