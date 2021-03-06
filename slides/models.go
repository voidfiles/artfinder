package slides

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Binding struct {
	In  chan Slide
	Out chan Slide
}

// Site contains identifying info about what parent site an image was found
type Site struct {
	Title string `json:"title,omitempty" binding:"required"`
	URL   string `json:"url,omitempty" binding:"required"`
}

// Page contains identifying info about where an image was found
type Page struct {
	Title     string     `json:"title,omitempty" binding:"required"`
	URL       string     `json:"url,omitempty" binding:"required"`
	Published *time.Time `json:"published,omitempty"`
	GUIDHash  string     `json:"guid_hash,omitempty"`
}

// ImageInfo is important information about a single image
type ImageInfo struct {
	URL         string `json:"url,omitempty" binding:"required"`
	Width       int    `json:"width,omitempty" binding:"required"`
	Height      int    `json:"height,omitempty" binding:"required"`
	ContentType string `json:"content_type,omitempty" binding:"required"`
	Filename    string `json:"filename,omitempty" binding:"required"`
}

// ArtistInfo is important information about the artist who made the image
type ArtistInfo struct {
	Name         string   `json:"name,omitempty" validate:"required"`
	ArtsyURL     string   `json:"artsy_url,omitempty"`
	WikipediaURL string   `json:"wikipedia_url,omitempty"`
	WebsiteURL   string   `json:"website_url,omitempty"`
	FeedURL      string   `json:"feed_url,omitempty"`
	InstagramURL string   `json:"instagram_url,omitempty"`
	TwitterURL   string   `json:"twitter_url,omitempty"`
	Description  string   `json:"description,omitempty"`
	Feeds        []string `json:"feeds,omitempty"`
	Sites        []string `json:"sites,omitempty"`
}

// WorkInfo is important information about the work
type WorkInfo struct {
	Name string `json:"name,omitempty"`
}

// RenderHash storage for various renders of the slide its self
type RenderHash struct {
	Blog  string `json:"blog,omitempty"`
	Slide string `json:"slide,omitempty"`
}

// Slide bundles together important information about a find
type Slide struct {
	Site           Site          `json:"site,omitempty"`
	Page           Page          `json:"page,omitempty"`
	Content        string        `json:"content,omitempty"`
	Edited         time.Time     `json:"edited,omitempty"`
	GUIDHash       string        `json:"guid_hash,omitempty" binding:"required"`
	SourceImageURL string        `json:"source_image_url,omitempty"`
	ArchivedImage  *ImageInfo    `json:"archived_image,omitempty"`
	ArtistInfo     *ArtistInfo   `json:"artist_info,omitempty"`
	ArtistsInfo    []*ArtistInfo `json:"artists,omitempty"`
	WorkInfo       *WorkInfo     `json:"work_info,omitempty"`
	RenderHash     RenderHash    `json:"render_hash,omitempty"`
}

func (s Slide) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Slide) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed.")
	}

	err := json.Unmarshal(source, &s)
	if err != nil {
		return err
	}

	return nil
}
