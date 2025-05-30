// Configs for the app runtime, env file path default to ./.env.local

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	SessionSecret      string `env:"SESSION_SECRET"`
	CSRFSecret         string `env:"CSRF_SECRET"`
	DomainName         string `env:"DOMAIN_NAME" envDefault:"localhost"`
	AppPort            int    `env:"APP_PORT" envDefault:"3000"`
	AppOuterPort       int    `env:"APP_OUTER_PORT" envDefault:"3000"`
	NginxPort          int    `env:"NGINX_PORT" envDefault:"80"`
	NginxSSLPort       int    `env:"NGINX_SSL_PORT" envDefault:"443"`
	Debug              bool   `env:"DEBUG" envDefault:"false"`
	BrandName          string `env:"BRAND_NAME"`
	BrandDomainName    string `env:"BRAND_DOMAIN_NAME"`
	Slogan             string `env:"SLOGAN"`
	DB                 *DBConfig
	ReplyDepthPageSize int
	AdminEmail         string `env:"ADMIN_EMAIL"`
	Redis              *RedisConfig
	SMTP               *SMTPConfig
	Testing            bool   `env:"TEST"`
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GithubClientID     string `env:"GITHUB_CLIENT_ID"`
	GithubClientSecret string `env:"GITHUB_CLIENT_SECRET"`
	CloudflareSiteKey  string `env:"CLOUDFLARE_SITE_KEY"`
	CloudflareSecret   string `env:"CLOUDFLARE_SECRET"`
}

func (ac *AppConfig) GetServerURL() string {
	protocol := "http"
	if ac.NginxSSLPort == 443 {
		protocol += "s"
	}

	port := ":" + strconv.Itoa(ac.NginxPort)
	if port == ":80" {
		port = ""
	}

	return fmt.Sprintf("%s://%s%s", protocol, ac.DomainName, port)
}

// Get app host as host:port
// func (ac *AppConfig) GetHost() string {
// 	return fmt.Sprintf("%s:%d", ac.DomainName, ac.AppOuterPort)
// }

type DBConfig struct {
	DBHost              string `env:"DB_HOST"`
	DBName              string `env:"DB_NAME"`
	DBPort              int    `env:"DB_PORT"`
	DBUser              string `env:"DB_USER"`
	DBPassword          string `env:"DB_PASSWORD"`
	UserDefaultPassword string `env:"USER_DEFAULT_PASSWORD"`
}

func (dbCfg *DBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbCfg.DBUser,
		dbCfg.DBPassword,
		dbCfg.DBHost,
		dbCfg.DBPort,
		dbCfg.DBName,
	)
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     string `env:"REDIS_PORT"`
	User     string `env:"REDIS_USER"`
	Password string `env:"REDIS_PASSWORD"`
}

type SMTPConfig struct {
	Server     string `env:"SMTP_SERVER"`
	ServerPort string `env:"SMTP_SERVER_PORT"`
	User       string `env:"SMTP_USER"`
	Sender     string `env:"SMTP_SENDER"`
	Password   string `env:"SMTP_PASSWORD"`
}

var Config *AppConfig

// var BrandName = "DizKaz"
var BrandName = "笛卡"
var BrandDomainName = "dizkaz.com"
var Slogan = ""
var ReplyDepthPageSize = 10

// For testing
const SuperCode = "686868"

func Init(envFile string) error {
	cfg, err := Parse(envFile)
	if err != nil {
		return err
	}
	Config = cfg
	return nil
}

func InitFromEnv() error {
	cfg, err := ParseFromEnv()
	if err != nil {
		return err
	}
	Config = cfg
	return nil
}

func NewTest() (*AppConfig, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("can't get caller info")
	}

	// fmt.Println("caller filename: ", filename)

	testingEnvFile := filepath.Join(filepath.Dir(filename), "./.env.testing")

	// fmt.Println("testing env file: ", testingEnvFile)

	if _, err := os.Stat(testingEnvFile); err != nil {
		return nil, err
	}
	// if err != nil {
	// 	return nil, err
	// }
	return Parse(testingEnvFile)
}

// Parse env file and generate AppConfig struct
func Parse(envFile string) (*AppConfig, error) {
	if err := godotenv.Load(envFile); err != nil {
		return nil, err
	}

	return ParseFromEnv()
}

func ParseFromEnv() (*AppConfig, error) {
	dbCfg := &DBConfig{}
	if err := env.Parse(dbCfg); err != nil {
		return nil, err
	}

	rdbCfg := &RedisConfig{}
	if err := env.Parse(rdbCfg); err != nil {
		return nil, err
	}

	smtpCfg := &SMTPConfig{}
	if err := env.Parse(smtpCfg); err != nil {
		return nil, err
	}

	cfg := &AppConfig{
		DB:                 dbCfg,
		BrandName:          BrandName,
		BrandDomainName:    BrandDomainName,
		Slogan:             Slogan,
		ReplyDepthPageSize: ReplyDepthPageSize,
		Redis:              rdbCfg,
		SMTP:               smtpCfg,
	}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
