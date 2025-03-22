package db

import (
    "database/sql"
	"encoding/json"
    "log"
)

type SongDBO struct {

    Id                   int       `json:"id"`
    Song_name            string    `json:"song_name"`
    Timestamps_encoding  string    `json:"timestamp_encoding"`
    Song                 []byte    `json:"song"`

}

func Find(slct string, where ...string) ([]SongDBO, []byte) {

    db := connect("djensen")
    defer db.Close()

    query := buildQuery(slct, where...)
    rows, err := db.Query(query)
    if err != nil {
        log.Fatal(err)
    }

    retDBO, ret := processRows(rows)
    return retDBO, ret

}

func buildQuery(slct string, where ...string) string {

    query := "SELECT " + slct + " FROM songs"
    if len(where) > 0 {

        query += " WHERE "
        for _, where_part := range where {
            query += where_part + " "
        }

    }

    return query

}

func Insert(song_name string, timestamps_encoding string, song []byte) {

    db := connect("djensen")
    defer db.Close()

    _, err := db.Query("INSERT INTO songs (song_name, timestamps_encoding, song) VALUES (?, ?, ?)", song_name, timestamps_encoding, song)

    if err != nil {
        log.Fatal(err)
    }
}

// processRows returns both the DBOs and the Marshalled DBOs; omit the return type you don't want
func processRows(rows *sql.Rows) ([]SongDBO, []byte) {

    var ret_array []SongDBO

    for rows.Next() {

       var b SongDBO
       err := rows.Scan(&b.Id, &b.Song_name, &b.Timestamps_encoding, &b.Song)
       if err != nil {
           log.Fatal(err)
       }
       ret_array = append(ret_array, b)

    }

    ret, err := json.Marshal(ret_array)
    if err != nil {
        log.Fatal(err)
    }

    return ret_array, ret

}
