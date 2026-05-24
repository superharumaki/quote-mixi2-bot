# quote-mixi2-bot

mixi2 に自動で名言を投稿する Bot です。

## 機能

* `quotes.json` から名言を読み込み
* 未投稿の名言をランダム選択
* 投稿済みを `state.json` に記録
* 全件投稿後に自動リセット
* GitHub Actions で毎日自動投稿

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
git add quotes.json
git commit -m "update quotes"
git push
```

---

# 注意

* 名言の順番はなるべく変更しない
* 新規名言は末尾に追加推奨
* 長すぎる文章は投稿エラーになる可能性あり
