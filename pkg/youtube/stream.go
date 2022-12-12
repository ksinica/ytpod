package youtube

import (
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	yt "github.com/kkdai/youtube/v2"
	"github.com/samber/lo"
)

func handleGetStream(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "videoID")

	c := yt.Client{
		HTTPClient: httpClient(r.Context()),
	}

	video, err := c.GetVideoContext(
		r.Context(),
		"https://youtube.com/watch?v="+id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	fmts := lo.Filter(video.Formats, func(x yt.Format, _ int) bool {
		return strings.HasPrefix(x.MimeType, "audio/mp4") && x.Width == 0
	})

	sort.Slice(fmts, func(i, j int) bool {
		return fmts[i].Bitrate > fmts[j].Bitrate
	})

	if len(fmts) > 0 {
		url, err := c.GetStreamURLContext(
			r.Context(),
			video,
			&fmts[0],
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
		return
	}

	http.Error(w, "no streams found", http.StatusNotFound)
}
