package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/alexsuraykin/SkillFactory_News/internal/models"
	storage "github.com/alexsuraykin/SkillFactory_News/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockFeedsRepository struct{}

func (m *mockFeedsRepository) Feeds(page int, pageSize int, title string, filter string) ([]models.Feeds, error) {
	return []models.Feeds{
		{Id: 1, Title: "Feed 1", Content: "Content 1"},
		{Id: 2, Title: "Feed 2", Content: "Content 2"},
	}, nil
}

func (m *mockFeedsRepository) StoreFeeds([]models.Feeds) error {
	return nil
}

func (m *mockFeedsRepository) FeedById(int) (*models.Feeds, error) {
	return nil, nil
}

func TestFeedsAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	repo := &mockFeedsRepository{}
	api := &API{
		storage: &storage.Storage{
			Feeds: repo,
		},
	}

	r.GET("/feeds", api.Feeds)

	limit := 5
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/feeds?"+"page="+strconv.Itoa(limit)+"&"+"title="+"&"+"filter=", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Feeds
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON response: %v", err)
	}

	expected := []models.Feeds{
		{Id: 1, Title: "Feed 1", Content: "Content 1"},
		{Id: 2, Title: "Feed 2", Content: "Content 2"},
	}

	assert.Equal(t, expected, response)
}
