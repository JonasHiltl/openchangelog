package analytics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jonashiltl/openchangelog/internal/store"
)

type Event struct {
	Time          time.Time `json:"time"`
	WID           string    `json:"wid"`
	CID           string    `json:"cid"`
	City          string    `json:"city"` // ISO 3166
	CountryCode   string    `json:"countryCode"`
	ContinentCode string    `json:"continentCode"`
	Lat           float64   `json:"lat"`
	Lng           float64   `json:"lng"`
	AccessDenied  bool      `json:"accessDenied"` // changelog protected
}

func NewEvent(r *http.Request, cl store.Changelog) Event {
	return newEvent(r, cl, false)
}

func NewAccessDeniedEvent(r *http.Request, cl store.Changelog) Event {
	return newEvent(r, cl, true)
}

func newEvent(r *http.Request, cl store.Changelog, denied bool) Event {
	return Event{
		CID:           cl.ID.String(),
		WID:           cl.WorkspaceID.String(),
		City:          r.Header.Get("cf-ipcity"),
		CountryCode:   r.Header.Get("cf-ipcountry"),
		ContinentCode: r.Header.Get("cf-ipcontinent"),
		Lat:           parseFloat(r.Header.Get("cf-iplatitude")),
		Lng:           parseFloat(r.Header.Get("cf-iplongitude")),
		AccessDenied:  denied,
		Time:          time.Now(),
	}
}

func parseFloat(str string) float64 {
	if str == "" {
		return 0
	}
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return f
}

type Emitter interface {
	Emit(Event)
}
