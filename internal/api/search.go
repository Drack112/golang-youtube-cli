package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Drack112/go-youtube/internal/models"

	"github.com/Drack112/go-youtube/pkg/logger"
	"github.com/Drack112/go-youtube/pkg/utils"
)

const youtubeSearchBase = "https://www.youtube.com/results?search_query="

var initialDataRegex = regexp.MustCompile(`var ytInitialData = (\{.*?\});`)

type SearchResponse struct {
	Results           []models.SearchResult
	ContinuationToken string
	HasMore           bool
}

func SearchVideos(input string) ([]models.SearchResult, error) {
	resp, err := SearchVideosWithPagination(input, "")
	if err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func SearchVideosWithPagination(input string, continuationToken string) (*SearchResponse, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}

	logger.Debug("[SearchVideos] raw input ", "input", input)

	if utils.IsYouTubeURL(input) {
		logger.Debug("[SearchVideosWithPagination] detected YouTube URL", "input", input)

		clean := utils.CleanYoutubeLink(input)
		id := utils.ExtractVideoID(clean)

		if id == "" {
			logger.Error("[SearchVideosWithPagination] failed to extract Id from ", "clean", clean)
			return nil, errors.New("invalid YouTube link")
		}

		// Try to fetch the watch page and extract initialPlayerResponse for richer metadata
		watchURL := "https://www.youtube.com/watch?v=" + id
		body, err := utils.Fetch(watchURL)
		if err != nil {
			logger.Warn("[SearchVideosWithPagination] failed to fetch watch page, falling back to basic data", "error", err)
			return &SearchResponse{
				Results: []models.SearchResult{{
					ID:          id,
					Title:       "(Fetching metadata...)",
					URL:         watchURL,
					Thumbnail:   "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg",
					Duration:    "",
					DurationSec: 0,
				}},
				ContinuationToken: "",
				HasMore:           false,
			}, nil
		}

		// look for ytInitialPlayerResponse object
		playerRegex := regexp.MustCompile(`ytInitialPlayerResponse\s*=\s*(\{.*?\});`)
		match := playerRegex.FindStringSubmatch(body)
		if len(match) < 2 {
			// fallback to basic
			logger.Warn("[SearchVideosWithPagination] initialPlayerResponse not found, using fallback")
			return &SearchResponse{
				Results: []models.SearchResult{{
					ID:          id,
					Title:       "(Fetching metadata...)",
					URL:         watchURL,
					Thumbnail:   "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg",
					Duration:    "",
					DurationSec: 0,
				}},
				ContinuationToken: "",
				HasMore:           false,
			}, nil
		}

		var resp map[string]any
		if err := json.Unmarshal([]byte(match[1]), &resp); err != nil {
			logger.Warn("[SearchVideosWithPagination] failed to unmarshal player response", "error", err)
			return &SearchResponse{
				Results: []models.SearchResult{{
					ID:          id,
					Title:       "(Fetching metadata...)",
					URL:         watchURL,
					Thumbnail:   "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg",
					Duration:    "",
					DurationSec: 0,
				}},
				ContinuationToken: "",
				HasMore:           false,
			}, nil
		}

		videoDetails := utils.DeepGet(resp, "videoDetails")
		title := ""
		author := ""
		lengthSec := 0
		isLive := false
		channelID := ""

		if vd, ok := videoDetails.(map[string]any); ok {
			title = utils.Str(vd["title"])
			author = utils.Str(vd["author"])
			if ls := utils.Str(vd["lengthSeconds"]); ls != "" {
				if n, err := strconv.Atoi(ls); err == nil {
					lengthSec = n
				}
			}
			if idv := utils.Str(vd["videoId"]); idv != "" {
				// ensure URL uses normalized id
				id = idv
			}
			if cid := utils.Str(vd["channelId"]); cid != "" {
				channelID = cid
			}
			if isLiveVal, ok := vd["isLiveContent"].(bool); ok {
				isLive = isLiveVal
			}
		}

		return &SearchResponse{
			Results: []models.SearchResult{{
				ID:          id,
				Title:       title,
				URL:         "https://www.youtube.com/watch?v=" + id,
				Thumbnail:   "https://i.ytimg.com/vi/" + id + "/hqdefault.jpg",
				Duration:    "",
				DurationSec: lengthSec,
				ChannelName: author,
				ChannelID:   channelID,
				IsLive:      isLive,
			}},
			ContinuationToken: "",
			HasMore:           false,
		}, nil
	}

	logger.Debug("[SearchVideosWithPagination] performing search for ", "input", input)

	url := "https://www.youtube.com/results?search_query=" + utils.URLEncode(input)
	body, err := utils.Fetch(url)
	if err != nil {
		logger.Error("[SearchVideosWithPagination] fetch failed", "error", err)
		return nil, err
	}

	jsonData, err := extractInitialData(body)
	if err != nil {
		logger.Error("[SearchVideosWithPagination] failed to extract ytInitialData", "error", err)
		return nil, err
	}

	results, continuation, err := parseSearchResultsWithContinuation(jsonData)
	if err != nil {
		logger.Error("[SearchVideosWithPagination] failed to parse results", "error", err)
		return nil, fmt.Errorf("extract error: %w", err)
	}

	return &SearchResponse{
		Results:           results,
		ContinuationToken: continuation,
		HasMore:           continuation != "",
	}, nil
}

func extractInitialData(html string) ([]byte, error) {
	match := initialDataRegex.FindStringSubmatch(html)
	if len(match) < 2 {
		return nil, errors.New("ytInitialData not found")
	}

	return []byte(match[1]), nil
}

func parseSearchResults(jsonData []byte) ([]models.SearchResult, error) {
	results, _, err := parseSearchResultsWithContinuation(jsonData)
	return results, err
}

func parseSearchResultsWithContinuation(jsonData []byte) ([]models.SearchResult, string, error) {
	var root map[string]any
	if err := json.Unmarshal(jsonData, &root); err != nil {
		return nil, "", err
	}

	contents := utils.DeepGet(
		root,
		"contents",
		"twoColumnSearchResultsRenderer",
		"primaryContents",
		"sectionListRenderer",
		"contents",
	)

	if contents == nil {
		return nil, "", errors.New("search results missing")
	}

	results, err := extractVideoRenderers(contents)
	if err != nil {
		return nil, "", err
	}

	continuation := ""
	if arr, ok := contents.([]any); ok {
		for _, block := range arr {
			if blockMap, ok := block.(map[string]any); ok {
				if contItem := blockMap["continuationItemRenderer"]; contItem != nil {
					if contItemMap, ok := contItem.(map[string]any); ok {
						if token := utils.DeepGet(contItemMap, "continuationEndpoint", "continuationCommand", "token"); token != nil {
							if tokenStr, ok := token.(string); ok {
								continuation = tokenStr
							}
						}
					}
				}
			}
		}
	}

	return results, continuation, nil
}

func extractVideoRenderers(contents any) ([]models.SearchResult, error) {
	arr, ok := contents.([]any)
	if !ok {
		return nil, errors.New("invalid result container")
	}

	var results []models.SearchResult

	for _, block := range arr {
		blockMap, ok := block.(map[string]any)
		if !ok {
			continue
		}

		if blockMap["continuationItemRenderer"] != nil {
			continue
		}

		if vr := utils.DeepGet(blockMap, "videoRenderer"); vr != nil {
			if parsed := parseVideoRenderer(vr.(map[string]any)); parsed != nil {
				results = append(results, *parsed)
			}
		}

		if sr := utils.DeepGet(blockMap, "reelItemRenderer"); sr != nil {
			if parsed := parseShortRenderer(sr.(map[string]any)); parsed != nil {
				results = append(results, *parsed)
			}
		}

		if rich := utils.DeepGet(blockMap, "richItemRenderer", "content", "videoRenderer"); rich != nil {
			if parsed := parseVideoRenderer(rich.(map[string]any)); parsed != nil {
				results = append(results, *parsed)
			}
		}

		if shelf := blockMap["shelfRenderer"]; shelf != nil {
			results = append(results, extractFromShelf(shelf.(map[string]any))...)
		}

		if itemSection := utils.DeepGet(blockMap, "itemSectionRenderer", "contents"); itemSection != nil {
			if sectionResults, err := extractVideoRenderers(itemSection); err == nil {
				results = append(results, sectionResults...)
			}
		}
	}

	return results, nil
}

func extractFromShelf(shelf map[string]any) []models.SearchResult {
	var results []models.SearchResult

	contents := utils.DeepGet(shelf, "content", "verticalListRenderer", "items")
	if contents == nil {
		return results
	}

	items, ok := contents.([]any)
	if !ok {
		return results
	}

	for _, it := range items {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}

		if vr := utils.DeepGet(m, "videoRenderer"); vr != nil {
			if parsed := parseVideoRenderer(vr.(map[string]any)); parsed != nil {
				results = append(results, *parsed)
			}
		}
	}

	return results
}

func parseVideoRenderer(m map[string]any) *models.SearchResult {
	id := utils.Str(m["videoId"])
	if id == "" {
		return nil
	}

	title := utils.GetText(m, "title", "runs", "text")
	if title == "" {
		title = utils.GetText(m, "title", "simpleText")
	}

	thumb := utils.GetThumbnail(m)
	channel := utils.GetText(m, "ownerText", "runs", "text")
	if channel == "" {
		channel = utils.GetText(m, "ownerText", "simpleText")
	}

	dur := utils.GetText(m, "lengthText", "simpleText")
	if dur == "" {
		dur = utils.GetText(m, "lengthText", "runs", "text")
	}

	channelID := ""
	channelURL := ""
	if ownerBrowseEndpoint := utils.DeepGet(m, "ownerText", "runs", "navigationEndpoint", "browseEndpoint"); ownerBrowseEndpoint != nil {
		if browseID, ok := ownerBrowseEndpoint.(map[string]any)["browseId"].(string); ok {
			channelID = browseID
			channelURL = "https://www.youtube.com/channel/" + browseID
		}
	}

	// Check if live stream
	isLive := utils.DeepGet(m, "badges", "liveBadgeRenderer") != nil
	if !isLive {
		// Alternative live indicator
		isLive = utils.DeepGet(m, "thumbnailOverlays", "thumbnailOverlayTimeStatusRenderer", "style") == "LIVE"
	}

	return &models.SearchResult{
		ID:          id,
		Title:       title,
		URL:         "https://www.youtube.com/watch?v=" + id,
		Thumbnail:   thumb,
		ChannelName: channel,
		ChannelID:   channelID,
		ChannelURl:  channelURL,
		Duration:    dur,
		DurationSec: utils.ParseDuration(dur),
		IsLive:      isLive,
		IsShort:     false,
	}
}

func parseShortRenderer(m map[string]any) *models.SearchResult {
	id := utils.Str(m["videoId"])
	if id == "" {
		return nil
	}

	title := utils.GetText(m, "headline", "simpleText")
	if title == "" {
		title = utils.GetText(m, "headline", "runs", "text")
	}

	thumb := utils.GetThumbnail(m)

	// Extract channel info for shorts
	channelName := utils.GetText(m, "shortBylineText", "simpleText")
	if channelName == "" {
		channelName = utils.GetText(m, "shortBylineText", "runs", "text")
	}

	channelID := ""
	channelURL := ""
	if ownerBrowseEndpoint := utils.DeepGet(m, "shortBylineText", "runs", "navigationEndpoint", "browseEndpoint"); ownerBrowseEndpoint != nil {
		if browseID, ok := ownerBrowseEndpoint.(map[string]any)["browseId"].(string); ok {
			channelID = browseID
			channelURL = "https://www.youtube.com/channel/" + browseID
		}
	}

	return &models.SearchResult{
		ID:          id,
		Title:       title,
		URL:         "https://www.youtube.com/shorts/" + id,
		Thumbnail:   thumb,
		ChannelName: channelName,
		ChannelID:   channelID,
		ChannelURl:  channelURL,
		Duration:    "SHORT",
		DurationSec: 0,
		IsLive:      false,
		IsShort:     true,
	}
}
