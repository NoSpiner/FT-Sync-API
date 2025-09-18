# FT-Sync-API
API for saving and syncronizing [FreeTube](https://freetubeapp.io/) playlists and subscriptions over the network.

[documentation](docs.md)

No client included, but with a bit tinkering should be possible. (probably like 4 curl commands in bash or 10 lines of python)

Remember to close Freetube if tinkering with playlists.db and profiles.db files to avoid file corruption. 
Also ffs don't open this to the internet. It's got no authethficiation.


Some launch options

    -addr string
 
      ip address that will be used. Localhost by default
 
    -port string
 
      port that will be used 9191, by default
