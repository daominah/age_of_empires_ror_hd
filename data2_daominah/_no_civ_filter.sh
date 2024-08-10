IFS=$'\n'  # make newlines the only separator

# remember to change target directory path in `ls` if needed
for i in $(ls ../data2/*.ai | grep -v -E "Assyria|Babylon|Carthage|Choson|Egypt|Greek|Hittite|Macedon|Minoa|Palmyra|Persia|Phoenicia|Rome|Shang|Sumeria|Yamato"); do
   echo "$i";
   # rm "$i";
done;

# Output:
#Archers Bronze.ai
#Archers Iron.ai
#Cav Archer Iron.ai
#Cavalry Bronze.ai
#Cavalry Iron.ai
#Default.ai
#Elephant Archer Iron.ai
#Infantry Bronze.ai
#Infantry Stone.ai
#Infantry Tool.ai
#Phalanx Bronze.ai
#Phalanx Iron.ai
#Priest Bronze.ai
#Priest Iron.ai
#War Elephant Iron.ai


# ls "../data2" | grep -ivE "water|death match" | grep -i choson
