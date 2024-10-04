CREATE DATABASE IF NOT EXISTS kinnikudb;

USE kinnikudb;

CREATE TABLE IF NOT EXISTS workout_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    exercise VARCHAR(255) NOT NULL,    -- 種目 (例: スクワット, ベンチプレス)
    reps INT NOT NULL,                 -- レップ数
    sets INT NOT NULL,                 -- セット数
    date DATE NOT NULL                 -- 日付
);

-- 初期データの挿入
INSERT INTO workout_records (exercise, reps, sets, date) VALUES
('Push-ups', 20, 3, '2024-10-01'),
('Squats', 15, 4, '2024-10-02'),
('Pull-ups', 10, 3, '2024-10-03');
