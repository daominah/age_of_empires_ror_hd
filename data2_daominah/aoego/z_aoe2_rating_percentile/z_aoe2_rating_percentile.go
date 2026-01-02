package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
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

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoego"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

//go:embed index_template.html
var indexTemplateHTML string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	const isForceReDownload = false

	today := time.Now().Format("2006-01-02")
	//today := "2025-11-01" // for testing with existing data only, should be commented out on normal run
	log.Printf("checking if data needs to be downloaded for today %v", today)

	// check whether data is downloaded
	projectRootDir, err := aoego.GetProjectRootGit()
	if err != nil {
		log.Fatalf("error GetProjectRootGit: %v", err)
	}
	goCodeDir := filepath.Join(projectRootDir, "data2_daominah", "aoego")
	todayOutputDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data", today)
	err = os.MkdirAll(todayOutputDir, 0755)
	if err != nil {
		log.Fatalf("error os.MkdirAll: %v", err)
	}

	// read the directory, download leaderboard if needed
	files, err := os.ReadDir(todayOutputDir)
	if err != nil {
		log.Fatalf("error os.ReadDir: %v", err)
	}
	if len(files) > 0 && !isForceReDownload {
		log.Printf("re-use existing ageofempires.com data, already have %d files", len(files))
	} else {
		log.Printf("downloading ageofempires.com data...")
		nDownloadedPages, err := DownloadAgeofempirescomData(todayOutputDir)
		if err != nil {
			log.Fatalf("error downloadAoe2insightsData: %v", err)
			return
		}
		log.Printf("downloaded %d pages of ageofempires.com data", nDownloadedPages)
	}

	// read all files in the directory and aggregate players
	players := make(map[int]AoEPlayer) // map key is "rlUserId"
	files, err = os.ReadDir(todayOutputDir)
	if err != nil {
		log.Fatalf("error os.ReadDir before aggregate players: %v", err)
	}
	for _, file := range files {
		filePath := filepath.Join(todayOutputDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("error os.ReadFile %v: %v", filePath, err)
		}
		var pageData AgeofempirescomDataResponse
		err = json.Unmarshal(data, &pageData)
		if err != nil {
			log.Fatalf("error json.Unmarshal %v: %v", filePath, err)
		}
		for _, player := range pageData.Items {
			players[player.RlUserId] = player
		}
	}

	// concise data for reduced size and re-use,
	// raw API responses size is about 16 MB, not good for GitHub Pages hosting,
	// concise JSON is 4MB, zipped is about xMB.
	var sortedPlayers []AoEPlayerLite
	for _, player := range players {
		sortedPlayers = append(sortedPlayers, player.ToLite())
	}
	sort.Slice(sortedPlayers, func(i, j int) bool {
		// highest rating player comes first
		return sortedPlayers[i].Elo > sortedPlayers[j].Elo
	})
	// save players concise data as zipped JSON
	liteDataBytes, err := json.MarshalIndent(sortedPlayers, "", "\t")
	if err != nil {
		log.Fatalf("error json.MarshalIndent liteData: %v", err)
	}
	fNameCompressedNoExt := fmt.Sprintf("all_players_%v", today)
	fPathCompressed := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_lite", fNameCompressedNoExt+".zip")
	zippedSize, err := saveToZip(fPathCompressed, fNameCompressedNoExt, liteDataBytes)
	if err != nil {
		log.Fatalf("error saveToZip fPathCompressed %v: %v", fPathCompressed, err)
	}
	log.Printf("wrote compressed players data to %v, size %v KiB", fNameCompressedNoExt, zippedSize/1024)

	log.Printf("-------------------------------------------------------")
	log.Printf("processing all zip files in data_lite and generate charts for each date")
	err = loopProcessAllZipFiles(goCodeDir)
	if err != nil {
		log.Fatalf("error loopProcessAllZipFiles: %v", err)
	}

	// final output "index.html" that can pick date to view a corresponding chart
	// from directory "output_charts"
	log.Printf("-------------------------------------------------------")
	log.Printf("combining all charts into index.html")
	err = generateIndexHTML(goCodeDir)
	if err != nil {
		log.Fatalf("error generateIndexHTML: %v", err)
	}
}

// readPlayersFromZip reads players data from a zip file
func readPlayersFromZip(zipPath string) ([]AoEPlayerLite, error) {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, fmt.Errorf("error opening zip file %v: %w", zipPath, err)
	}
	defer zipReader.Close()

	// Find the JSON file inside the zip
	var jsonFile *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, ".json") {
			jsonFile = f
			break
		}
	}
	if jsonFile == nil {
		return nil, fmt.Errorf("no JSON file found in zip %v", zipPath)
	}

	// Read the JSON file
	rc, err := jsonFile.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file in zip: %w", err)
	}
	defer rc.Close()

	jsonData, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data: %w", err)
	}

	var players []AoEPlayerLite
	err = json.Unmarshal(jsonData, &players)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return players, nil
}

// processPlayersData processes sorted players and generates chart data and CSV
func processPlayersData(sortedPlayers []AoEPlayerLite, dataDate string, goCodeDir string) ([]RatingBucket, []PercentileMarker, error) {
	// group players to Elo buckets, e.g. 1000-1099, 1100-1199, ..., 3100-3199
	const bucketSize float64 = 100
	bucketFunc := func(elo float64) string {
		roundDown := int(elo) / int(bucketSize) * int(bucketSize)
		roundUp := roundDown + int(bucketSize)
		return fmt.Sprintf("%04d→%04d", roundDown, roundUp)
	}
	ratingRanges := make(map[string][]AoEPlayerLite)
	for _, player := range sortedPlayers {
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
	totalPlayers := len(sortedPlayers)
	if totalPlayers == 0 {
		return nil, nil, fmt.Errorf("no players found in the data")
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

	// Ensure all buckets from 0→99 to 3100→3199 are present for consistent x-axis
	bucketsMap := make(map[string]RatingBucket)
	for _, bucket := range chartBars {
		bucketsMap[bucket.RatingRange] = bucket
	}
	const maxRating = 3200
	for rating := 0; rating < maxRating; rating += 100 {
		bucketKey := fmt.Sprintf("%04d→%04d", rating, rating+100)
		if _, exists := bucketsMap[bucketKey]; !exists {
			// Add empty bucket
			emptyBucket := RatingBucket{
				RatingRange:     bucketKey,
				RatingBoundLow:  rating,
				RatingBoundHigh: rating + 100,
				CountPlayers:    0,
				Percentile:      0,
				RankBoundLow:    0,
				RankBoundHigh:   0,
			}
			chartBars = append(chartBars, emptyBucket)
		}
	}
	// Sort chartBars by rating range
	sort.Slice(chartBars, func(i, j int) bool {
		return chartBars[i].RatingBoundLow < chartBars[j].RatingBoundLow
	})

	// write output CSV to a file
	outputCSVFileName := fmt.Sprintf("aoe2_rating_percentile_date_%v.csv", dataDate)
	outputCSVFilePath := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_summarized", outputCSVFileName)
	err := os.MkdirAll(filepath.Dir(outputCSVFilePath), 0755)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating CSV output directory: %w", err)
	}
	outputCSVFile, err := os.Create(outputCSVFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error os.Create outputCSVFilePath %v: %w", outputCSVFilePath, err)
	}
	csvWriter := csv.NewWriter(outputCSVFile)
	err = csvWriter.WriteAll(dataAsCSV)
	if err != nil {
		_ = outputCSVFile.Close()
		return nil, nil, fmt.Errorf("error csvWriter.WriteAll to %v: %w", outputCSVFilePath, err)
	}
	csvWriter.Flush()
	err = outputCSVFile.Close()
	if err != nil {
		return nil, nil, fmt.Errorf("error outputCSVFile.Close %v: %w", outputCSVFilePath, err)
	}
	//log.Printf("wrote rating percentile to %v", outputCSVFileName)

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
	}

	return chartBars, percentileMarkers, nil
}

// generateChartForDate generates a chart HTML file for a specific date
func generateChartForDate(
	goCodeDir string,
	dataISOStr string,
	chartBars []RatingBucket,
	percentileMarkers []PercentileMarker) error {
	outputChartFileName := fmt.Sprintf("chart_%v.html", dataISOStr)
	outputChartsDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "output_charts")
	err := os.MkdirAll(outputChartsDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating output_charts directory: %w", err)
	}
	outputChartFileFullPath := filepath.Join(outputChartsDir, outputChartFileName)

	//skip if output file already exists
	if _, err := os.Stat(outputChartFileFullPath); err == nil {
		log.Printf("skipping, chart file already exists %v", outputChartFileName)
		return nil
	}

	chartWidth, chartHeight := 1800, 800
	err = drawPercentilesChart(chartBars, percentileMarkers,
		dataISOStr, chartWidth, chartHeight, outputChartFileFullPath)
	if err != nil {
		return fmt.Errorf("error drawPercentilesChart: %w", err)
	}
	log.Printf("drew rating percentile chart to %v", outputChartFileName)

	// After page.Render(f) and closing the file, read-modify-write the HTML:
	htmlBytes, err := os.ReadFile(outputChartFileFullPath)
	if err != nil {
		return fmt.Errorf("error reading chart HTML: %w", err)
	}
	htmlStr := string(htmlBytes)

	// Inject CSS to remove vertical scrolling and disable animations
	cssStyle := `
	<style>
		html, body {
			overflow-y: hidden !important;
			height: 100%;
			margin: 0;
			padding: 0;
		}
		* {
			animation: none !important;
			transition: none !important;
		}
	</style>
	`
	// Inject CSS in the head section
	if strings.Contains(htmlStr, "</head>") {
		htmlStr = strings.Replace(htmlStr, "</head>", cssStyle+"</head>", 1)
	} else {
		// Fallback: inject before </body> if </head> not found
		htmlStr = strings.Replace(htmlStr, "<body>", "<body>"+cssStyle, 1)
	}

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
		return fmt.Errorf("error writing modified chart HTML: %w", err)
	}
	//log.Printf("injected overlay script into chart HTML %v", outputChartFileName)
	return nil
}

// generateIndexHTML creates an index.html file that lists all available charts and allows date selection
func generateIndexHTML(goCodeDir string) error {
	outputChartsDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "output_charts")

	// Read all files in output_charts directory
	files, err := os.ReadDir(outputChartsDir)
	if err != nil {
		return fmt.Errorf("error reading output_charts directory: %w", err)
	}

	// Extract dates from chart filenames
	var dates []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if !strings.HasPrefix(filename, "chart_") || !strings.HasSuffix(filename, ".html") {
			continue
		}
		// Extract date from "chart_yyyy-mm-dd.html"
		dateStr := strings.TrimPrefix(filename, "chart_")
		dateStr = strings.TrimSuffix(dateStr, ".html")

		// Validate date format
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("skipping file with invalid date format: %v", filename)
			continue
		}
		dates = append(dates, dateStr)
	}
	if len(dates) == 0 {
		return fmt.Errorf("no chart files found in output_charts directory")
	}
	// Sort dates ISO string in descending order (newest first)
	sort.Slice(dates, func(i, j int) bool {
		return dates[i] > dates[j]
	})

	// Generate date options HTML
	var dateOptions strings.Builder
	for i, date := range dates {
		selected := ""
		if i == 0 { // select the newest date by default
			selected = " selected"
		}
		dateOptions.WriteString(fmt.Sprintf(`				<option value="%s"%s>%s</option>
			`, date, selected, date))
	}

	// Replace placeholder in template with actual date options
	htmlContent := strings.Replace(indexTemplateHTML, "{{DATE_OPTIONS}}", dateOptions.String(), 1)
	// remove default <option value="">-- Select a date --</option>
	htmlContent = strings.Replace(htmlContent, `                <option value="">-- Select a date --</option>`, "", 1)

	// Write index.html to the z_aoe2_rating_percentile directory
	indexHTMLPath := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "index.html")
	err = os.WriteFile(indexHTMLPath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing index.html: %w", err)
	}

	log.Printf("generated index.html with %d available charts", len(dates))
	return nil
}

// loopProcessAllZipFiles draws charts for all zip files, each zip produces one HTML chart
func loopProcessAllZipFiles(goCodeDir string) error {
	dataLiteDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_lite")

	// Read all files in data_lite directory
	files, err := os.ReadDir(dataLiteDir)
	if err != nil {
		return fmt.Errorf("error reading data_lite directory: %w", err)
	}

	// Process each zip file
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".zip") {
			continue
		}

		// Extract date from filename: all_players_yyyy-mm-dd.zip
		zipName := file.Name()
		if !strings.HasPrefix(zipName, "all_players_") {
			log.Printf("skipping file with unexpected name format: %v", zipName)
			continue
		}
		dateStr := strings.TrimPrefix(zipName, "all_players_")
		dateStr = strings.TrimSuffix(dateStr, ".zip")

		// Validate date format
		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Printf("skipping file with invalid date format: %v", zipName)
			continue
		}

		log.Printf("processing file: %v", zipName)

		zipPath := filepath.Join(dataLiteDir, zipName)

		// Read players from zip
		sortedPlayers, err := readPlayersFromZip(zipPath)
		if err != nil {
			log.Printf("error reading zip file %v: %v", zipName, err)
			continue
		}

		// Sort players by rating (highest first)
		sort.Slice(sortedPlayers, func(i, j int) bool {
			return sortedPlayers[i].Elo > sortedPlayers[j].Elo
		})

		// Process data to generate chart bars and percentile markers
		chartBars, percentileMarkers, err := processPlayersData(sortedPlayers, dateStr, goCodeDir)
		if err != nil {
			log.Printf("error processing players data for %v: %v", dateStr, err)
			continue
		}

		// Generate chart
		err = generateChartForDate(goCodeDir, dateStr, chartBars, percentileMarkers)
		if err != nil {
			log.Printf("error generating chart for %v: %v", dateStr, err)
			continue
		}

		//log.Printf("successfully processed %v", dateStr)
	}

	return nil
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

type AoEPlayerLite struct {
	RlUserId int
	UserName string
	Elo      float64
	Rank     int
}

func (p AoEPlayer) ToLite() AoEPlayerLite {
	return AoEPlayerLite{
		RlUserId: p.RlUserId,
		UserName: p.UserName,
		Elo:      p.Elo,
		Rank:     p.Rank,
	}
}

func saveToZip(fileFullPath string, fileNameNoExt string, data []byte) (int, error) {
	zipFile, err := os.Create(fileFullPath)
	if err != nil {
		return 0, fmt.Errorf("os.Create: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	w, err := zipWriter.Create(fileNameNoExt + ".json")
	if err != nil {
		return 0, fmt.Errorf("zipWriter.Create: %w", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return 0, fmt.Errorf("w.Write: %w", err)
	}
	err = zipWriter.Close()
	if err != nil {
		return 0, fmt.Errorf("zipWriter.Close: %w", err)
	}
	info, err := os.Stat(fileFullPath)
	if err != nil {
		return 0, fmt.Errorf("os.Stat: %w", err)
	}
	zippedBytes := int(info.Size())
	return zippedBytes, nil
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

func printPercentileByHumanLevels(players []AoEPlayerLite) {
	fmt.Printf("______________________________________________________\n")
	group0To800 := 0  // beginner
	group0To1000 := 0 // novice
	group0To1200 := 0 // developing
	group0To1500 := 0 // competent
	group0To2000 := 0 // expert
	group0To2500 := 0 // elite
	group2500Up := 0  // professional
	totalPlayers := len(players)
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
	group2500Up = totalPlayers - group0To2500

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
	fmt.Printf("* rating 2500 and up: percentile 100%%  (rank #%5d → #    1)\n", group2500Up)
	// Output:
	//	* rating    0 →  800: percentile 26.1% (rank #43139 → #31878)
	//	* rating  800 → 1000: percentile 49.2% (rank #31877 → #21906)
	//	* rating 1000 → 1200: percentile 73.6% (rank #21905 → #11376)
	//	* rating 1200 → 1500: percentile 90.7% (rank #11375 → # 4017)
	//	* rating 1500 → 2000: percentile 98.8% (rank # 4016 → #  499)
	//	* rating 2000 → 2500: percentile 99.9% (rank #  498 → #   59)
	//	* rating 2500 and up: percentile 100%  (rank #   58 → #    1)
	fmt.Printf("______________________________________________________\n")

}

// drawPercentilesChart draws a chart and save as HTML to outputFilePath,
// the chart x-axis is rating buckets, y-axis is number of players in that bucket.
func drawPercentilesChart(
	bars []RatingBucket,
	percentileMarkers []PercentileMarker,
	dataDate string,
	chartWidth, chartHeight int,
	outputFilePath string) error {
	maxAxisX := 3200 // Always use 3200 for consistent x-axis across all charts
	maxAxisY := 0
	totalPlayers := 0
	for _, b := range bars {
		roundUpNPlayers := int(math.Ceil(float64(b.CountPlayers+500)/1000)) * 1000
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
		charts.WithAnimation(false),                                // Disable animation
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
		charts.WithAnimation(false), // Disable animation
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
