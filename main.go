package main

import (
    "database/sql"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "strings"
    _ "github.com/glebarez/go-sqlite"
)


func main() {
    port := flag.String("port", "9191", "port that will be used 9191, by default")
    addr := flag.String("addr", "", "ip address that will be used. Localhost by default")
    flag.Parse()
    connection := *addr + ":" + *port
    
    var err error
    db, err = sql.Open("sqlite", "./data.db")
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    log.Println("Connected to the SQLite database successfully.")
    createTables(db)
    
    mux := http.NewServeMux()
    
    mux.HandleFunc("GET /videos", getVideos)
    mux.HandleFunc("GET /videos/", getVideosID)
    mux.HandleFunc("GET /playlists", getPlaylists)
    mux.HandleFunc("GET /playlists/", getPlaylistID)
    mux.HandleFunc("POST /videos", uploadPlaylists)
    mux.HandleFunc("POST /playlists", uploadPlaylists)
    mux.HandleFunc("GET /playlistsDB", getPlaylistsDB)
    mux.HandleFunc("DELETE /playlists/", deletePlaylistID)
    
    mux.HandleFunc("POST /subscriptions", uploadSubscriptions)
    mux.HandleFunc("POST /channelgroups", uploadSubscriptions)
    mux.HandleFunc("GET /subscriptions", getSubs)
    mux.HandleFunc("GET /subscriptions/", getSubsID)
    mux.HandleFunc("GET /channelgroups", getChannelGroups)
    mux.HandleFunc("GET /channelgroups/", getChannelGroupsID)
    mux.HandleFunc("GET /channelgroupsDB", getChannelGroupsDB)
    mux.HandleFunc("DELETE /channelgroups/", deleteChannelGroupsID)
    
    server := &http.Server{
        Addr:    connection,
        Handler: mux,
    }
    
    log.Println("Server running on", connection)
    log.Fatal(server.ListenAndServe())
}

func getPlaylists(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(getAllPlaylists(db))
}

func getPlaylistID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/playlists/")
    for _, playlist := range getAllPlaylists(db) {
        if playlist.PlaylistName == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(playlist)
            return
        }
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"error": "Playlist not found"})
}

func getVideos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(getAllVideos(db))
}

func getPlaylistsDB(w http.ResponseWriter, r *http.Request) {
    var returnString string
    for _, playlist := range getAllPlaylists(db) {
        jsonBytes, err := json.Marshal(playlist)
        if err != nil {
            log.Println(err)
        }
        playlistString := string(jsonBytes) + "\n"
        returnString = returnString + playlistString
    }
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(returnString))
}

func getVideosID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/videos/")
    var videos []Video
    for _, video := range getAllVideos(db) {
        if video.Playlist == id {
            videos = append(videos, video)
        }
    }
    w.Header().Set("Content-Type", "application/json")
    if len(videos) > 0 {
        json.NewEncoder(w).Encode(videos)
    } else {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Playlist not found"})
    }
}

func updatePlaylist(w http.ResponseWriter, r *http.Request) {
    var updatedPlaylist Playlist
    if err := json.NewDecoder(r.Body).Decode(&updatedPlaylist); err != nil {
        // Ignore error as in original
    }
    processPlaylist(updatedPlaylist)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedPlaylist)
}

func uploadPlaylists(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // 32 MB max memory
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse form"})
        log.Println("upload failed:", err)
        return
    }
    
    file, handler, err := r.FormFile("file")
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
        log.Println("upload failed")
        return
    }
    defer file.Close()
    
    fileContent, _ := io.ReadAll(file)
    log.Println("uploaded:", handler.Filename)
    go importFtPlaylists(db, string(fileContent))
    w.Write([]byte(fmt.Sprintf("'%s' uploaded", handler.Filename)))
}

func deletePlaylistID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/playlists/")
    exists := false
    for _, playlist := range getAllPlaylists(db) {
        if playlist.PlaylistName == id {
            exists = true
        }
    }
    if !exists {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "No such playlist!"})
        return
    }
    deleteVideos(id)
    deletePlaylist(id)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Playlist deleted"})
}

func uploadSubscriptions(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse form"})
        log.Println("upload failed:", err)
        return
    }
    
    file, handler, err := r.FormFile("file")
    if err != nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save file"})
        log.Println("upload failed")
        return
    }
    defer file.Close()
    
    fileContent, _ := io.ReadAll(file)
    log.Println("uploaded:", handler.Filename)
    go importFtSubscriptions(string(fileContent))
    w.Write([]byte(fmt.Sprintf("'%s' uploaded", handler.Filename)))
}

func getSubs(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(getAllSubscriptions())
}

func getSubsID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
    var IdSubs []Subscription
    for _, sub := range getAllSubscriptions() {
        if sub.ChannelGroupName == id {
            IdSubs = append(IdSubs, sub)
        }
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(IdSubs)
    if len(IdSubs) == 0 {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Channel group not found"})
    }
}

func getChannelGroups(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(getAllChannelGroups())
}

func getChannelGroupsID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/channelgroups/")
    for _, group := range getAllChannelGroups() {
        if group.Name == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(group)
            return
        }
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotFound)
    json.NewEncoder(w).Encode(map[string]string{"error": "Channel group not found"})
}

func getChannelGroupsDB(w http.ResponseWriter, r *http.Request) {
    var returnString string
    for _, group := range getAllChannelGroups() {
        jsonBytes, err := json.Marshal(group)
        if err != nil {
            log.Println(err)
        }
        groupString := string(jsonBytes) + "\n"
        returnString = returnString + groupString
    }
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte(returnString))
}

func deleteChannelGroupsID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/channelgroups/")
    exists := false
    for _, group := range getAllChannelGroups() {
        if group.Name == id {
            exists = true
        }
    }
    if !exists {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": "No such channel group!"})
        return
    }
    deleteSubscriptions(id)
    deleteChannelGroup(id)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Channel group deleted"})
}
