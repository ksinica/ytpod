package youtube

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	initialDataPrefx  = "var ytInitialData ="
	initialDataSuffix = ";"
)

type initialData struct {
	Header struct {
		C4TabbedHeaderRenderer struct {
			ChannelID string `json:"channelId"`
			Avatar    struct {
				Thumbnails []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnails"`
			} `json:"avatar"`
		} `json:"c4TabbedHeaderRenderer"`
	} `json:"header"`
	Metadata struct {
		ChannelMetadataRenderer struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ChannelURL  string `json:"channelUrl"`
		} `json:"channelMetadataRenderer"`
	} `json:"metadata"`
}

func (d *initialData) avatarURL() string {
	if len(d.Header.C4TabbedHeaderRenderer.Avatar.Thumbnails) > 0 {
		ret := d.Header.C4TabbedHeaderRenderer.Avatar.Thumbnails[0]
		for _, thumb := range d.Header.C4TabbedHeaderRenderer.Avatar.Thumbnails {
			if thumb.Width > ret.Width {
				ret = thumb
			}
		}
		return ret.URL
	}
	return ""
}

func findInitialData(r io.Reader) (ret initialData, err error) {
	var q *goquery.Document
	q, err = goquery.NewDocumentFromReader(r)
	if err != nil {
		return
	}
	q.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		text := s.Text()
		if strings.HasPrefix(text, initialDataPrefx) {
			text = strings.TrimPrefix(text, initialDataPrefx)
			text = strings.TrimSuffix(text, initialDataSuffix)
			err = json.Unmarshal([]byte(text), &ret)
			return false
		}
		return true
	})
	return
}

func fetchInitialData(ctx context.Context, url string) (ret initialData, err error) {
	err = httpFetch(ctx, url, func(r io.Reader) error {
		ret, err = findInitialData(r)
		return err
	})
	return
}
