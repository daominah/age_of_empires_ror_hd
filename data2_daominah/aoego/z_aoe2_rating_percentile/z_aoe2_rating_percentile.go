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

	// step 1: download raw leaderboard pages for today
	err = stepDownload(todayOutputDir, isForceReDownload)
	if err != nil {
		log.Fatalf("error stepDownload: %v", err)
	}

	// step 2: merge pages into unique players
	players, err := stepMergePages(todayOutputDir)
	if err != nil {
		log.Fatalf("error stepMergePages: %v", err)
	}

	// step 3: save today's players as concise zipped JSON in "data_lite"
	err = stepSaveLite(players, goCodeDir, today)
	if err != nil {
		log.Fatalf("error stepSaveLite: %v", err)
	}

	// steps 4 to 6 for every date in "data_lite":
	// read zip, detect and fill rank gaps, summarize to CSV, draw chart
	log.Printf("-------------------------------------------------------")
	log.Printf("processing all zip files in data_lite and generate charts for each date")
	err = loopProcessAllZipFiles(goCodeDir)
	if err != nil {
		log.Fatalf("error loopProcessAllZipFiles: %v", err)
	}

	// step 7: final output "index.html" that can pick date to view a corresponding chart
	// from directory "output_charts"
	log.Printf("-------------------------------------------------------")
	log.Printf("combining all charts into index.html")
	err = generateIndexHTML(goCodeDir)
	if err != nil {
		log.Fatalf("error generateIndexHTML: %v", err)
	}

}

// stepDownload downloads leaderboard pages to todayOutputDir,
// re-uses existing pages unless isForceReDownload.
func stepDownload(todayOutputDir string, isForceReDownload bool) error {
	files, err := os.ReadDir(todayOutputDir)
	if err != nil {
		return fmt.Errorf("error os.ReadDir: %w", err)
	}
	if len(files) > 0 && !isForceReDownload {
		log.Printf("re-use existing ageofempires.com data, already have %d files", len(files))
		return nil
	}
	log.Printf("downloading ageofempires.com data...")
	nDownloadedPages, err := DownloadAgeofempirescomData(todayOutputDir)
	if err != nil {
		return fmt.Errorf("error DownloadAgeofempirescomData: %w", err)
	}
	log.Printf("downloaded %d pages of ageofempires.com data", nDownloadedPages)
	return nil
}

// stepMergePages reads all downloaded pages and returns unique players by RlUserId.
func stepMergePages(todayOutputDir string) (map[int]AoEPlayer, error) {
	players := make(map[int]AoEPlayer) // map key is "rlUserId"
	files, err := os.ReadDir(todayOutputDir)
	if err != nil {
		return nil, fmt.Errorf("error os.ReadDir: %w", err)
	}
	for _, file := range files {
		filePath := filepath.Join(todayOutputDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("error os.ReadFile %v: %w", filePath, err)
		}
		var pageData AgeofempirescomDataResponse
		err = json.Unmarshal(data, &pageData)
		if err != nil {
			return nil, fmt.Errorf("error json.Unmarshal %v: %w", filePath, err)
		}
		for _, player := range pageData.Items {
			players[player.RlUserId] = player
		}
	}
	return players, nil
}

// stepSaveLite saves players as concise zipped JSON sorted by Elo,
// raw API responses size is about 16 MB, not good for GitHub Pages hosting,
// concise JSON is 4MB, zipped is smaller.
// Correct in memory only: the zip stays as the API returned it (no estimated players).
func stepSaveLite(players map[int]AoEPlayer, goCodeDir string, today string) error {
	var sortedPlayers []AoEPlayerLite
	for _, player := range players {
		sortedPlayers = append(sortedPlayers, player.ToLite())
	}
	sort.Slice(sortedPlayers, func(i, j int) bool {
		// highest rating player comes first
		return sortedPlayers[i].Elo > sortedPlayers[j].Elo
	})
	liteDataBytes, err := json.MarshalIndent(sortedPlayers, "", "\t")
	if err != nil {
		return fmt.Errorf("error json.MarshalIndent liteData: %w", err)
	}
	fNameCompressedNoExt := fmt.Sprintf("all_players_%v", today)
	fPathCompressed := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_lite", fNameCompressedNoExt+".zip")
	zippedSize, err := saveToZip(fPathCompressed, fNameCompressedNoExt, liteDataBytes)
	if err != nil {
		return fmt.Errorf("error saveToZip %v: %w", fPathCompressed, err)
	}
	log.Printf("wrote compressed players data to %v, size %v KiB", fNameCompressedNoExt, zippedSize/1024)
	return nil
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
	dataAsCSV = append(dataAsCSV, []string{"RatingRange", "RatingLow", "RatingHigh", "CountPlayers", "Percentile", "RankLow", "RankHigh", "EstimatedPlayers"})
	var chartBars []RatingBucket
	for _, bucketKey := range sortedBucketKeys {
		playersInBucket := ratingRanges[bucketKey]
		cumulativeToCurrentBucket += len(playersInBucket)
		nEstimatedInBucket := 0
		for _, p := range playersInBucket {
			if p.IsEstimated {
				nEstimatedInBucket++
			}
		}
		ratingBucket := RatingBucket{
			RatingRange:      bucketKey,
			RatingBoundLow:   0, // will be set later based on RatingRange
			RatingBoundHigh:  0, // will be set later based on RatingRange
			CountPlayers:     len(playersInBucket),
			Percentile:       float64(cumulativeToCurrentBucket) / float64(totalPlayers) * 100.0,
			RankBoundLow:     totalPlayers - (cumulativeToCurrentBucket - len(playersInBucket)),
			RankBoundHigh:    totalPlayers - cumulativeToCurrentBucket + 1,
			EstimatedPlayers: nEstimatedInBucket,
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
			fmt.Sprintf("%v", ratingBucket.EstimatedPlayers),
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
	percentileMarkers []PercentileMarker,
	nEstimated int) error {
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
	err = drawPercentilesChart(chartBars, percentileMarkers, nEstimated,
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

// loopProcessAllZipFiles fills rank gaps, then writes the CSV summary and HTML chart of every zip.
// What goes wrong and the decisions: "fill-leaderboard-rank-gaps.md".
func loopProcessAllZipFiles(goCodeDir string) error {
	dataLiteDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_lite")
	zips, err := listDataLiteZips(dataLiteDir)
	if err != nil {
		return fmt.Errorf("error listDataLiteZips: %w", err)
	}

	// first pass: Elo by rank of every day,
	// so a gap can be filled from a reference day, the latest earlier day that has those ranks
	eloByRankOfDays := make(map[string]EloByRank)
	for _, z := range zips {
		players, err := readPlayersFromZip(z.Path)
		if err != nil {
			log.Printf("error reading zip file %v: %v", z.Path, err)
			continue
		}
		eloByRankOfDays[z.Date] = newEloByRank(players)
	}

	for _, z := range zips {
		log.Printf("processing file: %v", filepath.Base(z.Path))

		// step 4: read players from zip
		sortedPlayers, err := readPlayersFromZip(z.Path)
		if err != nil {
			log.Printf("error reading zip file %v: %v", z.Path, err)
			continue
		}

		// step 4.5a: detect rank gaps, blocks of ranks the API left out that day
		gaps := detectRankGaps(sortedPlayers)
		nEstimated := countMissing(gaps)

		// step 4.5b: fill gaps from a reference day, correct in memory only
		references := make([]EloByRank, len(gaps))
		nFromReference := 0
		for i, gap := range gaps {
			references[i] = findReferenceDay(gap, z.Date, eloByRankOfDays)
			if references[i] != nil {
				nFromReference += gap.RankLast - gap.RankFirst + 1
			}
		}
		if nEstimated > 0 {
			log.Printf("date %v: %v rank gaps, %v players missing from the API, %v estimated from a reference day",
				z.Date, len(gaps), nEstimated, nFromReference)
		}
		sortedPlayers = fillRankGaps(sortedPlayers, gaps, references)

		// Sort players by rating (highest first)
		sort.Slice(sortedPlayers, func(i, j int) bool {
			return sortedPlayers[i].Elo > sortedPlayers[j].Elo
		})

		// step 5: process data to generate chart bars and percentile markers
		chartBars, percentileMarkers, err := processPlayersData(sortedPlayers, z.Date, goCodeDir)
		if err != nil {
			log.Printf("error processing players data for %v: %v", z.Date, err)
			continue
		}

		// step 6: generate chart
		err = generateChartForDate(goCodeDir, z.Date, chartBars, percentileMarkers, nEstimated)
		if err != nil {
			log.Printf("error generating chart for %v: %v", z.Date, err)
			continue
		}
	}
	return nil
}

// listDataLiteZips returns the "all_players_yyyy-mm-dd.zip" files in dataLiteDir, sorted by date.
func listDataLiteZips(dataLiteDir string) ([]DataLiteZip, error) {
	files, err := os.ReadDir(dataLiteDir)
	if err != nil {
		return nil, fmt.Errorf("error os.ReadDir: %w", err)
	}
	var zips []DataLiteZip
	for _, file := range files {
		zipName := file.Name()
		if !strings.HasSuffix(zipName, ".zip") {
			continue
		}
		if !strings.HasPrefix(zipName, "all_players_") {
			log.Printf("skipping file with unexpected name format: %v", zipName)
			continue
		}
		dateStr := strings.TrimSuffix(strings.TrimPrefix(zipName, "all_players_"), ".zip")
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			log.Printf("skipping file with invalid date format: %v", zipName)
			continue
		}
		zips = append(zips, DataLiteZip{Date: dateStr, Path: filepath.Join(dataLiteDir, zipName)})
	}
	return zips, nil
}

// detectRankGaps returns blocks of ranks missing between known players.
// The ageofempires.com API sometimes serves a leaderboard without whole pages,
// or ranks shift while pages download minutes apart,
// and its Count shrinks to match, so gaps in Rank are the only signal.
// Players missing after the last rank cannot be detected.
func detectRankGaps(players []AoEPlayerLite) []RankGap {
	byRank := make([]AoEPlayerLite, len(players))
	copy(byRank, players)
	sort.Slice(byRank, func(i, j int) bool {
		return byRank[i].Rank < byRank[j].Rank
	})
	var gaps []RankGap
	for i := 1; i < len(byRank); i++ {
		above, below := byRank[i-1], byRank[i]
		if below.Rank-above.Rank <= 1 {
			continue
		}
		gaps = append(gaps, RankGap{
			RankFirst: above.Rank + 1,
			RankLast:  below.Rank - 1,
			EloAbove:  above.Elo,
			EloBelow:  below.Elo,
		})
	}
	return gaps
}

// fillRankGaps returns players plus one estimated player per missing rank,
// estimated players have IsEstimated and an Elo inside the gap bounds.
// The count and range of each gap are exact,
// only the spread of Elo inside the gap is estimated.
// Fill from a reference day, not a bell curve:
//   - references[i] not nil: copy the rating shape of that reference day at the same ranks,
//     rescaled so the gap's neighbors on that day land on today's neighbors.
//     The real rating distribution has a long high tail that no simple curve fits,
//     so a nearby day is the closest thing to the missing players.
//   - otherwise, the fallback: a bell curve (normal distribution) fitted to the known players,
//     the gap's Elo range is split into equal areas under the curve.
func fillRankGaps(players []AoEPlayerLite, gaps []RankGap, references []EloByRank) []AoEPlayerLite {
	if len(gaps) == 0 {
		return players
	}
	mean, stdDev := fitNormal(players)
	filled := make([]AoEPlayerLite, len(players), len(players)+countMissing(gaps))
	copy(filled, players)
	for i, gap := range gaps {
		var reference EloByRank
		if i < len(references) {
			reference = references[i]
		}
		nMissing := gap.RankLast - gap.RankFirst + 1
		cdfAbove := normalCDF(gap.EloAbove, mean, stdDev)
		cdfBelow := normalCDF(gap.EloBelow, mean, stdDev)
		// far in the tails the curve is flat in float64,
		// fall back to evenly spaced Elo
		isCurveUsable := stdDev > 0 && cdfAbove-cdfBelow > 1e-12
		for j := range nMissing {
			rank := gap.RankFirst + j
			// fraction of the way from EloAbove down to EloBelow,
			// the known neighbors sit at positions 0 and nMissing+1
			fraction := float64(j+1) / float64(nMissing+1)
			if reference != nil {
				refAbove, refBelow := reference[gap.RankFirst-1], reference[gap.RankLast+1]
				fraction = (refAbove - reference[rank]) / (refAbove - refBelow)
			}
			var elo float64
			if reference == nil && isCurveUsable {
				elo = normalQuantile(cdfAbove-fraction*(cdfAbove-cdfBelow), mean, stdDev)
			} else {
				elo = gap.EloAbove - fraction*(gap.EloAbove-gap.EloBelow)
			}
			elo = math.Max(gap.EloBelow, math.Min(gap.EloAbove, math.Round(elo)))
			filled = append(filled, AoEPlayerLite{Elo: elo, Rank: rank, IsEstimated: true})
		}
	}
	return filled
}

// findReferenceDay returns Elo by rank of the latest day before date
// that has every rank of the gap and both its neighbors,
// or nil if no earlier day qualifies.
// Fill from earlier days data only:
// a chart built on its own date gets the same estimates as a later rebuild.
func findReferenceDay(gap RankGap, date string, eloByRankOfDays map[string]EloByRank) EloByRank {
	var best EloByRank
	bestDate := ""
	// "yyyy-mm-dd" dates compare in time order as strings
	for otherDate, eloByRank := range eloByRankOfDays {
		if otherDate >= date || otherDate <= bestDate {
			continue
		}
		if !eloByRank.hasRanks(gap.RankFirst-1, gap.RankLast+1) {
			continue
		}
		if eloByRank[gap.RankFirst-1] <= eloByRank[gap.RankLast+1] {
			continue // no spread to copy
		}
		best, bestDate = eloByRank, otherDate
	}
	return best
}

// newEloByRank indexes players' Elo by rank, NaN for ranks with no player.
func newEloByRank(players []AoEPlayerLite) EloByRank {
	maxRank := 0
	for _, p := range players {
		maxRank = max(maxRank, p.Rank)
	}
	eloByRank := make(EloByRank, maxRank+1)
	for i := range eloByRank {
		eloByRank[i] = math.NaN()
	}
	for _, p := range players {
		if p.Rank > 0 {
			eloByRank[p.Rank] = p.Elo
		}
	}
	return eloByRank
}

// hasRanks reports whether every rank from first to last has a player.
func (e EloByRank) hasRanks(first int, last int) bool {
	if first < 1 || last >= len(e) {
		return false
	}
	for rank := first; rank <= last; rank++ {
		if math.IsNaN(e[rank]) {
			return false
		}
	}
	return true
}

func countMissing(gaps []RankGap) int {
	n := 0
	for _, gap := range gaps {
		n += gap.RankLast - gap.RankFirst + 1
	}
	return n
}

// fitNormal returns the mean and standard deviation of players' Elo.
func fitNormal(players []AoEPlayerLite) (float64, float64) {
	if len(players) == 0 {
		return 0, 0
	}
	sum := 0.0
	for _, p := range players {
		sum += p.Elo
	}
	mean := sum / float64(len(players))
	sumSquares := 0.0
	for _, p := range players {
		sumSquares += (p.Elo - mean) * (p.Elo - mean)
	}
	return mean, math.Sqrt(sumSquares / float64(len(players)))
}

func normalCDF(x float64, mean float64, stdDev float64) float64 {
	return 0.5 * (1 + math.Erf((x-mean)/(stdDev*math.Sqrt2)))
}

func normalQuantile(p float64, mean float64, stdDev float64) float64 {
	return mean + stdDev*math.Sqrt2*math.Erfinv(2*p-1)
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
			log.Printf("error downloading page %d: %v, continuing to next page", page, err)
			continue
		}
		var pageDataObj AgeofempirescomDataResponse
		err = json.Unmarshal(pageData, &pageDataObj)
		if err != nil {
			log.Printf("error json.Unmarshal page %d: %v, continuing to next page", page, err)
			continue
		}
		beautyJSON, err := json.MarshalIndent(pageDataObj, "", "\t")
		if err != nil {
			log.Printf("error json.MarshalIndent page %d: %v, continuing to next page", page, err)
			continue
		}
		outputFile := filepath.Join(todayOutputDir,
			fmt.Sprintf("ageofempirescom_leaderboard_page_%03d.json", page))
		err = os.WriteFile(outputFile, beautyJSON, 0644)
		if err != nil {
			log.Printf("error os.WriteFile page %d: %v, continuing to next page", page, err)
			continue
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
	// IsEstimated marks a player added in memory for a rank the API left out,
	// never saved to the zip
	IsEstimated bool `json:"-"`
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
	// EstimatedPlayers is how many of CountPlayers have ratings estimated, for ranks the API left out,
	// the last CSV column to disclose estimates
	EstimatedPlayers int
}

// DataLiteZip is one day of saved players in "data_lite".
type DataLiteZip struct {
	Date string // "yyyy-mm-dd"
	Path string
}

// EloByRank is one day's Elo indexed by rank, NaN where the rank has no player.
type EloByRank []float64

// RankGap is a block of consecutive ranks missing from the API response,
// bounded by the known players right above and right below it.
type RankGap struct {
	RankFirst int
	RankLast  int
	EloAbove  float64 // Elo of the player at RankFirst-1
	EloBelow  float64 // Elo of the player at RankLast+1
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
	nEstimated int,
	dataDate string,
	chartWidth, chartHeight int,
	outputFilePath string) error {
	maxAxisX := 3200 // Always use 3200 for consistent x-axis across all charts
	maxAxisY := 8000
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
	subtitle := fmt.Sprintf("Data from ageofempires.com leaderboards on %v with a total of %v players",
		dataDate, totalPlayers)
	if nEstimated > 0 {
		subtitle += fmt.Sprintf(" (ratings of %v players estimated from earlier days, the API left them out that day)", nEstimated)
	}
	newTitleOpts := func() opts.Title {
		return opts.Title{
			Title:    "AoE2DE rating distribution",
			Subtitle: subtitle,
			Left:     "center", Top: "0px",
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
