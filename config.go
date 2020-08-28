package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "bojan",
		Password: "bojan",
		Name:     "lenslocked",
	}
}

type Config struct {
	Port        int            `json:"port"`
	Env         string         `json:"env"`
	Pepper      string         `json:"pepper"`
	HMACKey     string         `json:"hmac_key"`
	BasePath    string         `json:"basepath`
	DatabaseURL string         `json:"databaseurl"`
	Database    PostgresConfig `json:"database"`
	Mailgun     MailgunConfig  `json:"mailgun"`
	Dropbox     OauthConfig    `json:"dropbox"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "secret-random-string",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
	}
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

type OauthConfig struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	AuthURL  string `json:"auth_url"`
	TokenURL string `json:"token_url"`
}

func LoadConfig(configReq bool) Config {
	f, err := os.Open(".env")
	if err != nil {
		if configReq {
			panic(err)
		}
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .env")
	return c
}

func LoadFromENV() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env")
		return Config{}, err
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return Config{}, err
	}
	env := os.Getenv("ENV")
	pepper := os.Getenv("PEPPER")
	hmac := os.Getenv("HMAC_KEY")
	basepath := os.Getenv("BASEPATH")
	dburl := os.Getenv("DATABASEURL")
	mgApiKey := os.Getenv("MAILGUN_APIKEY")
	mgpubkey := os.Getenv("MAILGUN_PUBLIC_API_KEY")
	mgdomain := os.Getenv("MAILGUN_DOMAIN")
	dropboxid := os.Getenv("DROPBOX_ID")
	dropboxSecret := os.Getenv("DROPBOX_SECRET")
	dropboxAuthUrl := os.Getenv("DROPBOX_AUTHURL")
	dropbboxTokenUrl := os.Getenv("DROPBOX_TOKENURL")
	return Config{
		Port:        port,
		Env:         env,
		Pepper:      pepper,
		HMACKey:     hmac,
		BasePath:    basepath,
		DatabaseURL: dburl,
		Mailgun: MailgunConfig{
			APIKey:       mgApiKey,
			Domain:       mgdomain,
			PublicAPIKey: mgpubkey,
		},
		Dropbox: OauthConfig{
			ID:       dropboxid,
			AuthURL:  dropboxAuthUrl,
			Secret:   dropboxSecret,
			TokenURL: dropbboxTokenUrl,
		},
	}, nil

}
