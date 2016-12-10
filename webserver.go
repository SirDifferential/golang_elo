package main

import (
    "os"
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "html/template"
    "errors"
    "time"
    "elolib/eloutils"
)

type HistoryData struct {
    Player1 string
    Player2 string
    Result string
    Change1 int
    Change2 int
    Date string
}

type PlayerData struct {
    Name string
    Rank int
    Rating int
}

type TemplateRenderData struct {
    Ratings []elolib.PlayerRating
    History []HistoryData
    Player PlayerData
}

type Page struct {
    Title string

    // Body is a byte slice as the IO libraries expect that
    Body  []byte
}

/**
* @brief: Returns true if the given string a contains the case insensitive string b
*/
func strContains(a string, b string) bool {
    a, b = strings.ToUpper(a), strings.ToUpper(b)
    return strings.Contains(a, b)
}

/**
* @brief: Renders a simple error code page
*/
func displayError(err error, w http.ResponseWriter, request string) {
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

/**
* @brief Reads a file from the file system and parses it as a Page type
* @param filename: Relative path to the file
* @return: A page instance, or nil with an error
*/
func loadPage(filename string) (*Page, error) {

    // TODO: Add a whitelist of paths that are allowed to be searched
    // And maybe limit the paths to only child dirs (no going up)
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: filename, Body: body}, nil
}

/**
* @brief Handles a single HTTP request to any non-defined request. Only index is served
* @param w: Object that writes the response to the requesting client
* @param r: The details about the request that was made to this server
*/
func handler(w http.ResponseWriter, r *http.Request) {
    // Parse out the full requested url path
    // If it has anything extra in it, don't return anything
    // The request should have gone to one of the designated handlers
    // if it was a request we should process

    path := r.URL.Path[1:]
    //fmt.Println("main handler ", path)
    if len(path) > 0 {
        displayError(errors.New("File not found"), w, path)
        return
    }

    file, err := os.Open("index.html")
    if err != nil {
        displayError(err, w, "index.html")
        return;
    }
    defer file.Close()

    ratings, err := elolib.GetRatings()

    if ratings == nil || err != nil {
        displayError(err, w, "index.html")
        return
    }

    history, err := elolib.GetHistory()
    if err != nil {
        displayError(err, w, "index.html")
        return
    }

    history_data := make([]HistoryData, 0)
    for i := range history {
        h := history[i]
        time_temp := time.Unix(h.EpochTime, 0)
        time_str := time_temp.Format("2006-01-02 15:04:05")
        history_data = append(history_data, HistoryData{Player1: h.Player1, Player2: h.Player2, Result: h.Result, Change1: h.NewRating_p1 - h.OldRating_p1, Change2: h.NewRating_p2 - h.OldRating_p2, Date: time_str})
    }

    template_render := TemplateRenderData{Ratings: ratings, History: history_data}
    t, _ := template.ParseFiles("index.html")
    err = t.Execute(w, template_render)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/**
* @brief Serves requests to URLs beginning with /player/
*/
func playerHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/player/"):]
    player_name := title

    for i := range player_name {
        if player_name[i] == '/' || player_name[i] == '.' {
            displayError(errors.New("Invalid character in player name"), w, "player")
            return
        }
    }

    //fmt.Println("serving player ", player_name)
    ratings, err := elolib.GetRatings()

    if ratings == nil || err != nil {
        displayError(err, w, "player")
        return
    }

    history, err := elolib.GetHistory()
    if err != nil {
        displayError(err, w, "player")
        return
    }

    history_data := make([]HistoryData, 0)
    for i := range history {
        h := history[i]
        if h.Player1 != player_name && h.Player2 != player_name {
            continue
        }
        time_temp := time.Unix(h.EpochTime, 0)
        time_str := time_temp.Format("2006-01-02 15:04:05")
        history_data = append(history_data, HistoryData{Player1: h.Player1, Player2: h.Player2, Result: h.Result, Change1: h.NewRating_p1 - h.OldRating_p1, Change2: h.NewRating_p2 - h.OldRating_p2, Date: time_str})
    }

    player := PlayerData{Name: player_name, Rating: 0, Rank: 0}

    found_player := false
    for i := range ratings {
        if ratings[i].Player == player_name {
            found_player = true

            player.Rating = ratings[i].Rating
            player.Rank = ratings[i].Rank
            player.Name = ratings[i].Player
            break
        }
    }

    if found_player == false {
        displayError(errors.New("Player not found"), w, "player")
        return
    }

    template_render := TemplateRenderData{Ratings: ratings, History: history_data, Player: player}
    t, _ := template.ParseFiles("player.html")
    err = t.Execute(w, template_render)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/**
* @brief: Serves any requested file to the client
*/
func genericFileHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[1:]

    //fmt.Println("serving ", path)

    // For now limit requested files to only the stylesheet
    if path != "web.css" && path != "./web.css" && path != "../web.css" {
        displayError(errors.New("File not found"), w, path)
        return;
    }

    f, err := ioutil.ReadFile(path)
    if err != nil {
        displayError(err, w, path)
        return
    }

    fmt.Fprintf(w, "%s", f);
}

func main() {
    // Handler for any URL not specified explicitly
    http.HandleFunc("/", handler)

    // Handler for the path "/player/", listing player history
    http.HandleFunc("/player/", playerHandler)

    // Serves the stylesheet
    http.HandleFunc("/web.css", genericFileHandler)
    http.ListenAndServe(":9005", nil)
}

