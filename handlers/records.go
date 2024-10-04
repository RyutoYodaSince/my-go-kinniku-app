package handlers

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "text/template"
    "time"

    "myapp/models"
)

var DB *sql.DB
var tmpl = template.Must(template.ParseFiles("/app/templates/index.html"))

func ShowRecords(w http.ResponseWriter, r *http.Request) {
    rows, err := DB.Query("SELECT id, exercise, reps, sets, date FROM workout_records")
    if err != nil {
        log.Printf("クエリ実行エラー: %v", err)
        http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var records []models.WorkoutRecord
    for rows.Next() {
        var record models.WorkoutRecord
        if err := rows.Scan(&record.ID, &record.Exercise, &record.Reps, &record.Sets, &record.Date); err != nil {
            log.Printf("行のスキャンエラー: %v", err)
            continue
        }
        records = append(records, record)
    }

    if err := tmpl.Execute(w, records); err != nil {
        log.Printf("テンプレート実行エラー: %v", err)
        http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
    }
}

func AddRecord(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    exercise := r.FormValue("exercise")
    reps := r.FormValue("reps")
    sets := r.FormValue("sets")
    date := time.Now().Format("2006-01-02")

    if exercise == "" || reps == "" || sets == "" {
        http.Error(w, "全ての項目を入力してください。", http.StatusBadRequest)
        return
    }

    _, err := DB.Exec("INSERT INTO workout_records (exercise, reps, sets, date) VALUES ($1, $2, $3, $4)", exercise, reps, sets, date)
    if err != nil {
        log.Printf("レコード挿入エラー: %v", err)
        http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
        return
    }

    fmt.Fprint(w, "トレーニング記録が追加されました！")
}