package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/components"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoego"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func main() {
	const isForceReDownload = false

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// check whether data is downloaded
	projectRootDir, err := aoego.GetProjectRootGit()
	if err != nil {
		log.Fatalf("error GetProjectRootGit: %v\n", err)
	}
	today := time.Now().Format("2006-01-02")
	todayOutputDir := filepath.Join(projectRootDir, "data2_daominah", "aoego",
		"z_aoe2_rating_percentile", "data", today)
	err = os.MkdirAll(todayOutputDir, 0755)
	if err != nil {
		log.Fatalf("error os.MkdirAll: %v\n", err)
	}

	// read the directory, download leaderboard if needed
	files, err := os.ReadDir(todayOutputDir)
	if err != nil {
		log.Fatalf("error os.ReadDir: %v\n", err)
	}
	if len(files) > 0 && !isForceReDownload {
		log.Printf("re-use existing ageofempires.com data, already have %d files\n", len(files))
	} else {
		log.Printf("downloading ageofempires.com data...")
		nDownloadedPages, err := DownloadAgeofempirescomData(todayOutputDir)
		if err != nil {
			log.Fatalf("error downloadAoe2insightsData: %v\n", err)
			return
		}
		log.Printf("downloaded %d pages of ageofempires.com data\n", nDownloadedPages)
	}

	// read all files in the directory and aggregate players
	players := make(map[int]AoEPlayer) // map key is "rlUserId"
	files, err = os.ReadDir(todayOutputDir)
	if err != nil {
		log.Fatalf("error os.ReadDir before aggregate players: %v\n", err)
	}
	for _, file := range files {
		filePath := filepath.Join(todayOutputDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error os.ReadFile %v: %v\n", filePath, err)
		}
		var pageData AgeofempirescomDataResponse
		err = json.Unmarshal(data, &pageData)
		if err != nil {
			log.Fatalf("error json.Unmarshal %v: %v\n", filePath, err)
		}
		for _, player := range pageData.Items {
			players[player.RlUserId] = player
		}
	}

	// group players to Elo buckets, e.g. 1000-1099, 1100-1199, ...
	const bucketSize float64 = 100
	bucketFunc := func(elo float64) string {
		roundDown := int(elo) / int(bucketSize) * int(bucketSize)
		roundUp := roundDown + int(bucketSize)
		return fmt.Sprintf("%04d→%04d", roundDown, roundUp)
	}
	ratingRanges := make(map[string][]AoEPlayer)
	for _, player := range players {
		bucket := bucketFunc(player.Elo)
		ratingRanges[bucket] = append(ratingRanges[bucket], player)
	}
	// sort the buckets by key (same as rating from low to high)
	sortedBucketKeys := make([]string, 0, len(ratingRanges))
	for key := range ratingRanges {
		sortedBucketKeys = append(sortedBucketKeys, key)
	}
	sort.Strings(sortedBucketKeys)
	// summarize then output as a CSV
	totalPlayers := len(players)
	if totalPlayers == 0 {
		log.Fatalf("no players found in the data\n")
	}
	cumulativeToCurrentBucket := 0
	var dataAsCSV [][]string
	dataAsCSV = append(dataAsCSV, []string{"RatingRange", "RatingLow", "RatingHigh", "CountPlayers", "Percentile", "RankLow", "RankHigh"})
	var chartBars []RatingBucket
	for _, bucketKey := range sortedBucketKeys {
		playersInBucket := ratingRanges[bucketKey]
		cumulativeToCurrentBucket += len(playersInBucket)
		ratingBucket := RatingBucket{
			RatingRange:     bucketKey,
			RatingBoundLow:  0, // will be set later based on RatingRange
			RatingBoundHigh: 0, // will be set later based on RatingRange
			CountPlayers:    len(playersInBucket),
			Percentile:      float64(cumulativeToCurrentBucket) / float64(totalPlayers) * 100.0,
			RankBoundLow:    totalPlayers - (cumulativeToCurrentBucket - len(playersInBucket)),
			RankBoundHigh:   totalPlayers - cumulativeToCurrentBucket + 1,
		}
		_ = ratingBucket.setRatingBound()
		dataAsCSV = append(dataAsCSV, []string{
			ratingBucket.RatingRange,
			fmt.Sprintf("%v", ratingBucket.RatingBoundLow),
			fmt.Sprintf("%v", ratingBucket.RatingBoundHigh),
			fmt.Sprintf("%v", ratingBucket.CountPlayers),
			fmt.Sprintf("%.3f", ratingBucket.Percentile),
			fmt.Sprintf("#%v", ratingBucket.RankBoundLow),
			fmt.Sprintf("#%v", ratingBucket.RankBoundHigh),
		})
		chartBars = append(chartBars, ratingBucket)
	}
	// write output CSV to a file
	outputCSVFileName := fmt.Sprintf("aoe2_rating_percentile_date_%v.csv", time.Now().Format("2006_01_02"))
	outputCSVFilePath := filepath.Join(projectRootDir, "data2_daominah", "aoego",
		"z_aoe2_rating_percentile", "data_processed", outputCSVFileName)
	outputCSVFile, err := os.Create(outputCSVFilePath)
	if err != nil {
		log.Fatalf("error os.Create outputCSVFilePath %v: %v\n", outputCSVFilePath, err)
	}
	csvWriter := csv.NewWriter(outputCSVFile)
	err = csvWriter.WriteAll(dataAsCSV)
	if err != nil {
		log.Fatalf("error csvWriter.WriteAll to %v: %v\n", outputCSVFilePath, err)
	}
	csvWriter.Flush()
	err = outputCSVFile.Close()
	if err != nil {
		log.Fatalf("error outputCSVFile.Close %v: %v\n", outputCSVFilePath, err)
	}
	log.Printf("wrote rating percentile to %v\n", outputCSVFileName)
	//log.Printf("AoE2 rating percentile:")
	//for _, row := range dataAsCSV {
	//	log.Printf("%12v %12v %12v %16v", row[0], row[1], row[2], row[3])
	//}

	fmt.Printf("______________________________________________________\n")
	group0To800 := 0
	group0To1000 := 0
	group0To1200 := 0
	group0To1500 := 0
	group0To2000 := 0
	group0To2500 := 0
	for _, player := range players {
		if player.Elo < 800 {
			group0To800++
		}
		if player.Elo < 1000 {
			group0To1000++
		}
		if player.Elo < 1200 {
			group0To1200++
		}
		if player.Elo < 1500 {
			group0To1500++
		}
		if player.Elo < 2000 {
			group0To2000++
		}
		if player.Elo < 2500 {
			group0To2500++
		}
	}
	fmt.Printf("* rating    0 →  800: percentile %.1f%% (rank #%5d → #%5d)\n",
		float64(group0To800)/float64(totalPlayers)*100.0,
		totalPlayers, totalPlayers-group0To800+1)
	fmt.Printf("* rating  800 → 1000: percentile %.1f%% (rank #%5d → #%5d)\n",
		float64(group0To1000)/float64(totalPlayers)*100.0,
		totalPlayers-group0To800, totalPlayers-group0To1000+1)
	fmt.Printf("* rating 1000 → 1200: percentile %.1f%% (rank #%5d → #%5d)\n",
		float64(group0To1200)/float64(totalPlayers)*100.0,
		totalPlayers-group0To1000, totalPlayers-group0To1200+1)
	fmt.Printf("* rating 1200 → 1500: percentile %.1f%% (rank #%5d → #%5d)\n",
		float64(group0To1500)/float64(totalPlayers)*100.0,
		totalPlayers-group0To1200, totalPlayers-group0To1500+1)
	fmt.Printf("* rating 1500 → 2000: percentile %.1f%% (rank #%5d → #%5d)\n",
		float64(group0To2000)/float64(totalPlayers)*100.0,
		totalPlayers-group0To1500, totalPlayers-group0To2000+1)
	fmt.Printf("* rating 2000 → 2500: percentile %.1f%% (rank #%5d → #%5d)\n",
		float64(group0To2500)/float64(totalPlayers)*100.0,
		totalPlayers-group0To2000, totalPlayers-group0To2500+1)
	fmt.Printf("* rating 2500 and up: percentile 100%%  (rank #%5d → #    1)\n", totalPlayers-group0To2500)
	// Output:
	//	* rating    0 →  800: percentile 26.1% (rank #43139 → #31878)
	//	* rating  800 → 1000: percentile 49.2% (rank #31877 → #21906)
	//	* rating 1000 → 1200: percentile 73.6% (rank #21905 → #11376)
	//	* rating 1200 → 1500: percentile 90.7% (rank #11375 → # 4017)
	//	* rating 1500 → 2000: percentile 98.8% (rank # 4016 → #  499)
	//	* rating 2000 → 2500: percentile 99.9% (rank #  498 → #   59)
	//	* rating 2500 and up: percentile 100%  (rank #   58 → #    1)
	fmt.Printf("______________________________________________________\n")

	// now sort players to find exact rating for markers 25%, 50%, 75%, 90%, 99%, 99.9%
	sortedPlayers := make([]AoEPlayer, 0, len(players))
	for _, player := range players {
		sortedPlayers = append(sortedPlayers, player)
	}
	sort.Slice(sortedPlayers, func(i, j int) bool {
		// highest rating player comes first
		return sortedPlayers[i].Elo > sortedPlayers[j].Elo
	})
	getRatingAtPercentile := func(percentile float64) (float64, int) {
		rankIndexFloat := math.Floor(float64(totalPlayers) * (1 - percentile))
		rankIndex := int(rankIndexFloat)
		if rankIndex < 0 {
			rankIndex = 0
		}
		if rankIndex >= totalPlayers {
			rankIndex = totalPlayers - 1
		}
		return sortedPlayers[rankIndex].Elo, sortedPlayers[rankIndex].Rank
	}
	percentileMarkers := make([]PercentileMarker, 0)
	for _, percentile := range []float64{0.25, 0.5, 0.75, 0.9, 0.99, 0.999} {
		ratingAtPercentile, rankPosition := getRatingAtPercentile(percentile)
		percentileMarker := PercentileMarker{
			Percentile:   percentile,
			Rating:       ratingAtPercentile,
			RankPosition: rankPosition,
		}
		percentileMarkers = append(percentileMarkers, percentileMarker)
		fmt.Printf("rating %4.0f is better than %.1f%% players.\n",
			ratingAtPercentile, percentile*100)
	}
	// Output:
	//	rating  787 is better than 25.0% players.
	//	rating 1005 is better than 50.0% players.
	//	rating 1214 is better than 75.0% players.
	//	rating 1480 is better than 90.0% players.
	//	rating 2027 is better than 99.0% players.
	//	rating 2546 is better than 99.9% players.
	fmt.Printf("______________________________________________________\n")

	outputChartFile := "index.html"
	outputChartFileFullPath := filepath.Join(projectRootDir, "data2_daominah", "aoego",
		"z_aoe2_rating_percentile", outputChartFile)
	chartWidth, chartHeight := 1800, 800
	err = drawPercentilesChart(chartBars, percentileMarkers,
		today, chartWidth, chartHeight, outputChartFileFullPath)
	if err != nil {
		log.Fatalf("error drawPercentilesChart: %v\n", err)
	}
	log.Printf("drew rating percentile chart to %v\n", outputChartFile)
	// Output: an HTML file with scripts, onload will render two charts as <svg>.
	// I want to overlay the two charts in the same place,
	// but go-echarts overlapping func does not work,
	// so I will edit the output HTML here:
	// code here
	// After page.Render(f) and closing the file, read-modify-write the HTML:
	htmlBytes, err := os.ReadFile(outputChartFileFullPath)
	if err != nil {
		log.Fatalf("error reading chart HTML: %v", err)
	}
	htmlStr := string(htmlBytes)

	// Inject a script before </body> to overlay the SVGs
	overlayScript := `
	<script>
	window.onload = function() {
	  // Find all SVGs in the page
	  var svgs = document.querySelectorAll('svg');
	  if (svgs.length >= 2) {
	    // Place the second SVG inside the first SVG's parent, absolutely positioned
	    var svg1 = svgs[0];
	    var svg2 = svgs[1];
	    svg2.style.position = 'absolute';
	    svg2.style.left = svg1.offsetLeft + 'px';
	    svg2.style.top = svg1.offsetTop + 'px';
	    svg2.style.pointerEvents = 'none'; // let mouse events pass through
	    svg1.parentNode.style.position = 'relative';
	    svg1.parentNode.appendChild(svg2);
	  }
	};
	</script>
	`
	htmlStr = strings.Replace(htmlStr, "</body>", overlayScript+"</body>", 1)
	err = os.WriteFile(outputChartFileFullPath, []byte(htmlStr), 0644)
	if err != nil {
		log.Fatalf("error writing modified chart HTML: %v", err)
	}
	log.Printf("injected overlay script into chart HTML %v\n", outputChartFile)
}

func DownloadAgeofempirescomData(todayOutputDir string) (int, error) {
	beginT := time.Now()
	firstPage, err := downloadAgeofempirescomData(1)
	endT := time.Now()
	if err != nil {
		return 0, fmt.Errorf("download first page error: %w", err)
	}
	var firstPageData AgeofempirescomDataResponse
	err = json.Unmarshal(firstPage, &firstPageData)
	if err != nil {
		return 0, fmt.Errorf("json.Unmarshal first page error: %w", err)
	}
	beautyJSON, err := json.MarshalIndent(firstPageData, "", "\t")
	if err != nil {
		return 0, fmt.Errorf("json.MarshalIndent first page error: %w", err)
	}
	outputFile := filepath.Join(todayOutputDir, "ageofempirescom_leaderboard_page_001.json")
	err = os.WriteFile(outputFile, beautyJSON, 0644)
	if err != nil {
		return 0, fmt.Errorf("os.WriteFile first page error: %w", err)
	}

	nTotalPlayers := firstPageData.Count
	nPages := nTotalPlayers/100 + 1
	log.Printf("total players: %d, total pages: %d", nTotalPlayers, nPages)
	nDownloadedPages := 1
	for page := 2; page <= nPages; page++ {
		log.Printf("begin download page %v, estimated time left: %v",
			page, time.Duration(nPages-page)*endT.Sub(beginT))
		pageData, err := downloadAgeofempirescomData(page)
		if err != nil {
			return 0, fmt.Errorf("download page %d error: %w", page, err)
		}
		var pageDataObj AgeofempirescomDataResponse
		err = json.Unmarshal(pageData, &pageDataObj)
		if err != nil {
			return 0, fmt.Errorf("json.Unmarshal page %d error: %w", page, err)
		}
		beautyJSON, err := json.MarshalIndent(pageDataObj, "", "\t")
		if err != nil {
			return 0, fmt.Errorf("json.MarshalIndent page %d error: %w", page, err)
		}
		outputFile := filepath.Join(todayOutputDir,
			fmt.Sprintf("ageofempirescom_leaderboard_page_%03d.json", page))
		err = os.WriteFile(outputFile, beautyJSON, 0644)
		if err != nil {
			return 0, fmt.Errorf("os.WriteFile page %d error: %w", page, err)
		}
		nDownloadedPages++
	}
	return nDownloadedPages, nil
}

func downloadAgeofempirescomData(page int) ([]byte, error) {
	u := `https://api.ageofempires.com/api/v2/ageii/Leaderboard`
	filter := map[string]any{
		"region":           "7",
		"matchType":        "3",
		"consoleMatchType": "15",
		"count":            "100", // maximum limit the API allows
		"sortColumn":       "rank",
		"sortDirection":    "ASC",
		"searchPlayer":     "",
		"page":             fmt.Sprintf("%v", page),
	}
	reqBody, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal req filter: %w", err)
	}
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("http.Post error: %w", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("io.ReadAll error: %w", err)
	}
	_ = resp.Body.Close()
	return respBody, nil
}

// AgeofempirescomDataResponse example Item:
//
//	{
//		"gameId": "age2",
//		"rlUserId": 199325,
//		"userName": "GL.Hera",
//		"elo": 2834,
//		"eloHighest": 3045,
//		"rank": 2,
//		"region": "3",
//		"wins": 4515,
//		"losses": 1640,
//		"winPercent": 73.35,
//		"winStreak": -1,
//		"totalGames": 6155,
//	}
type AgeofempirescomDataResponse struct {
	Id            string // random hex string
	Count         int    // always is total players, regardless of page, 43139 on 2025-10
	LeaderboardId int    // 0
	Region        int    // 0
	LastUpdated   string // "2006-01-02T15:04:05.999999999Z07:00"
	Items         []AoEPlayer
}

type AoEPlayer struct {
	GameId     string // const "age2"
	RlUserId   int    // key map to filter unique players
	UserName   string // e.g. "GL.Hera"
	AvatarUrl  string
	Elo        float64 // main data
	EloHighest float64
	Rank       int // can be calculated if we have all players rating, but they provide it
	Region     string
	Wins       int
	WinPercent float64
	Losses     int
	WinStreak  int
	TotalGames int
}

type RatingBucket struct {
	RatingRange     string // e.g. "0000→0099", "0100→0199", ...
	RatingBoundLow  int
	RatingBoundHigh int
	CountPlayers    int
	Percentile      float64 // from 0 to 100, rounded .999
	RankBoundLow    int     // e.g. #43139
	RankBoundHigh   int     // e.g. #43054
}

type PercentileMarker struct {
	Percentile   float64 // some milestones, e.g. 0.25, 0.5, 0.75, 0.9, 0.99, 0.999
	Rating       float64 // should be calculated from Percentile field
	RankPosition int     // can be ignored
}

func (rr *RatingBucket) setRatingBound() error {
	parts := strings.Split(rr.RatingRange, "→")
	if len(parts) != 2 {
		return fmt.Errorf("invalid RatingRange format: %v", rr.RatingRange)
	}
	lowStr := strings.TrimLeft(parts[0], "0")
	highStr := strings.TrimLeft(parts[1], "0")
	if lowStr == "" {
		lowStr = "0"
	}
	if highStr == "" {
		highStr = "0"
	}
	var err error
	rr.RatingBoundLow, err = strconv.Atoi(lowStr)
	if err != nil {
		return fmt.Errorf("error strconv.Atoi RatingBoundLow %v: %w", lowStr, err)
	}
	rr.RatingBoundHigh, err = strconv.Atoi(highStr)
	if err != nil {
		return fmt.Errorf("error strconv.Atoi RatingBoundHigh %v: %w", highStr, err)
	}
	return nil
}

// drawPercentilesChart draws a chart and save as HTML to outputFilePath,
// the chart x-axis is rating buckets, y-axis is number of players in that bucket.
func drawPercentilesChart(
	bars []RatingBucket,
	percentileMarkers []PercentileMarker,
	dataDate string,
	chartWidth, chartHeight int,
	outputFilePath string) error {
	maxAxisX, maxAxisY := 0, 0
	totalPlayers := 0
	for _, b := range bars {
		if b.RatingBoundHigh > maxAxisX {
			maxAxisX = b.RatingBoundHigh
		}
		roundUpNPlayers := int(math.Ceil(float64(b.CountPlayers)/1000)) * 1000
		if roundUpNPlayers > maxAxisY {
			maxAxisY = roundUpNPlayers
		}
		totalPlayers += b.CountPlayers
	}

	// barChart is the main chart,
	// displays number of players in each rating bucket (0→100, 100→200, ...)
	barChart := charts.NewBar()
	newInitializationOpts := func() opts.Initialization {
		return opts.Initialization{
			Width:     fmt.Sprintf("%vpx", chartWidth),
			Height:    fmt.Sprintf("%vpx", chartHeight),
			PageTitle: "AoE2DE rating distribution",
			Renderer:  "svg",
		}
	}
	newXAxisOptsDiscrete := func() opts.XAxis {
		return opts.XAxis{
			Name:      "rating",
			Type:      "category", // enum defined by go-echarts, for discrete buckets
			AxisLine:  &opts.AxisLine{Show: opts.Bool(true)},
			AxisLabel: &opts.AxisLabel{Interval: strconv.Itoa(0)}, // display all bars labels, regardless of space
			Min:       "0",
			Max:       strconv.Itoa(maxAxisX),
		}
	}
	newXAxisOptsContinuous := func() opts.XAxis {
		return opts.XAxis{
			//Name:      "percentile",
			Type:      "value",
			AxisLine:  &opts.AxisLine{Show: opts.Bool(true)},
			Min:       "0",
			Max:       strconv.Itoa(maxAxisX),
			Position:  "top",
			AxisTick:  &opts.AxisTick{Show: opts.Bool(false)},
			AxisLabel: &opts.AxisLabel{Show: opts.Bool(false)},
			SplitLine: &opts.SplitLine{Show: opts.Bool(false)},
		}
	}
	newYAxisOpts := func(showSplitLine bool) opts.YAxis {
		return opts.YAxis{
			Name:     "players count",
			Type:     "value",
			AxisLine: &opts.AxisLine{Show: opts.Bool(true)},
			Min:      "0", Max: strconv.Itoa(maxAxisY),
			SplitLine: &opts.SplitLine{Show: opts.Bool(showSplitLine)},
		}
	}
	newTitleOpts := func() opts.Title {
		return opts.Title{
			Title: "AoE2DE rating distribution",
			Subtitle: fmt.Sprintf("Data from ageofempires.com leaderboards on %v with a total of %v players",
				dataDate, totalPlayers),
			Left: "center", Top: "0px",
		}
	}
	newGridOpts := func() opts.Grid {
		return opts.Grid{Top: "120px"} // lower the chart to get space for the title
	}
	barChart.SetGlobalOptions(
		charts.WithInitializationOpts(newInitializationOpts()),
		charts.WithXAxisOpts(newXAxisOptsDiscrete()),
		charts.WithYAxisOpts(newYAxisOpts(true)),
		charts.WithTitleOpts(newTitleOpts()),
		charts.WithGridOpts(newGridOpts()),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}), // hide the only 1 legend button
	)

	xLabels := make([]string, 0, len(bars))
	yValues := make([]opts.BarData, 0, len(bars))
	for _, b := range bars {
		// use RatingBoundHigh instead of RatingRange to simulate continuous x-axis,
		// string space to align label to the right of the bar
		xLabels = append(xLabels, fmt.Sprintf("%16v", b.RatingBoundHigh))
		itemTooltip := &opts.Tooltip{
			Formatter: types.FuncStr(fmt.Sprintf(`
				Rating range: %d → %d<br/>
				Percentile: %.2f%%<br/>
				Rank range: #%d → #%d<br/>
				Players count: %d<br/>`,
				b.RatingBoundLow, b.RatingBoundHigh,
				b.Percentile,
				b.RankBoundLow, b.RankBoundHigh,
				b.CountPlayers,
			))}
		yValues = append(yValues, opts.BarData{
			Value:   b.CountPlayers,
			Tooltip: itemTooltip,
		})
	}
	barChart.SetXAxis(xLabels).
		AddSeries("CountPlayers", yValues).
		SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
			BarWidth: "98%", // almost no gap between bars
		}))

	// TODO: draw vertical markers from percentileMarkers data
	// where rating at 25%, 50%, 75%, 90%, 99%, 99.9% percentile.
	markerChart := charts.NewLine()
	markerHeight := maxAxisY
	markerChart.SetGlobalOptions(
		charts.WithInitializationOpts(newInitializationOpts()),
		charts.WithXAxisOpts(newXAxisOptsContinuous()),
		charts.WithYAxisOpts(newYAxisOpts(false)),
		charts.WithGridOpts(newGridOpts()),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(false)}),
	)

	// Add one vertical line series per percentile marker
	for _, marker := range percentileMarkers {
		markerChart.
			AddSeries(
				fmt.Sprintf("%.1fth percentile", marker.Percentile*100), // series name
				[]opts.LineData{
					{
						Value: []float64{marker.Rating, 0},
						//Name:  fmt.Sprintf("%.0f", marker.Rating), // bot data item name
					},
					{
						Value: []float64{marker.Rating, float64(markerHeight)},
						Name: fmt.Sprintf("p%.1f\n\nElo %v",
							marker.Percentile*100, marker.Rating), // top data item name
					},
				},
			).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					Symbol: "circle", SymbolSize: 1,
				}),
				charts.WithLabelOpts(opts.Label{
					Show:      opts.Bool(true),
					Position:  "top", // label shows on top of the symbol
					Formatter: "{b}", // to make label depend on data, "{b}" means data.Name
				}),
				charts.WithLineStyleOpts(opts.LineStyle{
					Color: "grey", Width: 1, Type: "dashed"}),
			)
	}

	page := components.NewPage()
	page.SetPageTitle("AoE2DE rating distribution")
	page.AddCharts(barChart)
	page.AddCharts(markerChart)

	// write the chart as HTML to target output file
	f, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("create output file %q: %w", outputFilePath, err)
	}
	err = page.Render(f)
	if err != nil {
		return fmt.Errorf("render chart to %q: %w", outputFilePath, err)
	}
	_ = f.Close()
	return nil
}

func generateChartBackgroundDataURI(width, height int) string {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	orange := color.RGBA{R: 255, G: 165, B: 0, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, orange)
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})
	base64Img := base64.StdEncoding.EncodeToString(buf.Bytes())
	return fmt.Sprintf("url('data:image/jpeg;base64,%s')", base64Img)
}
