# connect

フォルダ分割されたデザインモックアップHTML群に、ナビゲーション用の `<script>` を注入してページ間遷移を可能にするCLIツール。

## 背景

`modular-design` スキルで作成されるデザインモックアップは、ページ/機能ごとに独立したディレクトリに分かれている。

```
designs/
├── changes/
│   ├── index.html
│   └── styles.css
├── history/
│   ├── index.html
│   └── styles.css
├── cherry-pick/
│   ├── index.html
│   └── styles.css
├── tags/
│   ├── index.html
│   └── styles.css
└── ...
```

各HTMLは静的デザイン確認用であり、サイドバーやツールバーのボタンをクリックしても何も起きない。このツールは各HTMLに小さなJavaScriptを注入し、クリッカブルプロトタイプとして動作させる。

## 動作概要

```
connect ./designs
```

`designs/` 配下の各ディレクトリを走査し、各 `index.html` の `</body>` 直前にナビゲーション用の `<script>` タグを1つ挿入する。

## アプローチ: Script注入

HTMLのDOM構造は一切変更しない。各 `index.html` の `</body>` 直前に以下のような `<script>` を挿入する。

```javascript
<script data-connect>
(function() {
  var pages = ["changes","history","cherry-pick","tags","revert","reset","reflog","submodules","worktrees","search"];
  var sel = {"nav":".nav-item","toolbar":".toolbar-btn","activeClass":"active","modalClose":".modal-close","modalCloseText":".modal-footer .btn-secondary"};
  var customMapping = {};
  var toolbarMapping = {};

  function toPageId(text) {
    return text.trim().toLowerCase().replace(/\s+/g, '-');
  }

  function resolve(text, map) {
    if (map[text]) return map[text];
    return toPageId(text);
  }

  function navigate(pageId) {
    if (pages.indexOf(pageId) !== -1) {
      location.href = '../' + pageId + '/index.html';
    }
  }

  // サイドバー
  document.querySelectorAll(sel.nav).forEach(function(item) {
    if (item.classList.contains(sel.activeClass)) return;
    var text = item.textContent.replace(/\d+/g, '');
    var pageId = resolve(text, customMapping);
    if (pages.indexOf(pageId) !== -1) {
      item.style.cursor = 'pointer';
      item.addEventListener('click', function() { navigate(pageId); });
    }
  });

  // ツールバー
  document.querySelectorAll(sel.toolbar).forEach(function(btn) {
    var text = btn.textContent.trim();
    var pageId = resolve(text, toolbarMapping);
    if (pages.indexOf(pageId) !== -1) {
      btn.addEventListener('click', function() { navigate(pageId); });
    }
  });

  // モーダル閉じる
  if (sel.modalClose) {
    document.querySelectorAll(sel.modalClose).forEach(function(btn) {
      btn.style.cursor = 'pointer';
      btn.addEventListener('click', function() { history.back(); });
    });
  }
  if (sel.modalCloseText) {
    document.querySelectorAll(sel.modalCloseText).forEach(function(btn) {
      if (btn.textContent.trim() === 'Close') {
        btn.style.cursor = 'pointer';
        btn.addEventListener('click', function() { history.back(); });
      }
    });
  }
})();
</script>
```

`connect` は設定に応じて `sel`, `customMapping`, `toolbarMapping` の値を差し替えてスクリプトを生成する。`selectors` が未設定の場合は上記のデフォルト値が埋め込まれる。

### `pages` 配列の生成

`connect` はディレクトリ走査の結果から `pages` 配列を自動生成し、スクリプトに埋め込む。設定ファイルで上書きされたマッピングも反映される。

### `data-connect` 属性

注入済みの `<script>` には `data-connect` 属性を付与する。これにより:

- 再実行時に二重注入を防止できる（既存の `data-connect` スクリプトを検出したら削除してから再注入）
- `--clean` オプションで注入スクリプトだけを除去できる

## ナビゲーション対象

### サイドバー `.nav-item`

1. `.nav-item` のテキストコンテンツ（SVG・バッジの数字除く）を取得
2. 小文字化・スペースをハイフン化してページIDに変換（例: "Cherry-pick" → `cherry-pick`）
3. `pages` 配列に存在すれば `click` イベントで `../cherry-pick/index.html` に遷移
4. `.active` クラスが付いた項目（現在のページ）はスキップ

### ツールバー `.toolbar-btn`

1. `.toolbar-btn` のテキストコンテンツからページIDを特定
2. `pages` 配列に存在すれば `click` イベントで遷移

### モーダルの閉じるボタン

Tags等のモーダルページでは:

1. `.modal-close` ボタンのクリックで `history.back()`
2. `.modal-footer` 内の "Close" ボタンのクリックで `history.back()`

## 入力

- 引数: デザインディレクトリのパス（例: `./designs`）
- 前提: 各サブディレクトリに `index.html` が存在する

## ページの検出

1. 指定ディレクトリの直下サブディレクトリを列挙
2. 各サブディレクトリに `index.html` が存在するものをページとして認識
3. ディレクトリ名がページIDとなる（例: `changes`, `history`, `cherry-pick`）

## 出力モード

### `--in-place`（デフォルト）

元の `index.html` に `<script>` を注入する。

### `--out-dir <path>`

指定ディレクトリにコピーを作成し、コピー側に `<script>` を注入する。元ファイルは変更しない。ディレクトリ構造を再現し、`styles.css` 等の関連ファイルもコピーする。

### `--clean`

`data-connect` 属性を持つ `<script>` タグを全HTMLから除去する。注入前の状態に戻す。

## マッピングの上書き

テキストからページIDへの自動変換で対応できないケースのために、設定ファイルで明示的なマッピングを指定できる。

### `connect.json`（デザインディレクトリ直下に配置）

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
    "Stash": "stash",
    "Search": "search"
  },
  "toolbar": {
    "Tags": "tags",
    "Settings": "settings"
  },
  "ignore": ["mockup"]
}
```

- `selectors`: ナビゲーション要素のCSSセレクタ（後述）
- `mapping`: サイドバーのテキスト → ディレクトリ名
- `toolbar`: ツールバーボタンのテキスト → ディレクトリ名
- `ignore`: リンク対象から除外するディレクトリ名（例: `mockup` は元の単体モックアップなので除外）

すべてのフィールドはオプション。設定ファイル自体が存在しない場合はデフォルト値で動作する。

### セレクタのカスタマイズ

`selectors` でナビゲーション要素のCSSセレクタを指定できる。これにより、異なるHTML構造を持つデザインでもconnectを利用できる。

| キー             | デフォルト値                   | 説明                           |
| ---------------- | ------------------------------ | ------------------------------ |
| `nav`            | `.nav-item`                    | サイドバーのナビゲーション項目 |
| `toolbar`        | `.toolbar-btn`                 | ツールバーのボタン             |
| `activeClass`    | `active`                       | 現在のページを示すクラス名     |
| `modalClose`     | `.modal-close`                 | モーダルの×ボタン              |
| `modalCloseText` | `.modal-footer .btn-secondary` | モーダルのCloseテキストボタン  |

例えば `<nav>` 内の `<li>` でナビゲーションするデザインの場合:

```json
{
  "selectors": {
    "nav": "nav li",
    "toolbar": ".action-btn",
    "activeClass": "current"
  }
}
```

指定されなかったキーはデフォルト値が使われる。

### マッピングの埋め込み

マッピングは `pages` 配列の生成には影響しない。スクリプト内の `toPageId` 変換を補助するために、注入スクリプト内に追加のマッピングテーブルとして埋め込まれる:

```javascript
var customMapping = { Stash: "stash", Search: "search" };
var toolbarMapping = { Tags: "tags", Settings: "settings" };
```

設定ファイルが存在しない場合は、テキストの自動変換のみで動作する。

## CLIインターフェース

```
connect <designs-dir> [options]

Arguments:
  designs-dir          デザインディレクトリのパス

Options:
  --in-place           元ファイルに直接注入する（デフォルト）
  --out-dir <path>     コピーを作成してそちらに注入する
  --clean              注入済みスクリプトを除去する
  --config <path>      設定ファイルのパス（デフォルト: <designs-dir>/connect.json）
  --dry-run            変更内容を表示するが、ファイルは書き換えない
  --verbose            処理の詳細を表示する
  --help               ヘルプを表示する
```

## 実行例

```bash
# デフォルト（in-place注入）
connect ./designs

# 別ディレクトリにコピー＆注入
connect ./designs --out-dir ./designs-linked

# 注入済みスクリプトを除去
connect ./designs --clean

# dry-runで確認
connect ./designs --dry-run --verbose
```

### dry-run出力例

```
Detected pages: changes, cherry-pick, history, reflog, reset, revert, search, submodules, tags, worktrees
Skipping: mockup (ignored)

designs/changes/index.html:
  inject <script data-connect> (pages: 10 entries)

designs/history/index.html:
  inject <script data-connect> (pages: 10 entries)

designs/tags/index.html:
  inject <script data-connect> (pages: 10 entries)

...

10 files processed.
```

## 処理フロー

1. デザインディレクトリを走査し、`index.html` を持つサブディレクトリを列挙
2. `connect.json` があれば読み込み、`ignore` に含まれるディレクトリを除外
3. ページID一覧から `pages` 配列を生成
4. カスタムマッピングがあればスクリプトに埋め込む
5. 各 `index.html` について:
   a. 既存の `data-connect` スクリプトがあれば除去
   b. `</body>` の直前にスクリプトを挿入
   c. ファイルを書き出す（`--in-place`）またはコピー先に書き出す（`--out-dir`）

## 技術選定

Go で実装する。HTMLの `</body>` 位置の特定は文字列検索で行う（`golang.org/x/net/html` でのフルパースは不要）。

## ファイル構成

```
connect/
├── SPEC.md
├── go.mod
├── go.sum
├── main.go
├── connect/
│   ├── detect.go       # ページ検出
│   ├── inject.go       # スクリプト生成・注入
│   ├── clean.go        # スクリプト除去
│   ├── config.go       # 設定ファイル読み込み
│   └── script.go       # スクリプトテンプレート
└── testdata/
    └── ...
```
