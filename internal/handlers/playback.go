package handlers

import (
	"time"

	"github.com/Drack112/go-youtube/internal/api"
	"github.com/Drack112/go-youtube/internal/flags"
	"github.com/Drack112/go-youtube/internal/models"
	"github.com/Drack112/go-youtube/pkg/logger"
)

const maxRetries = 3

func SearchWithRetries(input *flags.Options) ([]models.SearchResult, error) {
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		var results []models.SearchResult

		logger.Debug("[Handler] Search attempt ", "attempt", attempt, "for", input.Input)
		results, err := api.SearchVideos(input.Input)
		if err == nil {
			return results, nil
		}

		logger.Warn("[Handler]", "attempt", attempt, "error", err)
		time.Sleep(300 * time.Millisecond)
	}

	return nil, err
}
