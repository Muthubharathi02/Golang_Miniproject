package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
)

func main() {
	log.Println("Connecting to server...")
	//connecting to the server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected")

	defer c.Logout()
	//providing email id and password
	if err := c.Login("kmuthupavithra@gmail.com", "9629461146"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")
	//selecting inbox
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	//selecting latest mail in inbox
	seqset := new(imap.SeqSet)
	seqset.AddRange(mbox.Messages, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchRFC822}, messages)
	}()
	fmt.Println(messages)
	for msg := range messages {
		for _, r := range msg.Body {

			entity, _ := message.Read(r)

			multiPartReader := entity.MultipartReader()
			count := 0
			for e, err := multiPartReader.NextPart(); err != io.EOF; e, err = multiPartReader.NextPart() {

				kind, _, cErr := e.Header.ContentType()

				if kind == "multipart/alternative" {
					continue
				}
				count++

				if cErr != nil {
					log.Fatal(cErr)
				}
				fmt.Println(kind)

				c, rErr := ioutil.ReadAll(e.Body)

				if rErr != nil {
					log.Fatal(rErr)
				}
				//downloading png format image attachments
				if kind == "image/png" {
					if fErr := ioutil.WriteFile("./output"+strconv.Itoa(count)+".png", c, 0777); fErr != nil {
						log.Fatal(fErr)
					}
				}
				//downloading jpeg format image attachments
				if kind == "image/jpeg" {
					if fErr := ioutil.WriteFile("./output"+strconv.Itoa(count)+".jpeg", c, 0777); fErr != nil {
						log.Fatal(fErr)
					}
				}
				//downloading pdf attachments
				if kind == "application/pdf" {
					if fErr := ioutil.WriteFile("./output"+strconv.Itoa(count)+".pdf", c, 0777); fErr != nil {
						log.Fatal(fErr)
					}
				}
				//saving html content of the mail in pdf format
				if kind == "text/html" {

					pdfg, err := wkhtml.NewPDFGenerator()
					if err != nil {
						return
					}
					htmlStr := string(c)
					pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(htmlStr)))
					err = pdfg.Create()
					if err != nil {
						log.Fatal(err)
					}

					err = pdfg.WriteFile("./output" + strconv.Itoa(count) + ".pdf")
					if err != nil {
						log.Fatal(err)
					}

				}

			}
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done")
}
