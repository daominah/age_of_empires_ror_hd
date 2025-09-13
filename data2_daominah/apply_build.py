import os
import shutil
from pathlib import Path


def main():
    sourceDir = Path(os.path.dirname(__file__))
    print(f"sourceDir: {sourceDir}")
    targetDir = os.path.join(sourceDir.parent, "data2")
    print(f"targetDir: {targetDir}")
    # return 0

    # mapBuilds maps strategy source file to targets it will be copied to,
    # for Assyria, Egypt, Greek, Minoa, Sumeria, Yamato only Immortal file is
    # meaningful, other strategies are not selected randomly.
    # remaining civilizations randomly select one of their strategies:
    # Babylon:2, Carthage:2, Choson:3, Hittite:3, Macedon:3,
    # Palmyra:3, Persia:3, Phoenicia:1, Rome:3, Shang:3
    mapBuilds = {
        "Assyria_Archer.ai": {
            # "Assyria Archer Bronze.ai",
            # "Assyria Archer Iron.ai",
            # "Assyria Ballista.ai",
            # "Assyria Bowmen.ai",
            # "Assyria Infantry Bronze.ai",
            "Immortal Assyria.ai",
        },
        "Babylon_Tower_Priest.ai": {
            "Babylon Scouts.ai",
            # "Babylon Swordsmen.ai",
        },
        "Babylon_Chariot.ai": {
            # "Babylon Scouts.ai",
            "Babylon Swordsmen.ai",
        },
        "Carthage_Helepolis.ai": {
            "Carthage Phalanx.ai",
            "Carthage War Elephant.ai",
        },
        "Choson_Swordsmen.ai": {
            "Choson Axemen.ai",
            # "Choson Priests.ai",
            "Choson Swordsmen.ai",
        },
        "Choson_Tower.ai": {
            # "Choson Axemen.ai",
            "Choson Priests.ai",
            # "Choson Swordsmen.ai",
        },
        "Egypt_Chariot_Priest.ai": {
            # "Egypt Archers Bronze.ai",
            # "Egypt Archers Iron.ai",
            # "Egypt Chariot Archer.ai",
            # "Egypt Chariots.ai",
            # "Egypt War Elephants.ai",
            "Immortal Egypt.ai",
        },
        "Greek_Centurion.ai": {
            # "Greek Phalanx.ai",
            # "Greek Priests.ai",
            # "Greek Siege.ai",
            "Immortal Greek.ai",
        },
        "Hittite_Horse_Archer.ai": {
            # "Hittite Bowmen.ai",
            "Hittite Elephant.ai",
            "Hittite Horse Archers.ai",
        },
        # very similar to Sumeria_Catapult.ai, focus on Massive Catapult
        "Hittite_Catapult.ai": {
            "Hittite Bowmen.ai",
        },
        "Macedon_Centurion.ai": {
            "Macedon Cavalry.ai",
            "Macedon Elephant.ai",
            "Macedon Phalanx.ai",
        },
        "Minoa_Bowmen_Helepolis.ai": {
            "Immortal Minoa.ai",
            # "Minoa Composite Bowmen.ai",
        },
        "Palmyra_Camel.ai": {
            "Palmyra Composite Bow.ai",
            "Palmyra Elephant.ai",
            # "Palmyra Horse Archer.ai",
        },
        "Palmyra_Horse_Archer.ai": {
            # "Palmyra Composite Bow.ai",
            # "Palmyra Elephant.ai",
            "Palmyra Horse Archer.ai",
        },
        "Persia_War_Elephant.ai": {
            "Persia Elephant Archers.ai",
            "Persia Priests.ai",
            "Persia War Elephant.ai",
        },
        "Phoenicia_Elephant_Archer.ai": {
            "Phoenicia Elephants.ai",
        },
        # "Rome_Chariot.ai": {
        #     "Rome Axemen.ai",
        # },
        "Rome_Legion.ai": {
            # "Rome Axemen.ai",
            "Rome Legion.ai",
            # "Rome Siege.ai",
        },
        "Rome_Siege.ai": {
            "Rome Axemen.ai",
            # "Rome Legion.ai",
            "Rome Siege.ai",
        },
        "Shang_Stable.ai": {
            "Shang Cavalry.ai",
            "Shang Clubmen.ai",
            # "Shang Heavy Cavalry.ai",
        },
        "Shang_Cataphract.ai": {
            # "Shang Cavalry.ai",
            # "Shang Clubmen.ai",
            "Shang Heavy Cavalry.ai",
        },
        "Sumeria_Catapult.ai": {
            "Immortal Sumeria.ai",
            # "Sumeria Catapults.ai",
            # "Sumeria Chariots.ai",
            # "Sumeria Scouts.ai",
        },
        "Yamato_Cavalry.ai": {
            "Immortal Yamato.ai",
            # "Yamato Heavy Cavalry.ai",
        },
    }

    print("___________________________________________________________________")
    usedSources = []
    unusedSources = []
    for src, targets in mapBuilds.items():
        srcPath = os.path.join(sourceDir, src)
        if not Path(srcPath).exists():
            print(f"error file {srcPath} does not exist")
            continue
        nCopiedTargets = 0
        for target in targets:
            targetPath = os.path.join(targetDir, target)
            try:
                shutil.copyfile(srcPath, targetPath)
                nCopiedTargets += 1
            except Exception as err:
                print(f"error copyfile '{srcPath}' to '{targetPath}': {err}")
        if nCopiedTargets > 0:
            usedSources.append(src)
        else:
            unusedSources.append(src)
        print(f"copied '{src}' to {nCopiedTargets} targets.")
    print("___________________________________________________________________")
    print(f"unused source files: {unusedSources}")
    print(f"used {len(usedSources)} source files: {usedSources}")


if __name__ == "__main__":
    main()
