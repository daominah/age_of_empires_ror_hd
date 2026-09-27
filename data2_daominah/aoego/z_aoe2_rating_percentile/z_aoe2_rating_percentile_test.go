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
