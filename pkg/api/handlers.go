package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/tenebresus/djensen_server/pkg/db"
)

type Song struct {
    Id int `json:"id"`
    Name string `json:"name"`
}

type SongInfo struct {
    Timestamps string `json:"timestamps"`
    Song []byte `json:"song"`
}

func Run() {

    http.HandleFunc("/api/songs", getSongs)
    http.HandleFunc("/api/songs/{song}", getSongInfo)
    http.HandleFunc("/api/songs/{song}/timestamps", getSongTimestamps)
    http.HandleFunc("/api/songs/{song}/song_data", getSongData)
    http.ListenAndServe(":8080", nil)

}

func getSongData(w http.ResponseWriter, r *http.Request) {

    log.Println("Received request for song data!")
    song_id := r.PathValue("song")
    songs, _ := db.Find("*", "id = " + song_id)

    ret := []byte(songs[0].Song)
    w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
    w.Write(ret)

    log.Println("Sent song data to client!")
}

func getSongTimestamps(w http.ResponseWriter, r *http.Request) {

    log.Println("Received request for song timestamps!")
    song_id := r.PathValue("song")
    songs, _ := db.Find("*", "id = " + song_id)

    ret := []byte(songs[0].Timestamps_encoding)
    w.Write(ret)

    log.Println("Sent timestamps to client!")

}

func getSongInfo(w http.ResponseWriter, r *http.Request) {

    song_id := r.PathValue("song")
    songs, _ := db.Find("*", "id = " + song_id)

    var song SongInfo
    song.Timestamps = songs[0].Timestamps_encoding
    song.Song = songs[0].Song

    ret, _ := json.Marshal(song)
    w.Write(ret)

}

func getSongs(w http.ResponseWriter, r *http.Request) {

    songs, _ := db.Find("*")

    var ret_arry []Song
    for _, song := range songs {

        var ret_entry Song
        ret_entry.Id = song.Id
        ret_entry.Name = song.Song_name
        ret_arry = append(ret_arry, ret_entry)

    }

    ret, _ := json.Marshal(ret_arry)
    w.Write(ret)

}
