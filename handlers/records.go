package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"myapp/models" // modelsパッケージをインポート

	_ "github.com/go-sql-driver/mysql" // MySQLドライバをインポート
)

// HTMLテンプレートをロード
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

// ShowRecordsは筋トレ記録を表示するハンドラです。
func ShowRecords(w http.ResponseWriter, r *http.Request) {
	// データベースに接続
	db, err := sql.Open("mysql", "root:strongpassword@tcp(tidb:4000)/kinnikudb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // 処理が終わったら必ずデータベース接続を閉じる

	// データベースからレコードを取得
	rows, err := db.Query("SELECT id, exercise, reps, sets, date FROM workout_records")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // 処理が終わったら必ず結果セットを閉じる

	// 結果を格納するスライスを定義
	var records []models.WorkoutRecord

	// 結果セットをループで処理して、各レコードをスライスに追加
	for rows.Next() {
		var record models.WorkoutRecord
		if err := rows.Scan(&record.ID, &record.Exercise, &record.Reps, &record.Sets, &record.Date); err != nil {
			log.Fatal(err)
		}
		records = append(records, record)
	}

	// テンプレートにデータを渡してHTMLを生成
	tmpl.Execute(w, records)
}

// AddRecordは新しい筋トレ記録を追加するハンドラです。
func AddRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// フォームデータから入力値を取得
		exercise := r.FormValue("exercise")
		reps := r.FormValue("reps")
		sets := r.FormValue("sets")
		date := time.Now().Format("2006-01-02") // 現在の日付をフォーマット

		// 入力値のバリデーション
		if exercise == "" || reps == "" || sets == "" {
			fmt.Fprint(w, "全ての項目を入力してください。")
			return
		}

		// データベースに接続
		db, err := sql.Open("mysql", "root:strongpassword@tcp(tidb:4000)/kinnikudb")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// 新しいレコードをデータベースに挿入
		_, err = db.Exec("INSERT INTO workout_records (exercise, reps, sets, date) VALUES (?, ?, ?, ?)", exercise, reps, sets, date)
		if err != nil {
			log.Fatal(err)
		}

		// 成功メッセージを表示
		fmt.Fprint(w, "トレーニング記録が追加されました！")
	} else {
		// POSTリクエストでない場合はメインページにリダイレクト
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
