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
