# article-scraper

ブログサイトをスクレイピングして Markdown 形式で出力するスクリプトです。
現時点で [Zenn](https://zenn.dev) に対応しています。

```bash
# Zenn
go run . zenn <id>
```

## 対応記法

- 見出し
- リスト（ul, ol）
- 画像
- 太字（strong）・強調（i）・打ち消し線（s）
- リンク
- 脚註
- アコーディオン（details）
