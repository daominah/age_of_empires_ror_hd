# Age of Empires: The Rise of Rome

## Download

From [github.com/daominah/age_of_empires_ror_hd](https://github.com/daominah/age_of_empires_ror_hd) (this page)
click on green button `Code` then `Download ZIP` (size about 100 MB).

## Run the game

* Init setup: run file `SETUP.EXE`

* For single player, run file `empires.exe`,
  cheat enabled, population limit from 200 to 1000.

* For multi player, run file `EMPIRESX.EXE`, if someone uses cheat code in the game,
  theirs units will be deleted then kicked from the game.

## Edit AI (computer player)

Guides to edit how computer player plays are in directory
[data2_daominah/doc](data2_daominah/doc/edit_computer_player.md).

## AoE version detail

### Farm replenish bug available

The game version is The Rise of Rome 1.0, pressing "S" when both a Farmer and
Farms are selected at the same time will replenish the Farm with full food
(so researching Domestication, Plow, Irrigation is not necessary,
the same as Sumerian and Minoan Farm bonus).

### Civilization bonuses

The result here comes from real tests and compares with the game document.

#### 1. Assyrian

* Villager move speed +18% (stated +30%):  
  +0.2 tiles/second, starting at 1.3 instead of 1.1.
* Archers attack speed +36% (stated +40%):  
  attack reload time reduced to 1.1 instead of 1.5.

#### 2. Babylonian

* Priest rejuvenation +0.75: starting at 2.75 instead of 2 (so +38% but stated +50%).  
  So Priest's rest duration is 36s instead of 50s,
  with Fanaticism, it is 24s instead of 29s.
* Stone Miner work rate +44%  and capacity +3 (stated 30%).
* Tower and Wall hit points x2.

#### 3. Carthaginian

* Elephant and Academy units hit points +25%.
* Light Transport move speed +25%, Heavy Transport move speed +43% (stated 30%),
  so Heavy Transport moves as fast as Heavy Horse Archer.
* Fire Galley attack +6 (24+12 instead of 24+6).

#### 4. Choson

* Priest cost -32% (stated -30%): 85 gold instead of 125.
* Iron Age Swordsmen +80 HP
  (Long Swordsman HP is 160 instead of 80, Legion 240 instead of 160).
* Towers +2 range.

#### 5. Egyptian

* Chariots +33% HP (Chariot Archer HP starts at 93 instead of 70,
  Scythe Chariot HP is 182 instead of 137).
* Priest +3 range (from 10+3 to 10+6).
* Gold Miner +44% work rate and capacity +2 (stated +20%).

#### 6. Greek

* Academy units move speed +33% (stated +30%):  
  +0.3 tiles/s, Hoplite starts at 1.2 instead of 0.9,
  with Aristocracy, the speed is 1.45 instead of 1.15.
* War ships move speed +17% (stated +30%).

#### 7. Hittite

* Archers +1 attack.
* Siege units x2 HP.
* Warship +4 range.

#### 8. Macedonian

* Academy units +2 pierce armor.
* Melee units +2 sight.
* Siege units cost -50%.
* All units are 4 times more resistant to conversion
  (not sure how conversion works, will have tests in the next chapter).

#### 9. Minoan

* Composite Bowman +2 Composite Bowman range
* Farm have +60 food (starting at 310 instead of 250).
* Ships cost -30%.

#### 10. Palmyran

* Forager, Hunter, Gold Miner, Stone Miner work rate +44%, Woodcutter +36%,  
  (Farmer and Builder work as same as normal Villager, stated Villager work rate +25%).
* Villager +1 armor and +1 pierce armor (in game only show +1 melee armor).
* Villager cost +50% (so 75 food instead of 50).
* Camel Rider move speed +25% (so as fast as Heavy Horse Archer).
* Free tribute
  (instead of 25% taxed, other civilizations lost 125 resource to give 100).

#### 11. Persian

* Hunter work rate +66% and capacity +3 (stated +30%).
* Elephants move speed +56% (stated +50%).
* Trireme attack speed +38% (stated 50%).

#### 12. Phoenician

#### 13. Roman

#### 14. Shang

#### 15. Sumerian

#### 16. Yamato

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
