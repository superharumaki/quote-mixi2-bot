# quote-mixi2-bot

mixi2 に名言を自動投稿する Bot です。

## 機能

* `quotes.json` から名言を読み込み
* 未投稿の名言をランダム選択
* 投稿済みを `state.json` に記録
* 全件投稿後に自動リセット
* GitHub Actions で毎日自動投稿

---

# 自動投稿時間

GitHub Actions の cron で毎日実行しています。

```yaml
cron: '17 21 * * *'
```

UTC基準なので、日本時間では朝6時17分ごろ実行されます。

---

# 必要なSecrets

GitHub：

Settings
→ Secrets and variables
→ Actions
→ New repository secret

必要なもの：

```text
CLIENT_ID
CLIENT_SECRET
```

---

# 手動実行

GitHub：

Actions
→ Quote Bot
→ Run workflow

---

# 名言追加方法

## 1. Excel に追加

Excel の一覧に新しい名言を追加します。

| text | author |
| ---- | ------ |
| 名言本文 | 出典     |

※ 既存の行は削除せず、そのまま残してください。

---

## 2. CSV 保存

Excel で：

* 名前を付けて保存
* 「Unicode テキスト」で保存

保存後：

```text
quotes.txt
↓
quotes.csv
```

に名前変更します。

---

## 3. Codespaces にアップロード

`quotes.csv` を Codespaces にドラッグ＆ドロップします。

---

## 4. quotes.json に変換

ターミナルで実行：

```bash
python3 convert.py
```

---

## 5. GitHub に反映

```bash
git pull --rebase
git add quotes.json
git commit -m "update quotes"
git push
```

---

# 投稿済み管理

`state.json` で投稿済み名言を管理しています。

例：

```json
{
  "posted_indexes": [0, 3, 5]
}
```

---

# 投稿リセット

最初から再投稿したい場合：

```json
{
  "posted_indexes": []
}
```

に変更。

その後：

```bash
git add state.json
git commit -m "reset posted quotes"
git push
```

---

# 投稿削除方法

mixi2 の投稿IDが分かっている場合、ターミナルから削除できます。

## 1. delete.go を作成

```go
package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"

	"github.com/mixigroup/mixi2-application-sdk-go/auth"
	application_apiv1 "github.com/mixigroup/mixi2-application-sdk-go/gen/go/social/mixi/application/service/application_api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	authenticator, err := auth.NewAuthenticator(
		os.Getenv("CLIENT_ID"),
		os.Getenv("CLIENT_SECRET"),
		os.Getenv("TOKEN_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}

	authCtx, err := authenticator.AuthorizedContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(
		os.Getenv("API_ADDRESS"),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := application_apiv1.NewApplicationServiceClient(conn)

	postID := "ここに投稿ID"

	_, err = client.DeletePost(authCtx, &application_apiv1.DeletePostRequest{
		PostId: postID,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("削除成功")
}
```

---

## 2. 環境変数設定

```bash
export CLIENT_ID=xxxxxxxx
export CLIENT_SECRET=xxxxxxxx
export TOKEN_URL=https://application-auth.mixi.social/oauth2/token
export API_ADDRESS=application-api.mixi.social:443
```

---

## 3. 実行

```bash
go run delete.go
```

成功すると：

```text
削除成功
```

と表示されます。

---

# 注意

* 名言の順番はなるべく変更しない
* 新規名言は末尾追加推奨
* 長すぎる文章は投稿エラーになる可能性あり
* GitHub Actions の cron は UTC 基準です
