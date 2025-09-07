#!/bin/bash

# Script to switch the data file (units, buildings, techs, civs stats) between:
# - Age of Empires: Rise of Rome v1.0 (default)
# - Age of Empires: Definitive Edition

# set to "true" to switch to Definitive Edition
is_change_to_definitive_edition=false
# is_change_to_definitive_edition=true

definitive_edition="data/empires_definitive_edition.dat"  # switch_dat_to_aoe_de
rise_of_rome="data/empires.dat.backup"                    # switch_dat_to_aoe_ror_v1.0

target="data/empires.dat"  # actual data file used by the game

if [ "$is_change_to_definitive_edition" = true ]; then
    if [ -f "$definitive_edition" ]; then
        cp -f "$definitive_edition" "$target"
        echo "Switched to Age of Empires: Definitive Edition data file."
    else
        echo "Definitive Edition data file not found: $definitive_edition"
        exit 1
    fi
else
    if [ -f "$rise_of_rome" ]; then
        cp -f "$rise_of_rome" "$target"
        echo "Switched to Age of Empires: Rise of Rome v1.0 data file."
    else
        echo "Rise of Rome v1.0 data file not found: $rise_of_rome"
        exit 1
    fi
fi
