package youtube

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type feedEntry struct {
	VideoID   string    `xml:"videoId"`
	ChannelID string    `xml:"channelId"`
	Title     string    `xml:"title"`
	Published time.Time `xml:"published"`
	Updated   time.Time `xml:"updated"`
	Group     struct {
		Thumbnail struct {
			URL string `xml:"url,attr"`
		} `xml:"thumbnail"`
		Content struct {
			URL string `xml:"url,attr"`
		} `xml:"content"`
		Description string `xml:"description"`
	} `xml:"group"`
}

type feed struct {
	XMLName xml.Name `xml:"feed"`
	Link    []struct {
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	ID     string `xml:"id"`
	Title  string `xml:"title"`
	Author struct {
		Name string `xml:"name"`
		URI  string `xml:"uri"`
	} `xml:"author"`
	Published time.Time   `xml:"published"`
	Entry     []feedEntry `xml:"entry"`
}

func fetchFeed(ctx context.Context, channelID string) (ret feed, err error) {
	err = httpFetch(
		ctx,
		fmt.Sprintf(
			"https://youtube.com/feeds/videos.xml?channel_id=%s",
			channelID,
		),
		func(r io.Reader) error {
			return xml.NewDecoder(r).Decode(&ret)
		},
	)
	return
}
