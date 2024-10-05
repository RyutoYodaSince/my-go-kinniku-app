package handlers // パッケージ名を宣言。Goのプログラムはパッケージ単位で構成される。

import (
    "database/sql"       // データベース操作に必要なパッケージ
    "fmt"                // 標準出力フォーマットに使用するパッケージ
    "log"                // ログを記録するためのパッケージ
    "net/http"           // HTTPサーバーやクライアントのためのパッケージ
    "text/template"      // テンプレートを処理するパッケージ
    "time"               // 時間を扱うためのパッケージ

    "myapp/models"       // カスタムデータ型（WorkoutRecord）を定義している自作のパッケージ
)

// Goでは型推論を使うことができ、`var`キーワードを使用して変数を宣言する。
// 型を指定する必要はないが、必要に応じて指定もできる。
// 例えば、この場合、`*sql.DB`はデータベース接続を扱うポインタ型の変数。
var DB *sql.DB // データベース接続用の変数をグローバルに宣言

// テンプレートファイルをパースし、`tmpl`という変数に格納。
// Goでは`template.Must`を使うことで、エラーが発生した場合にパニックを発生させることができる。
var tmpl = template.Must(template.ParseFiles("/app/templates/index.html")) // テンプレートファイルをパースして読み込み

// ShowRecords関数: データベースからトレーニング記録を取得して表示する
// `func`はGoで関数を定義する際に使うキーワード
func ShowRecords(w http.ResponseWriter, r *http.Request) {
    // `:=`はGo特有の書き方で、型推論を用いて変数を宣言し、値を代入する際に使う。
    // ここでは`rows`と`err`の2つの変数を宣言し、クエリの結果とエラーを受け取る。
    rows, err := DB.Query("SELECT id, exercise, reps, sets, date FROM workout_records")
    if err != nil {
        // エラーが発生した場合、`log.Printf`でエラー内容をログに出力し、HTTPエラーを返す。
        log.Printf("クエリ実行エラー: %v", err)
        http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
        return // エラー発生時は早期リターン
    }
    defer rows.Close() // `defer`は関数が終了する時に実行される。リソースリークを防ぐために使う。

    // `var`を使って変数を宣言。`records`はWorkoutRecord型のスライスとして宣言している。
    var records []models.WorkoutRecord
    // `for rows.Next()`は、クエリの結果が次に進む限り繰り返し処理を行う。
    for rows.Next() {
        var record models.WorkoutRecord // 1つのトレーニング記録を格納するための変数
        // `rows.Scan`はSQLクエリの結果を構造体のフィールドにマッピングする。
        if err := rows.Scan(&record.ID, &record.Exercise, &record.Reps, &record.Sets, &record.Date); err != nil {
            // スキャンに失敗した場合はログにエラーメッセージを出力し、次の行に進む。
            log.Printf("行のスキャンエラー: %v", err)
            continue // エラーがあってもスライスに追加せず、ループを継続する。
        }
        records = append(records, record) // 成功したらスライスに追加
    }

    // `tmpl.Execute`で、テンプレートを実行し、HTTPレスポンスに記録を表示する。
    if err := tmpl.Execute(w, records); err != nil {
        // テンプレート実行に失敗した場合、エラーログを出力し、エラーメッセージを返す。
        log.Printf("テンプレート実行エラー: %v", err)
        http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
    }
}

// AddRecord関数: トレーニング記録を追加する
// `func`で関数を宣言し、`w http.ResponseWriter`と`r *http.Request`でHTTPレスポンスとリクエストを受け取る。
func AddRecord(w http.ResponseWriter, r *http.Request) {
    // `r.Method`でHTTPメソッドを取得し、POSTメソッドでなければリダイレクト。
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return // POSTメソッドでなければ処理を終了する
    }

    // フォームから送信されたデータを取得
    // `r.FormValue`はフォームの値を取得するために使われるメソッド。
    exercise := r.FormValue("exercise")
    reps := r.FormValue("reps")
    sets := r.FormValue("sets")
    date := time.Now().Format("2006-01-02") // `time.Now()`で現在の日付を取得し、指定フォーマットで出力

    // 入力された値が空であればエラーメッセージを表示する。
    if exercise == "" || reps == "" || sets == "" {
        http.Error(w, "全ての項目を入力してください。", http.StatusBadRequest)
        return // フォームの値が空の場合、処理を終了する
    }

    // データベースに新しいトレーニング記録を挿入する
    // `DB.Exec`はクエリを実行し、データベースに変更を加える。
    _, err := DB.Exec("INSERT INTO workout_records (exercise, reps, sets, date) VALUES ($1, $2, $3, $4)", exercise, reps, sets, date)
    if err != nil {
        // 挿入に失敗した場合はエラーログを記録し、エラーメッセージを表示する。
        log.Printf("レコード挿入エラー: %v", err)
        http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
        return // エラーが発生した場合、処理を終了する
    }

    // 成功した場合、レスポンスにメッセージを表示
    fmt.Fprint(w, "トレーニング記録が追加されました！")
}
