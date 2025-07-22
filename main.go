package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Entry struct {
	XMLName xml.Name `xml:"entry"`
	Id      int      `xml:"id,attr"`
	Date    string   `xml:"date,attr"`
	Message string   `xml:"message,attr"`
}

type Log struct {
	XMLName xml.Name `xml:"log"`
	Entries []*Entry `xml:"entry"`
}

func (m Entry) String() string {
	return fmt.Sprintf("Message id=%v, date=%v, logmessage=%v", m.Id, m.Date, m.Message)
}

func (l Log) String() string {
	return fmt.Sprintf("Entries: %v", l.Entries)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func sync(installDir string) {
	FTP_SITE := os.Getenv("FTP_SITE")
	FTP_USERNAME := os.Getenv("FTP_USERNAME")
	FTP_PASSWORD := os.Getenv("FTP_PASSWORD")

	config := &ssh.ClientConfig{
		User: FTP_USERNAME,
		Auth: []ssh.AuthMethod{
			ssh.Password(FTP_PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, tcp_err := ssh.Dial("tcp", FTP_SITE, config)
	if tcp_err != nil {
		log.Fatal("failed to dial: ", tcp_err)
	}
	defer client.Close()

	sftp_client, _ := sftp.NewClient(client)
	defer sftp_client.Close()

	file, _ := sftp_client.Create("oplog.xml")
	file_data, _ := os.ReadFile(installDir)
	if _, err := file.Write(file_data); err != nil {
		log.Fatal(err)
	}

	file.Close()

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: oplog '<INSERT MEANINGFUL MESSAGE HERE>'")
		fmt.Println("\n\nin the future, when oplog is called with no args, there will a program similar to an empty git commit.\n a log in memory is opening in nano. you can modify the message, save and close to push.")
		os.Exit(1)
	}

	homeDir, _ := os.UserHomeDir()
	installDir := filepath.Join(homeDir, ".oplog")

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

	if message == "sync" {
		sync(installDir)
		os.Exit(0)
	}

	var entry Log
	timeStr := time.Now().Format("2006-01-02 3:04PM")

	data, fileNotExistErr := os.ReadFile(installDir)
	if os.IsNotExist(fileNotExistErr) {
		tempEntry := &Entry{Id: 1, Date: timeStr, Message: message}
		entry.Entries = []*Entry{tempEntry}
	} else {
		if err := xml.Unmarshal(data, &entry); err != nil {
			panic(err)
		}
		entry.Entries = append(entry.Entries, &Entry{Id: entry.Entries[len(entry.Entries)-1].Id + 1, Date: timeStr, Message: message})
	}

	out, _ := xml.MarshalIndent(entry, "", "\t")

	err := os.WriteFile(installDir, out, 0644)
	check(err)
}
