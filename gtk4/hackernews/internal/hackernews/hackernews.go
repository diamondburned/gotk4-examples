// Package hackernews provides an API wrapper around the HackerNews API.
// Some of these codes are taken from https://github.com/peterhellberg/hn.
// The API reference is in https://github.com/HackerNews/API.
package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"
)

// ItemID is the ID if a post (item).
type ItemID int

// Client is a HackerNews client.
type Client struct {
	http.Client
}

// DefaultClient is the default HackerNews client.
var DefaultClient = &Client{
	Client: *http.DefaultClient,
}

// BaseEndpoint is the base endpoint used by the Client.
const BaseEndpoint = "https://hacker-news.firebaseio.com/v0"

// Get gets the given path (concatenated after BaseEndpoint) and unmarshals the
// body into jsonVal if it's not nil.
func (c *Client) Get(ctx context.Context, path string, jsonVal interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	req, err := http.NewRequestWithContext(ctx, "GET", BaseEndpoint+path, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status %s given", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(jsonVal); err != nil {
		return fmt.Errorf("HackerNews returned invalid JSON: %w", err)
	}

	return nil
}

// StoryFile enumerates the possible types of stories to be fetched.
type StoryFile string

const (
	NewStories  StoryFile = "newstories.json"
	TopStories  StoryFile = "topstories.json"
	BestStories StoryFile = "beststories.json"
	AskStories  StoryFile = "askstories.json"
	ShowStories StoryFile = "showstories.json"
	JobStories  StoryFile = "jobstories.json"
)

// Name returns the human-friendly name of the story file instead of the JSON
// path name.
func (f StoryFile) Name() string {
	switch f {
	case NewStories:
		return "New Stories"
	case TopStories:
		return "Top Stories"
	case BestStories:
		return "Best Stories"
	case AskStories:
		return "Ask Stories"
	case ShowStories:
		return "Show Stories"
	case JobStories:
		return "Job Stories"
	default:
		return string(f)
	}
}

// Stories fetches a list of stories, returning their item IDs.
func (c *Client) Stories(ctx context.Context, file StoryFile) ([]ItemID, error) {
	var items []ItemID
	if err := c.Get(ctx, string(file), &items); err != nil {
		return nil, err
	}
	return items, nil
}

// Item represents an item.
type Item struct {
	ID      ItemID `json:"id"`
	Deleted bool   `json:"deleted"`
	// Type is the item's type.
	Type ItemType `json:"type"`
	// By is the author's username.
	By string `json:"by"`
	// Time is the creation time.
	Time UnixTime `json:"time"`
	// Text is the comment, story or poll text in HTML.
	Text string `json:"text"`
	// Dead is true if the item is dead.
	Dead bool `json:"dead"`
	// Parent is the comment's parent, which is another comment or relevant
	// story.
	Parent ItemID `json:"parent"`
	// Kids is the IDs of the item's comments in ranked display order.
	Kids []ItemID `json:"kids"`
	// URL is the URL of the story.
	URL string `json:"url"`
	// Score is the story's score or the votes for a pollopt.
	Score int `json:"score"`
	// Title is the title of the story in HTML.
	Title string `json:"title"`
	// Parts is the list of relateed pollopts, in display order.
	Parts []ItemID `json:"parts"`
	// Descendants is the total comment count in case of stories or polls.
	Descendants int `json:"descendants"`
	// Skipped: pollopt.poll
}

// ItemType describes the type of an item.
type ItemType string

const (
	JobItem     ItemType = "job"
	StoryItem   ItemType = "story"
	CommentItem ItemType = "comment"
	PollItem    ItemType = "poll"
	PollOptItem ItemType = "pollopt"
)

// UnixTime is a timestamp number in Unix epoch time.
type UnixTime int64

// Time returns the Unix time as a time.Time.
func (t UnixTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// Item fetches a single item by its ID.
func (c *Client) Item(ctx context.Context, id ItemID) (*Item, error) {
	var item Item
	if err := c.Get(ctx, path.Join("item", fmt.Sprintf("%d.json", id)), &item); err != nil {
		return nil, err
	}
	return &item, nil
}
