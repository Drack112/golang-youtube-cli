package utils

import (
	"net/url"
	"strings"
)

func IsYoutubeURL(text string) bool {
	u, err := url.Parse(text)
	if err != nil || u.Scheme == "" {
		return false
	}

	host := strings.ToLower(u.Host)
	return strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be")
}
