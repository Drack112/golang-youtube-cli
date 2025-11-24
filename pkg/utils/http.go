package utils

import (
	"io"
	"net/http"
	"net/url"
)

func Fetch(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 OPR/123.0.0.0 (Edition Yx 08)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(body), nil
}

func URLEncode(s string) string {
	return url.QueryEscape(s)
}
