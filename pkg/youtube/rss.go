package youtube

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/eduncan911/podcast"
	"github.com/go-chi/chi/v5"
)

const (
	rssAuthor = "YouTube"
	rssTTL    = 60 // 1 hour
)

func makeStreamURL(r *http.Request, videoID string) string {
	rctx := chi.RouteContext(r.Context())
	scheme := httpRequestScheme(r)
	if len(rctx.RoutePatterns) > 0 {
		var path string
		for _, x := range rctx.RoutePatterns[:len(rctx.RoutePatterns)-1] {
			path += strings.TrimRight(x, "/*")
		}
		return fmt.Sprintf(
			"%s://%s%s/stream/%s",
			scheme,
			r.Host,
			path,
			videoID,
		)
	}
	return fmt.Sprintf("%s://%s/stream/%s", scheme, r.Host, videoID)
}

func feedEntryToEnclosure(r *http.Request, entry feedEntry) *podcast.Enclosure {
	return &podcast.Enclosure{
		URL:  makeStreamURL(r, entry.VideoID),
		Type: podcast.M4A,
	}
}

func feedEntryToPodcastItem(r *http.Request, entry feedEntry) podcast.Item {
	return podcast.Item{
		GUID:        entry.VideoID,
		Title:       entry.Title,
		Description: entry.Group.Description,
		PubDate:     &entry.Published,
		Enclosure:   feedEntryToEnclosure(r, entry),
	}
}

func doHandleFeed(w http.ResponseWriter, r *http.Request) error {
	rctx := chi.RouteContext(r.Context())
	init, err := fetchInitialData(
		r.Context(),
		"https://youtube.com/"+strings.TrimPrefix(rctx.RoutePath, "/feed/"),
	)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(
		r.Context(),
		init.Header.C4TabbedHeaderRenderer.ChannelID,
	)
	if err != nil {
		return err
	}

	rss := podcast.New(
		feed.Title,
		init.Metadata.ChannelMetadataRenderer.ChannelURL,
		init.Metadata.ChannelMetadataRenderer.Description,
		&feed.Published,
		nil,
	)
	rss.AddImage(init.avatarURL())
	rss.IAuthor = rssAuthor
	rss.TTL = rssTTL

	for _, entry := range feed.Entry {
		if _, err := rss.AddItem(feedEntryToPodcastItem(r, entry)); err != nil {
			return err
		}
	}

	if len(rss.Items) > 0 {
		rss.AddPubDate(rss.Items[len(rss.Items)-1].PubDate)
	}

	w.Header().Set("Content-Type", "application/rss+xml")
	w.Header().Set("Content-Disposition", "attachment; filename=\"feed.xml\"")
	return rss.Encode(w)
}

func handleGetFeed(w http.ResponseWriter, r *http.Request) {
	if err := doHandleFeed(w, r); err != nil {
		if errors.Is(err, errNotFound) {
			http.Error(
				w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)
			return
		}

		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}
