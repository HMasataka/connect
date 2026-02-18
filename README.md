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

| オプション | 説明 |
|------------|------|
| `--in-place` | 元ファイルに直接注入（デフォルト） |
| `--out-dir <path>` | コピー先に注入 |
| `--clean` | スクリプトを除去 |
| `--config <path>` | 設定ファイルを指定 |
| `--dry-run` | 実行内容を表示のみ |
| `--verbose` | 詳細表示 |

## 設定ファイル

`connect.json` をデザインディレクトリ直下に配置すると自動で読み込まれます。

```json
{
  "selectors": {
    "nav": ".nav-item",
    "toolbar": ".toolbar-btn",
    "activeClass": "active",
    "modalClose": ".modal-close",
    "modalCloseText": ".modal-footer .btn-secondary"
  },
  "mapping": {
    "Stash": "stash"
  },
  "toolbar": {
    "Tags": "tags"
  },
  "ignore": ["mockup"]
}
```

| フィールド | 説明 |
|------------|------|
| `selectors` | ナビゲーション要素のCSSセレクタ |
| `mapping` | サイドバーのテキスト → ディレクトリ名 |
| `toolbar` | ツールバーボタンのテキスト → ディレクトリ名 |
| `ignore` | 除外するディレクトリ名 |

## 仕組み

各 `index.html` の `</body>` 直前に `<script data-connect>` を挿入します。このスクリプトがサイドバーやツールバーのクリックイベントを処理し、対応するページへ遷移させます。

再実行時は既存のスクリプトを自動で置き換えます。
