INIPATH="/etc/comskip/comskip.ini"

comskip --ini $INIPATH "$1"
echo "$1" | sed -e 's/.mp4//g' | xargs -I{} rm {}.txt
