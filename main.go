package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"
)

var token = flag.String("token", os.Getenv("LINE_ACCESS_TOKEN"), "LINE Access Token")

type Build struct {
	RepositoryId int    `json:"repository_id"`
	EventType    string `json:"event_type"`
	FinishedAt   string `json:"finished_at"`
	Number       string `json:"number"`
	State        string `json:"state"`
	Result       int    `json:"result"`
	Branch       string `json:"branch"`
	Duration     int    `json:"duration"`
	Commit       string `json:"commit"`
	Message      string `json:"message"`
	StartedAt    string `json:"started_at"`
	Id           int    `json:"id`
}

func main() {
	flag.Parse()

	idmap := map[int]bool{}

	var wg sync.WaitGroup
	for _, proj := range os.Args {
		wg.Add(1)
		go func(proj string) {
			first := true
			for {
				r, err := http.Get(fmt.Sprintf("https://api.travis-ci.org/repositories/%s/builds.json", proj))
				if err != nil {
					log.Println(err)
					time.Sleep(30 * time.Second)
					continue
				}
				defer r.Body.Close()

				var builds []Build
				json.NewDecoder(r.Body).Decode(&builds)

				for _, build := range builds {
					if _, ok := idmap[build.Id]; ok {
						continue
					}
					if !first {
						log.Printf("%s: #%d %s", proj, build.Id, build.State)
					}
					if build.State != "finished" {
						continue
					}
					idmap[build.Id] = true

					if first {
						continue
					}

					message := proj
					icon := ""
					if build.Result != 0 {
						message += "(failed): "
						icon = "https://raw.githubusercontent.com/mattn/travis-notify/gh-pages/failed.png"
					} else {
						message += "(success): "
						icon = "https://raw.githubusercontent.com/mattn/travis-notify/gh-pages/success.png"
					}
					message += fmt.Sprintf("https://travis-ci.org/%s/builds/%d", proj, build.Id)

					buf := bytes.Buffer{}
					mw := multipart.NewWriter(&buf)
					mw.WriteField("message", message)
					mw.WriteField("imageThumbnail", icon)
					mw.WriteField("imageFullsize", icon)
					mw.Close()
					req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", &buf)
					if err != nil {
						log.Print(err)
						continue
					}
					req.Header.Add("Content-Type", mw.FormDataContentType())
					req.Header.Set("Authorization", "Bearer "+*token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Print(err)
						continue
					}
					if b, err := ioutil.ReadAll(resp.Body); err == nil {
						log.Print(string(b))
					} else {
						log.Print(err)
					}
					resp.Body.Close()
				}
				first = false

				time.Sleep(30 * time.Second)
			}
			wg.Done()
		}(proj)
	}
	wg.Wait()
}
