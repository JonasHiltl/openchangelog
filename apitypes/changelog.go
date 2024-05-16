package apitypes

import (
	"encoding/json"
	"time"
)

// Used to marshal/unmarshal a domain Changelog to json
type Changelog struct {
	Id          int64     `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Title       string    `json:"title,omitempty"`
	Subtitle    string    `json:"subtitle,omitempty"`
	Logo        Logo      `json:"logo"`
	Source      Source    `json:"source"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (c *Changelog) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	// We'll store the error (if any) so we can return it if necessary
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	// Unmarshal all fields except for 'source' into the Changelog struct
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	if sourceRaw, ok := objMap["source"]; ok {
		var sourceMap map[string]*json.RawMessage
		err := json.Unmarshal(b, &objMap)
		if err != nil {
			return err
		}
		typeRaw, ok := sourceMap["type"]
		if !ok {
			// No source type specified, so no source is set.
			return nil
		}

		var Type string
		err = json.Unmarshal(*typeRaw, &Type)
		if err != nil {
			return err
		}

		switch Type {
		case string(GitHub):
			var ghSource GHSource
			err = json.Unmarshal(*sourceRaw, &ghSource)
			if err != nil {
				return err
			}
			c.Source = ghSource
		}
	}

	return nil
}

type Logo struct {
	Src    string `json:"src"`
	Link   string `json:"link"`
	Alt    string `json:"alt"`
	Height string `json:"height"`
	Width  string `json:"width"`
}
