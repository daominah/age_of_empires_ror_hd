# Age of Empires: The Rise of Rome

## Download

From [github.com/daominah/age_of_empires_ror_hd](https://github.com/daominah/age_of_empires_ror_hd) (this page)
click on green button `Code` then `Download ZIP` (size about 100 MB).

## Run the game

* Init setup: run file `SETUP.EXE`

* For single player, run file `empires.exe` (lowercase file name, not `EMPIRESX.EXE`),
  cheat enabled, population limit increased from 200 to 1000.

* For multiplayer, run file `EMPIRESX.EXE`.
  If someone uses a cheat code in the game,
  their units will be deleted, and they will be kicked from the game.

Guides on how to edit computer player behavior are in the directory
[data2_daominah/doc](data2_daominah/doc/edit_computer_player.md).

## AoE version detail

The game version is The Rise of Rome 1.0.

### Farm bug still works

The Vietnamese AoE 1 community plays with the Farm bug for a very long time,
accepting it as a feature that balances Farms being costly to set up with a lot of wood at first,
but a strong food source in the long run.
Fixing the bug fundamentally changes the game,
this is one of the reasons they do not want the Definitive Edition.

The Farm bug: Pressing "S" when both a Farmer and Farms are selected at the same time will
replenish the Farm with full food quantity, as same as a new Farm.
(so researching Domestication, Plow, or Irrigation is not necessary,
food quantity bonus of the Sumerian and Minoan Farm is not important).

### Civilization bonuses AoE-RoR v1.0

The results here come from real tests and the game dat file,
some values are slightly different from the documentation.

The balance between civilizations is bad.
Shang is absolutely the best.
Greek, Choson, and Carthaginian are too weak in the most common "rule" game mode,
Bronze Age without Wall (a.k.a. "Đời ba không thành").

#### 1. Assyrian

* Villager move speed +18% (stated +30%):  
  +0.2 tiles/second, so starts at 1.3 instead of 1.1.
* Archers attack speed +36% (stated +40%):  
  attack reload time reduced to 1.1 instead of 1.5.

#### 2. Babylonian

* Priest rejuvenation +0.75, so starts at 2.75 instead of 2 (+38%, stated +50%).  
  so Priest's rest duration is 36s instead of 50s.  
  With Fanaticism, it is 24s instead of 29s.
* Stone Miner work rate +44% and capacity +3 (stated +30%).
* Tower and Wall HP x2.
  Fully upgraded Wall has 960 HP, Tower has 480 HP.

#### 3. Carthaginian

* Elephant and Academy units HP +25%.
* Light Transport move speed +25%, Heavy Transport move speed +43% (stated +30%).  
  (so Heavy Transport moves as fast as a Heavy Horse Archer at 2.5 tiles/s,  
  this is the fastest move speed in the game, not counting cheat units.)
* Fire Galley attack +6 (so fully upgraded attack is 24+12 instead of 24+6).

#### 4. Choson

* Priest cost -32% (stated -30%), so 85 gold instead of 125.
* Iron Age Swordsman receive +80 HP, so:
  - Long Swordsman has 160 HP instead of 80.
  - Legion has 240 HP instead of 160.
* Towers range +2, so fully upgraded range is 7+5 instead of 7+3.

#### 5. Egyptian

* Chariots HP +33%, so:
  - Chariot Archer starts at 93 HP instead of 70.
  - Chariot starts at 133 HP instead of 100.
  - Scythe Chariot has 182 HP instead of 137 (Nobility upgraded).
* Priest range +3 (fully upgraded range is 10+6 instead of 10+3).
* Gold Miner work rate +44% and capacity +2 (stated +20%).

#### 6. Greek

* Academy units move speed +33% (stated +30%) by +0.3 tiles/s, so:
  - Start at 1.2 instead of 0.9.
  - With Aristocracy, the speed is 1.45 instead of 1.15.
* Warships move speed +17% (stated +30%).

#### 7. Hittite

* Archers attack +1.
* Siege units HP x2 (so Heavy Catapult has 300 HP).
* Warships range +4
  (but missing Trireme in the Iron Age).

#### 8. Macedonian

* Academy units pierce armor +2.
* Melee units sight +2.
* Siege units cost -50%.  
  (so Stone Thrower cost reduced to 90 wood, 40 gold.)
* All units are 4 times more resistant to conversion.  
  (so Cavalry usually cannot be converted by a Priest on open field.)

#### 9. Minoan

* Composite Bowman range +2
  (so fully upgraded range is 7+5 instead of 7+3).
* Farm food +60 (so starts at 310 instead of 250).
* Ships cost -30%.

#### 10. Palmyran

* Forager, Hunter, Gold Miner, Stone Miner work rate +44%;    
  Woodcutter +36%;  
  Farmer, Builder, Repairman work just as fast as others
  (stated Villager work rate +25%).
* Villager armor +1 and pierce armor +1  
  (game UI not show pierce armor, but it is there).
* Villager cost +50% (so 75 food instead of 50).
* Camel Rider move speed +25% (so moves as fast as Heavy Horse Archer).
* Free tribute
  (so while other civilizations lose 125 resources to give ally 100,
  Palmyran only lose 100 resources to give ally 100).

#### 11. Persian

* Hunter work rate +66% and capacity +3 (stated +30%).
* Elephants move speed +56% (stated +50%),
  so starts at 1.4 tiles/s instead of 0.9.
* Trireme attack speed +38% (stated 50%).

#### 12. Phoenician

* Woodcutter work rate +36% and capacity +3 (stated +30%).
* Elephants cost -25%.
* Catapult Trireme and Juggernaught attack speed +72%.

#### 13. Roman

* Buildings cost -15% (except Tower, Wall, Wonder).  
  Tower cost -50% (but missing Guard Tower upgrade).
* Swordsman attack speed +50%  
  (stated +33%, probably refers to attack reload time).

#### 14. Shang

* Villager cost 35 food instead of 50 food (-30%).
* Wall HP x2.

#### 15. Sumerian

* Villager +15 HP (so has 40 HP instead of 25).
* Siege units attack speed +43%
  (attack reload time reduced to 3.5s instead of 5s).
* Farm food +250 (so starts at 500 instead of 250).

#### 16. Yamato

* Villager move speed +18% (stated +30%, same as Assyrian).
* Mounted units cost -25%, so:
  - Scout cost reduced to 75 food (normally 100).
  - Cavalry cost reduced to 52 food, 60 gold (normally 70 food, 80 gold).
  - Horse Archer cost reduced to 37 food, 52 gold (normally 50 food, 70 gold).
* Ships HP +30%.

### Definitive Edition and Return of Rome balance

In-game unit stats are defined in the file `data\empires.dat`.
In the same directory, the file `empires_definitive_edition.dat` is a copy with
changes that mimic the official balance patches of the
[Age of Empires: Definitive Edition](https://ageofempires.fandom.com/wiki/Summary_of_changes_in_Age_of_Empires:_Definitive_Edition)
and [Return of Rome](https://ageofempires.fandom.com/wiki/Age_of_Empires_II:_Definitive_Edition_-_Return_of_Rome).

To play with these changes, replace `data\empires.dat` with `data\empires_definitive_edition.dat`.
To revert, restore `data\empires.dat.backup`.
Script to do this [switch_data_definitive_edition.sh](switch_data_definitive_edition.sh).

Because I do not want to nerf good things, I only apply the buffs in the following list.

#### Units and buildings update

Mainly buffs for weak or expensive but underperformed [units](https://ageofempires.fandom.com/wiki/Unit_(Age_of_Empires)#Town_Center):
Swordsman, Elephant Archer, Cataphract, Academy, Towers.
The Catapult line now has increased pierce armor but negative melee armor.

##### 1. Villager

- [x] [Villager](https://ageofempires.fandom.com/wiki/Villager_(Age_of_Empires)#Changelog)
  move 10% faster after advancing to the Tool Age.  
  (Tool Age effect now adds +0.1 to Civilian Movement,
  so Villager speed increased to 1.2,
  Wheel effect adjusted accordingly to +0.6, so Wheel Villager stays at 1.8).
- [x] Hunters work 5% faster than before  
  (previously they worked at the same rate as Forager 0.45 food/s, now 0.4725 food/s).
- [ ] Gold and stone miners work 15% faster (to compensate bug fixes to technologies).
- [x] Increased the quantity of Gold Mine to 450 (previously 400).
- [x] Increased the quantity of Stone Mine to 300 (previously 250).
- [x] Towers build time reduced to 65s (previously 80s).
- [x] Farm build time reduced to 24s (previously 30s).
- [x] Farm technologies improved, fully upgraded Farm now has 550 food (prev 475):
  - Domestication: cost reduced to 150 food, 50 wood (previously 200 food, 50 wood).
  - Plow: +100 Farm food (previously +75).
  - Irrigation: +125 Farm food (previously +75).

##### 2. Barracks

- [x] [Axeman](https://ageofempires.fandom.com/wiki/Axeman_(Age_of_Empires)#Changelog)
  pierce armor increased to 1 (previously 0).
- [ ] Axeman has +4 bonus attack vs mounted units.
- [x] Short Sword research is free upon reaching the Bronze Age:
  - previously costs Food 120, Gold 50, Research Time 50, require Battle Axe
  - now costs Food 0, Gold 0, Research Time 1, only require Bronze Age
- [x] [Broad Swordsman](https://ageofempires.fandom.com/wiki/Broad_Swordsman#Changelog)
  have +10 HP, so now have 80 HP.
- [x] Long Swordsman have +20 HP, so now have 100 HP.
- [ ] Long Sword upgrade cost increased to 240 food, 100 gold
  (previously 160 food, 50 gold).
- [ ] Legions lose 20 HP.
- [x] [Slinger](https://ageofempires.fandom.com/wiki/Slinger_(Age_of_Empires))
  pierce armor increased to 3 (previously 2).
- [x] Slinger line of sight increased to 6 (previously 5).
- [ ] Slinger bonus attack against all archers increased to +3 (previously +2).
- [ ] Slinger training time increased to 35s (previously 24s).

##### 3. Archery Range

- [ ] [Bowman](https://ageofempires.fandom.com/wiki/Bowman_(Age_of_Empires))
  has +2 bonus attack vs light infantry (previously 0).
- [ ] Chariot Archer bonus attack vs Priest removed (previously +7),
  also conversion resistance reduced to x2 (previously x8).
- [ ] Horse Archer and [Heavy Horse Archer](https://ageofempires.fandom.com/wiki/Heavy_Horse_Archer)
  pierce armor reduced to 1 (previously 2).
- [x] [Elephant Archer](https://ageofempires.fandom.com/wiki/Elephant_Archer_(Age_of_Empires))
  attack increased to 6 (previously 5)
  (cost unchanged 180 food, 60 gold).
- [x] Improved Bow research time reduced to 45s (previously 60s).

##### 4. Stable

- [x] [Scout](https://ageofempires.fandom.com/wiki/Scout_(Age_of_Empires))
  cost reduced to 90 food (previously 100).
- [x] Scout pierce armor increased to 1 (previously 0).
- [ ] Scout have bonus attack vs archers +3 (previously 0).
- [ ] Chariot bonus attack vs Priest removed (previously +7),
  also conversion resistance reduced to x2 (previously x8).
- [ ] Scythe Chariot melee armor reduced to 1 (previously 2).
- [ ] Scythe Chariot upgrade cost increased to 1400 wood, 1000 gold (previously 1200 wood, 800 gold).
- [x] [Cavalry](https://ageofempires.fandom.com/wiki/Cavalry_(Age_of_Empires))
  has 1 pierce armor (previously 0).
- [x] Heavy Cavalry has 2/2 melee/pierce armor (previously 1/1).
- [x] [Cataphract](https://ageofempires.fandom.com/wiki/Cataphract_(Age_of_Empires))
  has 240 HP, 5/3 melee/pierce armor (previously 180 HP, 3/1 armor).
- [x] Cataphract upgrade cost reduced to 1600 food, 600 gold
  (previously 2000 food, 850 gold).
- [ ] War Elephant has +25 bonus attack vs buildings (previously 0).
- [ ] Armored Elephant bonus attack vs all buildings reduced to +35
  (previously +40 vs buildings, additional +40/+80 vs Towers/Walls).
- [ ] Camel Rider added +4 bonus attack vs elephants
  (so now Camel have bonus attack vs all mounted units).

##### 5. Academy

- [x] Academy cost reduced to 150 wood (previously 200 wood).

##### 6. Siege Workshop

- [ ] Stone Thrower/Catapult/[Heavy Catapult](https://ageofempires.fandom.com/wiki/Heavy_Catapult)
  have 3/4/5 pierce armor (previously 0).
- [ ] Stone Thrower/Catapult/Heavy Catapult have -2 melee armor (previously 0).
- [ ] Catapult and Heavy Catapult have slightly smaller damage area radius.
- [ ] Siege weapon projectiles (stones and bolts) travel slightly faster,
  but are also slightly slower to reload (affects all stone and bolt firing units).
- [x] Ballista attack reload time reduced to 2.7s (previously 3).
- [ ] Helepolis attack reload time increased to 1.85s (previously 1.5s).
- [ ] Helepolis attack increased to 45 (previously 40).
- [x] Helepolis upgrade cost reduced to 1200 food, 1000 wood
  (previously 1500 food, 1000 wood).

##### 7. Temple

- [ ] Priest have +1 healing range.
- [x] Sacrifice costs 400 gold (previously 600 gold).

##### 8. Tower

- [x] [Towers](https://ageofempires.fandom.com/wiki/Tower#Age_of_Empires)
  build time reduced to 65s (previously 80s).
- [ ] Watch/Sentry/Guard/Ballista Towers HP increased to 125/185/240/240
  (previously 100/150/200/200 HP).
- [ ] Watch/Sentry/Guard/Ballista Towers attack increased to 5/6/8/20
  (previously 3/4/6/20 base pierce attack).
- [ ] Watch/Sentry/Guard/Ballista Towers range increased to 6/7/8/8
  (previously 5/6/7/7 range).
- [ ] Watch/Sentry/Guard/Ballista Towers pierce armor increased to 3/4/5/5
  (previously 3/4/4/4 pierce armor).

##### 9. Dock

- [ ] Catapult Trireme HP increased to 135 (previously 120).
- [ ] Catapult Trireme and Juggernaught move 10% faster.
- [ ] Catapult Trireme and Juggernaught cost reduced to 135 wood, 50 gold
  (previously 135 wood, 75 gold).
- [ ] Juggernaught upgrade cost reduced to 1300 food, 500 wood
  (previously: 2000 food, 900 wood).
- [ ] Juggernaught can no longer destroy trees.
- [ ] Trade Boat / Merchant Ship HP reduced to 120/200 (previously 200/250 HP).

#### Civilizations bonuses AoE-DE

##### Shared changes

Shared changes are applied to all civilizations:

- [x] Writing's Research Time reduced to 30s (previously 60s).
- [x] Wheel (TechID 28) is available for all (previously disabled for Persian and Macedonian).
- [x] Coinage (TechID 30) is available for all (previously disabled for Egyptian, Palmyran, Persian, Shang and Sumerian).
- [ ] Heavy Transport is available for all (previously disabled for Assyrian, Babylonian, Choson, Hittite, Palmyran, Shang and Sumerian).

##### Engine limitations

Genie Engine version used in AoE Rise of Rome is old, so has some limitations:

- For civilization-specific tech research bonuses,  
  in AoE2 we can use [Tech Cost Modifier](#tech-cost-modifier),  
  but in this old AoE-RoR, we use a workaround similar to [unique tech](#unique-tech).
- All team bonuses are implemented as normal bonuses,
  only affecting the civilization itself, not the whole team.

##### 1. [Assyrian](https://ageofempires.fandom.com/wiki/Assyrians#Civilization_bonuses)

- [x] Chain Mail (Infantry/Archers/Cavalry) upgrades are available
  (but still missing Heavy Horse Archer techID 38).
- [x] Alchemy and Engineering are available at the Government Center.
- [x] Siege Workshops work 20% faster (team bonus).
- [x] Siege units upgrades cost -50%.
- [ ] Villagers move 10% faster (previously 18% faster).
- [ ] Archers have -25% Attack Reload Time (so +33% fire rate, previously +36%).

##### 2. [Babylonian](https://ageofempires.fandom.com/wiki/Babylonians#Civilization_bonuses)

- [x] Chariot units have +1 pierce armor.
- [x] Metallurgy is available.
- [x] Chain Mail (Infantry/Archers/Cavalry) upgrades are available
  (but still missing Heavy Horse Archer).
- [ ] Market technologies cost -30%.
- [x] Builders work 10% faster (team bonus).

- [ ] Towers and walls have +60% HP (previously +100% HP).

##### 3. [Carthaginian](https://ageofempires.fandom.com/wiki/Carthaginians#Civilization_bonuses)

- [x] Start the game with +50 of all resources Wood, Food, Gold, Stone.
- [x] Camel Riders have +15% HP.
- [ ] Nobility cost is free (still requires Government Center).
- [x] Academy work 20% faster (team bonus).

- [ ] Transport ships move 25% faster (previously 43% faster for Heavy Transport).

##### 4. [Choson](https://ageofempires.fandom.com/wiki/Choson#Civilization_bonuses)

- [x] Axeman HP increased to 55 (previously 50 HP, same as normal).
- [x] Short/Broad/Long Swordsman and Legion have +15/+20/+60/+80 HP.
  (previously +0/+0/+80/+80 HP).
- [x] Get Nobility.
- [x] Storage Pit technologies cost -40%.  
  (so the total saving is 570 food and 222 gold;
  could be more if Choson had Iron Age armors, but they do not).
- [x] Buildings have +2 Line of Sight (team bonus).

##### 5. [Egyptian](https://ageofempires.fandom.com/wiki/Egyptians_(Age_of_Empires)#Civilization_bonuses)

- [x] Farms cost -20%, so their cost reduced to 60 wood (previously 75 wood).
- [x] Get Coinage.
- [x] Priests have +1 pierce armor (team bonus).
- [ ] Priests have +2/+3 range in the Bronze/Iron Age.
  (previously +3 range for all ages).

##### 6. [Greek](https://ageofempires.fandom.com/wiki/Greeks_(Age_of_Empires)#Civilization_bonuses)

- [x] Academy units cost -20% (in addition to the already present speed bonus).
  So Hoplite/Phalanx/Centurion cost 48 food, 32 gold (previously 60 food, 40 gold).
- [x] Town Centers work 10% faster starting in the Tool Age.
  (hard to implement, so just increased the work rate in all ages).
- [ ] Polytheism and Astrology free (requires Temple).
- [ ] Ships are 20% faster (previously 17% faster).
- [ ] Get Fire Galley.
- [x] Market cost -50% (team bonus).

##### 7. [Hittite](https://ageofempires.fandom.com/wiki/Hittites#Civilization_bonuses)

- [ ] Wheel -50% cost and research time.
- [x] Towers provide +4 population room (team bonus).

- [ ] Catapults have +50% HP (previously +100% HP).
- [ ] War ships (except Fire Galley) have +1/+2/+3 range in the Tool/Bronze/Iron Age
  (previously +4 range for all ages).
- [ ] Lose Centurion, Architecture and Irrigation.

##### 8. [Macedonian](https://ageofempires.fandom.com/wiki/Macedonians#Civilization_bonuses)

- [x] Get Wheel.
- [x] Get Catapult upgrade. (still missing Heavy Catapult).
- [x] Houses have +50 HP (team bonus),
  so they have 125 HP instead of 75 HP.

- [ ] Academy units have +1/+2 pierce armor in the Bronze/Iron Age
  (previously +2 pierce armor for all ages).
- [ ] Siege Workshop units cost -25% (previously -50%).

##### 9. [Minoan](https://ageofempires.fandom.com/wiki/Minoans#Civilization_bonuses)

- [x] Farmers work 10% faster.
- [x] Docks cost -20% (team bonus).

- [ ] Bonus Farms have +60 food removed.
- [ ] Improved Bowman line has +1/+2 range in the Bronze/Iron Age, respectively.
  (previously +2 for Composite Bowman for all ages).
  - [x] Only implement Improved Bowman +1 range for all ages.
- [ ] Ships cost -15%/-20%/-25%/-30% in the Stone/Tool/Bronze/Iron Age.
  (previously -30% for all ages).

##### 10. [Palmyran](https://ageofempires.fandom.com/wiki/Palmyrans#Civilization_bonuses)

- [x] Fix Palmyran Farmer, Builder, Repairman work rate to +25% (previously +0%).
- [x] Get Coinage and Plow.
- [x] Start the game with +75 food (previously have the same 200 food as normal).
- [ ] Trade units return +20% gold.
- [ ] Technologies researched 30% faster (team bonus).

- [ ] Villagers work 25% faster for all tasks
  (previously Forager, Hunter, Gold Miner, Stone Miner +44%; Woodcutter +36%).

##### 11. [Persian](https://ageofempires.fandom.com/wiki/Persians_(Age_of_Empires)#Civilization_bonuses)

- [x] Wheel, Artisanship, Plow, Coinage are available.
- [x] Ballistics is available.
- [x] Walls cost -20%, so 4 stone (previously 5 stone).
- [x] Stables work 20% faster (team bonus).

- [ ] All elephant units move 25% faster (previously 56% faster).
- [ ] Scout Ship line fires 18%/25%/33% faster in the Tool/Bronze/Iron Age.
  (previously 38% faster but only for Trireme in the Iron Age).

##### 12. [Phoenician](https://ageofempires.fandom.com/wiki/Phoenicians#Civilization_bonuses)

Seem like no buffs, but the shared Elephant Archer +1 attack is significant
due to the Phoenician low Elephants cost.

- [x] Docks have +150 HP, +4 line of sight.
- [x] Archers have +2 line of sight (team bonus).

- [ ] Catapult Trireme and Juggernaught fire 30% faster (previously 72% faster).
- [ ] Woodcutters work 15% faster and carry +2 wood (previously work rate +36%, carry +3).

##### 13. [Roman](https://ageofempires.fandom.com/wiki/Romans_(Age_of_Empires)#Civilization_bonuses)

- [x] Ballista and Helepolis have +1 range, so start at 9+1, max at 10+3.
- [x] Priests +50% heal speed (team bonus)  
  (implemented by set Roman `Heal Bonus (ResourceID 56)` to 1.5 (default is 0, works same as 1),
  along with change `Medicine` effect from `set Heal Bonus to 3` to `change Heal Bonus +3`,  
  if `Heal Bonus ≠ 0`, it will be multiplied with the base heal rate 3 HP/s;  
  Result a Roman Priest heal rate starts at 4.5 HP/s, with Medicine it is 13.5 HP/s, missing Astrology;  
  other civs Priest start at 3 HP/s, Astrology increases to 3.9 HP/s, combined with Medicine it is 11.7 HP/s, same as before).
- [ ] Buildings cost -10% (previously -15%)
- [ ] Towers cost -40% (previously -50%)

##### 14. [Shang](https://ageofempires.fandom.com/wiki/Shang#Civilization_bonuses)

- [x] Get Coinage.
- [x] Get Ballistics.
- [x] Cavalry, Heavy Cavalry, Cataphract attack 10% faster.
- [x] Town Center provide +4 population room (team bonus).

- [ ] Villagers cost 40 food (previously 35 food). Normal cost is 50 food.
- [ ] Start the game with -40 food.
- [ ] Walls have +60% HP (previously +100% HP).

##### 15. [Sumerian](https://ageofempires.fandom.com/wiki/Sumerians#Civilization_bonuses)

- [x] Get Coinage.
- [x] Camels have +1 pierce armor.
- [ ] Stone Thrower, Catapult, and Heavy Catapult fire 45% faster (previously 43%).
- [x] Town Centers cost -25%, so 150 wood instead of 200 wood (team bonus).

- [ ] Farms have +125 food (previously +250 food).

##### 16. [Yamato](https://ageofempires.fandom.com/wiki/Yamato#Civilization_bonuses)

- [ ] Stable and Archery Range upgrades cost -30%.
- [ ] Fishing Boats work 20% faster.
- [x] Stable and Archery Range cost -33% wood (team bonus).
  So their cost reduced to 100 wood (previously 150 wood).

- [ ] Mounted units cost -15% (previously -25%).
- [ ] Villagers move 10% faster (previously 18% faster).
- [ ] Ships have +10%/+15%/+20%/+25% HP in the Stone/Tool/Bronze/Iron Age
  (previously +30% for all ages).

##### 17. [Lac Viet](https://ageofempires.fandom.com/wiki/Lac_Viet#Civilization_bonuses)

TODO: Lac Viet is a new civilization,
need to modify the `empires.exe` to show the new civilization in the game,
not just edit the `data/empires.dat` file. Can consider to replace Carthaginian,
the civilization still bad after got buffs.

- [x] Foragers work 20% faster.
  (Forager work rate +0.1, so 0.55 food/s instead of 0.45 food/s).
- [x] Archers have +2 melee armor.
- [x] Ballista have +2 melee armor (but Helepolis is missing).
- [ ] All military units are created 25% faster.
- [ ] Houses and Farms are built 50% faster (team bonus).

#### Note on using Genie Editor

##### Open dat file

Click `Open` on the Genie Editor toolbar,
it will show a lot of paths that we need to configure,
can try button `Fill paths from registry` but usually not totally correct,
we can manually fill 3 necessary paths:

- Compressed data set (*.dat):
  `D:\game\age_of_empires_ror_hd\data\empires_definitive_edition.dat`
- Language file location:
  `D:\game\age_of_empires_ror_hd\language.dll`
- Language x1 file location:
  `D:\game\age_of_empires_ror_hd\languagex.dll`

##### Civilization Tech Tree

In tab `Civilizations`, left panel show list of civilizations,  
click on a civilization will show its corresponding Technology Tree, that is an EffectID,  
then switch to tab `Effects`, find that EffectID to edit the Effect Commands.

##### Unique Tech

Some civilization bonuses are having some cheaper tech costs,
they can be implemented by creating a new dummy unique Tech for that civilization,
disable non-corresponding civilizations to have that Tech (all techs are enabled by default),
then change the discounted Tech to require that unique Tech.

My dummy unique techs list added in `data/empires_definitive_edition.dat`:

- TechID 127: uniq_Assyrian_E81
- TechID 128: uniq_Babylonian_E82
- TechID 129: uniq_Carthaginian_E205
- TechID 130: uniq_Choson_E91
- TechID 131: uniq_Egyptian_E83
- TechID 132: uniq_Greek_E84
- TechID 133: uniq_Hittite_E85
- TechID 134: uniq_Macedonian_E206
- TechID 135: uniq_Minoan_E86
- TechID 136: uniq_Palmyran_E207
- TechID 137: uniq_Persian_E87
- TechID 138: uniq_Phoenician_E88
- TechID 139: uniq_Roman_E208
- TechID 140: uniq_Shang_E89
- TechID 141: uniq_Sumerian_E90
- TechID 142: uniq_Yamato_E92
- TechID 143: uniq_LacViet_E218_unused

##### Tech Cost Modifier

The following seems **not work** in this AoE version (The Rise of Rome 1.0),
but still keep this note here for future trying.

When adding a Command for an Effect
(usually the Effect that defines a civilization bonuses),
the dropdown to choose Command Type does not show all available commands,
the list does show the most common commands that are:

- Attribute Modifier (Command Type 0, 4, 5 for Set, Add, Multiply).
- Disable Tech (Command Type 102).

We have to manually type the **Command Type** to `101`,
which is [Tech Cost Modifier](https://agecommunity.fandom.com/wiki/Tech_Cost_Modifier_(Set/%2B/-)),
then need to fill in its arguments:

- Attribute A: **TechID** to modify. Example `27` for `Helepolis` upgrade.
- Attribute B: **ResourceID** to modify, the value can be:

  - `0` for Food
  - `1` for Wood
  - `2` for Stone
  - `3` for Gold
  - ... full values list [here](https://agecommunity.fandom.com/wiki/Resource_List)

- Atrtibute C: modifier **Mode**, the value can be:

  - `0` for Set to exact value.
  - `1` for Increase/Decrease to original cost.

- Attribute D: the **Value** to Increase/Decrease or Set at.

In the `Helepolis` upgrade example,
if the original cost is 1200 food 1000 wood,
and we want to -50% cost, we need to add two commands:

- Command 101, TechID 27, ResourceID 0, Mode 0, Value 600 (set food cost to 600)
- Command 101, TechID 27, ResourceID 1, Mode 0, Value 500 (set wood cost to 500)

---

### Appendix: Priest Conversion Test

Priest is the only unit that has 33% accuracy in AoE-RoR;
all other units have 100% accuracy.
However, it is unclear exactly how this 33% accuracy functions in practice.

In the following tests, 20 priests start to convert enemy unit at range 10,
the result number is the number of failed conversions.

#### Priest vs Cavalry

[14, 13, 14, 13, 11, 13, 15, 11, 16, 15, 10, 11, 14, 14, 7, 13, 13, 11, 13, 16]

#### Priest vs Cavalry

36% converted.

#### Mysticism Astrology Priest vs Cavalry

[7, 5, 6, 8, 8, 7, 5, 10, 5, 9, 6, 5, 4, 5, 4, 5, 6, 7, 8, 6]

68% converted.

#### Priest vs Centurion

[7, 6, 5, 4, 4, 8, 4, 3, 6, 3, 6, 9, 11, 3, 6, 3, 7, 4, 8, 4]

72% converted.

#### Mysticism Astrology Priest vs Centurion

[3, 2, 7, 5, 3, 4, 6, 2, 2, 4, 8, 2, 3, 5, 3, 3, 5, 5, 5, 6]

79% converted.

#### Priest vs Macedonian Centurion

[16, 14, 15, 19, 18, 16, 16, 19, 15, 14, 11, 15, 13, 15, 17, 11, 18, 15, 17, 13]

23% converted.

#### Mysticism Astrology Priest vs Macedonian Centurion

[15, 16, 17, 16, 15, 11, 9, 14, 12, 19, 17, 16, 13, 17, 16, 11, 17, 18, 16, 11]

26% converted.
