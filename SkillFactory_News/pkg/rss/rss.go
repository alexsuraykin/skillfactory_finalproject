package rss

import (
	"encoding/xml"
	"io"
	"net/http"

	"github.com/alexsuraykin/SkillFactory_News/internal/models"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/rs/zerolog/log"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// Parse links from config.JSON (./cmd)
func Parse(url string) ([]models.Feeds, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return nil, err
	}

	var f Feed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal xml")
		return nil, err
	}

	var data []models.Feeds
	for _, item := range f.Channel.Items {
		var n models.Feeds
		n.Title = item.Title
		n.Content = item.Description
		n.Content = strip.StripTags(n.Content)
		n.Link = item.Link
		n.PubDate = item.PubDate

		data = append(data, n)
	}

	return data, nil
}
