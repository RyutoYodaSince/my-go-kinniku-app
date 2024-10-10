package handlers

import (
    "encoding/csv"
    "net/http"
    "strconv"

    "myapp/models"
)

func DownloadCSV(w http.ResponseWriter, r *http.Request) {
    rows, err := DB.Query("SELECT id, exercise, weight, reps, sets, date FROM workout_records")
    if err != nil {
        http.Error(w, "データの取得に失敗しました", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    w.Header().Set("Content-Type", "text/csv")
    w.Header().Set("Content-Disposition", "attachment; filename=workout_records.csv")

    writer := csv.NewWriter(w)

    writer.Write([]string{"ID", "種目", "重量 (kg)", "回数", "セット数", "日付"})

    for rows.Next() {
        var record models.WorkoutRecord
        err := rows.Scan(&record.ID, &record.Exercise, &record.Weight, &record.Reps, &record.Sets, &record.Date)
        if err != nil {
            http.Error(w, "データの読み取りに失敗しました", http.StatusInternalServerError)
            return
        }

        writer.Write([]string{
            strconv.Itoa(record.ID),
            record.Exercise,
            strconv.Itoa(record.Weight),
            strconv.Itoa(record.Reps),
            strconv.Itoa(record.Sets),
            record.Date,
        })
    }

    writer.Flush()

    if err := writer.Error(); err != nil {
        http.Error(w, "CSVの書き込みに失敗しました", http.StatusInternalServerError)
        return
    }
}