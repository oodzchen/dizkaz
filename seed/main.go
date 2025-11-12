package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/oodzchen/dizkaz/config"
	"github.com/oodzchen/dizkaz/mocktool"
	"github.com/oodzchen/dizkaz/service"
	"github.com/oodzchen/dizkaz/store"
	"github.com/oodzchen/dizkaz/store/pgstore"
)

const timeoutDuration int = 300

var articleNum int
var categoryType string
var userNum int
var envFile string
var showHead bool

var mock *mocktool.Mock

var replyCmd *flag.FlagSet
var replyNum int
var replyNumOfReply int
var replyToId int
var maxLevel int

var usersCmd *flag.FlagSet
var usersFile string

var categoriesCmd *flag.FlagSet
var categoriesFile string

func init() {
	const defaultArticleNum int = 16
	const defaultUserNum int = 30
	const defaultEnvFile string = ".env.testing"

	flag.IntVar(&articleNum, "an", defaultArticleNum, "Create article with specific number")
	flag.StringVar(&categoryType, "c", "general", "Article category: general, internet, computer-sciense, hacker-news, qna, show, dizkaz")
	flag.IntVar(&userNum, "un", defaultUserNum, "User(goroutine) number")
	flag.StringVar(&envFile, "e", defaultEnvFile, "ENV file path, default to .env.testing")
	flag.BoolVar(&showHead, "h", false, "Show browser head")

	replyCmd = flag.NewFlagSet("reply", flag.ExitOnError)
	replyCmd.IntVar(&replyNum, "n", 10, "Number of replies")
	replyCmd.IntVar(&replyToId, "t", 0, "Article id to reply")
	replyCmd.IntVar(&maxLevel, "l", 0, "Max nested reply level")
	replyCmd.IntVar(&replyNumOfReply, "rnr", 1, "reply number of reply")

	usersCmd = flag.NewFlagSet("users", flag.ExitOnError)
	usersCmd.StringVar(&usersFile, "f", "../config/users.yml", "Users YAML file path")
	usersCmd.StringVar(&envFile, "e", defaultEnvFile, "ENV file path, default to .env.testing")

	categoriesCmd = flag.NewFlagSet("categories", flag.ExitOnError)
	categoriesCmd.StringVar(&categoriesFile, "f", "../config/categories.yml", "Categories YAML file path")
	categoriesCmd.StringVar(&envFile, "e", defaultEnvFile, "ENV file path, default to .env.testing")
}

func main() {
	// 检查是否有特定的子命令
	var subcommand string
	var subcommandArgs []string

	// 直接检查是否有已知的子命令
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "reply" || arg == "users" || arg == "categories" {
			subcommand = arg
			if i+1 < len(os.Args) {
				subcommandArgs = os.Args[i+1:]
			}
			break
		}
	}

	// 解析全局flag
	flag.Parse()

	// 如果是子命令，解析子命令的flag
	if subcommand != "" {
		switch subcommand {
		case "reply":
			replyCmd.Parse(subcommandArgs)
		case "users":
			usersCmd.Parse(subcommandArgs)
		case "categories":
			categoriesCmd.Parse(subcommandArgs)
		}
	}

	_, err := os.Stat(envFile)

	if os.IsNotExist(err) {
		err = config.InitFromEnv()
	} else {
		fmt.Printf("ENV file path: %s\n", envFile)
		err = config.Init(envFile)
	}

	if err != nil {
		log.Fatal(err)
	}
	cfg := config.Config

	// fmt.Println("App config: ", cfg)

	mock = mocktool.NewMock(cfg)

	// fmt.Println("Mock: ", mock)

	startTime := time.Now()

	// ----------------------- by headless brwoser -------------------------------
	// dir, err := os.MkdirTemp("", "chromedp-temp")
	// defer os.RemoveAll(dir)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// opts := append(chp.DefaultExecAllocatorOptions[:],
	// 	chp.DisableGPU,
	// 	chp.UserDataDir(dir),
	// 	chp.Flag("headless", !showHead),
	// )

	// allocCtx, cancel := chp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()

	// ----------------------- by database operations -------------------------------
	pg := pgstore.New(&pgstore.DBConfig{
		DSN: cfg.DB.GetDSN(),
	})

	err = pg.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer pg.CloseDB()

	err = pg.InitModules()
	if err != nil {
		log.Fatal(err)
	}

	dataStore := store.New(pg.Article, pg.User, pg.Role, pg.Permission, pg.Activity, pg.Message, pg.Category)

	var wg sync.WaitGroup
	policy := bluemonday.UGCPolicy()
	userSrv := &service.User{Store: dataStore, SantizePolicy: policy}
	articleSrv := &service.Article{Store: dataStore, SantizePolicy: policy, WG: &wg}

	// var wg sync.WaitGroup
	fmt.Println("os.Args", os.Args)
	fmt.Println("subcommand:", subcommand)
	if subcommand != "" {
		switch subcommand {
		case "reply":
			fmt.Println("Reply to article", replyToId, "with", replyNum, "replies")
			fmt.Println("Max level: ", maxLevel)
			fmt.Println("Reply number of reply: ", replyNumOfReply)
			if replyToId == 0 {
				log.Fatal("Article id is required")
			}

			wg.Add(replyNum)
			replyArticle(userSrv, articleSrv)
		case "users":
			fmt.Println("Seed users from file:", usersFile)
			seedUsers(userSrv, usersFile)
		case "categories":
			fmt.Println("Seed categories from file:", categoriesFile)
			categorySrv := service.NewCategory(dataStore)
			seedCategories(categorySrv, categoriesFile)
		default:
			wg.Add(articleNum)
			seedArticles(userSrv, articleSrv, startTime, categoryType)
		}
	} else {
		wg.Add(articleNum)
		seedArticles(userSrv, articleSrv, startTime, categoryType)
	}

	wg.Wait()
}
