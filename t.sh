# maybe more powerful
# for mac (sed for linux is different)
dir=`echo ${PWD##*/}`
grep "real-reading" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' "s#real-reading#$dir#g"
mv real-reading.ini $dir.ini

