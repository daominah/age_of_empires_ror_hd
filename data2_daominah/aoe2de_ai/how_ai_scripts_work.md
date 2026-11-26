# How the AoE2 DE AI scripts work

How three Age of Empires II: Definitive Edition AI scripts work, side by side:

- **The built-in Extreme AI**, the official AI that ships with the game.
- **Immortal**, a mod AI by TheMaximalBeing.
- **Rehoboam**, a mod AI by ThNightPanda.

Both mod AIs can be installed from the game's Browse Mods.

Every fact here was read from the game files and the AI scripts themselves;
forum posts were only hints.
The script files belong to their authors and Microsoft:
this doc cites them instead of copying them.
Each citation gives a name to search for, like a goal, a rule number or a section comment,
then the line numbers as of the versions in "Game and AI versions":
the built-in AI's lines move with every game patch, but its names mostly stay.

## Game and AI versions

All read on 2026-11-26.
The built-in AI and the lobby handicap change with game patches,
and the mod AIs change when their authors upload a new version:
when a version below differs, recheck the findings.

| AI or file                    | Version                                                                 |
|-------------------------------|-------------------------------------------------------------------------|
| Game                          | Steam build 26249352                                                    |
| Built-in AI                   | Shipped with that build, its files last changed by the 2026-10-28 patch |
| Lobby handicap (`Effects.xs`) | Shipped with that build, last changed by the 2026-10-28 patch           |
| Immortal                      | v0.10f, compiled 2024-05-11, mod upload of 2024-07-25                   |
| Immortal, newer               | v1.2c, compiled 2026-11-04, forum attachment only, not in Browse Mods   |
| Rehoboam                      | v1.80k, mod upload of 2026-10-11                                        |

- The mod upload dates come from the game's mod list:
  they tell when the upload changed, not what changed inside the script.
- Rehoboam's file is named v1.80k, but its header still says v1.80f and 2020-2022.

## Where AI files are

Paths use the terms of [how_to_create_mod.md](../aoe2de_mod/how_to_create_mod.md):

- `{AoE2DEInstallDir}` is the game's install folder,
  `steamapps\common\AoE2DE` inside the Steam library.
- `{AoE2DEModsDir}` is `%USERPROFILE%\Games\Age of Empires 2 DE\{Steam_ID}\mods\subscribed`,
  the subscribed mods.
  Our own mods go next to it, under `mods\local`.

| AI                  | Location                                                                            | Shape                                                                                                     |
|---------------------|-------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------|
| Built-in Extreme AI | `{AoE2DEInstallDir}\resources\_common\ai\Promisory\`                                | 38 files, about 8,100 rules, hand-written and commented                                                   |
| Immortal v0.10f     | `{AoE2DEModsDir}\264791_Immortal AI 0.10f\resources\_common\ai\Immortal v0d10f.per` | 1 file, 115,901 lines, 3.5 MB, about 9,700 rules                                                          |
| Immortal v1.2c      | Forum attachment `Immortal 1.2c.zip`, unpacks to `Immortal 1_2.per`                 | 1 file, 124,263 lines, 3.9 MB, about 9,000 rules                                                          |
| Rehoboam v1.80k     | `{AoE2DEModsDir}\44695_Rehoboam AI\resources\_common\ai\Rehoboam 1.80k.per`         | 1 file, 108,427 lines, 5.2 MB, about 9,300 rules, hand-written and commented, plus 3 micro test scenarios |

- **Extreme is one of the built-in AI's difficulty levels**, not a separate AI.
- The `.ai` file next to Immortal's `.per` is empty: that is normal,
  the game only needs it to list the AI.
- **The folder "Promisory" is named after its main author**, the AI scripter Promiskuitiv,
  brought in by the Definitive Edition developers.
  This comes from community knowledge: the game files name no author.
  The likely pun on "promissory" (containing a promise) is our guess.
- Immortal's author is TheMaximalBeing,
  original thread `https://forums.aiscripters.com/viewtopic.php?f=8&t=4156`
  (from the mod description in `info.json`).
- Rehoboam's author is ThNightPanda (The_Night_Panda in the script header),
  more information at `aoeaidatabase.pythonanywhere.com/?src=rh`
  (from its `info.json`).

## Built-in Extreme AI

- **Extreme is not a separate script**: it is the same rule set
  without weakening itself, and with some extra behaviors turned on
  (see "Difficulty levels" below).
- **No resource cheats unless asked**: a human typing "200" in chat
  turns on free resources, with the chat line "how do you turn this on"
  (goal `cheating`, `Promisory\interaction.per:1-60`).
  Off by default, and nothing in the script turns it on.
  A comment on the `cheating` goal (`Promisory\customConstants.per:1193`)
  says the old cheating option "has been obseleted by the handicap system".
- **Mostly steers the engine**: sets gatherer percentages,
  uses the engine's `build` placement and attack logic.
- **Adjusts kiting during the game**
  by how far ahead its army is
  (rules on `ability-to-kite`, `Promisory\init.per:1813-1853`).
- Largest files: `units.per`, `buildings.per`, `tsa.per`, `init.per`, `researches.per`.

### Difficulty levels

The built-in AI reads the difficulty once at game start and sets per-level constants
(the `#load-if-defined DIFFICULTY-...` blocks, `Promisory\customConstants.per:7383-7560`).
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
| Times it weakens itself when ahead of a human     |       4 |        3 |        2 |    1 |       0 |       0 |
| Extreme-only behaviors (below)                    |      no |       no |       no |   no |      no |     yes |

- **"Normal" is the villager target for the population cap**
  (constant `up-max-civ`, `Promisory\customConstants.per:7583-8137`):
  130 villagers at 200 population, 65 at 100, 200 at most.
  Below Hardest the target becomes that share, at least 6 villagers
  (constant `civilian-percent`, `Promisory\init.per:1218-1228`).
- **The villager cap is enforced**: it stops training villagers
  once its civilian count reaches the cap
  (goal `trainvillager`, `Promisory\units.per:11393-11402`).
  So Easiest grows from its 3 starting villagers to 7 and stops.
- **Easiest through Moderate also mirror the human**:
  after 12 villagers, they stop training more
  once they have more civilians than the player they target
  (goal `my-cpop`, `Promisory\units.per:11456-11467`).
- **The army cap, Easiest to Hard only**: a share of what the population cap
  leaves after villagers, 200 minus 130 is 70 at 200 population,
  at least 5 (goal `allowed-num-military`, `Promisory\init.per:1229-1248`).
  It counts military population, so a unit that takes 2 slots counts 2.
  Once reached, it stops training military units
  (goal `my-mpop`, `Promisory\units.per:12011-12030`).
- **Hardest and Extreme have no army cap**:
  it starts unlimited (`allowed-num-military 16384`, `Promisory\init.per:211`) and is never cut,
  so only the game's population cap limits them.
- **Easiest to Moderate also mirror the biggest army** early in the game,
  before their delay runs out and while not attacking or under attack
  (goal `allowed-num-military`, `Promisory\units.per:11987-12010`).
- **Hard and easier weaken themselves** when clearly ahead of a human,
  as many times as the table says; Hardest and Extreme never do
  (goal `difficulty-adjustment`, `Promisory\init.per:2113-2154`):
  - villagers work 15% slower;
  - buildings and siege fire 15% slower, Town Center and Castle arrows included;
  - each time is permanent and stacks on the last:
    at 4 times, on Easiest, villagers work at about half speed.
- "Lower is better" is read from the ordering:
  Easiest gets the highest value, Extreme the lowest.
  The scripts do not document the scale.
- **Kiting gets sloppier when the AI is winning**
  (rules on `ability-to-kite`, `Promisory\init.per:1813-1853`).
  The script rates its army against the enemy's on a scale from -3 to +4
  (section "Military Superiority", `Promisory\threats.per:1272-1400`),
  then raises the kiting value:
  - ahead by 2, or even and attacking: 81
  - ahead by 3 and attacking: 88
  - ahead by 4 and attacking: 94
  - kiting off (100) at max population

  It never goes below the difficulty's starting value,
  so Hardest stays at 81 or worse, and only Extreme reaches 75.
- Easiest and Standard also wait for a human to age up first,
  and drop any custom strategy for a plain one;
  Moderate drops only Feudal Age rush strategies, after its delay
  (section "Age advancement", `Promisory\researches.per:35-46`;
  goal `strategylocked`, `Promisory\init.per:9582-9606`).

**Hardest vs Extreme**: the economy and army caps are the same.
Extreme differs by:

- **Better missile dodging and kiting** (0 and 75 vs 25 and 81).
- **Extreme-only behaviors**, enabled by the flag `diff-fp`,
  which is 1 only on Extreme in the Definitive Edition (about 150 rules):
  - The newer attack system, group micro, melee micro,
    early militia rush micro, scout micro:
    goals `attackingsystem`, `group-micro`, `melee-micro`, `drushmicro`
    and `modern-scout-micro` (`Promisory\tsa.per:56`,
    `tsa.per:7256`, `tsa.per:10660`, `init.per:8670`, `scoutcontrol.per:14`).
  - Handling enemy mangonels and Town Centers in fights
    (section "Find enemy mangonels", `Promisory\tsa.per:12709-12762`).
  - Palisade and stone wall placement with gates
    (section "Center walling", `Promisory\buildings.per:2559-2752`).
  - Finer villager control: smarter food and wood splits
    (non-Extreme uses fixed `sn-...-gatherer-percentage` splits,
    `Promisory\gatherers.per:2189-2222`),
    straggler tree handling, a refined boar lure, scouting by direct unit control
    (comment "initial straggler chopper", `Promisory\gatherers.per:5151`;
    goal `boar-hunting-stage`, `boarhunting.per:267`;
    goal `ducscouting`, `general.per:6602`).
  - Repairing damaged Castles, Wonders and Town Centers:
    every 3 seconds, up to 8 nearby villagers go repair
    (section "Expedient repairs of key buildings", `Promisory\general.per:6480-6571`).
  - Placing buildings in a line near the forest and resource camps,
    so they wall off part of the base (`Promisory\extremebuildings3.per`,
    loaded at `Promisory\buildings.per:11451`).
  - In Death Match, a separate plan for military buildings and Houses,
    up to 6 Barracks and 6 Archery Ranges (`Promisory\extremebuildings4.per`).

### Army composition and counters

Mostly by countering what the enemy it targets has on the field,
counted every few seconds
(section "Archers", `Promisory\threats.per:569-603`):

- **Enemy archers → skirmishers**, about 2 per archer:
  3 archers allow up to 6 skirmishers, 10 allow 20, 40 allow 60
  (section "Skirmishers", `Promisory\units.per:5306-5390`).
- **Enemy infantry → archers**, unless the enemy already has
  many skirmishers, Huskarls or siege
  (section "Arbs", `Promisory\units.per:5745-6747`).
- **Enemy cavalry → pikemen and camels**
  (sections "Pikes" and "Camels", `Promisory\units.per:4747`, `units.per:7275`).
- **Out of gold in Imperial Age → skirmishers, pikemen and hussars**,
  the units that cost no gold (section "No gold left", `Promisory\units.per:8517`).

**Cavalry archers are rare**, though they are among the strongest units,
fast and ranged. The script trains them only in narrow cases
(goal `traincavarcher`, `Promisory\units.per:5808-6695`):

- **A cavalry archer rush**: weighted high only for Huns, Mongols and Tatars,
  and even then only while its army is behind and it has spare wood and gold.
- **As a counter**: to enemy infantry, capped at 15 cavalry archers
  until the enemy has 40 infantry, or to Teutonic Knights and war elephants.
- **As the default unit** only for Huns and Tatars (`Promisory\customConstants.per`).
  Most other civs default to their unique unit, a few to Champions, Knights or Mangonels.

**Hussars are rare too**, though they cost no gold and are strong late in the game.
The script treats them as a fallback, capped at 12 to 25 of the scout line
(section "Hussars", `Promisory\units.per:6747-6997`):

- **When it is short of gold**: spare gold under about 200, with spare food.
- **Only while the enemy has few pikemen and camels**: as few as 5 camels stop it.
- **Against enemy siege**, while it has few cavalry of its own.

Each civ also weighs its opening strategies (`Promisory\customConstants.per`):
across the civs, the Knight rush weighs the most, not an archer strategy.
So an Extreme that makes skirmishers and archers is mostly answering
an army of archers and infantry, or running out of gold.

### When it resigns

Rules on the goal `resign`, in `Promisory\resign.per:107-381`.
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
  like a taunt, and the script reads it (goal `fastresign`, `resign.per:711-735`):
  - "104": "Don't resign!"
  - "105": "You can resign again."
  - "205": no taunt text or sound in the game.
    It makes the AI resign sooner: roughly once its population
    is about half the enemy's or less (goal `fastresign`, `resign.per:750-825`).
- **One second chance against humans**: the first time it decides to resign
  before 42 minutes, while not under attack or with an ally,
  it holds off for 32 seconds and checks again (`resign.per:404-430`).
- **On resigning** it sends a random chat line, then resigns 6 seconds later
  (constant `enable-resign`, `resign.per:431-438`).

Before the population gap rules above can fire, all of these must hold
(`resign.per:113-200`):

- **Its population is under about 67** (villagers plus army).
- **It has lost at least 80 units and buildings** to enemies.
- **It is under attack right now.**
- **Its army is rated far behind**, on its own scale from -3 to +4:
  the lowest step, -3, for itself and for its team.

## Lobby Handicap Setting

The lobby handicap is a game setting per player, separate from any AI script.
It boosts the player it is set on.
**Confirmed from the game's own script**, the function `HandicapSetup` in
`{AoE2DEInstallDir}\resources\_common\xs\Effects.xs:139-203`:

- **Applied once, at game start**, by the "Dark Age effect"
  (function `EffectFunction2`, `Effects.xs:257-261`).
  It multiplies the base stats,
  so later flat upgrades like Loom's +15 hit points are not multiplied.
- **Does nothing at 100%**, the lowest value the lobby allows
  (in `HandicapSetup`, `Effects.xs:143`).
- **Per player**: it reads the multiplier for each player from the lobby,
  so an AI slot gets it only if set on that slot.

At 110%, everything below is × 1.1:

| Stat                | Applies to                                                                 |
|---------------------|----------------------------------------------------------------------------|
| Hit points          | All buildings, villagers                                                   |
| Work rate           | Villagers (gathering, building, repairing), Farms                          |
| Carry capacity      | Villagers                                                                  |
| Gathering bonuses   | Food, wood, gold, stone, shepherding, hunting, foraging                    |
| Training time ÷ 1.1 | Military units and monks: about 9% faster. Villagers train at normal speed |

**Not** changed by the current handicap script,
even if older patches, forum posts or the lobby tooltip say otherwise:

- **Military unit on-field stats** (hit points, attack, bonus damage, armor, speed).
- **Research time** and **villager training time**.
- **Starting resources**: the lobby tooltip says "increasing starting resources", which is outdated.

The built-in AI's comment that cheating was "obsoleted by the handicap system"
(on the `cheating` goal, `Promisory\customConstants.per:1193`) fits this:
a handicap on an AI slot is the official way to give it a boost.
The built-in AI detects a handicap on itself
by checking whether its villagers have more than 25 hit points,
then researches Loom early (section "Loom", `Promisory\researches.per:2801-2819`).

So at 110% on the Extreme AI's slot, against a human at 100%:

- **Its economy runs at least 10% faster**,
  with 10% tougher buildings and villagers, and about 9% faster military training.
  This makes up for some of its less efficient villager play,
  bringing its economy closer to a human's.
- **It notices the handicap**: its villagers have 27.5 hit points instead of 25,
  so it researches Loom early.
- **It still never weakens itself**, as on any Extreme game.
- **The human gets no boost.**
- **Immortal gets the same boosts on its slot**, but has no handicap checks at all,
  so its play does not change.

## Immortal

- **Compiled from C++ by the author's own compiler, source not shared.**
  The author on the forum thread, 2022-01-20:
  "I worte a compiler in C++. It runs C++ code to generate the script."
  No post in the thread shares the compiler or the C++ code.
  The file says "This script was compiled on: 2024-05-11" (rule 3575, line 48634).
  Signs:
  - Generated sections: constants (line 1), goals (line 2233),
    function names (line 17174), script (line 17231).
  - Machine variable names like `varTmp5547`,
    and return-address goals like `___ra1___`.
  - 3,236 jumps to exact rule numbers (the built-in AI has 21):
    inserting or deleting a rule shifts every target after it.
- **Mostly bypasses the engine**:
  - Sets the engine's civilian gatherer percentage to 0
    and assigns each villager itself.
  - Sets the engine's attack groups to 0
    and runs its own unit control per unit type.
  - Places Farms, camps, towers and Castles with its own code.
- **Ignores the difficulty setting**: no difficulty checks,
  never weakens itself, no cheats. It always plays at full strength.
- **About 30 named build orders**,
  for example "FF->Scout->Skirm" (fast Feudal, scouts, then skirmishers)
  or "FC->Knight" (fast Castle into knights).
  It picks one with random weights per civ and announces it in chat.
- **Uses spare strategic number slots as memory**:
  unit counts every loop, the chosen strategy,
  shared with allies along with "change strategy" chat messages.

### Immortal v1.2c vs v0.10f

The author came back on the forum thread in late 2026:
a work-in-progress update on 2026-10-30, then v1.2c on 2026-11-04,
posted in the thread only, not in Browse Mods.
Its chat reads "Compiled on 2026-11-04".
Read from both files:

| Question              | v0.10f                            | v1.2c                                                                                      |
|-----------------------|-----------------------------------|--------------------------------------------------------------------------------------------|
| Size                  | 115,901 lines, about 9,700 rules  | 124,263 lines, about 9,000 rules                                                           |
| Incas                 | 🔴 Never mentions the Settlement  | 🟢 Counts and uses Settlements (104 mentions), untested in game                            |
| Civs checked by name  | 45                                | 52: adds Jurchens, Khitans, Shu, Wei, Wu, Mapuche, Muisca, Tupi, Danes, Saxons, Varangians |
| Build orders          | About 30                          | 16, all named for Arabia (like "ff-arabia-scout-full-2stable"), plus "generic"             |
| Jumps to exact rules  | 3,236                             | 3,152                                                                                      |
| Difficulty and cheats | None                              | None                                                                                       |
| Engine turned off     | Villager gathering, attack groups | Same                                                                                       |
| Army micro            | Own, per unit type                | Same functions (melee, archer, siege, monk, Trebuchet dodging)                             |
| Economy and buildings | Own code                          | Rewritten: new Mill, Farm, drop-off and build order code                                   |
| Walls                 | No wall lines                     | Same                                                                                       |

- **Khitans are checked by name, but the author says they are not supported yet.**
- **The author reports much better late-game performance** in 4v4 than the old version.
- **The rewritten economy code matches the reported Mill bug**:
  the Mill placement and Farm functions are all new in v1.2c.
- **On walls, the author wrote in 2023**: "Yeah it's difficult for AIs to deal with walls properly."

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
- **Fixed styles**: Feudal 2 Ranges into Knights in 1v1 and 2v2;
  in team games, flanks rush archers and pockets go fast Castle into knights.
- **Claimed strength**: around 1150 Elo in 1v1 on the Definitive Edition,
  beats every non-cheating AI on Arabia except Immortal, "where games can go either way".
- **Known gaps**: does not handle walls or forward buildings.

What a few 1v1 games against it as Franks and Persians showed, 2026-11-27:

- **The same game every time**: 2 Archery Ranges, a big ball of archers,
  then Castle Age and a switch to Knights. It gets boring fast.
- **Inhuman archer micro**:
  fights with even Skirmishers are not worth taking,
  so the human ends up avoiding fights instead of naturally chasing the archers down.
  This fits the micro below: a target for each arrow, no wasted arrows, timed volleys, hit and run.
- **Its Knights dance against Ballistics and get shot down**:
  they keep dodging and never charge the archers.
  The Knight code has no dodging of its own: by default the Knights follow the archer group
  and stay within about 10 tiles of it
  (goal `gl-melee-army-mode`, lines 23147-23215; comment "Follow ranged army", lines 23572-23600).
  Against Ballistics the archer group flips its walking direction on every pass
  (goal `gl-army-direction`, lines 17870-17885), and the games show the Knights flipping with it.
  The flipping may help archers, but it plays Knights against Cavalry Archers badly:
  only the Knight being shot at should dodge, while the others charge in.
  Rehoboam dodges with the whole group, so its Knights never fight back.
  They only leave to raid villagers when the archer group is safe.

What the file shows:

- **Hand-written, not compiled**: an indented outline of comments,
  partly in Danish ("Fase" for phase, "ny" for new).
  It loads no other file.
- **Built as four phases**, each a block of rules that the others jump over
  (comments "Fase 1" to "Fase 4", lines 31881, 62680, 72142, 81219):

| Phase | Starts at line | What it does                                                                                                                                                          |
|-------|---------------:|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1     |          31881 | Fixed Dark Age opening, villagers under direct control: build orders, Houses, scouting, boar, sheep, deer luring, berries, Farms, wood, gold, anti militia rush micro |
| 2     |          62680 | Chosen Feudal strategy: fast Castle with 2 Town Centers into knights, or an archer and skirmisher rush                                                                |
| 3     |          72142 | Mostly for answering a tower rush: full economy, buildings, military, scouting                                                                                        |
| 4     |          81219 | Rest of the game: the 1v1 plan                                                                                                                                        |

- **Strategy is picked by team position**, not by randomness:
  flank or pocket in team games, a separate plan in 1v1 and 2v2
  (comment "Selecting strategy", line 8246).
- **Its own unit control**: ranged micro from the comment "Ranged unit micro" (line 11729)
  and melee micro from "Melee unit micro" (line 22465),
  about 14,000 lines together, plus scout micro and making Town Centers fire.
- **Archer and skirmisher micro, the strongest part of its army control**:
  - **Picks a target for each unit, one arrow at a time**
    (goal `gl-total-score`, lines 17371-17733):
    it works out the damage against each enemy in range
    (pierce armor, the skirmisher bonus against archers,
    25% more from higher ground and 25% less from lower ground),
    then the hits left to kill it.
  - **No wasted arrows**: after each order it subtracts the expected damage
    from that target's hit points, so the next unit picks someone not already dead.
    A target one hit from dying gets a big bonus
    (goal `gl-one-hit-remaining-score-bonus`, line 17365).
  - **Times its volleys**: against archers it waits until all its units can shoot,
    and fires when an enemy is about to shoot or after 5.7 seconds at most
    (comment "Finding out if enemy units are ready to attack", lines 16829-16915;
    goal `gl-time-since-last-shot`, lines 17272-17312).
    It does not wait when the enemy is close or its army is clearly ahead.
  - **Hit and run**: holds a range just beyond the enemy's,
    walks around them clockwise, and switches direction against Ballistics
    (goal `gl-army-direction`, lines 17860-17905;
    comment "Ajusting preffered distance to wether enemy has ballistics", lines 18512-18541).
    It turns away from trees and keeps to higher ground
    (comment "Change direction if tree is found", line 17989,
    and "Changing direction to stay on higher elevation", lines 18048-18230).
  - **Retreats by army strength**: it scores its group against the enemy's,
    counting enemy Knights, Scorpions and Mangonels as several archers each,
    and falls back to its towers or Town Center when the score drops below zero
    (goal `gl-group-superiority`, lines 12585-13609;
    comment "Defensive ranged army movements", lines 19894-20035).
  - **Avoids enemy Town Centers, Castles and towers** when moving
    (comment "moving away from enemy tc", lines 21286-21800).
- **Lag handling**: it measures how long each pass takes,
  spreads heavy code across passes, and skips work late in the game when slow
  (goal `gl-time-since-last-tick`, lines 2429-2551;
  comment "Skip current tick to reduce late game lag", line 8611).
- **A cheat that is off by default** (goal `gl-cheating`, lines 8719-8860):
  anyone typing "201" in chat turns it on, with a chat line saying so;
  "202" turns it off. Nothing in the script turns it on by itself.
  While on, every 15 seconds it adds 1 of each resource per villager gathering it,
  and clicks Feudal Age at 22 villagers.
- **Difficulty**: plays fully only on Hard and above;
  below Hard it turns off the training queue.
- **Mixes both approaches**: turns engine gathering and attack groups off in places
  for direct control, but still tunes many engine settings.

## How each AI handles Palisade walls

In games, the built-in Extreme AI is the most dangerous against a Palisade:
it attacks the wall when there is no gap,
while Immortal circles around looking for one.
The code explains both.

**Built-in Extreme AI: lets the engine attack walls**

- **It turns on the engine's wall targeting**, `sn-wall-targeting-mode`, at the start
  (`Promisory\init.per:282`),
  and keeps it on whenever it has siege
  (a ram, Bombard Cannon, Trebuchet, Petard or 2 Mangonels) or is out of houses
  (`finaling.per:1029-1039`).
  It turns it off only to defend, when your units come within 15 tiles of its town center
  (`finaling.per:998-1027`).
- **Its army moves as engine attack groups**:
  when a Palisade blocks the way, the group may attack the wall
  instead of searching far for a gap.
  A Palisade segment has only 250 hit points.
- **Target priority**, higher is attacked first
  (`up-set-offense-priority` and goal `attackpriorities`, `init.per:287-289`, `510`, `6384-6540`):
  - **Palisade Walls and Gates: 2. Stone Walls: 0.**
  - **Houses: 1**, so it does not go for Houses over bigger buildings.
  - **Buildings get 1 to 11 by its attack plan.**
    In its usual plan: Town Center 11, Lumber and Mining Camps 10, Mill 9,
    towers 7, military buildings 5, Castles 4.
  - It changes those numbers by plan:
    a Feudal rush puts the camps and Mill first,
    and enough siege or a maxed Imperial army switches back to the usual plan.

**Immortal: only targets units, so it walks around**

- **It turns the engine's wall targeting off** (rule 3590, `Immortal v0d10f.per:49171`).
- **Its unit control removes every building from the target list**:
  melee control (rule 1494, line 30714), the shared unit data (rule 1709, line 32546),
  and movement (rule 2494, line 39608).
  Its melee units and archers only ever choose units to attack.
- **It searches hard for a way in, but never checks for walls**:
  - **Attack angle** (rules 439-472, lines 21729-22032):
    it scores 25 to 50 spots on a full circle, 20 to 40 tiles around your main buildings,
    one spot every 5 tiles, and keeps the best.
    Spots near your Lumber Camps, Mining Camps, Mill or mule carts score higher,
    spots near your Castle, towers, Town Center or military buildings lower,
    and the spot it used last time scores lower too, so each new try comes from another side.
    Your Houses, Outposts and Palisades are not part of the score.
  - **Each move** (rules 2599-2608, lines 40598-40670):
    it tries up to 3 waypoints along the way, shifting one sideways when blocked,
    then checks up to 2 more against the walking distance:
    a waypoint is accepted only if walking there is less than 1.45 times the straight line.
  - **Against a complete Palisade, likely no waypoint passes**,
    and the next try comes from another side:
    this fits the army circling your base in games.
  - **Then it stops and the army stacks on one spot** (seen in games, against a wall with no gap).
    With no gap, the engine finds no walking path inside at all,
    so no waypoint ever passes and the group cannot make progress.
    A group that moves too long without progress switches to regroup:
    after 8 to 15 seconds in one movement state, and at most 90 seconds in another
    (rules 2252-2255, lines 37112-37158).
    Regroup orders every unit to the same point, in "stand ground" stance,
    and repeats the order every few seconds for any unit heading elsewhere
    (rules 2397-2406, lines 38679-38740).
    In stand ground, units stay put and only hit enemy units in reach, not your wall.
    Likely the regroup point sits by the wall and the group never finishes gathering,
    so it stays there; the compiled variable names leave this not fully proven.
- **Only its siege control targets buildings**, so without siege it never breaks in.
- **Its siege pushes are its biggest strength** (seen in games):
  - **Each attack comes from a new side**, so it feels like being outplayed by a good player.
  - **A lone ram or Trebuchet easily takes down a bare Castle** if you do not notice it in time.

  Its rams and Trebuchets share one control (rule 2797, line 42391),
  and each picks its target in this order (rules 2816-2837, lines 42522-42713):
  - enemy rams within 4 tiles,
  - then enemy siege nearby (Mangonels, Scorpions, Bombard Cannons, Organ Guns),
  - then defenses: Gates within 20 tiles, towers and Town Centers within 25 to 50 tiles,
    Castles within 30 to 60 tiles,
  - then any other building within 13 tiles, or 30 tiles once your population is small,
    except Houses, Lumber Camps, Mining Camps and mule carts.

  Palisade segments are not on the list: a ram goes for what stands behind them,
  and likely breaks the wall only where it blocks the way.
  The new side each time comes from the attack angle above:
  the last angle it used scores lower, so each push hits another side of your base.
- **v1.2c is the same**: wall targeting off (rule 3674, line 65972),
  and the same building removal in its unit control.
- **Fixing it means new rules in the middle of compiled code**,
  for example "when no enemy unit is reachable, attack the nearest Palisade":
  that shifts the jumps to exact rule numbers (see "Immortal").

**Rehoboam: booms and waits**

- **It warns in chat** it "will likely bug out against a fully walled opponent".
- **From 14 minutes on, it marks you as walled**
  with more than 2 Palisade Wall segments or more than 1 Gate,
  and none of your buildings in its town
  (goal `g-enemy-is-walled`, `Rehoboam 1.80k.per:76680-76687`).
  Any Palisade line passes the segment count;
  the Gate check names the plain gate, probably Stone Gates only, never Palisade Gates.
  The rules cited below all check `g-enemy-is-walled`.
- **Still in Feudal Age**: stops training archers once it has 5 military units (line 77940),
  skips the archer upgrades (lines 77336, 77474),
  stops at 13 Farms (line 75032), and builds no defensive tower
  unless you have more than 6 military units (line 76130).
- **In Castle Age**: adds Farms with a single Town Center (line 75084),
  keeps one Stable (line 76324),
  builds the Monastery and University without first having
  2 Stables or 3 Archery Ranges (lines 76359, 76974, 76425, 77068).
- **Holds its attacks against a walled enemy with a Castle**:
  it lowers its own army rating by 12 against an attack threshold of 14,
  until minute 37, or minute 47 without the Plate Barding Armor upgrade
  (`gl-MILITARY-SUPERIORITY` against `Defensiveness`, line 26841).
- **No code found for breaking into a walled base**:
  walls appear only in its tower rush and villager targeting.
  How well this works in a game is not checked.

**Their own walls**

- **Built-in Extreme AI walls early, Extreme difficulty only** (`Promisory\buildings.per`):
  - **Palisades near its camps** (goal `camp-walls`, lines 1954-2190), in the early game only:
    it stops by minute 17, in Castle Age, when you have 7 or more military units,
    when it has 4 or more, or 32 villagers,
    or when its plan goes archers, skirmishers, gunpowder, siege or monks.
    Likely the Palisade line seen in games, which stops you circling around its Town Center.
  - **A Palisade Gate when its own villager cannot walk out**
    (`palisade-gate1` checks, lines 10890-10950):
    it checks for no walking path, then puts a Gate in its wall.
    It does this only while it has no Palisade Gate yet,
    which fits its villagers sometimes walking the long way around.
  - **Builds a forward tower by your base and walls it in**
    (`watch-tower-building`, lines 8609-8728).
    In Feudal Age, when its ranged army is 8 to 16 tiles from your base
    and its army is rated ahead of yours, it builds a Watch Tower there,
    away from your Town Center and towers.
    Then, once per game, it surrounds that tower with Palisades,
    if the tower is more than 35 tiles from its own base and within 30 of yours.
  - **Its Dark Age wall around the map center is switched off**:
    "let's disable this until improvements are made"
    (constant `CENTER-WALL`, `customConstants.per:2278`).
- **Immortal builds no wall lines.**
  Its only Palisade code places single Palisade foundations
  on the tiles around one of its villagers going to build (rules 5670-5676, lines 70230-70315),
  likely a pathing trick; the purpose is not clear from the code. Same in v1.2c.
- **Rehoboam builds no walls.**
  It uses "can a Palisade be placed here" only as a test for free space,
  and builds a single Palisade only on a debug chat code (`Rehoboam 1.80k.per:37595-37599`).

## Script commands each AI uses

Counts of non-comment lines using each command, a rough picture of each design.

| Command                | Built-in | Immortal | Rehoboam | Reading                                         |
|------------------------|---------:|---------:|---------:|-------------------------------------------------|
| `up-jump-direct`       |       21 |    3,236 |      170 | Immortal: compiled function calls and loops     |
| `up-jump-rule`         |    1,903 |      145 |      506 | Built-in AI and Rehoboam: hand-written skips    |
| `up-get-object-data`   |      225 |    1,233 |      687 | Immortal reads single units most, Rehoboam next |
| `set-strategic-number` |    1,626 |        0 |    1,109 | Built-in AI and Rehoboam tune engine settings   |
| `research`             |    3,276 |      151 |    1,619 | Built-in AI spells out research per case        |
| `build`                |    6,035 |    1,106 |    4,517 | Built-in AI spells out building per case        |

## Built-in Extreme AI vs Immortal vs Rehoboam

All three at Extreme difficulty, as installed on 2026-11-26.
"Claimed" rows come from the AI's own description, not checked by play.

🟢 good, 🟡 less good, 🔴 clearly worse than the others;
no color where the row only describes or was not checked.
The Immortal column is v0.10f, the version in Browse Mods;
for what v1.2c changes, see "Immortal v1.2c vs v0.10f".
File sizes and rule counts are in "Where AI files are".

| Question                   | Built-in Extreme AI                                                                   | Immortal v0.10f                                             | Rehoboam v1.80k                                                       |
|----------------------------|---------------------------------------------------------------------------------------|-------------------------------------------------------------|-----------------------------------------------------------------------|
| Author                     | Official, shipped with the game (written by Promiskuitiv)                             | TheMaximalBeing                                             | ThNightPanda                                                          |
| Last changed               | 🟢 2026-10-28, changes with game patches                                              | 🟡 Compiled 2024-05-11 (author working on v1.2c in 2026)    | 🟡 Mod upload 2026-10-11                                              |
| Readable source            | 🟢 Yes, hand-written, English comments                                                | 🔴 No, compiler output, source not included                 | 🟢 Yes, hand-written, comments partly Danish                          |
| Safe to edit               | 🟢 Yes                                                                                | 🔴 Small value changes only: rule jumps are hard-coded      | 🟢 Yes                                                                |
| Game variety               | 🟢 High: plan per civ and map, army counters yours                                    | 🟢 High: about 30 build orders, picked at random            | 🔴 Low: one 1v1 plan, played as Persians or Franks                    |
| Civs                       | 🟢 All                                                                                | 🟡 All, but Incas break and post-2024 civs get generic play | 🔴 Persians or Franks (Teutons pass the check too)                    |
| Maps                       | 🟢 All                                                                                | 🟢 All                                                      | 🟡 Arabia and similar land maps                                       |
| Strategy choice            | 🟡 Rule tree per civ and map                                                          | 🟢 About 30 build orders, random weights per civ            | 🔴 Fixed: one plan for 1v1                                            |
| Villager control           | 🟡 Engine, steered by percentages                                                     | 🟢 Direct, per villager                                     | 🟡 Direct in the opening, mixed later                                 |
| Building placement         | 🔴 Mostly the engine: Castles often useless at the back                               | 🟢 Own code for Farms, camps, towers, Castles               | 🟡 Own build-line code for the opening buildings                      |
| Attack planning            | 🟡 Engine attack groups, targets by priority: Town Center, camps, Mill first          | 🟢 Strong: a new side each attack, like a good player       | 🟡 Moves as one archer ball; Knights follow it                        |
| Missile dodging and kiting | 🟡 Engine settings at their best on Extreme                                           | 🟢 Own micro                                                | 🟡 Own micro: good for archers, bad for Knights                       |
| Archer micro in play       | 🔴 Weakest: engine kiting only                                                        | 🟢 Clearly better than Extreme, close to human skill        | 🟡 Strong, but inhuman: a target per arrow, timed volleys             |
| Builds walls               | 🟢 Extreme only: early Palisades near its camps, so its Town Center cannot be circled | 🔴 No wall lines                                            | 🔴 None                                                               |
| Against your Palisades     | 🟢 Attacks the wall when there is no gap                                              | 🔴 Only siege hits walls, other units circle for a gap      | 🟡 Booms and holds attacks; no code to break in                       |
| Lag handling               | Not checked                                                                           | Not checked                                                 | 🟢 Measures pass time, spreads heavy code, skips work when slow       |
| Reacts to difficulty       | Yes, 6 levels                                                                         | No, always full strength                                    | Full strength on Hard and above                                       |
| Cheats                     | 🟢 None by default; "200" in chat turns them on                                       | 🟢 None                                                     | 🟢 None by default; "201" in chat turns them on                       |
| Claimed strength           | Not stated                                                                            | "Really powerful", plays all civs and settings              | About 1150 Elo in 1v1, even with Immortal, beats the others on Arabia |

## Limitations and bugs of each AI

**Built-in Extreme AI**:

- **Changes with every game patch**: findings hold only for the game build
  in "Game and AI versions".
- **Predictable army**: it mostly counters what the enemy fields
  (see the built-in AI's "Army composition and counters").
- **Kites worse when winning** (see the built-in AI's "Difficulty levels").
- **Relies on the engine** for gathering, placement and attacks: generic choices.
- **Poor Castle placement** (seen in games):
  its Castles often end up useless at the back of its base.
  The script mostly leaves the spot to the engine,
  placing around its own base with `sn-placement-zone-size`
  (`Promisory\buildings.per:10367-10442`, `buildings.per:11510`),
  and only one rule asks for a forward Castle (`buildings.per:10419`).
  Immortal, which places Castles with its own code, puts them in much better spots.

**Immortal v0.10f**:

- **Cannot play Incas** since the Settlement patch (seen in games):
  Incas now drop food, wood and gold at a Settlement
  instead of a Mill and camps (block `INCAN-CIV`, `Promisory\customConstants.per:3759-3813`),
  but Immortal never mentions the Settlement.
  It keeps trying to build Lumber Camps, Mining Camps and Mills, which Incas cannot build,
  so its economy gets stuck.
- **Does not know the civs released after May 2024**: it plays them with only its generic logic.
  They are Jurchens, Khitans, Shu, Wei, Wu, Mapuche, Muisca and Tupi;
  a Full Random AI can land on any of them that the game offers.
  Romans, from the Return of Rome DLC (2023), are covered:
  14 Roman checks, plus the Centurion and Legionary.
- **Stuck against walls** (seen in games): only its siege attacks walls,
  its other units walk around looking for a gap, even when the wall is complete
  (see "How each AI handles Palisade walls"). Same in v1.2c.
- **Only small value changes are safe to edit** (see "Immortal").
- **Replaced by v1.2c on the forum**, which should fix the Incas and adds most of the new civs
  (see "Immortal v1.2c vs v0.10f").

**Immortal v1.2c**, a work in progress:

- **Early bugs reported on the forum**: one or two Mills built per Farm,
  and Britons on Arena never building a University.
- **Khitans not supported yet**, by the author's post of 2026-11-04.
- **Fewer, Arabia-named build orders** (see "Immortal v1.2c vs v0.10f").
- **Not in Browse Mods**: it is installed by copying the files into the game's AI folder.
- **Same editing limit as v0.10f**: compiled, with 3,152 jumps to exact rule numbers.

**Rehoboam v1.80k**:

- **Few civs and one style**: Persians or Franks in 1v1, Feudal 2 Ranges into Knights;
  on other civs it warns "Wrong civ selected, this will lead to bugs" (see "Rehoboam").
- **Boring to play against**: the same archers into Knights every game,
  with inhuman archer micro (see "Rehoboam").
- **Weak against walls, by its own warning**: it booms and waits instead of breaking in
  (see "How each AI handles Palisade walls").
- **Arabia-like maps only**, with Conquest, low resources and Dark Age start.
- **Micro needs a fast computer**: a benchmark score of 1100 or more,
  and it warns in chat when the game lags.
- **Full Random civ does not fit**: a script cannot pick its own civ,
  so a random slot lands on Persians, Franks or Teutons only about 3 times in 50.
- **Unclear what its latest update changed**: the mod upload changed on 2026-10-11,
  but the file header still says v1.80f and 2020-2022.

## Where to start reading Immortal's code

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

- **The built-in AI, for game knowledge** (civs, maps, counter units):
  `units.per`, `researches.per`, `tsa.per`,
  and the `ability-to-kite` and `difficulty-adjustment` rules
  in `init.per:1800-2160` for the difficulty logic.
- **The forum thread** in "Where AI files are", for the source or the compiler name.

## Sources

- `{AoE2DEInstallDir}\resources\_common\ai\Promisory\`: the built-in AI.
- `{AoE2DEInstallDir}\resources\_common\xs\Effects.xs`: the lobby handicap.
- The Immortal and Rehoboam `.per` files and their mod `info.json`, see "Where AI files are".
