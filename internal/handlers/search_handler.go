package handlers

import (
	"time"

	"github.com/Drack112/go-youtube/internal/api"
	"github.com/Drack112/go-youtube/internal/flags"
	"github.com/Drack112/go-youtube/internal/models"
	"github.com/Drack112/go-youtube/internal/ui"

	"github.com/Drack112/go-youtube/pkg/logger"
)

func SearchWithRetries(value *flags.Options) string {
	maxAttempts := 3
	var results []models.SearchResult
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.Debug("[Handler] Search attempt", "attempt", attempt, "for", value.Input)

		results, err = api.SearchVideos(value.Input)
		if err == nil {
			break
		}

		logger.Warn("[Handler]", "attempt", attempt, "error", err)
		if attempt < maxAttempts {
			time.Sleep(time.Second * time.Duration(attempt))
		}
	}

	if err != nil {
		return ui.CreateErrorBox("Search Failed", err.Error())
	}

	if len(results) == 0 {
		return ui.CreateErrorBox("No Results", "No videos found for your search")
	}

	return ui.CreateSearchResultsView(results)
}
