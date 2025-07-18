package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: logcli '<INSERT MEANINGFUL MESSAGE HERE>'")
		fmt.Println("\n\nin the future, when logcli is called with no args, there will a program similar to an empty git commit.\n a log in memory is opening in nano. you can modify the message, save and close to push.")
		os.Exit(1)
	}

	homeDir, _ := os.UserHomeDir()
	installDir := filepath.Join(homeDir, ".logcli")

	message := os.Args[1]

	if message == "print" {
		data, err := os.ReadFile(installDir)
		check(err)

		println(string(data))
		os.Exit(0)
	}

	entryExample1 := &Entry{Id: 1, Date: time.Now(), Message: message}
	out, _ := xml.MarshalIndent(entryExample1, "", "\t")

	println(xml.Header, string(out))

	err := os.WriteFile(installDir, out, 0644)
	check(err)

}
