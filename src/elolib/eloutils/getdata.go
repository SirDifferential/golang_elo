package elolib

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "strconv"
    "errors"
    "sort"
)

type PlayerRating struct {
  Rank int
  Player string
  Rating int
}

type HistoryEntry struct {
  Player1 string
  Player2 string
  OldRating_p1 int
  OldRating_p2 int
  NewRating_p1 int
  NewRating_p2 int
  Expected float64
  EpochTime int64
}

type ByRating []PlayerRating
func (a ByRating) Len() int           { return len(a) }
func (a ByRating) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRating) Less(i, j int) bool { return a[i].Rating < a[j].Rating }

type ByEpoch []HistoryEntry
func (a ByEpoch) Len() int           { return len(a) }
func (a ByEpoch) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByEpoch) Less(i, j int) bool { return a[i].EpochTime < a[j].EpochTime }

func GetHistory() ([]HistoryEntry, error) {

    file, err := os.Open("history.txt")

    if err != nil {
        fmt.Println("Could not open history.txt")
        return nil, errors.New("Internal server error: Could not open ratings data")
    }

    scanner := bufio.NewScanner(file)

    var history = make([]HistoryEntry, 0)

    for scanner.Scan() {
        line := string(scanner.Text())
        delims_found := 0
        for i := 0; i < len(line); i++ {
            if rune(line[i]) == rune(':') {
                delims_found = delims_found + 1
            }
        }

        if delims_found != 7 {
            fmt.Println("invalid entry: ", line)
            continue
        }

        stringSlice := strings.Split(line, ":")

        player1 := stringSlice[0]
        player2 := stringSlice[1]

        old_r1, err := strconv.Atoi(stringSlice[2])
        if err != nil {
            fmt.Println("failed converting ", old_r1, " to int")
            continue
        }
        old_r2, err := strconv.Atoi(stringSlice[3])
        if err != nil {
            fmt.Println("failed converting ", old_r2, " to int")
            continue
        }

        new_r1, err := strconv.Atoi(stringSlice[4])
        if err != nil {
            fmt.Println("failed converting ", new_r1, " to int")
            continue
        }

        new_r2, err := strconv.Atoi(stringSlice[5])
        if err != nil {
            fmt.Println("failed converting ", new_r2, " to int")
            continue
        }

        expected, err := strconv.ParseFloat(stringSlice[6], 64)
        if err != nil {
            fmt.Println("failed converting ", expected, " to float")
            continue
        }

        epoch, err := strconv.ParseInt(stringSlice[7], 10, 64)
        if err != nil {
            fmt.Println("failed converting ", epoch, " to int")
            continue
        }

        history = append(history, HistoryEntry{Player1: player1, Player2: player2, OldRating_p1: old_r1, OldRating_p2: old_r2, NewRating_p1: new_r1, NewRating_p2: new_r2, Expected: expected, EpochTime: epoch})
    }

    sort.Sort(sort.Reverse(ByEpoch(history)))

    return history, nil
}

func GetRatings() ([]PlayerRating, error) {
    file, err := os.Open("ratings.txt")
    if err != nil {
        fmt.Println("Could not open ratings.txt")
        return nil, errors.New("Internal server error: Could not open ratings data")
    }

    scanner := bufio.NewScanner(file)

    var player_scores = make([]PlayerRating, 0)

    for scanner.Scan() {
        line := string(scanner.Text())
        delims_found := 0
        invalid := false
        for i := 0; i < len(line); i++ {
            if rune(line[i]) == rune(':') {
                delims_found = delims_found + 1
                if delims_found != 1 {
                    fmt.Println("invalid entry: ", line)
                    invalid = true
                    break
                }
            }
        }

        if invalid == true {
            continue
        }

        stringSlice := strings.Split(line, ":")

        player := stringSlice[0]

        score, err := strconv.Atoi(stringSlice[1])
        if err != nil {
            fmt.Println("failed converting ", score, " to int")
            continue
        }

        player_scores = append(player_scores, PlayerRating{Player: player, Rating: score})
    }

    sort.Sort(sort.Reverse(ByRating(player_scores)))
    index := 0
    player_len := len(player_scores)
    for index < player_len {
        player_scores[index].Rank = index+1
        index += 1
    }

    return player_scores, nil
}

