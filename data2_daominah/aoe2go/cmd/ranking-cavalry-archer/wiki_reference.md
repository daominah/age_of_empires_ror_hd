# Wiki reference for the cavalry archer ranking

Copied on 2026-11-05 from the Age of Empires Series Wiki (Fandom),
text under CC-BY-SA, only the sections that matter for the ranking:
"Technologies", "Civilization bonuses" and "Team bonuses".
Use it to double check `ranking_cavalry_archer.md` and `logic/effect.go`, it is not game data.

Source pages:

- https://ageofempires.fandom.com/wiki/Cavalry_Archer_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Heavy_Cavalry_Archer_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Mounted_Crossbowman
- https://ageofempires.fandom.com/wiki/Heavy_Mounted_Crossbowman
- https://ageofempires.fandom.com/wiki/Mangudai_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Camel_Archer_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Xianbei_Raider
- https://ageofempires.fandom.com/wiki/Kipchak
- https://ageofempires.fandom.com/wiki/Ratha
- https://ageofempires.fandom.com/wiki/Bolas_Rider
- https://ageofempires.fandom.com/wiki/War_Wagon_(Age_of_Empires_II)

Not captured:

- The Mounted Crossbowman and Heavy Mounted Crossbowman pages have a different layout,
  only the start of the unit information was read, so no technology list is stored.
- The elite variants (Elite Mangudai, Elite Camel Archer, Elite Kipchak, Elite War Wagon,
  Elite Bolas Rider) have no page of their own, and "Ratha" covers the ranged mode.

## Cavalry Archer and Heavy Cavalry Archer

- Hit points: Bloodlines +20, Sipahi +20 (Turks only).
- Attack:
  - Fletching +1, Bodkin Arrow +1, Bracer +1, Chemistry +1.
  - Parthian Tactics +2 vs Spearmen.
  - Recurve Bow +1, Magyars only.
- Range: Fletching +1, Bodkin Arrow +1, Bracer +1, Recurve Bow +1 (Magyars only).
- Accuracy: Thumb Ring, Ballistics.
- Armor:
  - Padded Archer Armor +1/+1, Leather Archer Armor +1/+1, Ring Archer Armor +1/+2.
  - Parthian Tactics +1 melee, +2 pierce.
  - Silk Armor +1/+1, Tatars only.
- Conversion resistance: Devotion, Faith, Heresy.
- Attack speed: Thumb Ring +11%.
- Movement speed: Husbandry +10%.
- Population cost: Aznauri Cavalry, Georgians only, -20% (Cavalry Archer), -15% (Heavy).
- Creation speed: Conscription +33%, Steppe Husbandry +100% (Cumans only).
- Aura enhancements: Liu Bei +45 hit points per minute (Shu), Sun Jian +15% movement speed
  (Wu, Cavalry Archer page only).

Civilization bonuses, per unit:

| Bonus                  | Cavalry Archer                | Heavy Cavalry Archer |
|------------------------|-------------------------------|----------------------|
| Georgians regeneration | 8/14 per minute Castle/Imp.   | 14 per minute        |
| Cumans movement speed  | +10%/+15% Castle/Imperial     | +15%                 |
| Huns cost              | -10%/-20% Castle/Imperial     | -20%                 |
| Georgians pop cost     | -20%                          | -15%                 |
| Upgrade                | Chinese -15% cost, Khitans -50% cost and Castle Age, Shu -25% cost | Khitans: available in the Castle Age |

Same for both:

- Attack: Japanese +2 vs archers except Skirmishers, Tatars +20% from higher elevation
  (+50% instead of +25%).
- Damage resistance: Georgians -20% (-40% instead of -25% on higher elevation),
  Jurchens -50% friendly fire.
- Attack speed: Mongols +25%, Jurchens +25%.
- Resource cost: Koreans -50% wood.

Technology bonuses:

- Bulgarians: Blacksmith technologies that benefit the unit cost -50% food.
- Burmese: Devotion and Faith are 50% cheaper.
- Chinese: technologies that benefit the unit are 5%/10%/15% cheaper in Feudal/Castle/Imperial
  (Cavalry Archer page only).
- Koreans: Padded, Leather and Ring Archer Armor are free.
- Persians: Parthian Tactics is available in the Castle Age.
- Shu: Thumb Ring and Blacksmith technologies that benefit the unit are 25% cheaper.
- Tatars: Thumb Ring and Parthian Tactics are free.
- Turks: Chemistry is free.
- Vietnamese: Conscription is free.

Team bonuses:

- Teutons conversion +3 min, +1 max, Tatars line of sight +2, Britons creation speed +10%,
  Magyars creation speed +25%, Britons and Portuguese upgrade speed.
- Britons: Thumb Ring and Parthian Tactics researched 10% faster.
- Bulgarians: Blacksmith technologies researched 80% faster.
- Huns: Bloodlines and Husbandry researched 20% faster.
- Lithuanians: Devotion, Faith and Heresy researched 20% faster.
- Malians: Ballistics and Chemistry researched 80% faster.
- Portuguese: technologies that benefit the unit researched 25% faster.

## Mangudai

- Hit points: Bloodlines +20.
- Attack: Fletching, Bodkin Arrow, Bracer, Chemistry (+1 each), Parthian Tactics +2 vs Spearmen.
- Range: Fletching, Bodkin Arrow, Bracer.
- Armor: Padded +1/+1, Leather +1/+1, Parthian Tactics +1/+2 (no Ring Archer Armor).
- Attack speed: Thumb Ring +18%.
- Movement speed: Husbandry +10%.
- Creation speed: Conscription +33%, Kasbah +25% with a Berber ally only.
- No civilization bonuses listed.
- Team bonuses: Britons Thumb Ring and Parthian Tactics 10% faster, Bulgarians Blacksmith 80% faster,
  Huns, Lithuanians, Malians, Portuguese as above, Teutons, Tatars and Magyars as above.

## Camel Archer

- Hit points: Bloodlines +20.
- Regeneration: Maghrebi Camels 15 hit points per minute.
- Attack: Fletching, Bodkin Arrow, Bracer, Chemistry (+1 each), no Parthian Tactics.
- Range: Fletching, Bodkin Arrow, Bracer.
- Armor: Padded +1/+1, Leather +1/+1, Ring +1/+2.
- Attack speed: Thumb Ring +18%.
- Movement speed: Husbandry +10%.
- Creation speed: Kasbah +25% (additional +25% with a Berber ally), Conscription +33%.
- Team bonuses: Hindustanis +2 vs standard buildings, Gurjaras creation speed +25%,
  Magyars creation speed +25%, plus the common ones.

## War Wagon

- No Bloodlines, no armor technologies, no Thumb Ring listed as a technology.
- Attack: Fletching, Bodkin Arrow, Bracer, Chemistry (+1 each).
- Range: Fletching, Bodkin Arrow, Bracer.
- Accuracy: Ballistics.
- Attack speed: Thumb Ring +11%.
- Movement speed: Husbandry +10%.
- Creation speed: Conscription +33%, Kasbah +25% with a Berber ally only.
- No civilization bonuses listed.
- Team bonuses: Britons Thumb Ring 10% faster, Bulgarians Fletching, Bodkin Arrow and Bracer 80% faster,
  Huns Husbandry 20% faster, Lithuanians Devotion and Faith, Malians, Portuguese,
  Teutons, Tatars, Magyars as above.

## Xianbei Raider

- Hit points: Bloodlines +20.
- Attack: Fletching, Bodkin Arrow, Bracer, Chemistry (+1 each), Parthian Tactics +2 vs Spearmen,
  affecting both regular and charged attacks.
- Range: Fletching, Bodkin Arrow, Bracer.
- Armor: Padded +1/+1, Leather +1/+1, Ring +1/+2, Parthian Tactics +1/+2, Ming Guang Armor +4 melee.
- Attack speed: Thumb Ring +18%.
- Movement speed: Husbandry +10%.
- Creation speed: Conscription +33%.
- Resource generation: Tuntian 1.8 food per minute.
- Aura: Cao Cao +8.75% attack speed.
- Team bonuses: Teutons, Tatars, Britons +10% creation, Magyars +25% creation.

## Kipchak (mercenary and Cuman unit)

- Hit points: Bloodlines +20 (+30 Varangians), Sipahi +20 (Turks only).
- Attack: Fletching, Bodkin Arrow, Chemistry (+1 each), Bracer (allies only),
  Parthian Tactics +2 vs Spearmen, Recurve Bow +1 (Magyars only).
- Range: Fletching, Bodkin Arrow, Bracer (allies only), Recurve Bow (Magyars only).
- Armor: Padded, Leather, Ring Archer Armor, Parthian Tactics +1/+2, Silk Armor (Tatars only).
- Attack speed: Thumb Ring +18%.
- Movement speed: Husbandry +10%, allies only.
- Creation speed: Conscription +33%, Kasbah +25% for Berbers and allies.
- Other: Aznauri Cavalry -20% population space (Georgians only).
- Civilization bonuses:
  - Hit points: Franks +20%.
  - Regeneration: Georgians 8/14 per minute in Castle/Imperial Age.
  - Attack: Japanese +2 vs Archers excluding Skirmishers, Tatars elevation bonus.
  - Armor: Teutons +1/+2 melee in Castle/Imperial Age.
  - Damage resistance: Georgians, Jurchens, Sicilians -40% bonus damage.
  - Attack speed: Jurchens +20%, Mongols +25%.
  - Movement speed: Cumans +10%/+15% in Castle/Imperial Age.
  - Creation speed: Aztecs +15%.
- Technology bonuses:
  - Bulgarians: Blacksmith technologies that benefit mercenary Kipchaks cost -50% food.
  - Burmese: Devotion and Faith 50% cheaper.
  - Burgundians: Husbandry 50% cheaper.
  - Celts: Elite Kipchaks can convert herdable animals.
  - Chinese: technologies cheaper 5%/10%/15% by age.
  - Koreans: armor technologies at the Blacksmith are free.
  - Mapuche: +3 gold per military unit killed.
  - Persians: Parthian Tactics in the Castle Age.
  - Poles: Bloodlines costs -50% food.
  - Spanish: Blacksmith technologies cost no gold, 20 gold each when researched.
  - Tatars: Parthian Tactics and Thumb Ring free.
  - Turks: Chemistry free.
  - Vietnamese: Conscription free.

## Ratha (ranged mode)

- Hit points: Bloodlines +20.
- Attack: Forging, Iron Casting, Blast Furnace for the melee mode,
  Fletching, Bodkin Arrow, Bracer, Chemistry for the ranged mode.
- Range: Fletching, Bodkin Arrow, Bracer (also raises melee mode line of sight and search radius).
- Accuracy: Ballistics (ranged mode).
- Armor: Padded +1/+1, Leather +1/+1, Ring +1/+2.
- Attack speed: Paiks +20%.
- Movement speed: Husbandry +10%.
- Creation speed: Conscription +33%, Kasbah +25% with a Berber ally only.
- Team bonuses: Wei +2 vs siege (melee mode), Teutons, Tatars, Magyars +25% creation,
  Portuguese upgrade 25% faster, Bulgarians, Huns, Malians as above.

## Bolas Rider

- No Bloodlines, no Thumb Ring, no Husbandry listed.
- Attack: Fletching, Bodkin Arrow, Chemistry (+1 each).
- Range: Fletching, Bodkin Arrow.
- Accuracy: Ballistics.
- Armor: Padded +1/+1, Leather +1/+1, Ring +1/+2.
- Creation speed: Conscription +33%.
- Other: Malon, 30% pass-through damage.
- Team unique technology: Butalmapu -15% resource cost (Mapuche).
- Team bonuses: Teutons, Tatars, Britons, Magyars, Portuguese, Bulgarians, Lithuanians, Malians.
