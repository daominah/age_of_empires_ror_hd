# Zips a mod folder for upload on the ageofempires.com mod edit page.
# The folder contents sit at the zip root,
# an extra top level folder breaks the mod.
# Usage: .\zip_mod.ps1 [ModDir]
# ModDir defaults to FullTechTreeKeepBonus next to this script.
# Output: FullTechTreeKeepBonus.zip, next to the mod folder.

param(
    [string]$ModDir = "$PSScriptRoot\FullTechTreeKeepBonus"
)

$ErrorActionPreference = "Stop"

$modPath = (Resolve-Path $ModDir).Path.TrimEnd("\")
if (-not (Test-Path $modPath -PathType Container)) {
    throw "not a directory: $modPath"
}
$zipPath = "$modPath.zip"

# Windows tar writes forward slashes in entry paths,
# unlike Compress-Archive in PowerShell 5.1 which writes backslashes.
# Use the full path because Git for Windows puts GNU tar (no zip support) on PATH.
$tar = "$env:SystemRoot\System32\tar.exe"

# Listing children instead of "." avoids "./" prefixes in entry paths.
$items = Get-ChildItem -LiteralPath $modPath -Force | ForEach-Object { $_.Name }
if ($items.Count -eq 0) {
    throw "empty directory: $modPath"
}

if (Test-Path $zipPath) {
    Remove-Item $zipPath -Confirm:$false
}
& $tar -a -c -f $zipPath -C $modPath @items
if ($LASTEXITCODE -ne 0) {
    throw "tar failed with exit code $LASTEXITCODE"
}

& $tar -t -v -f $zipPath
Write-Host "created $zipPath"
