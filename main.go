package main

import (
	"bytes"
	"fmt"
	"github.com/akamensky/argparse"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	slack "serverCheck/Slack"
)

type Sever struct {
	Servers []string `yaml:"servers"`
	Email   string   `yaml:"email"`
	Slack   string   `yaml:"slack"`
}

type Flags struct {
	port     int
	filename string
}

func getFileData(fileName string) (*Sever, error) {

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	data := &Sever{}
	err = yaml.Unmarshal(buf, data)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return data, nil
}

type test struct {
	text string
}

func sendSlack(ic <-chan string) {

	attachment1 := slack.Attachment{}
	var b bytes.Buffer
	// attachment1.AddAction(slack.Action { Type: "button", Text: "Book flights 🛫", Url: "https://flights.example.com/book/r123456", Style: "primary" })
	// attachment1.AddAction(slack.Action { Type: "button", Text: "Cancel", Url: "https://flights.example.com/abandon/r123456", Style: "danger" })
	for c := range ic {
		b.WriteString(c)
	}

	attachment1.AddField(slack.Field{Title: "Author", Value: b.String()}).AddField(slack.Field{Title: "Status", Value: "Completed"})
	payload := slack.Payload{
		Text:        "Server Chek",
		Username:    "robot",
		IconEmoji:   ":monkey_face:",
		Attachments: []slack.Attachment{attachment1},
	}
	url := "webhook url 입력"
	payload.SendSlack(url)
}

func CheckURLs(data *Sever) <-chan string {
	oc := make(chan string)

	go func() {
		for i := 0; i < len(data.Servers); i++ {
			url := data.Servers[i]
			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			oc <- fmt.Sprintf("%s - %d \n", url, resp.StatusCode)
		}
		close(oc)
	}()

	return oc
}

func ParserSetting(flags *Flags) {
	parser := argparse.NewParser("test_program", "test argparse program")
	f := parser.String("v", "value", &argparse.Options{Required: false, Help: "Value to print", Default: "ini.yaml"})
	p := parser.Int("p", "port", &argparse.Options{Required: false, Help: "this port is testing", Default: 3000})
	parsingArgs(parser)
	flags.filename = *f
	flags.port = *p
}

func parsingArgs(parser *argparse.Parser) {
	err := parser.Parse(os.Args)
	if err != nil {
		panic(parser.Usage(err))
	}
}

func main() {
	var flags Flags
	ParserSetting(&flags)
	data, err := getFileData(flags.filename)
	if err != nil {
		panic(err)
	}
	sendSlack(CheckURLs(data))
}
