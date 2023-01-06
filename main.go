package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"gopkg.in/urfave/cli.v1"
)

type Conf struct {
	ResultsDisplayCount int
	ApiUrl              string
	AuthToken           string
}

func (conf *Conf) DefaultConf() {
	if conf.ResultsDisplayCount == 0 {
		conf.ResultsDisplayCount = 5
	}

	if conf.ApiUrl == "" {
		conf.ApiUrl = "https://mastodon.social"
	}
}

// todo: move to Utils - PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

type Response []struct {
	ID                 string      `json:"id"`
	CreatedAt          time.Time   `json:"created_at"`
	InReplyToID        interface{} `json:"in_reply_to_id"`
	InReplyToAccountID interface{} `json:"in_reply_to_account_id"`
	Sensitive          bool        `json:"sensitive"`
	SpoilerText        string      `json:"spoiler_text"`
	Visibility         string      `json:"visibility"`
	Language           string      `json:"language"`
	URI                string      `json:"uri"`
	URL                string      `json:"url"`
	RepliesCount       int         `json:"replies_count"`
	ReblogsCount       int         `json:"reblogs_count"`
	FavouritesCount    int         `json:"favourites_count"`
	EditedAt           interface{} `json:"edited_at"`
	Content            string      `json:"content"`
	Reblog             interface{} `json:"reblog"`
	Account            struct {
		ID             string        `json:"id"`
		Username       string        `json:"username"`
		Acct           string        `json:"acct"`
		DisplayName    string        `json:"display_name"`
		Locked         bool          `json:"locked"`
		Bot            bool          `json:"bot"`
		Discoverable   bool          `json:"discoverable"`
		Group          bool          `json:"group"`
		CreatedAt      time.Time     `json:"created_at"`
		Note           string        `json:"note"`
		URL            string        `json:"url"`
		Avatar         string        `json:"avatar"`
		AvatarStatic   string        `json:"avatar_static"`
		Header         string        `json:"header"`
		HeaderStatic   string        `json:"header_static"`
		FollowersCount int           `json:"followers_count"`
		FollowingCount int           `json:"following_count"`
		StatusesCount  int           `json:"statuses_count"`
		LastStatusAt   string        `json:"last_status_at"`
		Emojis         []interface{} `json:"emojis"`
		Fields         []interface{} `json:"fields"`
	} `json:"account"`
	MediaAttachments []struct {
		ID               string      `json:"id"`
		Type             string      `json:"type"`
		URL              string      `json:"url"`
		PreviewURL       string      `json:"preview_url"`
		RemoteURL        string      `json:"remote_url"`
		PreviewRemoteURL interface{} `json:"preview_remote_url"`
		TextURL          interface{} `json:"text_url"`
		Meta             struct {
			Original struct {
				Width  int     `json:"width"`
				Height int     `json:"height"`
				Size   string  `json:"size"`
				Aspect float64 `json:"aspect"`
			} `json:"original"`
			Small struct {
				Width  int     `json:"width"`
				Height int     `json:"height"`
				Size   string  `json:"size"`
				Aspect float64 `json:"aspect"`
			} `json:"small"`
		} `json:"meta"`
		Description interface{} `json:"description"`
		Blurhash    string      `json:"blurhash"`
	} `json:"media_attachments"`
	Mentions []interface{} `json:"mentions"`
	Tags     []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"tags"`
	Emojis []interface{} `json:"emojis"`
	Card   interface{}   `json:"card"`
	Poll   interface{}   `json:"poll"`
}

func main() {
	conf := Conf{}
	configs_file, err := os.Open("conf.json")
	if os.IsNotExist(err) {
		conf.DefaultConf()
	} else {
		defer configs_file.Close()
		decoder := json.NewDecoder(configs_file)
		err := decoder.Decode(&conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Name = "mastodonctl"
	app.Usage = "commandline client for a Mastodon social media user"

	app.Authors = append(app.Authors, cli.Author{Name: "socraticDev", Email: "socraticdev@gmail.com"})
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:      "userinfos",
			ShortName: "id",
			Usage:     "Retrieve Mastodon Account ID by username",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "username",
					Usage: "Username of the targetted Mastodon Account",
					Value: "socdev",
				},
			},
			Action: func(c *cli.Context) error {
				var token_val string
				if len(conf.AuthToken) > 0 {
					token_val = fmt.Sprintf("Bearer %s", conf.AuthToken)
				} else {
					token_val = fmt.Sprintf("Bearer %s", os.Getenv("BEARER_TOKEN"))
				}

				accounts, err := GetAccounts(InAccounts{
					Username:     c.String("username"),
					AuthToken:    token_val,
					ApiUrl:       conf.ApiUrl,
					ResultsCount: conf.ResultsDisplayCount,
				})
				if err != nil {
					log.Fatal(err)
				}

				headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
				columnFmt := color.New(color.FgHiBlue).SprintfFunc()

				tbl := table.New("id", "username", "displayname", "URL", "follower count", "following count")
				tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

				for _, r := range accounts {

					tbl.AddRow(r.ID, r.UserName, r.DisplayName, r.URL, r.FollowersCount, r.FollowingCount)
				}

				tbl.Print()

				return nil
			},
		}, {
			Name:      "hashtag",
			ShortName: "tag",
			Usage:     "Will get latest post informations about a specific hashtag",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "Keyword (hashtag) word to search for - without # sign!",
					Value: "cat",
				},
			},
			Action: func(c *cli.Context) error {
				hashtag := c.String("name")

				if len(hashtag) <= 0 {
					fmt.Println("Error: must provide a hashtag value to look for!")
					return nil
				}

				uri := fmt.Sprintf("%s/api/v1/timelines/tag/%s?limit=%d", conf.ApiUrl, hashtag, conf.ResultsDisplayCount)

				resp, err := http.Get(uri)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body) // response body is []byte

				var results Response
				if err := json.Unmarshal(body, &results); err != nil { // Parse []byte to go struct pointer
					log.Fatal(err)
				}

				headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
				columnFmt := color.New(color.FgHiBlue).SprintfFunc()

				tbl := table.New("hashtag", "username", "media url")
				tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

				for _, r := range results {
					url := "{no media}"

					if len(r.MediaAttachments) > 0 {
						url = r.MediaAttachments[len(r.MediaAttachments)-1].URL
					}

					tbl.AddRow(hashtag, r.Account.Username, url)
				}

				tbl.Print()

				return nil
			},
		},
	}

	app.Run(os.Args)
}
