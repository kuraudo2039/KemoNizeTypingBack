FROM golang:1.22.2-bookworm as builder

# パッケージマネージャの更新とロケールの設定
RUN apt update && apt install -y curl locales && locale-gen ja_JP.UTF-8

# ワークディレクトリの設定
WORKDIR /app

# ソースコードのコピー
COPY ../workdir/gin_test .
RUN chown -R root:root ./
RUN ls -la
RUN pwd

# 依存関係のダウンロード
RUN go env -w GO111MODULE=on
RUN go mod download

# アプリケーションのビルド
RUN go build -buildvcs=false -o main . 

# 実行ステージ
FROM golang:1.22.2-bookworm
WORKDIR /app
COPY --from=builder /app/main .

# アプリケーションの実行
CMD ["./main"]
