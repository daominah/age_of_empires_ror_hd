package main

import (
	"encoding/csv"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// newLinearPlayers returns players ranked 1..nPlayers, Elo = 2000 - rank,
// so each 100-Elo bucket from 1000 to 2000 holds exactly 100 players.
func newLinearPlayers(nPlayers int) []AoEPlayerLite {
	players := make([]AoEPlayerLite, 0, nPlayers)
	for rank := 1; rank <= nPlayers; rank++ {
		players = append(players, AoEPlayerLite{
			RlUserId: rank,
			UserName: "player" + strconv.Itoa(rank),
			Elo:      float64(2000 - rank),
			Rank:     rank,
		})
	}
	return players
}

// removeRanks returns players without ranks in [rankFirst, rankLast].
func removeRanks(players []AoEPlayerLite, rankFirst int, rankLast int) []AoEPlayerLite {
	var kept []AoEPlayerLite
	for _, p := range players {
		if p.Rank >= rankFirst && p.Rank <= rankLast {
			continue
		}
		kept = append(kept, p)
	}
	return kept
}

// readBucketCounts reads a summarized CSV into a map from RatingLow to CountPlayers.
func readBucketCounts(t *testing.T, csvPath string) map[int]int {
	t.Helper()
	return readBucketColumn(t, csvPath, 3)
}

// readBucketColumn reads a summarized CSV into a map from RatingLow to the integer column.
func readBucketColumn(t *testing.T, csvPath string, column int) map[int]int {
	t.Helper()
	f, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("error os.Open: %v", err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("error csv ReadAll: %v", err)
	}
	counts := make(map[int]int)
	for _, row := range rows[1:] {
		ratingLow, _ := strconv.Atoi(row[1])
		count, _ := strconv.Atoi(row[column])
		counts[ratingLow] = count
	}
	return counts
}

func TestLoopProcessAllZipFiles_RankGapsFilled(t *testing.T) {
	// GIVEN a day where the API left out two pages of players:
	// ranks 301-400 all inside the 1600-1700 bucket,
	// and ranks 451-550 spanning the 1400-1500 and 1500-1600 buckets
	players := newLinearPlayers(1000)
	players = removeRanks(players, 301, 400)
	players = removeRanks(players, 451, 550)
	goCodeDir := t.TempDir()
	dataLiteDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_lite")
	if err := os.MkdirAll(dataLiteDir, 0755); err != nil {
		t.Fatalf("error os.MkdirAll: %v", err)
	}
	liteBytes, err := json.Marshal(players)
	if err != nil {
		t.Fatalf("error json.Marshal: %v", err)
	}
	const date = "2026-01-01"
	zipPath := filepath.Join(dataLiteDir, "all_players_"+date+".zip")
	if _, err := saveToZip(zipPath, "all_players_"+date, liteBytes); err != nil {
		t.Fatalf("error saveToZip: %v", err)
	}

	// WHEN the daily summary and chart are generated from the saved players
	if err := loopProcessAllZipFiles(goCodeDir); err != nil {
		t.Fatalf("error loopProcessAllZipFiles: %v", err)
	}

	// THEN the summary counts every ranked player, including the missing ones
	csvPath := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_summarized",
		"aoe2_rating_percentile_date_"+date+".csv")
	counts := readBucketCounts(t, csvPath)
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != 1000 {
		t.Errorf("total players: got %v, want 1000", total)
	}
	// THEN a gap inside one bucket is restored exactly
	if counts[1600] != 100 {
		t.Errorf("bucket 1600-1700: got %v, want 100", counts[1600])
	}
	// THEN a gap spanning two buckets is split close to the true 50/50
	if counts[1400]+counts[1500] != 200 {
		t.Errorf("buckets 1400-1600: got %v, want 200", counts[1400]+counts[1500])
	}
	for _, low := range []int{1400, 1500} {
		if counts[low] < 95 || counts[low] > 105 {
			t.Errorf("bucket %v: got %v, want about 100", low, counts[low])
		}
	}
	// THEN the summary tells readers how many players in each bucket are estimated
	estimated := readBucketColumn(t, csvPath, 7)
	if estimated[1600] != 100 || estimated[1400]+estimated[1500] != 100 {
		t.Errorf("estimated players: got %v, want 100 in 1600-1700 and 100 across 1400-1600", estimated)
	}
	// THEN the chart tells readers how many players are estimated
	chartPath := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "output_charts", "chart_"+date+".html")
	chartBytes, err := os.ReadFile(chartPath)
	if err != nil {
		t.Fatalf("error os.ReadFile: %v", err)
	}
	if !strings.Contains(string(chartBytes), "(ratings of 200 players estimated from earlier days, the API left them out that day)") {
		t.Errorf("chart subtitle does not mention the 200 estimated players")
	}
}

func TestDetectRankGaps(t *testing.T) {
	// GIVEN a leaderboard missing ranks 301-400 and 451-550,
	// served in a shuffled order
	players := newLinearPlayers(1000)
	players = removeRanks(players, 301, 400)
	players = removeRanks(players, 451, 550)
	players[0], players[len(players)-1] = players[len(players)-1], players[0]

	// WHEN gaps are detected
	gaps := detectRankGaps(players)

	// THEN each missing block is reported with the Elo of its neighbors
	want := []RankGap{
		{RankFirst: 301, RankLast: 400, EloAbove: 1700, EloBelow: 1599},
		{RankFirst: 451, RankLast: 550, EloAbove: 1550, EloBelow: 1449},
	}
	if len(gaps) != len(want) {
		t.Fatalf("gaps: got %+v, want %+v", gaps, want)
	}
	for i := range want {
		if gaps[i] != want[i] {
			t.Errorf("gap %v: got %+v, want %+v", i, gaps[i], want[i])
		}
	}
}

func TestDetectRankGaps_Complete(t *testing.T) {
	// GIVEN a complete leaderboard
	players := newLinearPlayers(1000)

	// WHEN gaps are detected
	gaps := detectRankGaps(players)

	// THEN no gap is reported
	if len(gaps) != 0 {
		t.Errorf("gaps: got %+v, want none", gaps)
	}
}

func TestFillRankGaps(t *testing.T) {
	// GIVEN a leaderboard missing ranks 301-400
	players := removeRanks(newLinearPlayers(1000), 301, 400)
	gaps := detectRankGaps(players)

	// WHEN the gaps are filled
	filled := fillRankGaps(players, gaps, nil)

	// THEN there is exactly one player per rank
	if len(filled) != 1000 {
		t.Fatalf("players: got %v, want 1000", len(filled))
	}
	seenRanks := make(map[int]bool)
	for _, p := range filled {
		seenRanks[p.Rank] = true
	}
	if len(seenRanks) != 1000 {
		t.Errorf("unique ranks: got %v, want 1000", len(seenRanks))
	}
	// THEN estimated players are marked, stay within the gap's Elo bounds,
	// and a better rank never has a lower Elo
	prevElo := math.Inf(1)
	for _, p := range filled[len(players):] {
		if !p.IsEstimated {
			t.Errorf("rank %v: estimated player is not marked IsEstimated", p.Rank)
		}
		if p.Elo > 1700 || p.Elo < 1599 {
			t.Errorf("rank %v: Elo %v outside gap bounds 1599-1700", p.Rank, p.Elo)
		}
		if p.Elo > prevElo {
			t.Errorf("rank %v: Elo %v higher than the rank above (%v)", p.Rank, p.Elo, prevElo)
		}
		prevElo = p.Elo
	}
	// THEN the input slice is not modified
	if len(players) != 900 {
		t.Errorf("input players: got %v, want 900", len(players))
	}
}

func TestFillRankGaps_FarTail(t *testing.T) {
	// GIVEN a gap so far above the average that the fitted curve is flat
	players := []AoEPlayerLite{
		{RlUserId: 1, Elo: 90000, Rank: 1},
		{RlUserId: 5, Elo: 80000, Rank: 5},
	}
	for rank := 6; rank <= 1000; rank++ {
		players = append(players, AoEPlayerLite{RlUserId: rank, Elo: 1000, Rank: rank})
	}
	gaps := detectRankGaps(players)

	// WHEN the gap is filled
	filled := fillRankGaps(players, gaps, nil)

	// THEN the estimated players are spread evenly between the neighbors
	var elos []float64
	for _, p := range filled[len(players):] {
		elos = append(elos, p.Elo)
	}
	want := []float64{87500, 85000, 82500}
	if len(elos) != len(want) {
		t.Fatalf("estimated Elo: got %v, want %v", elos, want)
	}
	for i := range want {
		if elos[i] != want[i] {
			t.Errorf("estimated Elo: got %v, want %v", elos, want)
			break
		}
	}
}

func TestLoopProcessAllZipFiles_RealData20260925(t *testing.T) {
	// GIVEN the real leaderboard saved on 2026-09-25,
	// when the API left out 5200 of 46168 ranked players in 6 blocks
	goCodeDir := newGoCodeDirWithFixtures(t, "2026-09-25")
	players, err := readPlayersFromZip(filepath.Join(goCodeDir,
		"z_aoe2_rating_percentile", "data_lite", "all_players_2026-09-25.zip"))
	if err != nil {
		t.Fatalf("error readPlayersFromZip: %v", err)
	}

	// WHEN the missing blocks are detected
	gaps := detectRankGaps(players)

	// THEN all 6 blocks are found
	wantGaps := [][2]int{{2201, 2400}, {12801, 14600}, {19601, 21600}, {23001, 23600}, {24601, 25000}, {28401, 28600}}
	if len(gaps) != len(wantGaps) {
		t.Fatalf("gaps: got %+v, want ranks %v", gaps, wantGaps)
	}
	for i, want := range wantGaps {
		if gaps[i].RankFirst != want[0] || gaps[i].RankLast != want[1] {
			t.Errorf("gap %v: got ranks %v-%v, want %v-%v", i, gaps[i].RankFirst, gaps[i].RankLast, want[0], want[1])
		}
	}

	// WHEN the daily summary is generated
	if err := loopProcessAllZipFiles(goCodeDir); err != nil {
		t.Fatalf("error loopProcessAllZipFiles: %v", err)
	}
	counts := readSummary(t, goCodeDir, "2026-09-25")

	// THEN every ranked player is counted
	total := 0
	for _, c := range counts {
		total += c
	}
	if total != 46168 {
		t.Errorf("total players: got %v, want 46168", total)
	}
	// THEN buckets whose gaps stay inside one bucket are restored exactly:
	// 1100-1200 had 2996 players plus 1800 missing (Elo 1154-1196)
	if counts[1100] != 4796 {
		t.Errorf("bucket 1100-1200: got %v, want 4796", counts[1100])
	}
	// THEN gaps crossing a boundary keep their exact total across the two buckets:
	// 1600-1800 had 861+620 players plus 200 missing (Elo 1681-1706),
	// 900-1100 had 4853+3569 players plus 2000+600+400+200 missing
	if counts[1600]+counts[1700] != 1681 {
		t.Errorf("buckets 1600-1800: got %v, want 1681", counts[1600]+counts[1700])
	}
	if counts[900]+counts[1000] != 11622 {
		t.Errorf("buckets 900-1100: got %v, want 11622", counts[900]+counts[1000])
	}
	// THEN the chart is bell-shaped again, peaking at 1000-1100
	for low, count := range counts {
		if count > counts[1000] {
			t.Errorf("bucket %v (%v players) is higher than the 1000-1100 peak (%v)", low, count, counts[1000])
		}
	}
}

func TestLoopProcessAllZipFiles_RealDataConsecutiveDays(t *testing.T) {
	// GIVEN the real leaderboards of two consecutive days:
	// 2026-09-24 missing 400 players, 2026-09-25 missing 5200 players
	goCodeDir := newGoCodeDirWithFixtures(t, "2026-09-24", "2026-09-25")

	// WHEN both daily summaries are generated
	if err := loopProcessAllZipFiles(goCodeDir); err != nil {
		t.Fatalf("error loopProcessAllZipFiles: %v", err)
	}
	day1 := readSummary(t, goCodeDir, "2026-09-24")
	day2 := readSummary(t, goCodeDir, "2026-09-25")

	// THEN every well-populated bucket is within 3% of the previous day,
	// the rating distribution barely moves in one day,
	// before the fill 1000-1100 differed by 42%
	for low, count1 := range day1 {
		if count1 < 500 {
			continue // small top buckets move a lot in relative terms
		}
		count2 := day2[low]
		diffPercent := math.Abs(float64(count2-count1)) / float64(count1) * 100
		if diffPercent > 3 {
			t.Errorf("bucket %v: 2026-09-24 has %v, 2026-09-25 has %v, differ %.1f%%",
				low, count1, count2, diffPercent)
		}
	}
}

// newGoCodeDirWithFixtures returns a temporary goCodeDir
// whose data_lite holds copies of this repo's data_lite zips of the given dates,
// so tests never write CSVs or charts into the repo.
func newGoCodeDirWithFixtures(t *testing.T, dates ...string) string {
	t.Helper()
	goCodeDir := t.TempDir()
	dataLiteDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_lite")
	if err := os.MkdirAll(dataLiteDir, 0755); err != nil {
		t.Fatalf("error os.MkdirAll: %v", err)
	}
	for _, date := range dates {
		fixture, err := os.ReadFile(filepath.Join("data_lite", "all_players_"+date+".zip"))
		if err != nil {
			t.Fatalf("error os.ReadFile: %v", err)
		}
		err = os.WriteFile(filepath.Join(dataLiteDir, "all_players_"+date+".zip"), fixture, 0644)
		if err != nil {
			t.Fatalf("error os.WriteFile: %v", err)
		}
	}
	return goCodeDir
}

// readSummary reads the bucket counts of a date generated in goCodeDir.
func readSummary(t *testing.T, goCodeDir string, date string) map[int]int {
	t.Helper()
	return readBucketCounts(t, filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_summarized",
		"aoe2_rating_percentile_date_"+date+".csv"))
}

func TestFillRankGaps_ReferenceDay(t *testing.T) {
	// GIVEN today is missing ranks 4 to 7, between Elo 2000 (rank 3) and 1000 (rank 8)
	players := []AoEPlayerLite{
		{RlUserId: 1, Elo: 2600, Rank: 1},
		{RlUserId: 2, Elo: 2300, Rank: 2},
		{RlUserId: 3, Elo: 2000, Rank: 3},
		{RlUserId: 8, Elo: 1000, Rank: 8},
	}
	gaps := detectRankGaps(players)
	// GIVEN a nearby day had those ranks, with most of them close to the lower neighbor
	reference := EloByRank{math.NaN(), 2610, 2310, 2100, 1900, 1300, 1200, 1150, 1100}

	// WHEN the gap is filled from the nearby day
	filled := fillRankGaps(players, gaps, []EloByRank{reference})

	// THEN the missing players keep that day's shape, rescaled to today's neighbors
	var elos []float64
	for _, p := range filled[len(players):] {
		elos = append(elos, p.Elo)
	}
	want := []float64{1800, 1200, 1100, 1050}
	if len(elos) != len(want) {
		t.Fatalf("estimated Elo: got %v, want %v", elos, want)
	}
	for i := range want {
		if elos[i] != want[i] {
			t.Errorf("estimated Elo: got %v, want %v", elos, want)
			break
		}
	}
}

func TestFindReferenceDay(t *testing.T) {
	// GIVEN a gap at ranks 3 to 4 on 2026-07-10,
	// a day before missing rank 4 too, two earlier days with all ranks,
	// and a later day with all ranks, nearer than any complete earlier day
	gap := RankGap{RankFirst: 3, RankLast: 4, EloAbove: 1900, EloBelow: 1600}
	nan := math.NaN()
	days := map[string]EloByRank{
		"2026-07-10": {nan, 2000, 1900, nan, nan, 1600},
		"2026-07-09": {nan, 2000, 1900, 1800, nan, 1600},
		"2026-07-07": {nan, 2000, 1900, 1800, 1700, 1600},
		"2026-07-05": {nan, 2000, 1900, 1750, 1650, 1600},
		"2026-07-11": {nan, 2000, 1900, 1850, 1750, 1600},
	}

	// WHEN the reference day is chosen
	reference := findReferenceDay(gap, "2026-07-10", days)

	// THEN the latest earlier day that has every rank of the gap is used,
	// never a later day, so a chart built on its own date matches a later rebuild
	if reference == nil || reference[3] != 1800 || reference[4] != 1700 {
		t.Errorf("reference: got %v, want the 2026-07-07 day", reference)
	}

	// WHEN no earlier day has every rank of the gap
	reference = findReferenceDay(gap, "2026-07-10", map[string]EloByRank{
		"2026-07-10": days["2026-07-10"],
		"2026-07-09": days["2026-07-09"],
		"2026-07-11": days["2026-07-11"],
	})

	// THEN there is no reference day
	if reference != nil {
		t.Errorf("reference: got %v, want nil", reference)
	}
}

func TestFillRankGaps_CurveOnSteepSide(t *testing.T) {
	// GIVEN a bell-shaped leaderboard (mean 1000, standard deviation 200)
	// missing every player rated 1300 to 1500, on the steep side of the bell,
	// and no nearby day to copy from
	const nPlayers = 10000
	var players []AoEPlayerLite
	for rank := 1; rank <= nPlayers; rank++ {
		quantile := 1 - (float64(rank)-0.5)/nPlayers
		elo := math.Round(1000 + 200*math.Sqrt2*math.Erfinv(2*quantile-1))
		if elo > 1300 && elo < 1500 {
			continue
		}
		players = append(players, AoEPlayerLite{RlUserId: rank, Elo: elo, Rank: rank})
	}
	gaps := detectRankGaps(players)

	// WHEN the gap is filled
	filled := fillRankGaps(players, gaps, nil)

	// THEN like the real bell, clearly more missing players sit near 1300 than near 1500,
	// an even spread would split them about equally
	nBelow1400, nAbove1400 := 0, 0
	for _, p := range filled[len(players):] {
		if p.Elo < 1400 {
			nBelow1400++
		} else {
			nAbove1400++
		}
	}
	if nBelow1400 < 2*nAbove1400 {
		t.Errorf("estimated players: %v below 1400, %v at or above, want at least twice as many below",
			nBelow1400, nAbove1400)
	}
}

func TestNormalQuantile(t *testing.T) {
	// GIVEN a normal distribution with mean 1000 and standard deviation 200
	const mean, stdDev = 1000.0, 200.0

	// WHEN the cumulative probability is taken at a rating, then turned back into a rating
	// THEN the mean splits players in half, and the round trip returns the same rating
	if got := normalCDF(mean, mean, stdDev); math.Abs(got-0.5) > 1e-12 {
		t.Errorf("normalCDF at mean: got %v, want 0.5", got)
	}
	for _, elo := range []float64{400, 900, 1000, 1500, 2200} {
		got := normalQuantile(normalCDF(elo, mean, stdDev), mean, stdDev)
		if math.Abs(got-elo) > 1e-3 { // ratings are whole numbers
			t.Errorf("round trip of %v: got %v", elo, got)
		}
	}
}

func TestLoopProcessAllZipFiles_RealDataHighTail(t *testing.T) {
	// GIVEN the real leaderboard of 2026-07-04, missing ranks 201 to 1400 (Elo 2279 to 1813),
	// and the day before, which has those ranks
	goCodeDir := newGoCodeDirWithFixtures(t, "2026-07-03", "2026-07-04")

	// WHEN both daily summaries are generated
	if err := loopProcessAllZipFiles(goCodeDir); err != nil {
		t.Fatalf("error loopProcessAllZipFiles: %v", err)
	}
	day1 := readSummary(t, goCodeDir, "2026-07-03")
	day2 := readSummary(t, goCodeDir, "2026-07-04")

	// THEN the high-rating buckets keep the real long tail, within 10% of the day before,
	// a bell curve fill put 38 players in 2100-2200 against 144 the day before
	for _, low := range []int{1800, 1900, 2000, 2100, 2200} {
		diffPercent := math.Abs(float64(day2[low]-day1[low])) / float64(day1[low]) * 100
		if diffPercent > 10 {
			t.Errorf("bucket %v: 2026-07-03 has %v, 2026-07-04 has %v, differ %.1f%%",
				low, day1[low], day2[low], diffPercent)
		}
	}
}

func TestGenerateTotalPlayersHTML(t *testing.T) {
	// GIVEN two daily summaries: an older one named "yyyy_mm_dd" without the EstimatedPlayers column,
	// and a newer one with 150 players, 30 of them estimated,
	// plus a stale "yyyy_mm_dd" copy of the older day under the newer day's date
	goCodeDir := t.TempDir()
	summarizedDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "data_summarized")
	if err := os.MkdirAll(summarizedDir, 0755); err != nil {
		t.Fatalf("error os.MkdirAll: %v", err)
	}
	oldCSV := "RatingRange,RatingLow,RatingHigh,CountPlayers,Percentile,RankLow,RankHigh\n" +
		"0000→0100,0,100,40,0.4,#100,#61\n" +
		"0100→0200,100,200,60,1,#60,#1\n"
	newCSV := "RatingRange,RatingLow,RatingHigh,CountPlayers,Percentile,RankLow,RankHigh,EstimatedPlayers\n" +
		"0000→0100,0,100,50,0.333,#150,#101,10\n" +
		"0100→0200,100,200,100,1,#100,#1,20\n"
	files := map[string]string{
		"aoe2_rating_percentile_date_2025_10_30.csv": oldCSV,
		"aoe2_rating_percentile_date_2026-09-25.csv": newCSV,
		"aoe2_rating_percentile_date_2026_09_25.csv": oldCSV,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(summarizedDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("error os.WriteFile: %v", err)
		}
	}

	// WHEN the daily totals are read
	totals, err := readDailyTotals(summarizedDir)
	if err != nil {
		t.Fatalf("error readDailyTotals: %v", err)
	}

	// THEN each day has its total including players the API left out, and how many were left out,
	// oldest first, and the stale copy is ignored
	want := []DailyTotal{
		{Date: "2025-10-30", NPlayers: 100, NEstimated: 0},
		{Date: "2026-09-25", NPlayers: 150, NEstimated: 30},
	}
	if len(totals) != len(want) || totals[0] != want[0] || totals[1] != want[1] {
		t.Errorf("totals: got %+v, want %+v", totals, want)
	}

	// WHEN the total players page is generated
	if err := generateTotalPlayersHTML(goCodeDir); err != nil {
		t.Fatalf("error generateTotalPlayersHTML: %v", err)
	}

	// THEN the page shows the totals by date,
	// and the players returned by the API (120 on 2026-09-25) as a line hidden until clicked,
	// on a count axis starting from 0, as the lowest day has fewer than 1000 players,
	// with a button back to the rating percentile chart
	page, err := os.ReadFile(filepath.Join(goCodeDir, "z_aoe2_rating_percentile", "total_players.html"))
	if err != nil {
		t.Fatalf("error os.ReadFile: %v", err)
	}
	for _, text := range []string{"2025-10-30", "2026-09-25", "100", "150", "120", `"selected":{"Returned by the API":false}`, `"min":0`, `href="index.html"`} {
		if !strings.Contains(string(page), text) {
			t.Errorf("total_players.html: missing %q", text)
		}
	}
}

func TestTotalPlayersAxisMin(t *testing.T) {
	// GIVEN the lowest daily count on the page
	// WHEN the count axis minimum is chosen
	// THEN it is 20000, or lower rounded down to a thousand when a day has fewer players
	cases := []struct{ lowest, want int }{
		{lowest: 44363, want: 20000},
		{lowest: 20000, want: 20000},
		{lowest: 19999, want: 19000},
		{lowest: 4673, want: 4000},
		{lowest: 100, want: 0},
	}
	for _, c := range cases {
		if got := totalPlayersAxisMin(c.lowest); got != c.want {
			t.Errorf("totalPlayersAxisMin(%v): got %v, want %v", c.lowest, got, c.want)
		}
	}
}

func TestDetectRankTail(t *testing.T) {
	// GIVEN the previous day had ranks 1 to 100, rated 1000 down to 10
	previous := make(EloByRank, 101)
	previous[0] = math.NaN()
	for rank := 1; rank <= 100; rank++ {
		previous[rank] = float64(1010 - 10*rank)
	}

	// WHEN today stops at rank 80, short by more than 10%
	tail, isSkipped := detectRankTail(previous[:81], previous)

	// THEN ranks 81 to 99 are a gap between today's worst player and an estimated worst player at rank 100,
	// rated like the previous day's worst
	if isSkipped || tail == nil {
		t.Fatalf("tail: got %+v, isSkipped %v, want a tail", tail, isSkipped)
	}
	wantGap := RankGap{RankFirst: 81, RankLast: 99, EloAbove: 210, EloBelow: 10}
	if tail.Gap != wantGap || tail.Bottom.Rank != 100 || tail.Bottom.Elo != 10 || !tail.Bottom.IsEstimated {
		t.Errorf("tail: got %+v, want gap %+v and an estimated rank 100 rated 10", tail, wantGap)
	}

	// WHEN today stops at rank 95, short by less than 10%
	tail, isSkipped = detectRankTail(previous[:96], previous)

	// THEN today counts as a normal day
	if tail != nil || isSkipped {
		t.Errorf("tail: got %+v, isSkipped %v, want none", tail, isSkipped)
	}

	// WHEN today stops at rank 40, less than half of the previous day
	tail, isSkipped = detectRankTail(previous[:41], previous)

	// THEN today is not the usual leaderboard and is skipped
	if tail != nil || !isSkipped {
		t.Errorf("tail: got %+v, isSkipped %v, want skipped", tail, isSkipped)
	}

	// WHEN there is no previous day
	tail, isSkipped = detectRankTail(previous[:41], nil)

	// THEN nothing can be detected
	if tail != nil || isSkipped {
		t.Errorf("tail: got %+v, isSkipped %v, want none", tail, isSkipped)
	}
}

func TestLoopProcessAllZipFiles_RealDataApiStoppedEarly(t *testing.T) {
	// GIVEN the real leaderboard of 2026-04-09, where the API stopped at rank 40700
	// and returned no player under 500 Elo,
	// and the day before, with last rank 45541
	goCodeDir := newGoCodeDirWithFixtures(t, "2026-04-08", "2026-04-09")

	// WHEN both daily summaries are generated
	if err := loopProcessAllZipFiles(goCodeDir); err != nil {
		t.Fatalf("error loopProcessAllZipFiles: %v", err)
	}
	day1 := readSummary(t, goCodeDir, "2026-04-08")
	day2 := readSummary(t, goCodeDir, "2026-04-09")

	// THEN 2026-04-09 is as long as the day before
	total1, total2 := 0, 0
	for _, count := range day1 {
		total1 += count
	}
	for _, count := range day2 {
		total2 += count
	}
	if total2 != 45541 {
		t.Errorf("2026-04-09 total: got %v, want 45541 like the last rank of 2026-04-08 (which has %v)",
			total2, total1)
	}
	// THEN the low-rating buckets are back, within 10% of the day before
	for _, low := range []int{0, 100, 200, 300, 400} {
		diffPercent := math.Abs(float64(day2[low]-day1[low])) / float64(day1[low]) * 100
		if diffPercent > 10 {
			t.Errorf("bucket %v: 2026-04-08 has %v, 2026-04-09 has %v, differ %.1f%%",
				low, day1[low], day2[low], diffPercent)
		}
	}
}

func TestLoopProcessAllZipFiles_RealDataWrongLeaderboard(t *testing.T) {
	// GIVEN the real leaderboard of 2026-05-01, where the API served only 4673 players
	// with a different top player, and the day before with 47262 players
	goCodeDir := newGoCodeDirWithFixtures(t, "2026-04-30", "2026-05-01")

	// WHEN the daily summaries are generated
	if err := loopProcessAllZipFiles(goCodeDir); err != nil {
		t.Fatalf("error loopProcessAllZipFiles: %v", err)
	}

	// THEN 2026-05-01 gets no summary and no chart, 2026-04-30 still does
	outputDir := filepath.Join(goCodeDir, "z_aoe2_rating_percentile")
	for _, path := range []string{
		filepath.Join(outputDir, "data_summarized", "aoe2_rating_percentile_date_2026-05-01.csv"),
		filepath.Join(outputDir, "output_charts", "chart_2026-05-01.html"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("%v: want no file, got err %v", filepath.Base(path), err)
		}
	}
	if day := readSummary(t, goCodeDir, "2026-04-30"); len(day) == 0 {
		t.Errorf("2026-04-30: want a summary")
	}
}
