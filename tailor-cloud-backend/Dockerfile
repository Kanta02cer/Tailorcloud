# マルチステージビルドでGoアプリケーションをビルド
FROM golang:1.24-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /app

# 依存関係をコピーしてダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# 本番環境用の軽量イメージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates tzdata

# 作業ディレクトリを設定
WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/main .

# フォントファイルをコピー（PDF生成用）
COPY --from=builder /app/fonts ./fonts

# ポートを公開
EXPOSE 8080

# 環境変数のデフォルト値
ENV PORT=8080

# アプリケーションを実行
CMD ["./main"]
