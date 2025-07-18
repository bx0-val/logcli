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

type Log struct {
	XMLName xml.Name `xml:"log"`
	Entries []*Entry `xml:"entry"`
}

func (m Entry) String() string {
	return fmt.Sprintf("Message id=%v, date=%v, logmessage=%v", m.Id, m.Date, m.Message)
}

func (l Log) String() string {
	return fmt.Sprintf("Entries: ", l.Entries)
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
		if os.IsNotExist(err) {
			println("No log file found. Make a log entry to create it.")
			os.Exit(0)
		}

		println(string(data))
		os.Exit(0)
	}

	entry := &Log{}
	tempEntry := &Entry{Id: 1, Date: time.Now(), Message: message}
	entry.Entries = []*Entry{tempEntry}
	out, _ := xml.MarshalIndent(entry, "", "\t")

	var in Log
	if err := xml.Unmarshal(out, &in); err != nil {
		panic(err)
	}

	// println(in)

	println("Inserted " + in.String() + " at " + installDir)

	err := os.WriteFile(installDir, out, 0644)
	check(err)

}
