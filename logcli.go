package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type Entry struct {
	XMLName xml.Name  `xml:"entry"`
	Id      int       `xml:"id,attr"`
	Date    time.Time `xml:"date"`
	Message string    `xml:"message"`
}

func (m Entry) String() string {
	return fmt.Sprintf("Message id=%v, date=%v, logmessage=%v", m.Id, m.Date, m.Message)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("here will be the git commit like program that opens a log in memory like it does in nano. lets you modify the message. review git commit for more details.")
	} else {
		message := os.Args[1]

		entryExample1 := &Entry{Id: 1, Date: time.Now(), Message: message}
		out, _ := xml.MarshalIndent(entryExample1, "", "\t")

		println(xml.Header, string(out))
	}

	// err := os.WriteFile("~/.logcli")
}
