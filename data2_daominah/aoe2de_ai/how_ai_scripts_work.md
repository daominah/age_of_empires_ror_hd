# How the AoE2 DE AI scripts work

How three Age of Empires II: Definitive Edition AI scripts work, side by side:

- **The built-in AI that ships with the game**, played as Extreme, its hardest difficulty.
  Its script files sit in a folder named "Promisory",
  so this doc also calls it Promisory.
- **Immortal**, a Workshop AI by TheMaximalBeing.
- **Rehoboam**, a Workshop AI by The_Night_Panda.

Every fact here was read from the game files and the AI scripts themselves,
cited as `file:line`; forum posts were only hints.
The script files belong to their authors and Microsoft:
this doc cites them by file and line instead of copying them.

## Versions these findings come from

All read on 2026-11-26.
The built-in AI and the lobby handicap change with game patches,
and the Workshop AIs change when their authors upload a new version:
when a version below differs, recheck the findings.

| AI or file                    | Version                                                                 |
|-------------------------------|-------------------------------------------------------------------------|
| Game                          | Steam build 26249352                                                    |
| Built-in AI (Promisory)       | Shipped with that build, its files last changed by the 2026-10-28 patch |
| Lobby handicap (`Effects.xs`) | Shipped with that build, last changed by the 2026-10-28 patch           |
| Immortal                      | v0.10f, compiled 2024-05-11, Workshop upload of 2024-07-25              |
| Rehoboam                      | v1.80k, Workshop upload of 2026-10-11                                   |

- The Workshop upload dates come from the game's mod list:
  they tell when the upload changed, not what changed inside the script.
- Rehoboam's file is named v1.80k, but its header still says v1.80f and 2020-2022.

## Where each AI lives

Paths use the terms of [how_to_create_mod.md](../aoe2de_mod/how_to_create_mod.md):

- `{AoE2DEInstallDir}` is the game's install folder,
  `steamapps\common\AoE2DE` inside the Steam library.
- `{AoE2DEModsDir}` is `%USERPROFILE%\Games\Age of Empires 2 DE\{Steam_ID}\mods\subscribed`,
  the subscribed Workshop mods.
  Our own mods go next to it, under `mods\local`.

| AI                           | Location                                                                            | Shape                                                                                                     |
|------------------------------|-------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------|
| Built-in Extreme (Promisory) | `{AoE2DEInstallDir}\resources\_common\ai\Promisory\`                                | 38 files, about 8,100 rules, hand-written and commented                                                   |
| Immortal v0.10f              | `{AoE2DEModsDir}\264791_Immortal AI 0.10f\resources\_common\ai\Immortal v0d10f.per` | 1 file, 115,901 lines, 3.5 MB, about 9,700 rules                                                          |
| Rehoboam v1.80k              | `{AoE2DEModsDir}\44695_Rehoboam AI\resources\_common\ai\Rehoboam 1.80k.per`         | 1 file, 108,427 lines, 5.2 MB, about 9,300 rules, hand-written and commented, plus 3 micro test scenarios |

- **Extreme is one of Promisory's difficulty levels**, not a separate AI.
- The `.ai` file next to Immortal's `.per` is empty: that is normal,
  the game only needs it to list the AI.
- **"Promisory" is named after its main author**, the AI scripter Promiskuitiv,
  brought in by the Definitive Edition developers.
  This comes from community knowledge: the game files name no author.
  The likely pun on "promissory" (containing a promise) is our guess.
- Immortal's author is TheMaximalBeing,
  original thread `https://forums.aiscripters.com/viewtopic.php?f=8&t=4156`
  (from the Workshop description in `info.json`).
- Rehoboam's author is The_Night_Panda (Workshop name ThNightPanda),
  more information at `aoeaidatabase.pythonanywhere.com/?src=rh`
  (from its `info.json`).

## Script and engine

An AI has two layers:

- **The script** (`.per` files): plain-text rules, "if this, then do that".
  Promisory and Immortal are both scripts.
- **The engine**: built-in AI code compiled into the game program `AoE2DE_s.exe`,
  closed source. It does the detailed work:
  - picking which villager goes to which tree or mine
  - finding a spot for a new building
  - forming and moving attack groups
  - choosing targets
  - pathfinding

A script either steers the engine or bypasses it:

- **Steering** (mostly Promisory): the script sets values called strategic numbers,
  for example "40% of villagers on food",
  and the engine picks the actual villagers and resources.
  Simpler and more stable, but the engine's choices are generic:
  distant resources, poor building spots, poorly timed attacks.
- **Bypassing** (mostly Immortal): the script turns those engine systems off,
  reads single units, and orders exact units to exact map points.
  It beats the engine defaults, at the cost of a lot of code.

## Built-in Extreme AI (Promisory)

- **Extreme is not a separate script**: it is the same rule set
  with the self-weakening turned off and some extra behaviors turned on
  (see "Built-in AI difficulty levels" below).
- **No resource cheats unless asked**: a human typing "200" in chat
  turns on free resources, with the chat line "how do you turn this on"
  (`Promisory\interaction.per:1-60`). Off by default, and nothing in the script turns it on.
  A comment in `Promisory\customConstants.per:1193` says the old cheating option
  "has been obseleted by the handicap system".
- **Mostly steers the engine**: sets gatherer percentages,
  uses the engine's `build` placement and attack logic.
- **Adjusts kiting during the game**
  by how far ahead its army is (`Promisory\init.per:1813-1853`).
- Largest files: `units.per`, `buildings.per`, `tsa.per`, `init.per`, `researches.per`.

## Built-in AI difficulty levels

Promisory reads the difficulty once at game start and sets per-level constants
(`Promisory\customConstants.per:7383-7560`).
The Definitive Edition loads the block under `DE-AVAILABLE` for Hard and Hardest;
Extreme falls through to the last block.
Names as shown in the game:
the script's `DIFFICULTY-EASY` is "Standard" in the lobby.

| Setting                                           | Easiest | Standard | Moderate | Hard | Hardest | Extreme |
|---------------------------------------------------|--------:|---------:|---------:|-----:|--------:|--------:|
| Villager cap, share of the normal 130 villagers   |      5% |      13% |      28% |  53% |    100% |    100% |
| Villager cap at 200 population                    |       7 |       17 |       36 |   69 |     130 |     130 |
| Army cap, share of the space left after villagers |      6% |      14% |      32% |  60% |       ∞ |       ∞ |
| Army cap at 200 population, population slots      |       5 |       10 |       22 |   42 |       ∞ |       ∞ |
| Missile dodging, scale 0 to 100, lower is better  |     100 |      100 |       75 |   50 |      25 |       0 |
| Kiting at start, scale 0 to 100, lower is better  |     100 |      100 |       94 |   88 |      81 |      75 |
| Max Castles                                       |       1 |        2 |        3 |    4 |       ∞ |       ∞ |
| Max Markets                                       |       1 |        1 |        2 |    3 |       ∞ |       ∞ |
| Max Feitorias                                     |       0 |        0 |        1 |    2 |       ∞ |       ∞ |
| Delay before extra production buildings, second   |    4800 |     3600 |     2400 | 1200 |       0 |       0 |
| Self-weakening steps when ahead of a human        |       4 |        3 |        2 |    1 |       0 |       0 |
| Extreme-only behaviors (below)                    |      no |       no |       no |   no |      no |     yes |

- **"Normal" is the villager target for the population cap**
  (`Promisory\customConstants.per:7583-8137`):
  130 villagers at 200 population, 65 at 100, 200 at most.
  Below Hardest the target becomes that share, at least 6 villagers
  (`Promisory\init.per:1218-1228`).
- **The villager cap is enforced**: it stops training villagers
  once the civilian count, fishing ships and trade carts included, reaches the cap
  (`Promisory\units.per:11393-11402`).
  So Easiest grows from its 3 starting villagers to 7 and stops.
- **Easiest through Moderate also mirror the human**:
  after 12 villagers, they stop training more
  once they have more civilians than the player they target
  (`Promisory\units.per:11456-11467`).
- **The army cap, Easiest to Hard only**: a share of what the population cap
  leaves after villagers, 200 minus 130 is 70 at 200 population,
  at least 5 (`Promisory\init.per:1229-1248`).
  It counts military population, so a unit that takes 2 slots counts 2.
  Once reached, it stops training military units (`Promisory\units.per:12011-12030`).
- **Hardest and Extreme have no army cap**:
  it starts unlimited (`Promisory\init.per:211`) and is never cut,
  so only the game's population cap limits them.
- **Easiest to Moderate also mirror the biggest army** early in the game,
  before their delay runs out and while not attacking or under attack
  (`Promisory\units.per:11987-12010`).
- "Lower is better" is read from the ordering:
  Easiest gets the highest value, Extreme the lowest.
  The scripts do not document the scale.
- **Kiting gets sloppier when the AI is winning** (`Promisory\init.per:1813-1853`).
  The script rates its army against the enemy's on a scale from -3 to +4
  (`Promisory\threats.per:1272-1400`), then raises the kiting value:
  - ahead by 2, or even and attacking: 81
  - ahead by 3 and attacking: 88
  - ahead by 4 and attacking: 94
  - kiting off (100) at max population, or with only fire and ram ships

  It never goes below the difficulty's starting value,
  so Hardest stays at 81 or worse, and only Extreme reaches 75.
- Easiest and Standard also wait for a human to age up first,
  and drop any custom strategy for a plain one;
  Moderate drops only Feudal Age rush strategies, after its delay
  (`Promisory\researches.per:35-46`, `Promisory\init.per:9582-9606`).

**Hardest against Extreme**: the economy and army caps are the same.
Extreme differs by:

- **Better missile dodging and kiting** (0 and 75 against 25 and 81).
- **Extreme-only behaviors**, enabled by the flag `diff-fp`,
  which is 1 only on Extreme in the Definitive Edition (about 150 rules):
  - The newer attack system, water micro, group micro, melee micro,
    early militia rush micro, scout micro (`Promisory\tsa.per:56`,
    `tsa.per:7256`, `tsa.per:10660`, `init.per:8670`, `scoutcontrol.per:14`).
  - Handling enemy mangonels and Town Centers in fights (`Promisory\tsa.per:12709-12762`).
  - Guarding builders of Docks with the scout,
    instead of researching Loom early to resist harassment
    (`Promisory\orb.per:43-74`, `Promisory\researches.per:2794`).
  - Palisade and stone wall placement with gates (`Promisory\buildings.per:2559-2752`).
  - Finer villager control: smarter food and wood splits
    (non-Extreme uses fixed splits, `Promisory\gatherers.per:2189-2222`),
    straggler tree handling, a refined boar lure, scouting by direct unit control
    (`Promisory\gatherers.per:5151`, `boarhunting.per:267`, `general.per:6602`).
  - Repairing damaged Castles, Wonders and Town Centers:
    every 3 seconds, up to 8 nearby villagers go repair
    (`Promisory\general.per:6480-6571`).
  - Placing buildings in a line near the forest and resource camps,
    so they wall off part of the base (`Promisory\extremebuildings3.per`,
    loaded at `Promisory\buildings.per:11451`).
  - In Death Match, a separate plan for military buildings and Houses,
    up to 6 Barracks and 6 Archery Ranges (`Promisory\extremebuildings4.per`).

## How the built-in AI picks its army

Mostly by countering what the enemy it targets has on the field,
counted every few seconds (`Promisory\threats.per:569-603`):

- **Enemy archers → skirmishers**, about 2 per archer:
  3 archers allow up to 6 skirmishers, 10 allow 20, 40 allow 60
  (`Promisory\units.per:5306-5390`).
- **Enemy infantry → archers**, unless the enemy already has
  many skirmishers, Huskarls or siege (`Promisory\units.per:5745-6747`).
- **Enemy cavalry → pikemen and camels** (`Promisory\units.per:4747`, `units.per:7275`).
- **Out of gold in Imperial Age → skirmishers, pikemen and hussars**,
  the units that cost no gold (`Promisory\units.per:8517`).

**Cavalry archers are rare**, though they are among the strongest units,
fast and ranged. The script trains them only in narrow cases (`Promisory\units.per:5808-6695`):

- **A cavalry archer rush**: weighted high only for Huns, Mongols and Tatars,
  and even then only while its army is behind and it has spare wood and gold.
- **As a counter**: to enemy infantry, capped at 15 cavalry archers
  until the enemy has 40 infantry, or to Teutonic Knights and war elephants.
- **As the default unit** only for Huns and Tatars (`Promisory\customConstants.per`).
  Most other civs default to their unique unit, a few to Champions, Knights or Mangonels.

**Hussars are rare too**, though they cost no gold and are strong late in the game.
The script treats them as a fallback, capped at 12 to 25 of the scout line
(`Promisory\units.per:6747-6997`):

- **When it is short of gold**: spare gold under about 200, with spare food.
- **Only while the enemy has few pikemen and camels**: as few as 5 camels stop it.
- **Against enemy siege**, while it has few cavalry of its own.

Each civ also weighs its opening strategies (`Promisory\customConstants.per`):
across the civs, the Knight rush weighs the most, not an archer strategy.
So an Extreme that makes skirmishers and archers is mostly answering
an army of archers and infantry, or running out of gold.

## Built-in AI self-weakening, Hard and easier only

Skip this section if you only care about Extreme:
Hardest and Extreme never weaken themselves.

On Hard and the easier levels, Promisory changes its own stats mid-game
when it is clearly ahead of a human enemy (`Promisory\init.per:2113-2154`):

- **Villagers work 15% slower.**
- **Buildings and siege fire 15% slower**, including Town Center and Castle arrows.
- **Permanent and stacking**: each trigger multiplies on top of the last,
  up to 4 times on Easiest, 3 on Standard, 2 on Moderate and 1 on Hard.
  At 4 times, villagers work at about half speed.

## When the built-in AI resigns

Rules in `Promisory\resign.per:107-381`.
They compare total population, villagers and army together, never the villager count alone:
`min-number-vills` (65 at 200 population) is not used there.

- **Never while its population is about a third of the cap or more**
  on Hardest and Extreme: about 67 at 200 population (`resign.per:113-131`).
- **Never while it or an ally holds all relics or a Wonder** (`resign.per:143-151`).
- **Then, against a human with no ally to lean on**, it resigns when, for example:
  - its population is under 30 and the enemy has 75 or more;
  - its population is under 60, the enemy has 120 or more, and its army is clearly weaker;
  - its population is under 80, the enemy has 144 or more, and its army is clearly weaker;
  - its population is under 60 and the enemy army alone is 60 or more.
- **It also resigns with nothing left to rebuild**:
  no villager, no monk, no Town Center, no Monastery (`resign.per:152-162`).
- **Chat codes change its resigning**: a human types the number in chat,
  like a taunt, and the script reads it (`resign.per:711-735`):
  - "104": "Don't resign!"
  - "105": "You can resign again."
  - "205": no taunt text or sound in the game.
    It makes the AI resign sooner: roughly once its population
    is about half the enemy's or less (`resign.per:750-825`).
- **One second chance against humans**: the first time it decides to resign
  before 42 minutes, while not under attack or with an ally,
  it holds off for 32 seconds and checks again (`resign.per:404-430`).
- **On resigning** it sends a random chat line, then resigns 6 seconds later
  (`resign.per:431-438`).

Before the population gap rules above can fire, all of these must hold
(`resign.per:113-200`):

- **Its population is under about 67** (villagers plus army).
- **It has lost at least 80 units and buildings** to enemies.
- **It is under attack right now.**
- **Its army is rated far behind**, on its own scale from -3 to +4:
  the lowest step, -3, for itself and for its team.

## Lobby handicap

The lobby handicap is a game setting per player, separate from any AI script.
It boosts the player it is set on.
**Confirmed from the game's own script**
`{AoE2DEInstallDir}\resources\_common\xs\Effects.xs:139-203`:

- **Applied once, at game start**, by the "Dark Age effect" (`Effects.xs:257-261`).
  It multiplies the base stats, so later flat upgrades like Loom's +15 hit points are not multiplied.
- **Does nothing at 100%**, the lowest value the lobby allows (`Effects.xs:143`).
- **Per player**: it reads the multiplier for each player from the lobby,
  so an AI slot gets it only if set on that slot.

At 110%, everything below is × 1.1:

| Stat                | Applies to                                                                                |
|---------------------|-------------------------------------------------------------------------------------------|
| Hit points          | Buildings, walls, gates, towers, Farms, villagers, trade carts, trade cogs, fishing ships |
| Work rate           | Villagers (gathering, building, repairing), Farms, trade carts, trade cogs                |
| Carry capacity      | Villagers, fishing ships                                                                  |
| Gathering bonuses   | Food, wood, gold, stone, fishing, shepherding, hunting, foraging                          |
| Training time ÷ 1.1 | Military units, monks and warships: about 9% faster. Villagers train at normal speed      |

**Not** changed by this handicap setting script:

- **Military unit on-field stats** (hit points, attack, bonus damage, armor, speed).
- **Research time** and **villager training time**.
- **Starting resources**: the lobby tooltip says "increasing starting resources", which is outdated.

Promisory's comment that cheating was "obseleted by the handicap system"
(`Promisory\customConstants.per:1193`) fits this:
a handicap on an AI slot is the official way to give it a boost.
Promisory detects a handicap on itself
by checking whether its villagers have more than 25 hit points,
then researches Loom early (`Promisory\researches.per:2801-2819`).

So at 110% on a human player's slot against Extreme:

- **That player's economy runs at least 10% faster**,
  with 10% tougher buildings and villagers, and about 9% faster military training.
- **The Extreme AI gets no boost, and it does not weaken itself.**
- **Neither AI script reacts to the human's handicap**:
  Promisory only checks its own villagers,
  Immortal has no handicap or difficulty checks at all.

## Immortal

- **Compiled from a higher-level language, source not included.**
  The file says "This script was compiled on: 2024-05-11" (line 48634).
  Signs:
  - Generated sections: constants (line 1), goals (line 2233),
    function names (line 17174), script (line 17231).
  - Machine variable names like `varTmp5547`,
    and return-address goals like `___ra1___`.
  - 3,236 jumps to exact rule numbers (Promisory has 21):
    inserting or deleting a rule shifts every target after it.
- **Mostly bypasses the engine**:
  - Sets the engine's civilian gatherer percentage to 0
    and assigns each villager itself.
  - Sets the engine's attack groups to 0
    and runs its own unit control per unit type.
  - Places Farms, camps, towers and Castles with its own code.
- **Ignores the difficulty setting**: no difficulty checks,
  no self-weakening, no cheats. It always plays at full strength.
- **About 30 named build orders**,
  for example "FF->Scout->Skirm" (fast Feudal, scouts, then skirmishers)
  or "FC->Knight" (fast Castle into knights).
  It picks one with random weights per civ and announces it in chat.
- **Uses spare strategic number slots as memory**:
  unit counts every loop, the chosen strategy,
  shared with allies along with "change strategy" chat messages.

## Rehoboam

What its own description says (`info.json` and the file header, lines 1-15),
not checked by play:

- **Narrow by design**: tuned for Arabia, 200 population, low resources, Dark Age start,
  Extreme difficulty, game speed 2.0 or lower ("speedhack will severly weaken micro").
- **Only a few civs**: in 1v1, its own slot must be Persians (preferred) or Franks;
  Teutons pass the check too. Mayans only on the flank in 3v3 and 4v4.
  The human's civ is free.
  On another civ it still runs, but warns in chat
  "Wrong civ selected, this will lead to bugs" (`Rehoboam 1.80k.per:4189-4225`).
  Typing "199" in chat lists its supported settings (lines 3649-3665).
- **Fixed styles**: Feudal flush into knights in 1v1 and 2v2;
  in team games, flanks rush archers and pockets go fast Castle into knights.
- **Claimed strength**: around 1150 Elo in 1v1 on the Definitive Edition,
  beats every non-cheating AI on Arabia except Immortal, "where games can go either way".
- **Known gaps**: does not handle walls or forward buildings.

What the file shows:

- **Hand-written, not compiled**: an indented outline of comments,
  partly in Danish ("Fase" for phase, "ny" for new).
  It loads no other file.
- **Built as four phases**, each a block of rules that the others jump over
  (lines 31881, 62680, 72142, 81219):

| Phase | Starts at line | What it does                                                                                                                                                          |
|-------|---------------:|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1     |          31881 | Fixed Dark Age opening, villagers under direct control: build orders, Houses, scouting, boar, sheep, deer luring, berries, Farms, wood, gold, anti militia rush micro |
| 2     |          62680 | Chosen Feudal strategy: fast Castle with 2 Town Centers into knights, or an archer and skirmisher rush                                                                |
| 3     |          72142 | Mostly for answering a tower rush: full economy, buildings, military, scouting                                                                                        |
| 4     |          81219 | Rest of the game: the 1v1 plan                                                                                                                                        |

- **Strategy is picked by team position**, not by randomness:
  flank or pocket in team games, a separate plan in 1v1 and 2v2 (line 8246).
- **Its own unit control**: ranged micro from line 11729 and melee micro from line 22465,
  about 14,000 lines together, plus scout micro and making Town Centers fire.
- **Lag handling**: it measures how long each pass takes,
  spreads heavy code across passes, and skips work late in the game when slow
  (lines 2429-2551, 8611).
- **A cheat that is off by default** (lines 8719-8860):
  anyone typing "201" in chat turns it on, with a chat line saying so;
  "202" turns it off. Nothing in the script turns it on by itself.
  While on, every 15 seconds it adds 1 of each resource per villager gathering it,
  and clicks Feudal Age at 22 villagers.
- **Difficulty**: plays fully only on Hard and above;
  below Hard it turns off the training queue.
- **Mixes both approaches**: turns engine gathering and attack groups off in places
  for direct control, but still tunes many engine settings.

## Command usage, side by side

Counts of non-comment lines using each command, a rough picture of each design.

| Command                | Promisory | Immortal | Rehoboam | Reading                                         |
|------------------------|----------:|---------:|---------:|-------------------------------------------------|
| `up-jump-direct`       |        21 |    3,236 |      170 | Immortal: compiled function calls and loops     |
| `up-jump-rule`         |     1,903 |      145 |      506 | Promisory and Rehoboam: hand-written skips      |
| `up-get-object-data`   |       225 |    1,233 |      687 | Immortal reads single units most, Rehoboam next |
| `set-strategic-number` |     1,626 |        0 |    1,109 | Promisory and Rehoboam tune engine settings     |
| `research`             |     3,276 |      151 |    1,619 | Promisory spells out research per case          |
| `build`                |     6,035 |    1,106 |    4,517 | Promisory spells out building per case          |

## Built-in Extreme against Immortal against Rehoboam

All three at Extreme difficulty, as installed on 2026-11-26.
"Claimed" rows come from the AI's own description, not checked by play.

🟢 good, 🟡 less good, 🔴 clearly worse than the others;
no color where the row only describes or was not checked.

| Question                    | Built-in Extreme (Promisory)                            | Immortal v0.10f                                             | Rehoboam v1.80k                                                       |
|-----------------------------|---------------------------------------------------------|-------------------------------------------------------------|-----------------------------------------------------------------------|
| Author                      | Promiskuitiv, shipped with the game                     | TheMaximalBeing                                             | The_Night_Panda                                                       |
| Last changed                | 🟢 2026-10-28, changes with game patches                | 🔴 Compiled 2024-05-11                                      | 🟡 Workshop 2026-10-11; header says 2020-2022                         |
| Size                        | 38 files, about 8,100 rules                             | 1 file, about 9,700 rules                                   | 1 file, about 9,300 rules                                             |
| Readable source             | 🟢 Yes, hand-written, English comments                  | 🔴 No, compiler output, source not included                 | 🟢 Yes, hand-written, comments partly Danish                          |
| Safe to edit                | 🟢 Yes                                                  | 🔴 Small value changes only: rule jumps are hard-coded      | 🟢 Yes                                                                |
| Game variety                | 🟢 High: plan per civ and map, army counters yours      | 🟢 High: about 30 build orders, picked at random            | 🔴 Low: one 1v1 plan, played as Persians or Franks                    |
| Civs                        | 🟢 All                                                  | 🟡 All, but Incas break and post-2024 civs get generic play | 🔴 Persians or Franks (Teutons pass the check too)                    |
| Maps                        | 🟢 All, land and water                                  | 🟢 All, with water build orders                             | 🟡 Arabia and similar land maps                                       |
| Strategy choice             | 🟡 Rule tree per civ and map                            | 🟢 About 30 build orders, random weights per civ            | 🔴 Fixed: one plan for 1v1                                            |
| Typical style               | Varies by civ                                           | Varies: Feudal rushes, fast Castle, fast Imperial           | Feudal flush into knights                                             |
| Villager control            | 🟡 Engine, steered by percentages                       | 🟢 Direct, per villager                                     | 🟡 Direct in the opening, mixed later                                 |
| Building placement          | 🔴 Mostly the engine: Castles often useless at the back | 🟢 Own code for Farms, camps, towers, Castles               | 🟡 Own build-line code for the opening buildings                      |
| Army control                | 🟡 Engine attack groups plus Extreme-only micro         | 🟢 Own control per unit type: melee, archer, siege, monk    | 🟢 Own ranged and melee micro, scout micro                            |
| Missile dodging and kiting  | 🟡 Engine settings at their best on Extreme             | 🟢 Own micro                                                | 🟢 Own micro, weaker above game speed 2.0                             |
| Walls and forward buildings | 🟢 Builds walls with gates                              | Not checked                                                 | 🔴 Claimed: cannot play against them                                  |
| Lag handling                | Not checked                                             | Not checked                                                 | 🟢 Measures pass time, spreads heavy code, skips work when slow       |
| Reacts to difficulty        | Yes, 6 levels                                           | No, always full strength                                    | Full strength on Hard and above                                       |
| Cheats                      | None by default; "200" in chat turns them on            | None                                                        | None by default; "201" in chat turns them on                          |
| Claimed strength            | Not stated                                              | "Really powerful", plays all civs and settings              | About 1150 Elo in 1v1, even with Immortal, beats the others on Arabia |

## Limitations and bugs

**Built-in Extreme (Promisory)**:

- **Changes with every game patch**: findings hold only for the game build
  in "Versions these findings come from".
- **Predictable army**: it mostly counters what the enemy fields
  (see "How the built-in AI picks its army").
- **Kites worse when winning** (see "Built-in AI difficulty levels").
- **Relies on the engine** for gathering, placement and attacks: generic choices.
- **Poor Castle placement** (seen in games):
  its Castles often end up useless at the back of its base.
  The script mostly leaves the spot to the engine,
  placing around its own base (`Promisory\buildings.per:10367-10442`, `buildings.per:11510`),
  and only one rule asks for a forward Castle (`buildings.per:10419`).
  Immortal, which places Castles with its own code, puts them in much better spots.

**Immortal v0.10f**:

- **Cannot play Incas** since the Settlement patch (seen in games):
  Incas now drop food, wood and gold at a Settlement
  instead of a Mill and camps (`Promisory\customConstants.per:3759-3813`),
  but Immortal never mentions the Settlement.
  It keeps trying to build Lumber Camps, Mining Camps and Mills, which Incas cannot build,
  so its economy gets stuck.
- **Does not know the civs released after May 2024**: it plays them with only its generic logic.
  They are Achaemenids, Spartans, Jurchens, Khitans, Shu, Wei, Wu, Mapuche, Muisca and Tupi;
  a Full Random AI can land on any of them that the game offers.
  Romans, from the Return of Rome DLC (2023), are covered:
  14 Roman checks, plus the Centurion and Legionary.
- **Frozen**: compiled on 2024-05-11 with no newer release,
  and only small value changes are safe to edit (see "Immortal").

**Rehoboam v1.80k**:

- **Few civs and one style**: Persians or Franks in 1v1, Feudal flush into knights;
  on other civs it warns "Wrong civ selected, this will lead to bugs" (see "Rehoboam").
- **Breaks against walls and forward buildings**: warns in chat when it sees walls.
- **Arabia-like maps only**, with Conquest, low resources and Dark Age start.
- **Micro needs a fast computer and a slow game**:
  a benchmark score of 1100 or more, game speed 2.0 or lower.
- **Full Random civ does not fit**: a script cannot pick its own civ,
  so a random slot lands on Persians, Franks or Teutons only about 3 times in 50.
- **Unclear what its latest update changed**: the Workshop upload changed on 2026-10-11,
  but the file header still says v1.80f and 2020-2022.

## Reading list for Immortal

For Immortal, read in this order:

- **Function index**: `Immortal v0d10f.per:17174-17230`,
  the 56 function names with their start rule numbers.
- **Main loop**: rule 3573 at line 48605,
  the body run every pass, which calls the functions.
- **Unit counting per loop**: rule 3629 around line 49593.
- **Strategy picker**: rule 4263 around line 56226,
  where the chosen strategy is announced.
- **Functions, by start line**:

| Function                       | Start rule |  Line |
|--------------------------------|-----------:|------:|
| Find enemy units in town       |          1 | 17237 |
| Defend ally                    |        300 | 20516 |
| Initial build (DUC)            |        351 | 21011 |
| Smart ungarrison               |        401 | 21399 |
| Place building                 |        473 | 22032 |
| Train (DUC)                    |        569 | 22807 |
| Load build order               |       1048 | 26844 |
| Micro: unit tag info           |       1108 | 27421 |
| Micro: melee state machine     |       1237 | 28449 |
| Micro: get data                |       1562 | 31344 |
| Micro: archer state machine    |       2624 | 40813 |
| Micro: trebuchet state machine |       2797 | 42391 |
| Micro: monk state machine      |       3065 | 44542 |
| Micro: mangonel state machine  |       3195 | 45579 |
| Shift gatherers closer         |       3307 | 46457 |
| Manual villager reassign       |       3408 | 47320 |
| Villager hopping               |       3519 | 48176 |

- **Promisory, for game knowledge** (civs, maps, counter units):
  `units.per`, `researches.per`, `tsa.per`,
  and `init.per:1800-2160` for the difficulty logic.
- **The forum thread** in "Where each AI lives", for the source or the compiler name.

## Sources

- `{AoE2DEInstallDir}\resources\_common\ai\Promisory\`: the built-in AI.
- `{AoE2DEInstallDir}\resources\_common\xs\Effects.xs`: the lobby handicap.
- The Immortal and Rehoboam `.per` files and their Workshop `info.json`, see "Where each AI lives".
