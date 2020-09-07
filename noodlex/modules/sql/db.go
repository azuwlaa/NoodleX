
package sql

import (
	"log"

	"github.com/NoodleSoup/NoodleX/noodlex"
	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/error_handling"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

var SESSION *gorm.DB

func init() {
	conn, err := pq.ParseURL(noodlex.BotConfig.SqlUri)
	error_handling.FatalError(err)

	db, err := gorm.Open("postgres", conn)
	error_handling.FatalError(err)

	if noodlex.BotConfig.DebugMode == "True" {
		SESSION = db.Debug()
		log.Println("[INFO][Database] Using database in debug mode.")
	} else {
		SESSION = db
	}

	db.DB().SetMaxOpenConns(100)

	log.Println("[INFO][Database] Database connected")

	// Create tables if they don't exist
	SESSION.AutoMigrate(&User{}, &Chat{}, &ChatMember{}, &Warns{}, &WarnFilters{}, &WarnSettings{}, &BlackListFilters{}, &Federation{},
		&FedChat{}, &FedAdmin{}, &FedBan{}, &Note{}, &Button{}, &Welcome{}, &WelcomeButton{}, &MutedUser{}, &Rules{})
	log.Println("Auto-migrated database schema")

}
