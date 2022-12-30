package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gopkg.in/urfave/cli.v1"
)

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

// PrettyPrint to print struct in a readable way
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func main() {
	app := cli.NewApp()
	app.Usage = "commandline client for a Mastadon social media user"
	app.Commands = []cli.Command{
		{
			Name:      "hashtag",
			ShortName: "#",
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

				uri := fmt.Sprintf("https://mastodon.social/api/v1/timelines/tag/%s?limit=1", hashtag)

				resp, err := http.Get(uri)
				if err != nil {
					fmt.Println("Error querying Mastodon.social")
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body) // response body is []byte

				var results Response
				if err := json.Unmarshal(body, &results); err != nil { // Parse []byte to go struct pointer
					fmt.Println("Can not unmarshal JSON")
				}

				for _, r := range results {
					url := r.MediaAttachments[len(r.MediaAttachments)-1].URL
					fmt.Println(fmt.Sprintf("Latest %s pic at this URL: %s", hashtag, url))
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}
