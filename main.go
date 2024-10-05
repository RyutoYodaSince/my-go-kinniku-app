package main // Goプログラムのエントリーポイントとして使うパッケージ

import (
    "database/sql"      // データベース操作のためのパッケージ
    "fmt"               // 文字列フォーマット用のパッケージ
    "log"               // ログを出力するためのパッケージ
    "net/http"          // HTTPサーバーやクライアントを作成するためのパッケージ
    "os"                // 環境変数などのシステム操作を扱うためのパッケージ
    "time"              // 時間を扱うためのパッケージ

    "myapp/handlers"    // ハンドラ関数やDB接続の処理を定義している自作パッケージ
    _ "github.com/lib/pq" // PostgreSQLドライバをインポート（使用しない変数やパッケージに対して`_`を使う）
)

func main() {
    // 環境変数を利用してデータベース接続情報を設定
    // `fmt.Sprintf`は指定したフォーマットに従って文字列を生成する関数。
    // `os.Getenv`は環境変数を取得する関数で、データベースのホスト、ユーザー、パスワード、DB名を指定。
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"))

    // データベース接続の初期化
    var err error // エラーを格納するための変数
    // データベース接続を最大30回試行
    for i := 0; i < 30; i++ {
        // データベース接続を開く。`sql.Open`はDBに接続せず、実際の接続はPingメソッドで確認する。
        handlers.DB, err = sql.Open("postgres", dsn)
        if err == nil { 
            // DB接続が成功しているか確認するためにPingを実行
            err = handlers.DB.Ping()
            if err == nil {
                // 成功したらループを抜ける
                break
            }
        }
        // 接続に失敗した場合、再試行するためにログを出力
        log.Printf("データベース接続を試行中... (%d/30)", i+1)
        // 2秒待ってから再試行
        time.Sleep(2 * time.Second)
    }
    // 30回試行しても接続できなければエラーを出力しプログラムを終了
    if err != nil {
        log.Fatalf("データベース接続に失敗しました: %v", err)
    }
    defer handlers.DB.Close() // プログラム終了時にデータベース接続を閉じる

    // サーバーの開始ログを出力
    log.Println("サーバーがポート8080で起動しています...")
    // URLパスとハンドラ関数を関連付ける
    http.HandleFunc("/", handlers.ShowRecords)  // "/"パスにアクセスがあった時、ShowRecordsを実行
    http.HandleFunc("/add", handlers.AddRecord) // "/add"パスにアクセスがあった時、AddRecordを実行
    // HTTPサーバーをポート8080で開始。エラーが発生した場合はログに記録して終了。
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
