package rss

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexsuraykin/SkillFactory_News/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	xmlData := `
	<rss version="2.0">
		<channel>
			<title>Sample Feed</title>
			<description>Sample RSS feed for testing</description>
			<link>http://example.com</link>
			<item>
				<title>Item 1</title>
				<description>Description of Item 1</description>
				<pubDate>Wed, 16 Jul 2024 00:00:00 +0000</pubDate>
				<link>http://example.com/item1</link>
			</item>
		</channel>
	</rss>`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, xmlData)
	}))
	defer mockServer.Close()

	url := mockServer.URL
	feeds, err := Parse(url)
	require.NoError(t, err, "Parse function returned an error")

	expectedFeeds := []models.Feeds{
		{
			Title:   "Item 1",
			Content: "Description of Item 1",
			Link:    "http://example.com/item1",
			PubDate: "Wed, 16 Jul 2024 00:00:00 +0000",
		},
	}

	assert.Equal(t, expectedFeeds, feeds, "Parsed feeds do not match expected")
}
