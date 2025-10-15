#!/bin/bash
#A simple bash script to test things out
DATAFOLDER=path/to/freetube
# 9191 is the default port
IP=<server_ip>:9191
if [ "$1" == "push" ]; then
curl -X POST -F "file=@ $DATAFOLDER/profiles.db" $IP/subscriptions
printf "\n"
curl -X POST -F "file=@ $DATAFOLDER/playlists.db" $IP/playlists
fi

if [ "$1" == "pull" ]; then
curl -X GET $IP/channelgroupsDB > $DATAFOLDER/profiles.db
curl -X GET $IP/playlistsDB > $DATAFOLDER/playlists.db
echo "pull"
fi
