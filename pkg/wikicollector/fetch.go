package wikicollector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"wikibot/pkg/config"

	"cgt.name/pkg/go-mwclient"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// WikiClient is a WikiClientConfiguration derivated type
type WikiClient config.WikiClientConfig

// NewWikiClient returns a new WikiClient
func NewWikiClient(wikiClientCfg config.WikiClientConfig) WikiClient {
	return WikiClient(wikiClientCfg)
}

func (wc *WikiClient) fetchrc() ([]RecentAction, error) {
	// Initialize a *Client with New(), specifying the wiki's API URL
	// and your HTTP User-Agent. Try to use a meaningful User-Agent.
	w, err := mwclient.New(wc.URL+"/api.php", wc.UserAgent)
	if err != nil {
		return nil, err
	}

	// Log in.
	err = w.Login(wc.Username, wc.Password)
	if err != nil {
		return nil, err
	}

	// Specify parameters to send.
	parameters := map[string]string{
		"format":  "json",
		"rcdir":   "newer",
		"rcprop":  "title|ids|sizes|flags|user|loginfo|redirect|tags|timestamp|comment",
		"list":    "recentchanges",
		"action":  "query",
		"rclimit": "100",
		"rcstart": fmt.Sprint(time.Now().Add(-wc.RefreshDelay).Unix()),
		"rcend":   fmt.Sprint(time.Now().Unix()),
	}

	// Make the request.
	resp, err := w.Get(parameters)
	if err != nil {
		return nil, err
	}

	// Print the *jason.Object
	queryResp := resp.Map()["query"]
	rcs, err := queryResp.Marshal()
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	json.Unmarshal(rcs, &m)
	rcsRev := m["recentchanges"].([]interface{})
	recentActions := make([]RecentAction, 0, len(rcsRev))
	for _, rc := range rcsRev {
		result := RecentAction{
			WikiURL: wc.URL,
		}
		_ = mapstructure.Decode(rc.(map[string]interface{}), &result)
		recentActions = append(recentActions, result)
	}

	return recentActions, nil
}

// RecentAction is the structure holding all the data of a recent action of mediawiki
type RecentAction struct {
	WikiURL   string
	Bot       bool     `mapstructure:"bot"`
	Comment   string   `mapstructure:"comment"`
	Minor     bool     `mapstructure:"minor"`
	New       bool     `mapstructure:"new"`
	NewLen    int      `mapstructure:"newlen"`
	Ns        int      `mapstructure:"ns"`
	OldRevID  int      `mapstructure:"old_revid"`
	OldLen    int      `mapstructure:"oldlen"`
	PageID    int      `mapstructure:"pageid"`
	RcID      int      `mapstructure:"rcid"`
	Redirect  bool     `mapstructure:"redirect"`
	RevID     int      `mapstructure:"revid"`
	Tags      []string `mapstructure:"tags"`
	Timestamp string   `mapstructure:"timestamp"`
	Title     string   `mapstructure:"title"`
	User      string   `mapstructure:"user"`
	Type      string   `mapstructure:"type"`
}

// StartFetchingActions fetches recent action from a wiki
func (wc *WikiClient) StartFetchingActions(ctx context.Context, mode string, raCh chan []RecentAction) {
	results, err := wc.fetchrc()
	if err != nil {
		logrus.Fatal(err)
	}
	if len(results) > 0 {
		raCh <- results
	}

	if mode == "job" {
		return
	}

	ticker := time.NewTicker(wc.RefreshDelay)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			results, err = wc.fetchrc()
			if len(results) > 0 {
				raCh <- results
			}
		}
	}
}
