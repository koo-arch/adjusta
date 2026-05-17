FROM golang:1.23.0-alpine

WORKDIR /backend

# Gitなどのツールのインストール
RUN apk update && apk add --no-cache git

# 必要なツールのインストール
RUN go install github.com/air-verse/air@v1.52.2

# ローカルのファイルをコンテナにコピー
COPY . .

# モジュールのダウンロード
RUN go mod download

# `air`の設定ファイルをコピー
COPY .air.toml /backend

# ビルド
RUN go build -o main .

# ポートの公開
EXPOSE 8080

# コンテナ起動時のコマンド
CMD ["air"]