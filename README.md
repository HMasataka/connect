# connect

デザインモックアップHTML群にナビゲーション用スクリプトを注入し、クリッカブルプロトタイプとして動作させるCLIツール。

## インストール

```bash
go install github.com/HMasataka/connect@latest
```

## 使い方

```bash
# デザインディレクトリにスクリプトを注入
connect ./designs

# 注入済みスクリプトを除去
connect --clean ./designs

# 別ディレクトリにコピーしてから注入
connect --out-dir ./designs-linked ./designs

# 実行内容を確認（ファイルは変更しない）
connect --dry-run --verbose ./designs
```

## オプション

| オプション         | 説明               |
| ------------------ | ------------------ |
| `--out-dir <path>` | コピー先に注入     |
| `--clean`          | スクリプトを除去   |
| `--config <path>`  | 設定ファイルを指定 |
| `--dry-run`        | 実行内容を表示のみ |
| `--verbose`        | 詳細表示           |

## 設定ファイル

`connect.json` をデザインディレクトリ直下に配置すると自動で読み込まれます。

```json
{
  "links": [
    {
      "selector": ".nav-item:not(.active)",
      "mapping": { "Stash": "stash" }
    },
    {
      "selector": ".toolbar-btn",
      "mapping": { "Tags": "tags", "AI": "ai-assist" }
    },
    {
      "selector": ".titlebar-btn[title='Settings']",
      "match": "title",
      "target": "settings-appearance"
    }
  ],
  "modals": ["tags", "branches-dialog"],
  "close": [
    ".modal-close",
    { "selector": ".btn-secondary", "match": "text", "value": "Close" }
  ],
  "ignore": ["mockup"]
}
```

### links

クリック時の遷移先を定義するルール。

| フィールド | 説明                                  |
| ---------- | ------------------------------------- |
| `selector` | CSSセレクタ                           |
| `match`    | `text`（デフォルト）, `title`, `auto` |
| `mapping`  | テキスト → ページID                   |
| `target`   | 固定の遷移先                          |

マッピングがなければテキストを自動変換（`Cherry-pick` → `cherry-pick`）。

### modals / close

モーダルページと閉じるボタンの定義。閉じるボタンは最後に訪れた非モーダルページに遷移。

```json
{
  "modals": ["tags", "branches-dialog"],
  "close": [".modal-close"]
}
```

### ignore

ページ検出から除外するディレクトリ名。

## 仕組み

各 `index.html` の `</body>` 直前に `<script data-connect>` を挿入します。このスクリプトがクリックイベントを処理し、対応するページへ遷移させます。

再実行時は既存のスクリプトを自動で置き換えます。
