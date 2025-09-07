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

### Civilization bonuses

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
* Fire Galley attack +6 (24+12 instead of 24+6).

#### 4. Choson

* Priest cost -32% (stated -30%),  
  so 85 gold instead of 125.
* Iron Age Swordsmen receive +80 HP.  
  Long Swordsman has 160 HP instead of 80;  
  Legion has 240 HP instead of 160.
* Towers range +2: fully upgraded range is 7+5 instead of 7+3.

#### 5. Egyptian

* Chariots HP +33%  
  (so Chariot Archer starts at 93 HP instead of 70,  
  Scythe Chariot has 182 HP instead of 137).
* Priest range +3 (fully upgraded range is 10+6 instead of 10+3).
* Gold Miner work rate +44% and capacity +2 (stated +20%).

#### 6. Greek

* Academy units move speed +33% (stated +30%):  
  +0.3 tiles/s, so starts at 1.2 instead of 0.9.  
  With Aristocracy, the speed is 1.45 instead of 1.15.
* Warships move speed +17% (stated +30%).

#### 7. Hittite

* Archers attack +1.
* Siege units HP x2.  
  (Heavy Catapult has 300 HP.)
* Warships range +4.

#### 8. Macedonian

* Academy units pierce armor +2.
* Melee units sight +2.
* Siege units cost -50%.  
  (Stone Thrower cost reduced to 90 wood, 40 gold.)
* All units are 4 times more resistant to conversion.  
  (so Cavalry usually cannot be converted by a Priest on open field.)

#### 9. Minoan

* Composite Bowman range +2.  
  Fully upgraded range is 7+5 instead of 7+3.
* Farm food +60 (starts at 310 instead of 250).
* Ships cost -30%.

#### 10. Palmyran

* Forager, Hunter, Gold Miner, Stone Miner work rate +44%;    
  Woodcutter +36%;  
  Farmer and Builder work rates are normal  
  (stated Villager work rate +25%).
* Villager armor +1 and pierce armor +1  
  (in-game not show pierce armor, but it is there).
* Villager cost +50% (so 75 food instead of 50).
* Camel Rider move speed +25% (as fast as a Heavy Horse Archer).
* Free tribute
  (so while other civilizations lose 125 resources to give ally 100,
  Palmyran gives 100 resources without losing anything).

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
  Tower cost -50%.
* Swordsmen attack speed +50%
  (stated +33%; this probably refers to attack reload time).

#### 14. Shang

* Villager cost 35 food instead of 50 food (-30%).
* Wall HP x2.

#### 15. Sumerian

* Villager HP +15 (so 40 instead of 25).
* Siege units attack speed +43%.
  (reload time reduced to 3.5 instead of 5).
* Farm food +250 (starts at 500 instead of 250).

#### 16. Yamato

* Villager move speed +18% (stated +30%, same as Assyrian).
* Mounted units cost -25%  
  (Cavalry cost reduced to 52 food, 60 gold,  
  Horse Archer cost reduced to 37 food, 52 gold).
* Ships HP +30%.

### Balance data to Definitive Edition and Return of Rome

In-game unit stats are defined in the file `data\empires.dat`.
In the same directory, the file `empires_definitive_edition.dat` is a copy with
changes that mimic the official balance patches of the
[Age of Empires: Definitive Edition](https://ageofempires.fandom.com/wiki/Summary_of_changes_in_Age_of_Empires:_Definitive_Edition)
and a few changes from [Return of Rome](https://ageofempires.fandom.com/wiki/Age_of_Empires_II:_Definitive_Edition_-_Return_of_Rome).

To play with these changes, replace `data\empires.dat` with `data\empires_definitive_edition.dat`.
To revert, restore `data\empires.dat.backup`.
Script to do this [switch_data_definitive_edition.sh](switch_data_definitive_edition.sh).

Because I do not want to nerf good things, I only apply the buffs in the following list.

#### Units and buildings update

##### 1. Villager

- [x] Hunters work 5% faster than before
  (previously they were as slow as Foragers 0.45 food/s, now 0.4725 food/s).
- [x] Fix Palmyran Farmer Work Rate not increased.
- [x] Villagers move 10% faster after advancing to the Tool Age.
  (Tool Age effect now adds +0.1 to Civilian Movement, so the speed increased to 1.2,
  Wheel effect adjusted accordingly to +0.6, so Wheel Villager stays at 1.8).
- [ ] Gold and stone miners work 15% faster (to compensate bug fixes to technologies).
- [x] Increased the quantity of Gold Mine to 450 (previously 400).
- [x] Increased the quantity of Stone Mine to 300 (previously 250).
- [x] Tower build time reduced to 72s (previously 80s).
- [x] Farm build time reduced to 24s (previously 30s).
- [x] Farm technologies improved, fully upgraded Farm now has 550 food (prev 475):
  - Domestication: cost reduced to 150 food, 50 wood (previously 200 food, 50 wood).
  - Plow: +100 Farm food (previously +75).
  - Irrigation: +125 Farm food (previously +75).

##### 2. Barracks

- [x] Short Sword research is free upon reaching the Bronze Age:
  - before: costs Food 120, Gold 50, Research Time 50, require Battle Axe
  - after: costs Food 0, Gold 0, Research Time 1, only require Bronze Age
- [x] Broad Swordsmen have +10 HP, so now have 80 HP.
- [x] Long Swordsmen have +20 HP, so now have 100 HP.
- [ ] Long Sword upgrade cost increased to 240 food, 100 gold
  (previously 160 food, 50 gold).
- [ ] Legions lose 20 HP.
- [ ] Slinger bonus attack vs mounted archers increased to +4  
  (bonus vs foot archers unchanged at +2).
- [ ] Slinger training time increased from 28 to 35 seconds.

##### 3. Archery Range

- [ ] Horse Archer and Heavy Horse Archer pierce armor reduced to 1 (previously 2).
- [x] Elephant Archer attack increased to 6 (previously 5)  
  (cost unchanged 180 food, 60 gold).
- [x] Improved Bow research time reduced to 45s (previously 60s).

##### 4. Stable

- [x] Scout cost reduced to 90 food (previously 100).
- [x] Cavalry has 0/1 melee/pierce armor (previously 0/0).
- [x] Heavy Cavalry has 2/2 melee/pierce armor (previously 1/1).
- [x] Cataphract has 240 HP, 5/3 melee/pierce armor (previously 180 HP, 3/1 armor).
- [x] Cataphract upgrade cost reduced to 1600 food, 600 gold (previously 2000 food, 850 gold).
- [ ] Camel Riders now have bonus attack against all mounted units. Missing +4 bonus attack vs elephants added.
- [ ] All chariots have no bonus damage against Priests, since last update. Their conversion resistance was changed to 2x (previously 8x).
- [ ] Scythe Chariot upgrade cost increased to 1400 wood, 1000 gold (previously: 1200 wood, 800 gold). Their melee armor reduced to 1 (previously 2).
- [ ] Armored Elephant bonus attack vs buildings decreased by 1.
- [ ] Elephant and Scythe Chariot trample damage area reduced just enough so that they no longer damage enemy units/buildings on the other side of a wall.

##### 5. Academy

- [x] Academy cost reduced to 150 wood (previously 200 wood).

##### 6. Siege Workshop

- [ ] Catapult and Heavy Catapult have slightly smaller damage area radius.
- [ ] Siege weapon projectiles (stones and bolts) travel slightly faster (they were very slow), but are also slightly slower to reload (affects all stone and bolt firing units).
- [ ] Helepolis has 23% slower fire rate (previously it fired as fast as archers). Gets +5 attack.
- [x] Helepolis upgrade cost reduced to 1200 food, 1000 wood (previously 1500 food, 1000 wood).

##### 7. Temple

- [ ] Priest have +1 healing range.
- [x] Sacrifice costs 400 gold (previously 600 gold).

##### 8. Tower

- [ ] Watch Towers have 125 HP (previously 100 HP).
- [ ] Sentry Towers have 185 HP (previously 150 HP).
- [ ] Guard/Ballista Towers have 240 HP (previously 200 HP).

##### 9. Dock

- [ ] Catapult Triremes now have 135 HP (previously 120 HP).
- [ ] Catapult Triremes and Juggernaughts move 10% faster, cost reduced to 135 wood, 50 gold (previously 135 wood, 75 gold).
- [ ] Juggernaughts can no longer destroy trees. Upgrade cost reduced to 1,300 food, 500 wood (previously: 2,000 food, 900 wood).
- [ ] Trade Boats to 120 HP (previously 200); Merchant Ships to 200 HP (previously 250).

#### Civilizations update

Shared changes to all civilizations:

- [x] Writing's Research Time reduced to 30s (previously 60s).
- [x] Wheel is available for all (previously disabled for Persian and Macedonian).
- [x] Coinage is available for all (previously disabled for Egyptian, Palmyran, Persian, Shang and Sumerian).
- [ ] Heavy Transport is available for all (previously disabled for Assyrian, Babylonian, Choson, Hittite, Palmyran, Shang and Sumerian).

End of shared changes. The following are civilization-specific changes.

##### 1. Assyrian

##### 2. Babylonian

##### 3. Carthaginian

##### 4. Choson

##### 5. Egyptian

##### 6. Greek

##### 7. Hittite

##### 8. Macedonian

##### 9. Minoan

##### 10. Palmyran

##### 11. Persian

##### 12. Phoenician

##### 13. Roman

##### 14. Shang

##### 15. Sumerian

##### 16. Yamato

##### 17. Lac Viet (new civilization)

TODO: need to modify the `empires.exe` to show the new civilization in the game,
not just edit the `data/empires.dat` file.

- [x] Foragers work 20% faster.
  (Forager work rate +0.1, so 0.55 food/s instead of 0.45 food/s).
- [x] Archers have +2 melee armor.
- [x] Ballista and Helepolis have +2 melee armor.

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
