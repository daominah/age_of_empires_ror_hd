# Game data extracted to JSON

`dat_to_json.py` extracts what the Go code needs
from the game files in [../dat_original_from_game](../dat_original_from_game),
one item per line, so a game patch shows as a small, readable git diff:

- `civilizations.json`: the `.dat` civilizations, in the `.dat` order (the ID).
  - Names: `dat_name` is the `.dat` name, some are old, e.g. "French" for Franks;
    `name_string_id` is the English name in the game strings.
  - `era` is "base" for AoE2 and "antiquity" for Chronicles.
  - What each one can use, from the game Technology Tree files `CivTechTrees/*.json`:
    `unit_ids`, `building_ids`, `technology_ids`,
    and `upgrade_technology_ids`, the technology that upgrades to a unit, e.g. Hussar.
- `units.json`: base stats of units and buildings,
  from the Gaia copy in the `.dat`, before any civilization bonus.
  Attacks and armors are amounts by armor class.
- `technologies.json`: name, effect, requirements, cost.
  `civilization_id` -1 is any civilization, if its Technology Tree allows it.
- `effects.json`: the commands of each effect, as `[type, A, B, C, D]`,
  see `logic/effect.go` for their meaning.
  Civilization bonuses are technologies too, named like "C-Bonus, ...".

The `.dat` stores numbers as 32-bit floats, e.g. 0.35 is stored as 0.3499999940395355:
the JSON files hold the shortest decimal that is the same 32-bit float, e.g. 0.35.

The Go code embeds the first three (`data/data.go`).
`logic/effect.go` copies by hand the effects the rankings need.

## How we know how to read the `.dat`

From [aoe2techtree](https://github.com/SiegeEngineers/aoe2techtree),
the data behind [aoe2techtree.net](https://aoe2techtree.net).
Its [scripts/generateDataFiles.py](https://github.com/SiegeEngineers/aoe2techtree/blob/master/scripts/generateDataFiles.py)
reads the same game files with the Python library `genieutils-py`:

- What a civilization can use comes from `CivTechTrees/*.json`:
  every node whose "Node Status" is not "NotAvailable".
- Unit stats come from the Gaia civilization copy of the units.
- A unit upgrade technology is the "Trigger Tech ID" of a node,
  skipped when it has no cost.
- A few IDs outside every Technology Tree are added by hand,
  e.g. the Feudal Age Town Center 71: we add the same ones in `cmd/generate-id`.

Checked on Update 188329: the civilization lists and unit stats here
matched its `data.json` for all 56 AoE2 civilizations.

## After a game patch

- Copy the new game files to `../dat_original_from_game`.
- Extract again (needs Python 3 and `genieutils-py`, in a venv outside this repository):

  ```bash
  python -m venv ~/.venvs/aoe2go
  ~/.venvs/aoe2go/Scripts/pip install genieutils-py
  ~/.venvs/aoe2go/Scripts/python dat_to_json.py
  ```

- Check the git diff of the JSON files, and edit `logic/effect.go` for the changed effects.
- Generate the Go ID constants again: `go generate ./logic`.
