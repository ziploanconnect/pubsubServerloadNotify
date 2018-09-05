NATS_ADDR=127.0.0.1:4222
LOADAVG=$(cat /proc/loadavg | cut -f1 -d" ")
NPROC=$(getconf _NPROCESSORS_ONLN)
SUBJECT="stats.loadaverage"
PAYLOAD=$(echo $(hostname) $LOADAVG $NPROC)
MESSAGE="PUB $SUBJECT ${#PAYLOAD}\r\n${PAYLOAD}\r\n"
printf "$MESSAGE" | /srv/nats/bin/catnats -q --raw --addr $NATS_ADDR --user user1 --pass pass1
