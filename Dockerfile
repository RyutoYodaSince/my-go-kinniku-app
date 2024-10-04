# ベースイメージとしてGoを使用
FROM golang:1.18-alpine

# 作業ディレクトリを設定
WORKDIR /app

# 依存パッケージをコピーしてインストール
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# Goアプリケーションをビルド
RUN go build -o main .

# アプリケーションを実行
CMD ["./main"]
