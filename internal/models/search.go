package models

type SearchResult struct {
	ID          string
	Title       string
	URL         string
	Thumbnail   string
	Duration    string // "1-:22"
	DurationSec int
	ChannelName string
	ChannelID   string
	ChannelURl  string
	IsLive      bool
	IsShort     bool
}
