package load

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetQueryIDs(t *testing.T) {
	tables := []struct {
		url   string
		hxURL string
		wID   string
		cID   string
	}{
		{
			url: "/?wid=ws_1&cid=cl_1",
			wID: "ws_1",
			cID: "cl_1",
		},
		{
			url: "/",
			wID: "",
			cID: "",
		},
		{
			hxURL: "http://localhost:6001/?wid=ws_2&cid=cl_2",
			wID:   "ws_2",
			cID:   "cl_2",
		},
	}

	for _, table := range tables {
		u, err := url.Parse(table.url)
		if err != nil {
			t.Error(err)
		}

		r := &http.Request{URL: u, Header: http.Header{}}
		if table.hxURL != "" {
			r.Header.Set("HX-Current-URL", table.hxURL)
		}

		wID, cID := GetQueryIDs(r)
		if wID != table.wID {
			t.Errorf("Expected %s to equals %s", wID, table.wID)
		}
		if cID != table.cID {
			t.Errorf("Expected %s to equals %s", cID, table.cID)
		}
	}
}
