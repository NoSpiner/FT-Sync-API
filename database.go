package main

import (
    "database/sql"
    "log"
		"encoding/json"
		"strings"
		"sync"
    _ "github.com/glebarez/go-sqlite"
)


type Video struct {
		VideoId string `json:"videoId"`
		Title string   `json:"title"`
		Author string  `json:"author"`
		AuthorId string `json:"authorId"`
		Published int `json:"published"`
		LengthSeconds int `json:"lengthSeconds"`
		TimeAdded int `json:"timeAdded"`
		PlaylistItemId string  `json:"playlistItemId"`
		Playlist string `json:"playlistName"`
}


type Playlist struct {
		PlaylistName string `json:"playlistName"`
		CreatedAt int  `json:"createdAt"`
		Description string `json:"description"`
		Playlist_id string `json:"_id"`
		LastUpdatedAt int `json:"lastUpdatedAt"`
		LastPlayedAt int `json:"lastPlayedAt"`
		Videos []*Video `json:"videos"`
}

//{"id","name","thumbnail","selected"}
type Subscription struct {
		Id string `json:"id"`
		Name string  `json:"name"`
		Thumbnail string `json:"thumbnail"`
		Selected bool `json:"selected"`
		ChannelGroupName string
}

//"name","bgColor","textColor", "_id", subscriptions
type ChannelGroup struct {
		Name string `json:"name"`
		BgColor string  `json:"bgColor"`
		TextColor string `json:"textColor"`
		ChannelGroup_id string `json:"_id"`
		Subscriptions []*Subscription `json:"subscriptions"`
}

func createTables(db *sql.DB){
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS playlists (
						PlaylistName TEXT NOT NULL UNIQUE,
            CreatedAt INTEGER,
            LastUpdatedAt INTEGER,
						LastPlayedAt INTEGER,
            Description TEXT,
						PlaylistInternalId TEXT
        )
    `)
    if err != nil {
        log.Println(err)
			}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS videos (
						VideoId TEXT NOT NULL,
            Title TEXT NOT NULL,
            Author TEXT,
            AuthorId TEXT,
            Published INTEGER DEFAULT (0),
            LengthSeconds INTEGER,
            TimeAdded INTEGER,
            Type TEXT,
						Playlist TEXT,
						PlaylistItemId TEXT PRIMARY KEY,
						FOREIGN KEY (Playlist) REFERENCES playlists (PlaylistName)
        )
    `)
    if err != nil {
        log.Println(err)
    }
//"name","bgColor","textColor", "_id", subscriptions
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS channelgroups (
						Name TEXT NOT NULL UNIQUE,
            BgColor TEXT,
            TextColor TEXT,
						ChannelGroupInternal_id TEXT UNIQUE
        )
    `)
    if err != nil {
        log.Println(err)
			}

//{"id","name","thumbnail","selected"}
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS subscriptions (
						ChannelId TEXT NOT NULL,
            Name TEXT NOT NULL,
            Thumbnail TEXT,
            Selected BOOL,
						ChannelGroup Text NOT NULL,
						FOREIGN KEY (ChannelGroup) REFERENCES channelgroups (Name)
						UNIQUE (ChannelId, ChannelGroup)
        )
    `)
    if err != nil {
        log.Println(err)
    }
	log.Println("database tables updated successfully")
}

func addPlaylist(db *sql.DB, playlist Playlist){
//Playlist_id,PlaylistName, CreatedAt, LastUpdatedAt, Description
//{"videoId":"KbTVBt8mLek","title":"万能青年旅店 [omnipotent youth society] - 冀西南林路行 (inside the cable temple) [Full Album]","author":"sunnlight","authorId":"UC74QiQdCsHuc5Q6sNSTlUoQ","published":"","lengthSeconds":2685,"timeAdded":1748850732898,"type":"video","playlistItemId":"9d85786c-778e-43f2-ad6c-41b091aa5777"}
//"videoId","title","author","authorId","published","lengthSeconds","timeAdded","type","playlistItemId"
    query := `INSERT OR IGNORE INTO playlists (PlaylistName, CreatedAt,LastUpdatedAt, LastPlayedAt, Description, PlaylistInternalId) VALUES (?,?,?,?,?,?) ON CONFLICT(PlaylistName) DO UPDATE SET LastUpdatedAt = excluded.LastUpdatedAt, LastPlayedAt = excluded.lastPlayedAt, Description = excluded.Description;`
    _, err := db.Exec(query, playlist.PlaylistName, playlist.CreatedAt, playlist.LastUpdatedAt, playlist.LastPlayedAt, playlist.Description, playlist.Playlist_id)
    if err != nil {
        log.Println("failed to add playlist:", err)
			}
}

func addVideo(tx *sql.Tx,video *Video, playlist string){
	query := "INSERT OR IGNORE INTO videos (VideoId,Title,Author, AuthorId, Published, LengthSeconds, TimeAdded, Type, playlistItemId, Playlist) VALUES (?,?,?,?,?,?,?,?,?,?)"
	Mutex.Lock()
	_,err := tx.Exec(query,video.VideoId,video.Title,video.Author,video.AuthorId,video.Published,video.LengthSeconds,video.TimeAdded,"video",video.PlaylistItemId,playlist)
	Mutex.Unlock()
	if err != nil{
	log.Println(err)
	}
}

func deleteVideos(playlistName string){
	query := "DELETE FROM videos WHERE Playlist=?"
	_,err := db.Exec(query,playlistName)
	if err != nil{
	log.Println(err)
}
}

func deletePlaylist(playlistName string){
	query := "DELETE FROM playlists WHERE PlaylistName=?"
	_,err := db.Exec(query,playlistName)
	if err != nil{
	log.Println(err)
}
}
func importFtPlaylists(db *sql.DB ,data string){
	//data, err := os.ReadFile(path)
	//if err != nil{log.Fatal(err)}
	temp := strings.ReplaceAll(string(data), "}\n{", "}SEPARATOR{")
	playlistJsonList := strings.Split(temp,"SEPARATOR")
	for _,i := range playlistJsonList{
		var playlist Playlist
		err := json.Unmarshal([]byte(i), &playlist)
		if err != nil {log.Println("Issue in unmarshaling json", err)}
		processPlaylist(playlist)
	}

}

func processPlaylist(playlist Playlist){
		addPlaylist(db,playlist)
		deleteVideos(playlist.PlaylistName)
	tx, err := db.Begin()
  if err != nil {
      log.Fatal(err)
  }
	defer tx.Rollback()

		for _,video_item := range playlist.Videos{
			addVideo(tx,video_item,playlist.PlaylistName)	
    }
		
	err = tx.Commit()
  if err != nil {
      log.Fatal(err)
  }
}

func getAllVideos(db *sql.DB)[]Video{
		var videos []Video
    query := "SELECT VideoId, Title, Author, AuthorId, Published, LengthSeconds, TimeAdded, playlistItemId, Playlist FROM videos"
    rows, err := db.Query(query)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var video Video
        err := rows.Scan(&video.VideoId, &video.Title, &video.Author, &video.AuthorId, &video.Published, &video.LengthSeconds, &video.TimeAdded, &video.PlaylistItemId, &video.Playlist)
        if err != nil {
            log.Println(err)
        }
        videos = append(videos, video)
    }
    
    // Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        log.Println(err)
    }
		return videos
}


func getAllPlaylists(db *sql.DB)[]Playlist{
    //Playlist_id,PlaylistName, CreatedAt, LastUpdatedAt, Description
		var playlists []Playlist
		var videos = getAllVideos(db)
    query := "SELECT PlaylistName, CreatedAt, LastUpdatedAt, Description, PlaylistInternalId FROM playlists"
    rows, err := db.Query(query)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var playlist Playlist
        err := rows.Scan(&playlist.PlaylistName, &playlist.CreatedAt, &playlist.LastUpdatedAt, &playlist.Description, &playlist.Playlist_id)
        if err != nil {
            log.Println(err)
        }
				var playlistVideos []*Video
				for _, video := range videos{
				if playlist.PlaylistName == video.Playlist{
					playlistVideos=append(playlistVideos, &video)
				}

				}
				playlist.Videos = playlistVideos
        playlists = append(playlists, playlist)
    }
    
    // Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        log.Println(err)
    }
		return playlists
}

func importFtSubscriptions(data string){
//{"id","name","thumbnail","selected"}
// name, bgColor, textColor, subscriptions, _id
	temp := strings.ReplaceAll(string(data), "}\n{", "}SEPARATOR{")
	subscriptionsJsonList := strings.Split(temp,"SEPARATOR")

	for _,i := range subscriptionsJsonList{
		var channelGroup ChannelGroup
		err := json.Unmarshal([]byte(i), &channelGroup)
		if err != nil {log.Println("Issue in unmarshaling json", err)}
		processChannelgroup(channelGroup)
	}

}

func processChannelgroup(channelgroup ChannelGroup){
		addChannelgroup(channelgroup)
		deleteSubscriptions(channelgroup.Name)
	tx, err := db.Begin()
	if err != nil {
	log.Fatal(err)
	}
	defer tx.Rollback()

		for _,subscription := range channelgroup.Subscriptions{
			addSub(tx, subscription,channelgroup.Name)	
    }

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func addChannelgroup(channelgroup ChannelGroup){
    query := `INSERT OR IGNORE INTO channelgroups (Name, BgColor, TextColor, ChannelGroupInternal_id) VALUES (?,?,?,?) ON CONFLICT(Name) DO UPDATE SET BgColor = excluded.BgColor, TextColor = excluded.TextColor;`

    _, err := db.Exec(query, channelgroup.Name, channelgroup.BgColor, channelgroup.TextColor, channelgroup.ChannelGroup_id)
    if err != nil {
        log.Println("failed to add playlist:", err)
			}

}

func addSub(tx *sql.Tx, subscription *Subscription, channeGroupName string){
	query := "INSERT OR IGNORE INTO subscriptions (Name, ChannelId, Thumbnail, Selected, ChannelGroup ) VALUES (?,?,?,?,?)"
	Mutex.Lock()
	_,err := tx.Exec(query,subscription.Name, subscription.Id, subscription.Thumbnail, subscription.Selected, channeGroupName)
	Mutex.Unlock()
	if err != nil{
	log.Println(err)
	}

}


func getAllSubscriptions()[]Subscription{
		var subs []Subscription
    query := "SELECT Name, ChannelId, Thumbnail, Selected, ChannelGroup  FROM subscriptions"
    rows, err := db.Query(query)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var sub Subscription
        err := rows.Scan(&sub.Name, &sub.Id, &sub.Thumbnail, &sub.Selected, &sub.ChannelGroupName)
        if err != nil {
            log.Println(err)
        }
        subs = append(subs, sub)
    }
    
    // Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        log.Println(err)
    }
		return subs
}


func getAllChannelGroups()[]ChannelGroup{
		var channelGroups []ChannelGroup
		var allSubs = getAllSubscriptions()
    query := "SELECT Name, BgColor, TextColor, ChannelGroupInternal_id FROM channelgroups"
    rows, err := db.Query(query)
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var group ChannelGroup
        err := rows.Scan(&group.Name, &group.BgColor, &group.TextColor, &group.ChannelGroup_id)
        if err != nil {
            log.Println(err)
        }
				var groupSubs []*Subscription
				for _, sub := range allSubs{
				if group.Name == sub.ChannelGroupName{
					groupSubs=append(groupSubs, &sub)
				}

				}
				group.Subscriptions = groupSubs
         channelGroups= append(channelGroups, group)
    }
    
    // Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        log.Println(err)
    }
		return channelGroups
}

func deleteSubscriptions(channeGroupName string){
	query := "DELETE FROM subscriptions WHERE ChannelGroup=?"
	_,err := db.Exec(query,channeGroupName)
	if err != nil{
	log.Println(err)
}
}

func deleteChannelGroup(channelGroupName string){
	query := "DELETE FROM channelgroups WHERE Name=?"
	_,err := db.Exec(query,channelGroupName)
	if err != nil{
	log.Println(err)
}
}


// Connect to the SQLite database
var db,_ = sql.Open("sqlite", "./data.db")
var Mutex sync.Mutex
