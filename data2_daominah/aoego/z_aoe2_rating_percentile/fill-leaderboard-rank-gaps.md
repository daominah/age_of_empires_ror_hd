# Fill leaderboard rank gaps

## What goes wrong

The ageofempires.com leaderboard API sometimes serves an incomplete leaderboard,
and its `Count` shrinks to match.
Our download trusts `Count` and the page responses,
so it saves the incomplete leaderboard without any error.
Across 243 days with a zip (2025-11-01 to 2026-09-26,
pruned zips recovered from the Pages repo git history), plus one leftover CSV:

- Rank gaps: blocks of ranks are missing,
  whole 100-player pages plus smaller gaps that look like ranks shifting while pages download.
  134 days miss at least one rank,
  the worst (2026-06-08 to 2026-06-25) miss about 20000 of 50000 players.
  A hole at ranks 2201 to 2600 persists over many days,
  so even charts that looked normal had a dip at 1600 to 1700.
- API stopped early: 2026-04-09 ends at rank 40700 against 45541 the day before,
  with no player under 500 Elo.
  The ranks have no hole, so gap detection cannot see it.
- Wrong leaderboard: 2026-05-01 has 4673 players
  with a different top player (2722 Elo against 2936 the day before).
  Most likely another leaderboard was returned by mistake,
  we cannot tell whether the API or our request caused it.
- Our download stopped early: 2025-12-23 has only ranks 1 to 7000 (70 pages) against 44363 before.
  A failed page download exits and leaves the earlier pages on disk,
  then a local rerun reuses any existing pages without checking they are complete.
  The workflow always starts from an empty folder, so only local runs are exposed.

## Decisions

Grouped by the kind of bad day they address.

- All kinds:
  - Correct in memory only, between reading a zip and bucketing:
    the zip stays as the API returned it.
  - No re-download when data looks wrong:
    retrying against an already unstable server can look like an attack,
    so fill best-effort from data we already have.
  - Disclose estimates: the chart subtitle says how many players have ratings estimated from earlier days,
    a new last CSV column `EstimatedPlayers` counts them per bucket.
  - Total players page: a cross-check that totals stay close from one day to the next.
    A second line, hidden until clicked in the legend, shows only the players the API returned.
- Rank gaps:
  - Fill from a reference day, not a bell curve:
    the count and Elo range of each gap are exact, only the spread inside is estimated,
    copied from the latest earlier day that has those ranks.
    The bell curve is only the fallback.
  - Fill from earlier days data only: a chart built on its own date gets the same estimates as a later rebuild.
  - Discard rare edge cases: reversed gap neighbors (4 gaps, off by 1 Elo)
    and duplicate ranks from the API are left as is.
- API stopped early:
  - Short day, the player count (read from the last rank)
    drops more than 10% below the previous day's (normal days change less than 3%):
    fill the missing bottom ranks like a gap.
    A drop of less than 10% at the bottom is not detected.
- Wrong leaderboard, and our download stopped early:
  - Too little data, the player count drops more than 50% below the previous day's:
    drop the day, no chart and no CSV, never a reference.
    Filling most of a day from the day before would only copy that day.
