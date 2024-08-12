import os
import shutil
from pathlib import Path


def main():
    sourceDir = Path(os.path.dirname(__file__))
    print(f"sourceDir: {sourceDir}")
    targetDir = os.path.join(sourceDir.parent, "data2")
    print(f"targetDir: {targetDir}")
    # return 0

    # mapBuilds maps source file to targets it will be copied to
    mapBuilds = {
        "Assyria_Archer.ai": {
            "Assyria Archer Bronze.ai",
            "Assyria Archer Iron.ai",
            "Assyria Ballista.ai",
            "Assyria Bowmen.ai",
            "Assyria Infantry Bronze.ai",
            "Immortal Assyria.ai",
        },
        "Choson_Swordsmen.ai": {
            "Choson Axemen.ai",
            "Choson Priests.ai",
            "Choson Swordsmen.ai",
        },
        "Greek_Centurion.ai": {
            "Greek Phalanx.ai",
            "Greek Priests.ai",
            "Greek Siege.ai",
            "Immortal Greek.ai",
        },
        "Hittite_Horse_Archer.ai": {
            "Hittite Bowmen.ai",
            "Hittite Elephant.ai",
            "Hittite Horse Archers.ai",
        },
        "Macedon_Centurion.ai": {
            "Macedon Cavalry.ai",
            "Macedon Elephant.ai",
            "Macedon Phalanx.ai",
        },
        "Persia_War_Elephant.ai": {
            "Persia Elephant Archers.ai",
            "Persia Priests.ai",
            "Persia War Elephant.ai",
        },
        "Phoenicia_Elephant_Archer.ai": {
            "Phoenicia Elephants.ai",
        },
        # "Phoenicia_Bronze.ai": {
        #     "Phoenicia Elephants.ai",
        # },
        "Yamato_Cavalry.ai": {
            "Immortal Yamato.ai",
            "Yamato Heavy Cavalry.ai",
        },
    }
    for src, targets in mapBuilds.items():
        srcPath = os.path.join(sourceDir, src)
        if not Path(srcPath).exists():
            print(f"file {srcPath} does not exist")
            continue
        print(f"copying '{src}'")
        for target in targets:
            targetPath = os.path.join(targetDir, target)
            try:
                shutil.copyfile(srcPath, targetPath)
                print(f"    to '{target}'")
            except Exception as err:
                print(f"error copyfile '{srcPath}' to '{targetPath}': {err}")


if __name__ == "__main__":
    main()
    print("main returned")
