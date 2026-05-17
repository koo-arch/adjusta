# ── 1. builder stage ─────────────────────────
FROM golang:1.23-alpine AS builder
WORKDIR /app

# ビルドに必要なパッケージをインストール
RUN apk update && apk add --no-cache git

# 依存を先にキャッシュ
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーして ent generate → ビルド
COPY . .
# ent のコード生成（intercept 機能を有効に）
RUN go run -mod=mod entgo.io/ent/cmd/ent generate \
    --feature intercept,schema/snapshot \
    ./ent/schema

# Linux amd64 向けにクロスコンパイルする
ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64
RUN go build -mod=mod -o server .

# ── 2. runtime stage ─────────────────────────
FROM alpine:3.18
WORKDIR /app

# ビルド済バイナリだけコピー
COPY --from=builder /app/server .

# 必要なら環境変数や証明書をここでセット
ENV GIN_MODE=release

# ポート
EXPOSE 8080

ENTRYPOINT ["./server"]