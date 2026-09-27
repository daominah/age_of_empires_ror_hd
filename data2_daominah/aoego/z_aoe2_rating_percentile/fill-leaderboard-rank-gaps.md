# Fill leaderboard rank gaps

## What goes wrong

The ageofempires.com leaderboard API sometimes serves an incomplete leaderboard,
and its `Count` shrinks to match.
Our download trusts `Count` and the page responses,
so it saves the incomplete leaderboard without any error.
Across 243 days with a zip (2025-11-01 to 2026-09-26,
pruned zips recovered from the Pages repo git history):

- Rank gaps: blocks of ranks are missing,
  whole 100-player pages plus smaller gaps that look like ranks shifting while pages download.
  134 days miss at least one rank,
  the worst (2026-06-08 to 2026-06-25) miss about 20000 of 50000 players.
  A hole at ranks 2201 to 2600 persists over many days,
  so even charts that looked normal had a dip at 1600 to 1700.

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
- Rank gaps:
  - Fill from a reference day, not a bell curve:
    the count and Elo range of each gap are exact, only the spread inside is estimated,
    copied from the latest earlier day that has those ranks.
    The bell curve is only the fallback.
  - Fill from earlier days data only: a chart built on its own date gets the same estimates as a later rebuild.
  - Discard rare edge cases: reversed gap neighbors (4 gaps, off by 1 Elo)
    and duplicate ranks from the API are left as is.
