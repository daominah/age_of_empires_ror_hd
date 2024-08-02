IFS=$'\n'  # make newlines the only separator

for i in $(ls ../data2/*.ai | grep -v -E "Assyria|Babylon|Carthage|Choson|Egypt|Greek|Hittite|Macedon|Minoa|Palmyra|Persia|Phoenicia|Rome|Shang|Sumeria|Yamato"); do
    echo "$i";
    rm "$i";
done;
