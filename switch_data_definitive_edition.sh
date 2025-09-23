#!/bin/bash

# Script to switch the data file (units, buildings, techs, civs stats) between:
# - Age of Empires: Rise of Rome v1.0
# - Age of Empires: Definitive Edition and Return of Rome balance
# - Overpowered Mace: civ Macedonian has combined bonuses from all civilizations


TARGET="data/empires.dat"  # actual data file used by the game

BRANCH_RISE_OF_ROME="data/empires.dat.backup"               # Rise of Rome v1.0
BRANCH_DEFINITIVE="data/empires_definitive_edition.dat"     # Definitive Edition & Return of Rome
BRANCH_MACE_OVERPOWERED="data/empires_mace_overpowered.dat" # Overpowered Mace


# SET HERE BEFORE RUNNING THE SCRIPT one of the three options:
data_branch="$BRANCH_MACE_OVERPOWERED"

if [ "$data_branch" = "$BRANCH_RISE_OF_ROME" ]; then
    if [ -f "$BRANCH_RISE_OF_ROME" ]; then
        cp -f "$BRANCH_RISE_OF_ROME" "$TARGET"
        echo "$(date --iso=s) switched to Rise of Rome v1.0 data file."
    else
        echo "ERROR: file not found: $BRANCH_RISE_OF_ROME"
        exit 1
    fi
elif [ "$data_branch" = "$BRANCH_DEFINITIVE" ]; then
    if [ -f "$BRANCH_DEFINITIVE" ]; then
        cp -f "$BRANCH_DEFINITIVE" "$TARGET"
        echo "$(date --iso=s) switched to Definitive Edition data file."
    else
        echo "ERROR: file not found: $BRANCH_DEFINITIVE"
        exit 1
    fi
elif [ "$data_branch" = "$BRANCH_MACE_OVERPOWERED" ]; then
    if [ -f "$BRANCH_MACE_OVERPOWERED" ]; then
        cp -f "$BRANCH_MACE_OVERPOWERED" "$TARGET"
        echo "$(date --iso=s) switched to Overpowered Macedonian data file."
    else
        echo "ERROR: file not found: $BRANCH_MACE_OVERPOWERED"
        exit 1
    fi
else
    echo "unknown data branch: $data_branch, make sure it is a file path"
    exit 1
fi

# - use the following command ignore change to empires.dat while you are editing:
#   git update-index --skip-worktree data/empires.dat
# - to track changes again:
#   git update-index --no-skip-worktree data/empires.dat
# - to list all files marked as skip:
#   git ls-files -v . | grep ^S
