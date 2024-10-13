package apitypes

import (
	"encoding/json"
	"time"
)

// Represents the Changelog returned by the API via json encoding.
// Implements json un-/marshaling.
type Changelog struct {
	ID            string
	WorkspaceID   string
	Subdomain     string
	Domain        NullString
	Title         NullString
	Subtitle      NullString
	ColorScheme   ColorScheme
	HidePoweredBy bool
	Logo          Logo
	Source        Source
	CreatedAt     time.Time
}

type ColorScheme string

const (
	Dark   ColorScheme = "dark"
	Light  ColorScheme = "light"
	System ColorScheme = "system"
)

type FullChangelog struct {
	Changelog       Changelog `json:"changelog"`
	Articles        []Article `json:"articles"`
	HasMoreArticles bool      `json:"hasMoreArticles"`
}

type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"publishedAt"`
	Tags        []string  `json:"tags"`
	HTMLContent string    `json:"htmlContent"`
}

func (l Changelog) MarshalJSON() ([]byte, error) {
	obj := struct {
		ID            string     `json:"id"`
		WorkspaceID   string     `json:"workspaceId"`
		Subdomain     string     `json:"subdomain,omitempty"`
		Title         string     `json:"title,omitempty"`
		Domain        string     `json:"domain,omitempty"`
		Subtitle      string     `json:"subtitle,omitempty"`
		ColorScheme   string     `json:"colorScheme,omitempty"`
		HidePoweredBy bool       `json:"hidePoweredBy"`
		Logo          *Logo      `json:"logo,omitempty"`
		Source        Source     `json:"source,omitempty"`
		CreatedAt     *time.Time `json:"createdAt,omitempty"`
	}{
		ID:            l.ID,
		WorkspaceID:   l.WorkspaceID,
		Subdomain:     l.Subdomain,
		Domain:        l.Domain.V(),
		Title:         l.Title.V(),
		Subtitle:      l.Subtitle.V(),
		ColorScheme:   string(l.ColorScheme),
		HidePoweredBy: l.HidePoweredBy,
		Source:        l.Source,
	}

	if l.Logo.IsValid() {
		obj.Logo = &l.Logo
	}

	if !l.CreatedAt.IsZero() {
		obj.CreatedAt = &l.CreatedAt
	}

	return json.Marshal(obj)
}

func (c *Changelog) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	// We'll store the error (if any) so we can return it if necessary
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	if idRaw, ok := objMap["id"]; ok {
		err = json.Unmarshal(*idRaw, &c.ID)
		if err != nil {
			return err
		}
	}

	if workspaceIdRaw, ok := objMap["workspaceId"]; ok {
		err = json.Unmarshal(*workspaceIdRaw, &c.WorkspaceID)
		if err != nil {
			return err
		}
	}

	if subdomainRaw, ok := objMap["subdomain"]; ok {
		err = json.Unmarshal(*subdomainRaw, &c.Subdomain)
		if err != nil {
			return err
		}
	}

	if domainRaw, ok := objMap["domain"]; ok {
		err = json.Unmarshal(*domainRaw, &c.Domain)
		if err != nil {
			return err
		}
	}

	if titleRaw, ok := objMap["title"]; ok {
		err = json.Unmarshal(*titleRaw, &c.Title)
		if err != nil {
			return err
		}
	}

	if subtitleRaw, ok := objMap["subtitle"]; ok {
		err = json.Unmarshal(*subtitleRaw, &c.Subtitle)
		if err != nil {
			return err
		}
	}

	if colorSchemeRaw, ok := objMap["colorScheme"]; ok {
		err = json.Unmarshal(*colorSchemeRaw, &c.ColorScheme)
		if err != nil {
			return err
		}
	}

	if hidePoweredByRaw, ok := objMap["hidePoweredBy"]; ok {
		err = json.Unmarshal(*hidePoweredByRaw, &c.HidePoweredBy)
		if err != nil {
			return err
		}
	}

	if logoRaw, ok := objMap["logo"]; ok {
		err = json.Unmarshal(*logoRaw, &c.Logo)
		if err != nil {
			return err
		}
	}

	if createdAtRaw, ok := objMap["createdAt"]; ok {
		err = json.Unmarshal(*createdAtRaw, &c.CreatedAt)
		if err != nil {
			return err
		}
	}

	if sourceRaw, ok := objMap["source"]; ok && sourceRaw != nil {
		c.Source = DecodeSource(*sourceRaw)
	}

	return nil
}

func DecodeSource(in json.RawMessage) Source {
	var sourceMap map[string]json.RawMessage
	err := json.Unmarshal(in, &sourceMap)
	if err != nil {
		return nil
	}

	typeRaw, ok := sourceMap["type"]
	if !ok {
		// No source type specified, so no source is set.
		return nil
	}

	var Type string
	err = json.Unmarshal(typeRaw, &Type)
	if err != nil {
		return nil
	}

	switch Type {
	case string(GitHub):
		var ghSource GHSource
		err = json.Unmarshal(in, &ghSource)
		if err != nil {
			return nil
		}
		return ghSource
	}
	return nil
}

type Logo struct {
	Src    NullString
	Link   NullString
	Alt    NullString
	Height NullString
	Width  NullString
}

// omits fields if they are empty
func (l Logo) MarshalJSON() ([]byte, error) {
	data := make(map[string]NullString)

	if !l.Src.IsZero() {
		data["src"] = l.Src
	}
	if !l.Link.IsZero() {
		data["link"] = l.Link
	}
	if !l.Alt.IsZero() {
		data["alt"] = l.Alt
	}
	if !l.Height.IsZero() {
		data["height"] = l.Height
	}
	if !l.Width.IsZero() {
		data["width"] = l.Width
	}
	return json.Marshal(data)
}

// Returns true if at least one field is valid
func (l Logo) IsValid() bool {
	return l.Src.IsValid() || l.Link.IsValid() || l.Alt.IsValid() || l.Height.IsValid() || l.Width.IsValid()
}

type CreateChangelogBody struct {
	Title         NullString  `json:"title"`
	Subtitle      NullString  `json:"subtitle"`
	Logo          Logo        `json:"logo"`
	Domain        NullString  `json:"domain"`
	ColorScheme   ColorScheme `json:"colorScheme"`
	HidePoweredBy bool        `json:"hidePoweredBy"`
}

type UpdateChangelogBody struct {
	Title         NullString  `json:"title"`
	Subtitle      NullString  `json:"subtitle"`
	Logo          Logo        `json:"logo"`
	Domain        NullString  `json:"domain"`
	ColorScheme   ColorScheme `json:"colorScheme"`
	Subdomain     NullString  `json:"subdomain"`
	HidePoweredBy *bool       `json:"hidePoweredBy,omitempty"`
}
