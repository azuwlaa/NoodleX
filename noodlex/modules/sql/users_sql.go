
package sql

import (
	"encoding/json"
	"strings"

	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/caching"
	"github.com/PaulSonOfLars/gotgbot"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId   int    `gorm:"primary_key" json:"user_id"`
	UserName string `json:"user_name"`
}

type Chat struct {
	ChatId   string `gorm:"primary_key" json:"chat_id"`
	ChatName string `json:"chat_name"`
}

type ChatMember struct {
	PrivChatId int    `gorm:"primary_key;AUTO_INCREMENT" json: "priv_chat_id"`
	ChatId     string `json: "chat_id"`
	UserId     int    `json: "user_id"`
}

func EnsureBotInDb(u *gotgbot.Updater) {
	// Insert bot user only if it doesn't exist already
	botUser := &User{UserId: u.Dispatcher.Bot.Id, UserName: u.Dispatcher.Bot.UserName}
	SESSION.Save(botUser)
	cacheUser()
}

func UpdateUser(userId int, username string, chatId string, chatName string) {
	username = strings.ToLower(username)
	tx := SESSION.Begin()

	// upsert user
	user := &User{UserId: userId, UserName: username}
	tx.Save(user)

	if chatId == "nil" || chatName == "nil" {
		tx.Commit()
		return
	}

	// upsert chat
	chat := &Chat{ChatId: chatId, ChatName: chatName}
	tx.Save(chat)
	tx.Commit()
	cacheUser()
}

func GetUserIdByName(username string) *User {
	username = strings.ToLower(username)

	userJson, err := caching.CACHE.Get("users")
	var users []User
	if err != nil {
		users = cacheUser()
	}

	_ = json.Unmarshal(userJson, &users)

	for _, user := range users {
		if user.UserName == username {
			return &user
		}
	}

	return nil
}

func UpdateChatMember(userId int, chatId string) {
	tx := SESSION.Begin()

	user := &ChatMember{UserId: userId, ChatId: chatId}
	if err := tx.Where("chat_id = ? AND user_id = ?", chatId, userId).First(&user).Error; !gorm.IsRecordNotFoundError(err) {
		return
	}

	// upsert user
	chatmember := &ChatMember{UserId: userId, ChatId: chatId}
	tx.Save(chatmember)

	if chatId == "nil" {
		tx.Commit()
		return
	}

	tx.Commit()
}

func GetUserChats(userId int) int {
	var user []ChatMember
	SESSION.Model(&ChatMember{}).Where("user_id = ?", userId).Find(&user)
	return len(user)
}

func cacheUser() []User {
	var users []User
	SESSION.Model(&User{}).Find(&users)
	// userJson, _ := jettison.Marshal(users)
	// err := caching.CACHE.Set("users", userJson)
	// error_handling.HandleErr(err)
	return users
}
