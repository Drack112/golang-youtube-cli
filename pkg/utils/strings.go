package utils

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func IsYouTubeURL(text string) bool {
	u, err := url.Parse(text)
	if err != nil || u.Scheme == "" {
		return false
	}

	host := strings.ToLower(u.Host)
	return strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be")
}

func CleanYoutubeLink(link string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9./?=&%_-]")
	cleanedLink := reg.ReplaceAllString(link, "")
	return cleanedLink
}

func ParseDuration(dur string) int {
	if dur == "" || dur == "LIVE" || dur == "SHORT" {
		return 0
	}

	if strings.Contains(strings.ToUpper(dur), "LIVE") {
		return 0
	}

	parts := strings.Split(dur, ":")
	var hours, minutes, seconds int

	switch len(parts) {
	case 3:
		hours, _ = strconv.Atoi(parts[0])
		minutes, _ = strconv.Atoi(parts[1])
		seconds, _ = strconv.Atoi(parts[2])
	case 2:
		minutes, _ = strconv.Atoi(parts[0])
		seconds, _ = strconv.Atoi(parts[1])
	case 1:
		seconds, _ = strconv.Atoi(parts[0])
	default:
		return 0
	}
	return hours*3600 + minutes*60 + seconds
}

func GetThumbnail(m map[string]any) string {
	if thumbnails := DeepGet(m, "thumbnail", "thumbnails"); thumbnails != nil {
		if thumbs, ok := thumbnails.([]any); ok && len(thumbs) > 0 {
			lastThumb := thumbs[len(thumbs)-1]
			if thumbMap, ok := lastThumb.(map[string]any); ok {
				if url, ok := thumbMap["url"].(string); ok {
					if strings.HasPrefix(url, "/") {
						return "https:" + url
					}
					if strings.HasPrefix(url, "/") {
						return "https://www.youtube.com" + url
					}
					return url
				}
			}
		}
	}

	if videoID := Str(m["videoID"]); videoID != "" {
		return "https://i.ytimg.com/vi/" + videoID + "/hqdefault.jpg"
	}

	return ""
}

func ExtractVideoID(link string) string {
	patterns := []struct {
		pattern string
		index   int
	}{
		{`(?:v=|v/|embed/|youtu\.be/|/v/|/embed/)([^&?#/]+)`, 1},
		{`^([a-zA-Z0-9_-]{11})$`, 1}, // Direct video ID
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p.pattern)
		matches := re.FindStringSubmatch(link)
		if len(matches) > p.index {
			return matches[p.index]
		}
	}
	return ""
}
