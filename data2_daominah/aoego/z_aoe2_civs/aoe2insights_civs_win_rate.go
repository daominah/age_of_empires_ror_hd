package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoego"
)

// main prints out AoE2 civilizations from AoE2Insights website data in the last 180 days,
// sorted by win rate, if you want to sort by popularity, change the const at the beginning
func main() {
	const isForceReDownload = false
	const isSortByPopularity = false // default sort by win rate

	//const ratingRange RatingRange = Rating1000To1200
	const ratingRange RatingRange = Rating1200To1900
	//const ratingRange RatingRange = Rating1900Up

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// check whether data is downloaded
	projectRootDir, err := aoego.GetProjectRootGit()
	if err != nil {
		log.Fatalf("error GetProjectRootGit: %v\n", err)
	}
	outputFilePath := filepath.Join(projectRootDir, "data2_daominah", "aoego", "z_aoe2_civs",
		"aoe2insights_civs_win_rate.json")
	data, err := os.ReadFile(outputFilePath)
	if err == nil && len(data) > 0 && !isForceReDownload {
		log.Printf("re-use existing aoe2insights data")
	} else {
		log.Printf("downloading aoe2insights data...")
		data, err = downloadAoe2insightsData(ratingRange)
		if err != nil {
			log.Fatalf("error downloadAoe2insightsData: %v\n", err)
			return
		}

		var tmpObj map[string]any
		err = json.Unmarshal(data, &tmpObj)
		if err != nil {
			log.Fatalf("error json.Unmarshal: %v\n", err)
		}
		beautyData, err := json.MarshalIndent(tmpObj, "", "\t")
		if err != nil {
			log.Fatalf("error json.MarshalIndent: %v\n", err)
		}
		err = os.WriteFile(outputFilePath, beautyData, 0644)
		if err != nil {
			log.Fatalf("error os.WriteFile: %v\n", err)
		}
		log.Printf("dowloaded and wrote data to %v\n", outputFilePath)
	}

	var aoe2insightsResp Aoe2insightsResponse
	err = json.Unmarshal(data, &aoe2insightsResp)
	if err != nil {
		log.Fatalf("error json.Unmarshal to Aoe2insightsResponse: %v\n", err)
	}

	civs := ToCivStats(aoe2insightsResp)
	if isSortByPopularity {
		sort.Sort(SortByPopularity(civs))
		// in 2025-10, tops are:
		//1 . Mongols          WinRate:  50.77%  Popularity:   6.01% (19499 matches)
		//2 . Persians         WinRate:  51.01%  Popularity:   3.20% (10394 matches)
		//3 . Khitans          WinRate:  56.15%  Popularity:   3.13% (10147 matches)
		//4 . Magyars          WinRate:  47.72%  Popularity:   3.06% (9939 matches)
		//5 . Khmer            WinRate:  48.30%  Popularity:   2.50% (8110 matches)
		//6 . Japanese         WinRate:  50.22%  Popularity:   2.50% (8100 matches)
		//7 . Huns             WinRate:  52.87%  Popularity:   2.44% (7923 matches)
		//8 . Lithuanians      WinRate:  48.97%  Popularity:   2.42% (7840 matches)
		//9 . Spanish          WinRate:  50.77%  Popularity:   2.35% (7640 matches)
		//10. Hindustanis      WinRate:  51.64%  Popularity:   2.33% (7562 matches)
	} else {
		sort.Sort(SortByWinrate(civs))
		// in 2025-10, tops are:
		//1 . Khitans          WinRate:  56.15%  Popularity:   3.13% (10147 matches)
		//	Pasture is better Farm, good army, consider to try their Heavy Cavalry Archer.
		//2 . Romans           WinRate:  54.92%  Popularity:   1.49% (4831 matches)
		//	Villagers works 5% faster, good Infantry.
		//3 . Wu               WinRate:  54.41%  Popularity:   1.35% (4374 matches)
		//	military buildings +55 food, good anti-archer Infantry.
		//4 . Shu              WinRate:  54.03%  Popularity:   1.22% (3968 matches)
		//	Lumberjacks gives 9% food too, good Archer+Siege rush.
		//5 . Malay            WinRate:  52.90%  Popularity:   1.73% (5631 matches)
		//	Age up only 60% time, probably good Elephant.
		//6 . Huns             WinRate:  52.87%  Popularity:   2.44% (7923 matches)
		//7 . Vikings          WinRate:  52.52%  Popularity:   1.68% (5438 matches)
		//8 . Malians          WinRate:  52.36%  Popularity:   2.09% (6784 matches)
		//9 . Hindustanis      WinRate:  51.64%  Popularity:   2.33% (7562 matches)
		//10. Mayans           WinRate:  51.53%  Popularity:   2.04% (6614 matches)
	}
	for i, civ := range civs {
		fmt.Printf("%-2v. %-15v  WinRate: %6.2f%%  Popularity: %6.2f%% (%v matches)\n",
			i+1, civ.Civilization, civ.WinRate*100, civ.Popularity*100, civ.NMatches)
	}
}

func downloadAoe2insightsData(ratingRange RatingRange) ([]byte, error) {
	baseURL := "https://www.aoe2insights.com/stats/api/match-ups/"
	params := url.Values{}
	params.Add("filter", `{"field":"ladder","operator":"equals","value":"3"}`)
	params.Add("filter", `{"field":"started_time","operator":"greater_than_or_equals_relative","value":[-180,"days"]}`)
	params.Add("aggregation", "agg_win_rate")
	params.Add("aggregation", "agg_n_matches")
	params.Add("by", "civilization")
	params.Add("order", "-agg_win_rate")
	params.Add("start", "0")
	params.Add("limit", "100")
	switch ratingRange {
	case Rating1000To1200:
		// Elo group from 1000 to 1200, developing players,
		// percentiles: 1000 (~50%), 1100 (~64%)
		params.Add("filter", `{"field":"elo_group","operator":"equals","value":"1000"}`)
		params.Add("filter", `{"field":"elo_group_opponent","operator":"equals","value":"1000"}`)
	case Rating1900Up:
		// Elo group from 1900 and up, top-tier players,
		// percentiles:
		// 1900 (~98%), 2000 (~99%)
		// 2400 (top 100 players, ~99.9%)
		params.Add("filter", `{"field":"elo_group","operator":"equals","value":"1900"}`)
		params.Add("filter", `{"field":"elo_group_opponent","operator":"equals","value":"1900"}`)
	default:
		// Elo group from 1200 to 2000, mid-tier players,
		// percentiles:
		// 1200 (~75%), 1300 (~81%), 1400 (~87%), 1500 (~91%),
		// 1600 (~94%), 1900 (~98%), 2000 (~99%)
		params.Add("filter", `{"field":"elo_group","operator":"equals","value":"1200"}`)
		params.Add("filter", `{"field":"elo_group_opponent","operator":"equals","value":"1200"}`)
	}

	fullURL := baseURL + "?" + params.Encode()
	log.Printf("data URL: %v", fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("http.Get error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll error: %w", err)
	}
	return body, nil
}

type Aoe2insightsResponse struct {
	Count int `json:"count"`
	Start int `json:"start"`
	Limit int `json:"limit"`
	Rows  struct {
		AggWinRate   []float64 `json:"agg_win_rate"`
		AggNMatches  []int     `json:"agg_n_matches"`
		Civilization []string  `json:"civilization"`
	} `json:"rows"`
}

// RatingRange are pre-defined by aoe2insights website (cannot customize)
type RatingRange string

const (
	Rating1000To1200 RatingRange = "1000"
	Rating1200To1900 RatingRange = "1200"
	Rating1900Up     RatingRange = "1900"
)

type CivStat struct {
	Civilization string
	WinRate      float64
	NMatches     int
	Popularity   float64
}

type SortByPopularity []CivStat
type SortByWinrate []CivStat

func (a SortByPopularity) Len() int           { return len(a) }
func (a SortByPopularity) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByPopularity) Less(i, j int) bool { return a[i].NMatches > a[j].NMatches } // Descending

func (a SortByWinrate) Len() int           { return len(a) }
func (a SortByWinrate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByWinrate) Less(i, j int) bool { return a[i].WinRate > a[j].WinRate } // Descending

func ToCivStats(resp Aoe2insightsResponse) []CivStat {
	stats := make([]CivStat, 0, len(resp.Rows.Civilization))
	totalMatches := 0
	for _, n := range resp.Rows.AggNMatches {
		totalMatches += n
	}
	if totalMatches == 0 {
		return nil
	}
	for i, civ := range resp.Rows.Civilization {
		if strings.ToLower(civ) == "unknown" {
			continue
		}
		stats = append(stats, CivStat{
			Civilization: civ,
			WinRate:      resp.Rows.AggWinRate[i],
			NMatches:     resp.Rows.AggNMatches[i],
			Popularity:   float64(resp.Rows.AggNMatches[i]) / float64(totalMatches),
		})
	}
	return stats
}
