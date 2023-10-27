package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type AppSettings struct {
	Paths    AppPaths  `json:"paths"`
	Server   AppServer `json:"server"`
	Amqp     Amqp      `json:"amqp"`
	Postgres Postgres  `json:"postgres"`
	Task     AppTask   `json:"task"`
}

type AppPaths struct {
	Database string `json:"database"`
	Queue    string `json:"queue"`
	Images   string `json:"images"`
}

type AppServer struct {
	Protocol string        `json:"protocol"`
	Host     string        `json:"host"`
	Urls     AppServerUrls `json:"urls"`
}

type AppServerUrls map[string]string

type AppTask struct {
	Items   []AppTaskItems `json:"items"`
	Timeout AppTaskTimeout `json:"timeout"`
	Filter  ArticleFilter  `json:"filter"`
}

type AppTaskItems struct {
	Region  string `json:"Region"`
	PageNum uint16 `json:"PageNum"`
}

type AppTaskTimeout struct {
	LapInterval     string `json:"lapInterval"`
	ReqInterval     string `json:"reqInterval"`
	ReqTTL          string `json:"reqTTL"`
	ReqEconnRefused string `json:"reqEconnRefused"`
	ReqError        string `json:"reqError"`
	ReqEmptyList    string `json:"reqEmptyList"`
}

type ArticleFilter struct {
	Id       map[string]uint64 `json:"Id"`
	RegionId map[string]uint8  `json:"RegionId"`
	Title    map[string]uint8  `json:"Title"`
	Message  map[string]uint8  `json:"Message"`
	Images   map[string]uint8  `json:"Images"`
	Phones   map[string]uint8  `json:"Phones"`
	Hash     map[string]uint8  `json:"Hash"`
	Link     map[string]uint8  `json:"Link"`
	Date     map[string]uint16 `json:"Date"`
}

type Amqp struct {
	Protocol string          `json:"protocol"`
	Host     string          `json:"host"`
	Port     int             `json:"port"`
	Vhost    string          `json:"virtual-host"`
	Creds    AmqpCredentials `json:"creds"`
}

type AmqpCredentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Postgres struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

func (a *Amqp) UrlWithCreds(user, password string) string {
	connectionLink := a.String()
	protocol := "amqp://"

	if user != "" && password != "" {
		hostCreds := protocol + user + ":" + password + "@"
		return strings.Replace(connectionLink, protocol, hostCreds, 1)
	}

	return connectionLink
}

func (a *Amqp) String() string {
	var (
		vhost string = a.Vhost
		host  string = a.Host
		port  int    = a.Port
	)

	if port == 0 {
		port = 5672
	}

	if host == "" {
		host = "localhost"
	}

	if vhost == "" {
		vhost = "/"
	}

	return fmt.Sprintf("amqp://%s:%d/%s", host, port, vhost)
}

func (a *Postgres) UrlWithCreds(user, password string) string {
	connectionLink := a.String()
	protocol := a.Protocol + "://"

	if protocol == "://" {
		protocol += "postgres"
	}

	if user != "" && password != "" {
		hostCreds := protocol + user + ":" + password + "@"
		return strings.Replace(connectionLink, protocol, hostCreds, 1)
	}

	return connectionLink
}

func (a *Postgres) String() string {
	var (
		SSLMode  string = a.SSLMode
		database string = a.Database
		protocol string = a.Protocol
		host     string = a.Host
		port     int    = a.Port
	)

	if port == 0 {
		port = 5432
	}

	if host == "" {
		host = "localhost"
	}

	if protocol == "" {
		protocol = "postgres"
	}

	if database == "" {
		database = "postgres"
	}

	if SSLMode == "" {
		database = "disable"
	}

	return fmt.Sprintf("%s://%s:%d/%s?sslmode=%s", a.Protocol, host, port, database, SSLMode)
}

func MustLoadSettings(jsonFilename string) (settings *AppSettings) {
	file, err := os.Open(jsonFilename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(contents, &settings); err != nil {
		panic(err)
	}

	return
}
