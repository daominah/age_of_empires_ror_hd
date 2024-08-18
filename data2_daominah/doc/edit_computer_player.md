# Edit how computer play

Computer players follow build orders defined in `.ai` files.  
Their play is affected by `.per` and `.ply` files too, but `.ai` files are the
most important, focus on them can be enough to make computer build a good army.

## AI file

### .ai file syntax

An AI file defines a list of buildings, units, and researches that the computers
will follow in sequence.
They have extension `.ai` and are located in the `data2` folder.

Except in a special case, you will see they research Tool Age before building
enough villagers in the `.ai` files, that is because `SNUpgradeToToolAgeASAP`
is enabled in [Random Map.per](../../data2/Random Map.per).
We talk more about this in `.per` file section.

List of all buildings, units, and researches: [ai_file_build_order.pdf](./ai_file_build_order.pdf)  
(source https://aoe.heavengames.com/siegeworkshop/ai/)

### apply an edited AI file

Each civilization has some AI files corresponding to them.
Not sure if they are selected by name prefix,
so you can add any number of AI files for a civilization.

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

Should edit in dir [data2_daominah](..) then copy to files in data2 with the
corresponding civilization (can do by running [apply_build.py](../apply_build.py)).

## Helper tool

### speed editor

[speed_aoe.exe](../tool/speed_aoe.exe) can be used to change the speed while
testing AI build orders, better than `steroids`.
Recommended to set speed to 40 (double the fastest speed), not too fast so that
you can see what is happening.

### .ai file editor

[ai_edit.exe](../tool/ai_edit.exe)
(downloaded from `https://aoe.heavengames.com/dl-php/showfile.php?fileid=1669`).

Not easy to search string, so better to use normal text editor to edit then
check syntax and give summary of a build order with this tool.

Known bug: when check condition to research `Heavy Horse Archer`,
it require `Chain Mail Cavalry`, but it should be `Chain Mail Archer`.

### empires.dat editor

[AdvancedGenieEditor3.exe](../tool/genie_engine_editor/AdvancedGenieEditor3.exe)
can be used to check units stats from file `data/empires.dat`.
Should view on a cloned dir, avoiding accidentally save.

Example output for Bowman:

* hit point: 35
* move speed: 1.2
* range: 5
* sight: 7
* attack: 3
* attack reload: 1.4
* armor: 0; armor vs Slinger: -2
* train time: 30
* cost: 40 food, 20 wood

Downloaded from https://github.com/Tapsa/AGE

## PER file

Not work when start with Random Map, but work with scenarios.

Control personalities of computer players,
e.g. an option like `SNUpgradeToToolAgeASAP=1` will make computer skip some
build orders to research Tool Age first, then back to what skipped.

## PLY file

"I made some of my own plays and tested them, and yes, the did work.
It was very rewarding to see the AI come in with 3 groups of units, one made of priests.
The two groups of units split up to attempt to flank my army,
then attacked my fastest moving units (cavalry).
After taking some damage, they retreated to a point, where they were healed by
the group of priests, and then attacked again, making sure to take out the same
unit they attacked last time. Quite interesting"

Cons: the `.ply` file works for ALL computer players, in all scenarios, and random maps.
So if you make a play that is useful only in one scenario, it will be used in
all aspects of the game, and could cause some extremely stupid AI attacks.

https://aoe.heavengames.com/siegeworkshop/ply/

## References

[aoe.heavengames.com](https://aoe.heavengames.com/) is a good source,
has a lot of information (but some minor infos are wrong).
