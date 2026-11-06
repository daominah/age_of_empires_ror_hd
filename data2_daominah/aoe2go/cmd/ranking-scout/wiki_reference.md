# Wiki reference for the scout ranking

Copied on 2026-11-05 from the Age of Empires Series Wiki (Fandom),
text under CC-BY-SA, only the sections that matter for the ranking:
"Technologies", "Civilization bonuses", "Team bonuses" and the unit statistics.
Use it to double check `ranking_scout.md` and `logic/effect.go`, it is not game data.

Source pages:

- https://ageofempires.fandom.com/wiki/Scout_Cavalry_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Light_Cavalry_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Hussar_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Winged_Hussar_(Age_of_Empires_II)
- https://ageofempires.fandom.com/wiki/Magyar_Huszar
  (the Elite Magyar Huszar has no page of its own)

## Light Cavalry statistics

- Available to all civilizations except Aztecs, Inca, Mapuche, Maya, Muisca, Teutons, Tupi.
- Castle Age, trained at the Stable, 80 food, 30 seconds.
- Hit points 60, melee attack 7, attack bonus +10 vs Monk, reload 2.
- Melee armor 0, pierce armor 2, armor class Cavalry, speed 1.5.
- Line of sight 8, 10 in the Imperial Age.
- Conversion resistance: +3 min and +1 max conversion intervals, 19% chance per interval.
- Upgrade from Scout Cavalry costs 150 food 50 gold.
- Upgrade to Hussar costs 500 food 600 gold.
- Upgrade to Winged Hussar (Lithuanians and Poles) costs 600 food 800 gold.
- It is the only unit slower than its predecessor (Scout Cavalry speed 1.55).
- Trivia: Malian Light Cavalry with Farimba has 14 attack.

## Technologies, shared by the scout line

Differences between Scout Cavalry, Light Cavalry and Hussar are listed in the next section.

- Hit points:
  - Bloodlines +20, +30 for Varangians.
  - Ordo Cavalry regeneration in combat, Khitans only.
- Attack:
  - Forging +1 (+2 for Khitans).
  - Iron Casting +1 (+2 for Khitans).
  - Blast Furnace +2.
  - Manipur Cavalry +4 against Archers, Burmese only.
  - Farimba +5, Malians only (not on the Hussar page).
- Armor:
  - Scale Barding +1 melee, +1 pierce.
  - Chain Barding +1 melee, +1 pierce.
  - Plate Barding +1 melee, +2 pierce.
  - Silk Armor +1 melee, +1 pierce, Tatars only.
  - Ming Guang Armor +4 melee, Wei only.
- Conversion resistance:
  - Devotion +1 min, +1 max.
  - Faith +4 min, +4 max.
  - Heresy: dies upon getting converted.
  - First Crusade +4 min, +4 max, Sicilians only (not on the Hussar page).
- Attack speed: Stirrups +33%, Bulgarians only.
- Movement speed: Husbandry +10%.
- Resource cost: Kshatriyas -25%, Gurjaras only.
- Population cost: Aznauri Cavalry -20%, Georgians only.
- Creation speed:
  - Conscription +33%.
  - Steppe Husbandry +100%, Cumans only.
  - Chivalry +40%, Franks only (not on the Hussar page).
- Resource generation: Tuntian 1.8 food per minute, Wei only.
- Ability (not on the Hussar page):
  - Lechitic Legacy: 33% trample damage in 0.5 tile radius, Poles only.
  - Wootz Steel: attack ignores armor, Dravidians only.
- Upgrading: Conscription +33% faster, Chivalry +40% faster for Franks.
- Aura enhancements:
  - Regeneration: Liu Bei 45 hit points per minute, Shu only (not on the Hussar page).
  - Attack speed: Cao Cao +8.75%, Wei only.
  - Movement speed: Sun Jian +15%, Wu only.

## Civilization bonuses, per unit

Values differ by unit, because a bonus can be staggered by age.

| Bonus                     | Scout Cavalry                 | Light Cavalry             | Hussar              |
|---------------------------|-------------------------------|---------------------------|---------------------|
| Ordo Cavalry regeneration | 67.5 per minute               | 90 per minute             | 118.42 per minute   |
| Franks hit points         | +20% from Feudal Age          | +20%                      | not listed          |
| Mongols hit points        | +20%/+30% Castle/Imperial     | +20%/+30% Castle/Imperial | +30%                |
| Georgians regeneration    | 2/8/14 Feudal/Castle/Imperial | 8/14 Castle/Imperial      | 14 per minute       |
| Gurjaras bonus damage     | +20%/+30%/+40% F/C/I          | +30%/+40% Castle/Imperial | +40%                |
| Cumans movement speed     | +5%/+10%/+15% F/C/I           | +10%/+15% Castle/Imperial | +15%                |
| Berbers cost              | -15%/-20% Castle/Imperial     | -15%/-20% Castle/Imperial | -20%                |
| Jurchens attack speed     | +25% from Feudal Age          | +25%                      | +25%                |
| Upgrade cost              | see below                     | see below                 | Khitans +15% faster |

Bonuses that are the same for all three units:

- Attack:
  - Bengalis +2 attack vs Skirmishers (not on the Hussar page).
  - Tatars +20% damage from higher elevation (+50% instead of +25%).
- Armor:
  - Teutons +1/+2 melee in Castle/Imperial Age (Scout Cavalry only, Light Cavalry is not theirs).
  - Turks +1 pierce.
- Damage resistance:
  - Georgians -20% damage, -40% instead of -25% on higher elevation.
  - Jurchens -50% friendly fire damage.
  - Sicilians -40% bonus damage (not on the Hussar page).
- Resource cost: Magyars -15%.
- Creation speed: Khitans +15%.
- Resource generation: Danes 25% of target cost when destroying buildings (not on the Hussar page).
- Upgrading cost:
  - Scout Cavalry: Burgundians -50%, Chinese -10%/-15% Castle/Imperial, Khitans +15% faster,
    Poles -50% food, Turks free.
  - Light Cavalry: Burgundians -50%, Khitans +15% faster, Poles -50% food, Turks free.

Technology bonuses (text of the civilization bonus lists):

- Bulgarians: Blacksmith technologies that benefit the unit cost -50% food.
- Burgundians: Husbandry is 50% cheaper.
- Burmese: Devotion and Faith are 50% cheaper.
- Celts: can convert herdable animals even if enemy units are next to them.
- Chinese: technologies that benefit the unit are 5%/10%/15% cheaper in Feudal/Castle/Imperial
  (not on the Hussar page).
- Khitans: receive double the effect from Forging and Iron Casting.
- Magyars: Forging, Iron Casting and Blast Furnace are free.
- Poles: Bloodlines costs -50% food (not on the Hussar page).
- Spanish: Blacksmith technologies that benefit the unit cost no gold,
  and each such technology provides 20 gold when researched.
- Vietnamese: Conscription is free (not on the Hussar page).

## Team bonuses, scout line

- Attack: Hindustanis +2 vs standard buildings, Poles +1 vs archers, Wei +2 vs siege.
- Conversion resistance: Teutons +3 min, +1 max.
- Line of sight: Mongols +2.
- Creation speed: Huns +20%.
- Upgrading: Huns +20% faster, Portuguese +25% faster.
- Bulgarians: Blacksmith technologies that benefit the unit are researched 80% faster.
- Huns: Bloodlines and Husbandry are researched 20% faster.
- Lithuanians: Devotion, Faith and Heresy are researched 20% faster.
- Portuguese: technologies that benefit the unit are researched 25% faster.
- Wei: Cavalry units +2 attack vs siege weapons.

## Winged Hussar (Lithuanians and Poles)

- Strong vs archers, siege weapons, Monks and gunpowder units.
- Hit points: Bloodlines +20.
- Attack: Forging +1, Iron Casting +1, Blast Furnace +2 (Poles only).
- Armor: Scale Barding +1/+1, Chain Barding +1/+1, Plate Barding +1/+2 (Lithuanians only).
- Conversion resistance: Devotion, Faith, Heresy (Lithuanians only).
- Movement speed: Husbandry +10%.
- Creation speed: Conscription +33%.
- Ability: Lechitic Legacy, 33% blast damage in 0.5 tile radius, Poles only.
- Civilization bonus: Poles, Bloodlines costs -50% food.
- Team bonuses: same as the scout line above.

## Magyar Huszar

- Strong vs siege weapons, archers, other light cavalry, Ghulams, Chakram Throwers.
- Hit points: Bloodlines +20.
- Armor: Scale, Chain and Plate Barding.
- Conversion resistance: Devotion, Heresy.
- Movement speed: Husbandry +10%.
- Resource cost: Corvinian Army (costs 80 food).
- Creation speed: Conscription +33%, Kasbah +25% with a Berber ally only.
- Team bonuses: Hindustanis +2 vs standard buildings, Wei +2 vs siege, Teutons conversion.
- The page lists no Forging, Iron Casting or Blast Furnace.
