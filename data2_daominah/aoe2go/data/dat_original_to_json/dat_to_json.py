#!/usr/bin/env python3
"""Extract what the Go code needs from the AoE2 DE game files in ../dat_original_from_game
into JSON files next to this script:

- civilizations.json: .dat civilizations, with what each one can use
  from the game Technology Tree files CivTechTrees/*.json.
- units.json: base stats of units and buildings (the Gaia copy, before any bonus).
- technologies.json: technologies, with their effect.
- effects.json: the commands of each effect.

One item per line, so after a game patch the git diff is small and readable.

Usage: python dat_to_json.py [path to the dat_original_from_game directory]
"""
import json
import struct
import sys
from pathlib import Path

from genieutils.datfile import DatFile

script_dir = Path(__file__).resolve().parent
default_game_dir = script_dir.parent / 'dat_original_from_game'

# resource types of the .dat costs
resource_names = {0: 'food', 1: 'wood', 2: 'stone', 3: 'gold'}
# unit types of the .dat: 70 is a creatable unit, 80 a building
exported_unit_types = {70, 80}


def main():
    game_dir = Path(sys.argv[1]) if len(sys.argv) > 1 else default_game_dir
    dat_dir = game_dir / 'resources/_common/dat'
    dat = DatFile.parse(dat_dir / 'empires2_x2_p1.dat')

    write_items('civilizations.json', civilizations(dat, dat_dir))
    write_items('units.json', units(dat))
    write_items('technologies.json', technologies(dat))
    write_items('effects.json', effects(dat))


def civilizations(dat, dat_dir):
    """The game civilizations.json lists the civilizations in the .dat order."""
    game_list = json.loads((dat_dir / 'civilizations.json').read_text(encoding='utf-8'))['civilization_list']
    if len(game_list) != len(dat.civs):
        raise SystemExit(f'civilizations.json has {len(game_list)} civilizations, the .dat {len(dat.civs)}')
    result = []
    for i, (civ, game) in enumerate(zip(dat.civs, game_list)):
        item = {
            'id': i,
            'dat_name': civ.name,  # some are old, e.g. "French" for Franks
            'internal_name': game['internal_name'],  # e.g. "Indians" for Hindustanis
            'name_string_id': game['name_string_id'],
            'era': game['era'],  # "base" for AoE2, "antiquity" for Chronicles
            'tech_tree_effect_id': civ.tech_tree_id,
            'team_bonus_effect_id': civ.team_bonus_id,
        }
        item.update(available_ids(dat_dir / 'CivTechTrees' / f'{game["tech_tree_name"]}.json'))
        result.append(item)
    return result


def available_ids(tree_path):
    """What the civilization can use, from its Technology Tree file:
    every node not "NotAvailable", like the in-game Technology Tree.
    upgrade_technology_ids maps a unit to the technology that upgrades to it, e.g. Hussar."""
    ids = {'unit_ids': set(), 'building_ids': set(), 'technology_ids': set()}
    upgrades = {}
    if not tree_path.exists():  # Gaia has no Technology Tree
        return {k: [] for k in ids} | {'upgrade_technology_ids': {}}
    tree = json.loads(tree_path.read_text(encoding='utf-8'))
    for node in tree['civ_techs_buildings'] + tree['civ_techs_units']:
        if node.get('Trigger Tech ID', -1) > -1:
            upgrades[str(node['Node ID'])] = node['Trigger Tech ID']
        if node['Node Status'] == 'NotAvailable':
            continue
        key = {'Unit': 'unit_ids', 'Building': 'building_ids', 'Tech': 'technology_ids'}[node['Use Type']]
        ids[key].add(node['Node ID'])
    result = {k: sorted(v) for k, v in ids.items()}
    result['upgrade_technology_ids'] = dict(sorted(upgrades.items(), key=lambda kv: int(kv[0])))
    return result


def units(dat):
    result = []
    for unit in dat.civs[0].units:  # Gaia: the base stats, civilization bonuses are effects
        if not unit or unit.type not in exported_unit_types:
            continue
        item = {
            'id': unit.id,
            'internal_name': unit.name,
            'name_string_id': unit.language_dll_name,
            'type': unit.type,
            'class': unit.class_,
            'hp': unit.hit_points,
            'speed': unit.speed,
            'line_of_sight': unit.line_of_sight,
        }
        if unit.type_50:
            item['range'] = unit.type_50.max_range
            if unit.type_50.min_range:  # most units have none
                item['min_range'] = unit.type_50.min_range
            item['reload_time'] = unit.type_50.reload_time
            item['attacks'] = {str(a.class_): a.amount for a in unit.type_50.attacks}
            item['armors'] = {str(a.class_): a.amount for a in unit.type_50.armours}
        if unit.creatable and unit.type_50:
            item.update(projectiles(dat, unit.creatable))
        if unit.creatable:
            item['cost'] = cost(unit.creatable.resource_costs)
            locations = [t for t in unit.creatable.train_locations if t.unit_id >= 0]
            item['train_location_ids'] = [t.unit_id for t in locations]
            # the first slot holds the time even without a building, e.g. Town Center 109
            slots = unit.creatable.train_locations
            item['train_time'] = slots[0].train_time if slots else 0
        result.append(item)
    return result


def projectiles(dat, creatable):
    """The arrows of a unit attack, only what differs from a single arrow:
    - total_projectiles: arrows fired by each attack, e.g. 4 for the Elite Kipchak.
    - max_total_projectiles: arrows of a charged volley, e.g. 6 for the Xianbei Raider.
    - secondary_projectile: the arrow besides the main one, with its own attacks,
      e.g. 3 pierce for the Kipchak, while the main arrow uses the unit attacks.
    - charge_projectile and charge_time: the strong arrow fired once charged, and the seconds to charge,
      e.g. the Bolas Rider bolas every 30 seconds."""
    result = {}
    if creatable.total_projectiles != 1:
        result['total_projectiles'] = creatable.total_projectiles
    if creatable.max_total_projectiles != 1:
        result['max_total_projectiles'] = creatable.max_total_projectiles
    secondary = projectile(dat, creatable.secondary_projectile_unit)
    if secondary:
        result['secondary_projectile'] = secondary
    charge = projectile(dat, creatable.charge_projectile_unit)
    if charge and creatable.recharge_rate > 0:
        result['charge_projectile'] = charge
        result['charge_time'] = round(creatable.max_charge / creatable.recharge_rate, 2)
    return result


def projectile(dat, unit_id):
    """A projectile unit ID with its attacks, None if there is none."""
    units = dat.civs[0].units
    if not 0 <= unit_id < len(units) or not units[unit_id] or not units[unit_id].type_50:
        return None
    return {'id': unit_id, 'attacks': {str(a.class_): a.amount for a in units[unit_id].type_50.attacks}}


def technologies(dat):
    result = []
    for i, tech in enumerate(dat.techs):
        if not tech.name and tech.effect_id < 0:
            continue
        locations = [r for r in tech.research_locations if r.location_id >= 0]
        result.append({
            'id': i,
            'internal_name': tech.name,
            'name_string_id': tech.language_dll_name,
            'effect_id': tech.effect_id,
            'civilization_id': tech.civ,  # -1 is any civilization (if its Technology Tree allows it)
            'required_technology_ids': [t for t in tech.required_techs if t >= 0],
            'required_technology_count': tech.required_tech_count,
            'research_location_ids': [r.location_id for r in locations],
            'research_time': locations[0].research_time if locations else 0,
            'cost': cost(tech.resource_costs),
        })
    return result


def effects(dat):
    """An effect command is [type, a, b, c, d], see logic/effect.go for the meaning."""
    return [
        {
            'id': i,
            'name': effect.name,
            'commands': [[c.type, c.a, c.b, c.c, c.d] for c in effect.effect_commands],
        }
        for i, effect in enumerate(dat.effects)
        if effect.effect_commands
    ]


def cost(resource_costs):
    return {resource_names[c.type]: c.amount for c in resource_costs if c.type in resource_names and c.amount}


def shortest_float32(value):
    """The .dat stores floats as 32-bit, read back as 64-bit with noise, e.g. 0.35 as 0.3499999940395355.
    Returns the shortest decimal that is the same 32-bit float, e.g. 0.35."""
    stored = struct.pack('<f', value)
    for precision in range(1, 10):
        candidate = float(f'{value:.{precision}g}')
        if struct.pack('<f', candidate) == stored:
            return candidate
    return value


def clean_floats(item):
    """Applies shortest_float32 to every float of a JSON item."""
    if isinstance(item, float):
        return shortest_float32(item)
    if isinstance(item, dict):
        return {k: clean_floats(v) for k, v in item.items()}
    if isinstance(item, list):
        return [clean_floats(v) for v in item]
    return item


def write_items(file_name, items):
    """Write a JSON list with one item per line, and one effect command per line."""
    path = script_dir / file_name
    items = clean_floats(items)
    with open(path, 'w', encoding='utf-8', newline='\n') as f:
        f.write('[\n')
        for i, item in enumerate(items):
            separator = ',' if i < len(items) - 1 else ''
            if 'commands' not in item:
                f.write(f'  {json.dumps(item, ensure_ascii=False)}{separator}\n')
                continue
            head = {k: v for k, v in item.items() if k != 'commands'}
            f.write(f'  {json.dumps(head, ensure_ascii=False)[:-1]}, "commands": [\n')
            for j, command in enumerate(item['commands']):
                command_separator = ',' if j < len(item['commands']) - 1 else ''
                f.write(f'    {json.dumps(command)}{command_separator}\n')
            f.write(f'  ]}}{separator}\n')
        f.write(']\n')
    print(f'wrote {path}: {len(items)} items')


if __name__ == '__main__':
    main()
