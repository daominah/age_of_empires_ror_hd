# Guide References

## AI file explanation

In short, an AI file defines buildings, units, and researches that the computers will
follow in sequence.

They have extension `.ai` and are located in the `data2` folder.

[ai_file_build_order.pdf](./ai_file_build_order.pdf)  
(online: https://aoe.heavengames.com/siegeworkshop/ai/)

## Apply an edited AI file

Each civilization has some AI files corresponding to them.

Example `Hittite` has 3 AI files:

- Hittite Bowmen.ai
- Hittite Elephant.ai
- Hittite Horse Archers.ai

There are some AI files that are shared among civilizations, such as:

- Archers Bronze.ai
- Archers Iron.ai
- Cav Archer Iron.ai
- Cavalry Bronze.ai
- Cavalry Iron.ai
- Default.ai
- Elephant Archer Iron.ai
- Infantry Bronze.ai
- Infantry Stone.ai
- Infantry Tool.ai
- Phalanx Bronze.ai
- Phalanx Iron.ai
- Priest Bronze.ai
- Priest Iron.ai
- War Elephant Iron.ai

Computer with randomly select a suitable AI file to follow.
Shared build orders are not good,
so I removed them to force computer to use civilization-specific AI files.
(deleted files have backup in dir [data2_no_edit_20240801](../../data2_no_edit_20240801)).

Should edit in dir data2_daominah then copy to files in data2 with the
corresponding civilization. Example: 3 files for Hittite were copied from
  [data2_daominah](../../data2_daominah).

Example 1 file `data2_daominah/Hittite_Horse_Archer.ai` is copied to 3 following files in dir `data2`:  

All the Hittite computers play the same way, we can add more build later.

## Helper tool

[Age of Empires AI Editor](./AIEDIT.exe)
(written by 1999 Stoyan Ratchev,
downloaded from `https://aoe.heavengames.com/dl-php/showfile.php?fileid=1669`)
can be helpful to edit AI files.
It can check syntax and give summary of a build order.

````text
Example statistics from Persian_War_Elephants.ai:
Units: 107 total
    48 War Elephant
    4 Elephant Archer
    4 Camel Rider
    1 Scout
    50 Villager
Buildings: 90 total
    5 Town Center
    1 Granary (auto build)
    1 Storage Pit (auto build)
    1 Barracks
    34 Farm
    10 Tower
    6 Stable
    1 Archery Range
    1 Market
    1 Government Center
    1 Temple
Research: 16 total
````

Known bug: when check condition to research `Heavy Horse Archer`,
it require `Chain Mail Cavalry`, but it should be `Chain Mail Archer`.
