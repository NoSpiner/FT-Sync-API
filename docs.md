## Uploading

Examples are given using curl. Python request library works as well.

> curl -X POST -F "file=@./profiles.db" localhost:9191/subscriptions 

Upload subscriptions and channel groups with POST request


> curl -X POST -F "file=@./playlists.db" localhost:9191/playlists

Uploads playlists 


Both also update and replacce the existing data. This works with any playlist or subscription file or valid .json equivalent.

## Viewing data in browser or other GET request


>localhost:9191/playlists

Shows a list of all the saved playlists in json format


>localhost:9191/playlists/Cooking

Shows only a playlist called "Cooking"


>localhost:9191/videos

List all saved videos in json format


>localhost:9191/videos/Cooking

List all saved videos in json format in the playlist "Cooking"


--------------------------------

>localhost:9191/channelgroups

List all channel groups in json format


>localhost:9191/channelgroups/Cooking

Returns only the channelgroup called cooking

>localhost:9191/subscriptions

Returns all subscribed channels from all channel groups

>localhost:9191/subscriptions

Returns all subscribed channels from all channel groups

>localhost:9191/subscriptions/Cooking

Returns all subscribed channels from channelgroup called Cooking

## import data back to freetube .db format

>curl -X GET localhost:9191/playlistsDB

or

>localhost:9191/playlistsDB

Returns all the saved playlists in string format that makes a valid .db file

>curl -X GET localhost:9191/channelgroupsDB

or

>localhost:9191/channelgroupsDB

Returns all the saved channelgroups and subscriptions in string format that makes a valid .db file

## deleting playlists and channelgroups

curl -X DELETE 'localhost:9191/channelgroups/Food'  

Deletes channelgroup called Food and all it's subscriptions

curl -X DELETE 'localhost:9191/playlists/Music'

Deletes playlist called Music and all it's related videos
