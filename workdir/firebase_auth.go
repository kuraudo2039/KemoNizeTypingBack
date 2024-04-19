// コンテキストとは
// https://zenn.dev/hsaki/books/golang-context/viewer/definition
// TODO_2024年4月18日_モジュール分割してfirestore連携APIデプロイ

import (
  "log"

  firebase "firebase.google.com/go"
  "google.golang.org/api/option"
)

// Use a service account
ctx := context.Background()
sa := option.WithCredentialsFile("secrets/kemonizetyping-322c43ee8d50.json")
app, err := firebase.NewApp(ctx, nil, sa)
if err != nil {
  log.Fatalln(err)
}

client, err := app.Firestore(ctx)
if err != nil {
  log.Fatalln(err)
}
defer client.Close()