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
Pressing "S" when both a Farmer and Farms are selected at the same time will
replenish the Farm with full food quantity, as same as a new Farm.
(so researching Domestication, Plow, or Irrigation is not necessary,
food quantity bonus of the Sumerian and Minoan Farm is not important).

### Civilization bonuses

The results here come from real tests and the game dat file,
some values are slightly different from the documentation.

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

### Priest conversion test

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
