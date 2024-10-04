package models

// WorkoutRecordは筋トレ記録の構造体です。
type WorkoutRecord struct {
	ID       int    // レコードID
	Exercise string // 種目（例: スクワット）
	Reps     int    // レップ数（回数）
	Sets     int    // セット数
	Date     string // 記録の日付
}
