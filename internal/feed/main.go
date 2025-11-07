package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not make request")
	}
	req.Header.Set("User-Agent", "gator")

	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not read from '%v'", feedURL)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code for %s: %d %s", feedURL, resp.StatusCode, resp.Status)
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil for %s (Status: 200). This indicates a highly unusual network error or invalid server response structure.", feedURL)
	}
	body, err := io.ReadAll(resp.Body)
	if len(body) == 0 {
		return nil, fmt.Errorf("fetched empty response body for %s (Status: %d). The URL likely returned a blank page or the network connection closed prematurely", feedURL, resp.StatusCode)
	}
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("could not marchal xml: %w", err)
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Items {
		item := &feed.Channel.Items[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &feed, nil
}
