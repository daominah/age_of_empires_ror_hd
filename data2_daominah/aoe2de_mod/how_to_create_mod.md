# How to create an AoE2DE data mod

For installing a mod, see [how_to_use_mod.md](how_to_use_mod.md).

`{AoE2DEInstallDir}` is `D:\game\SteamLibrary\steamapps\common\AoE2DE`.

## Steps

- Create the mod folder skeleton:
  - `info.json`, with `Title` matching the folder name.
  - `thumbnail.jpg`, picture shown in the mod list.
  - `resources\_common\dat\empires2_x2_p1.dat`,
    copied from the same path under `{AoE2DEInstallDir}`,
    so we get the latest data from Microsoft patches.
    Re-applying the same data edits for each Microsoft dat patch
    is the pain of maintaining a data mod.
- Edit the dat copy with AdvancedGenieEditor3.
- Copy the mod folder into
  `%USERPROFILE%\Games\Age of Empires 2 DE\{Steam_ID}\mods\local`.
- Test in a Single-player Skirmish,
  selecting the mod under `Game Settings` > `Data Mods`.
- Publish from the in-game Mods browser.

## Update an existing published mod

Use this after a Microsoft patch changes the official dat.
Example: [Full Tech Tree Keep Bonus](https://www.ageofempires.com/mods/details/549269),
edited in this repo at `FullTechTreeKeepBonus`.

The in-game `My Mods` > `Update Mod` only updates title and description:
its `Select Mod Folder` picker freezes, so upload the data from the website instead.
No mod ID needs to be stored in the mod files:
the website edit page decides which mod gets updated.

- Replace the mod's `resources\_common\dat\empires2_x2_p1.dat`
  with a fresh copy from the same path under `{AoE2DEInstallDir}`.
- Re-apply the data edits with AdvancedGenieEditor3.
- Zip the mod folder contents (`info.json`, thumbnail, `resources`) so they sit at the zip root,
  by running `.\zip_mod.ps1 FullTechTreeKeepBonus` in this directory
  (the mod folder argument defaults to `FullTechTreeKeepBonus`)
  (outputs the sibling `FullTechTreeKeepBonus.zip`).
  Zipping the mod folder itself adds an extra top level folder and breaks the mod.
- Log in with Steam on the mod edit page,
  for example `https://www.ageofempires.com/mods/details/549269/edit/`,
  drop the zip on "Drop mod file here, or click to attach",
  optionally describe the patch in the change list, then submit.
