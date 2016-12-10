package main

import (
    "os"
    "fmt"
    "bufio"
    "strings"
    "strconv"
    "math"
    "bytes"
    "time"
)

func main() {

    if len(os.Args) < 4 {
        fmt.Println("usage: player1_name player2_name result, where result is either wins, draws or loses")
        return
    }

    player1 := string(os.Args[1])
    player2 := string(os.Args[2])
    result := string(os.Args[3])

    //fmt.Println(player1, player2, result)
    if result != "wins" && result != "loses" && result != "draws" {
        fmt.Println("result must be either wins, loses or draws")
        return
    }

    file, err := os.Open("ratings.txt")
    if err != nil {
        fmt.Println("Could not open ratings.txt")
        return
    }

    var player_scores = make(map[string]int)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := string(scanner.Text())
        delims_found := 0
        for i := 0; i < len(line); i++ {
            if rune(line[i]) == rune(':') {
                delims_found = delims_found + 1
                if delims_found != 1 {
                    fmt.Println("invalid entry: ", line)
                    continue
                }
            }
        }

        stringSlice := strings.Split(line, ":")

        player := stringSlice[0]

        score, err := strconv.Atoi(stringSlice[1])
        if err != nil {
            fmt.Println("failed converting ", score, " to int")
            continue
        }
        
        player_scores[player] = int(score)
    }

    //for k, v := range player_scores {
    //    fmt.Println("player: ", k, " score: ", v)
    //}

    rating1, ok1 := player_scores[player1]
    if ok1 == false {
        fmt.Println("did not find player in ratings: ", player1)
        return
    }

    rating2, ok2 := player_scores[player2]
    if ok2 == false {
        fmt.Println("did not find player in ratings: ", player2)
        return
    }

    trans_p1 := math.Pow(10, float64(rating1) / 400.0)
    trans_p2 := math.Pow(10, float64(rating2) / 400.0)
    expected_p1 := trans_p1 / (trans_p1 + trans_p2)
    expected_p2 := trans_p2 / (trans_p1 + trans_p2)

    s1 := 0.0
    s2 := 0.0

    if result == "wins" {
        s1 = 1.0
        s2 = 0.0
    } else if result == "loses" {
        s1 = 0.5
        s2 = 0.5
    } else {
        s1 = 0.0
        s2 = 1.0
    }

    k := 32.0

    rating1_f := float64(rating1)
    rating2_f := float64(rating2)

    new_r1 := rating1_f + k * (s1 - expected_p1)
    new_r2 := rating2_f + k * (s2 - expected_p2)

    fmt.Printf("expected: %.3f against %.3f\n", expected_p1, expected_p2)
    fmt.Printf("%s %v -> %v (change %v)\n", player1, rating1, int(new_r1), int(new_r1) - rating1)
    fmt.Printf("%s %v -> %v (change %v)\n", player2, rating2, int(new_r2), int(new_r2) - rating2)

    player_scores[player1] = int(new_r1)
    player_scores[player2] = int(new_r2)

    file_out, err := os.Create("ratings.txt")
    if err != nil {
        fmt.Println("failed opening ratings.txt for writing")
        return
    }
    defer file_out.Close()

    file_out_history, err := os.OpenFile("history.txt", os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        fmt.Println("failed opening history.txt for writing")
        return
    }
    defer file_out_history.Close()

    file_writer := bufio.NewWriter(file_out)
    file_writer_history := bufio.NewWriter(file_out_history)

    for k, v := range player_scores {
        var buffer_out bytes.Buffer
        buffer_out.WriteString(k)
        buffer_out.WriteString(":")
        buffer_out.WriteString(strconv.Itoa(v))
        _, err := fmt.Fprintln(file_writer, buffer_out.String())
        if err != nil {
            fmt.Println("error writing line to file: ", k, ", ", v)
        }
    }

    now := time.Now()

    var buffer_out_history bytes.Buffer
    buffer_out_history.WriteString(player1)
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(player2)
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(result)
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(strconv.Itoa(rating1))
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(strconv.Itoa(rating2))
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(strconv.Itoa(int(new_r1)))
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(strconv.Itoa(int(new_r2)))
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(strconv.FormatFloat(expected_p1, 'f', 6, 64))
    buffer_out_history.WriteString(":")
    buffer_out_history.WriteString(strconv.FormatInt(now.Unix(), 10))

    _, err = fmt.Fprintln(file_writer_history, buffer_out_history.String())
    if err != nil {
        fmt.Println("error writing history to file: ", err)
    }

    file_writer.Flush()
    file_writer_history.Flush()

    return
}

