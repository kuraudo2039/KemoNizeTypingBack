FROM golang:1.22.2-bookworm

# パッケージマネージャの更新とロケールの設定
RUN apt update && apt install -y curl locales && locale-gen ja_JP.UTF-8

# Go mod設定 https://qiita.com/Syoitu/items/f221b52231703cebe8ff
RUN go env -w GO111MODULE=on
