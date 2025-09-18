package main

import (
    "log"
    _ "github.com/glebarez/go-sqlite"
    "github.com/gin-gonic/gin"
    "net/http"
		"io"
		"fmt"
		"encoding/json"
		"flag"
)


func main() {
	port := flag.String("port","9191", "port that will be used 9191, by default")
	addr := flag.String("addr","", "ip address that will be used. Localhost by default")
	flag.Parse()
	connection := *addr+":"+*port
  log.Println("Connected to the SQLite database successfully.")
	defer db.Close()
	createTables(db)
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/videos", getVideos)
	router.GET("/videos/:id", getVideosID)
	router.GET("/playlists", getPlaylists)
	router.GET("/playlists/:id", getPlaylistID)
	router.POST("/videos", uploadPlaylists)
	router.POST("/playlists", uploadPlaylists)
	router.GET("/playlistsDB", getPlaylistsDB)
	router.DELETE("/playlists/:id", deletePlaylistID)

	router.POST("/subscriptions", uploadSubscriptions)
	router.POST("/channelgroups", uploadSubscriptions)
	router.GET("/subscriptions", getSubs)
	router.GET("/subscriptions/:id", getSubsID)
	router.GET("/channelgroups", getChannelGroups)
	router.GET("/channelgroups/:id", getChannelGroupsID)
	router.GET("/channelgroupsDB", getChannelGroupsDB)
	router.DELETE("/channelgroups/:id",deleteChannelGroupsID )
	

	log.Println("Server running on", connection)
	log.Fatal(router.Run(connection))

	}


func getPlaylists(c *gin.Context){
		c.JSON(http.StatusOK, getAllPlaylists(db))

}

func getPlaylistID(c *gin.Context){
	  id := c.Param("id")
    for _, playlist := range getAllPlaylists(db) {
        if playlist.PlaylistName == id {
   					c.JSON(http.StatusOK, playlist)
						return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
}

func getVideos(c *gin.Context) {
    c.JSON(http.StatusOK, getAllVideos(db))
}

func getPlaylistsDB(c *gin.Context) {
		var returnString string
    for _, playlist := range getAllPlaylists(db) {
			jsonBytes, err := json.Marshal(playlist)
    if err != nil {
        log.Println(err)
			}
		playlistString := string(jsonBytes) +"\n"
		returnString = returnString+playlistString
		}
       c.String(http.StatusOK, returnString)
}



func getVideosID(c *gin.Context){
	  id := c.Param("id")
    var videos []Video
    for _, video := range getAllVideos(db) {
        if video.Playlist == id {
						videos = append(videos, video)
        }
    }
   	if len(videos) >0 {c.JSON(http.StatusOK, videos)
		}else{
    c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})}
}

// I don't know why this even exists :(
func updatePlaylist(c *gin.Context) {
    var updatedPlaylist Playlist
    if err := c.BindJSON(&updatedPlaylist); err != nil {
        //c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        //return
    }
				processPlaylist(updatedPlaylist)
      	c.JSON(http.StatusOK, updatedPlaylist)
        return
}

func uploadPlaylists (c *gin.Context) {
    // single file
    file, err := c.FormFile("file")
    if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        log.Println("upload failed")
				return
    }
		openedFile,_ := file.Open()
		fileContent, _ := io.ReadAll(openedFile)
		log.Println("uploaded:", file.Filename)
    importFtPlaylists(db,string(fileContent))
    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded", file.Filename))
  }

func deletePlaylistID(c *gin.Context){
	  id := c.Param("id")
		exists := false
    for _, playlist := range getAllPlaylists(db) {
        if playlist.PlaylistName == id {
   					exists = true
        }
    }
		if ! exists{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No such playlist!"})
				return
			}
		deleteVideos(id)
		deletePlaylist(id)
		c.JSON(http.StatusOK, gin.H{"message":"Playlist deleted"})
}

func uploadSubscriptions(c *gin.Context){
    file, err := c.FormFile("file")
    if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        log.Println("upload failed")
				return
    }
		openedFile,_ := file.Open()
		fileContent, _ := io.ReadAll(openedFile)
		log.Println("uploaded:", file.Filename)
    importFtSubscriptions(string(fileContent))
    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded", file.Filename))
}

func getSubs(c *gin.Context){
		c.JSON(http.StatusOK, getAllSubscriptions())
}

func getSubsID(c *gin.Context){
	  id := c.Param("id")
		var IdSubs []Subscription
    for _, sub := range getAllSubscriptions() {
        if sub.ChannelGroupName == id {
					IdSubs = append(IdSubs,sub)
        }
    }
   	c.JSON(http.StatusOK, IdSubs)
		if len(IdSubs) ==0{
    c.JSON(http.StatusNotFound, gin.H{"error": "Channel group not found"})}
}


func getChannelGroups(c *gin.Context){
		c.JSON(http.StatusOK, getAllChannelGroups())

}

func getChannelGroupsID(c *gin.Context){
	  id := c.Param("id")
    for _, group := range getAllChannelGroups(){
        if group.Name == id {
   					c.JSON(http.StatusOK, group)
						return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"error": "Channel group not found"})
}

func getChannelGroupsDB(c *gin.Context){
		var returnString string
    for _, group := range getAllChannelGroups() {
			jsonBytes, err := json.Marshal(group)
    if err != nil {
        log.Println(err)
			}
		groupString := string(jsonBytes) +"\n"
		returnString = returnString+groupString
		}
       c.String(http.StatusOK, returnString)

}


func deleteChannelGroupsID(c *gin.Context){
	  id := c.Param("id")
		exists := false
    for _, group := range getAllChannelGroups() {
        if group.Name == id {
   					exists = true
        }
    }
		if ! exists{
				c.JSON(http.StatusInternalServerError, gin.H{"error": "No such channel group!"})
				return
			}
		deleteSubscriptions(id)
		deleteChannelGroup(id)
		c.JSON(http.StatusOK, gin.H{"message":"Channel group deleted"})
}
