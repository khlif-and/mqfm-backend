package livestream

import "encoding/xml"

type AtomFeed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	VideoID   string `xml:"videoId"`
	ChannelID string `xml:"channelId"`
	Title     string `xml:"title"`
	Link      struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
}