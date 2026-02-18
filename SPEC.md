# connect

デザインモックアップHTML群にナビゲーション用スクリプトを注入するCLIツール。

## コンセプト

1. **すべてはリンクルール**: セレクタとマッピングのペアで、クリック時の遷移先を定義
2. **自動変換**: マッピングがなければテキストからページIDを自動生成（`Cherry-pick` → `cherry-pick`）
3. **マッチ方式**: テキスト、title属性、または両方のフォールバック

## 設定ファイル

`connect.json`（デザインディレクトリ直下）

```json
{
  "links": [
    {
      "selector": ".nav-item:not(.active)",
      "match": "text",
      "mapping": {
        "Stash": "stash"
      }
    },
    {
      "selector": ".toolbar-btn",
      "match": "text",
      "mapping": {
        "Tags": "tags",
        "AI": "ai-assist"
      }
    },
    {
      "selector": ".titlebar-btn[title='Settings']",
      "match": "title",
      "mapping": {
        "Settings": "settings-appearance"
      }
    },
    {
      "selector": ".branch-selector",
      "match": "text",
      "target": "branches-dialog"
    }
  ],
  "modals": ["tags", "branches-dialog"],
  "close": [
    ".modal-close",
    {
      "selector": ".modal-footer .btn-secondary",
      "match": "text",
      "value": "Close"
    }
  ],
  "ignore": ["mockup"]
}
```

### links

リンクルールの配列。各ルールは以下のフィールドを持つ。

| フィールド | 必須 | 説明                                                                           |
| ---------- | ---- | ------------------------------------------------------------------------------ |
| `selector` | Yes  | CSSセレクタ                                                                    |
| `match`    | No   | マッチ方式: `text`（デフォルト）, `title`, `auto`（textを優先、なければtitle） |
| `mapping`  | No   | テキスト→ページIDのマッピング                                                  |
| `target`   | No   | 固定の遷移先（mappingより優先）                                                |

**遷移先の決定順序:**

1. `target` が指定されていればそれを使用
2. `mapping` にマッチすればその値を使用
3. テキストを小文字化・スペースをハイフン化してページIDに変換

### modals

モーダルページのリスト。クローズボタンは最後に訪れた非モーダルページに戻る。

```json
"modals": ["tags", "open-repository", "branches-dialog", "settings-appearance"]
```

**動作:**

1. ページ遷移時、モーダルでなければ `sessionStorage` に現在のページを保存
2. クローズボタンは `sessionStorage` のページに遷移

これにより、ページA → モーダルX → モーダルY → クローズ → ページA という遷移が可能。

### close

モーダルの閉じるボタン。クリックで最後の非モーダルページに遷移。

`modals` と組み合わせて使用する。`modals` が未定義の場合は何もしない。

```json
"close": [
  ".modal-close",
  ".titlebar-btn[title='Close']",
  { "selector": ".btn-secondary", "match": "text", "value": "Close" },
  { "selector": ".icon-btn", "match": "title", "value": "Close" }
]
```

- 文字列: セレクタのみ（マッチした全要素が対象）
- オブジェクト: セレクタ + マッチ条件でフィルタ

| フィールド | 必須 | 説明                                 |
| ---------- | ---- | ------------------------------------ |
| `selector` | Yes  | CSSセレクタ（属性セレクタも使用可）  |
| `match`    | No   | マッチ方式: `text`, `title`          |
| `value`    | No   | マッチさせる値（指定時のみフィルタ） |

**例:**

```json
{
  "modals": ["tags", "branches-dialog"],
  "close": [
    ".modal-close",
    {
      "selector": ".modal-footer .btn-secondary",
      "match": "text",
      "value": "Close"
    }
  ]
}
```

### ignore

ページ検出から除外するディレクトリ名。

## 生成されるスクリプト

```javascript
<script data-connect>
(function() {
  var pages = ["changes", "history", ...];
  var modals = ["tags", "branches-dialog", ...];
  var links = [
    { selector: ".nav-item:not(.active)", match: "text", mapping: {...}, target: null },
    ...
  ];
  var close = [
    { selector: ".modal-close", match: null, value: null },
    { selector: ".modal-footer .btn-secondary", match: "text", value: "Close" }
  ];

  // 現在のページを取得
  var currentPage = location.pathname.split('/').slice(-2)[0];

  // 非モーダルページなら記録
  if (modals.indexOf(currentPage) === -1) {
    sessionStorage.setItem('connect-last-page', currentPage);
  }

  function toPageId(text) {
    return text.trim().toLowerCase().replace(/\s+/g, '-');
  }

  function getText(el, match) {
    if (match === 'title') return el.getAttribute('title') || '';
    if (match === 'auto') return el.textContent.trim() || el.getAttribute('title') || '';
    return el.textContent.trim();
  }

  function resolve(text, mapping, target) {
    if (target) return target;
    if (mapping && mapping[text]) return mapping[text];
    return toPageId(text);
  }

  function navigate(pageId) {
    if (pages.indexOf(pageId) !== -1) {
      location.href = '../' + pageId + '/index.html';
    }
  }

  function closeModal() {
    var lastPage = sessionStorage.getItem('connect-last-page') || pages[0];
    location.href = '../' + lastPage + '/index.html';
  }

  // リンク
  links.forEach(function(rule) {
    document.querySelectorAll(rule.selector).forEach(function(el) {
      var text = getText(el, rule.match || 'text');
      var pageId = resolve(text, rule.mapping, rule.target);
      if (pages.indexOf(pageId) !== -1) {
        el.style.cursor = 'pointer';
        el.addEventListener('click', function(e) {
          e.preventDefault();
          navigate(pageId);
        });
      }
    });
  });

  // 閉じる
  close.forEach(function(rule) {
    var selector = typeof rule === 'string' ? rule : rule.selector;
    var match = typeof rule === 'string' ? null : rule.match;
    var value = typeof rule === 'string' ? null : rule.value;
    document.querySelectorAll(selector).forEach(function(el) {
      if (value) {
        var actual = match === 'title' ? el.getAttribute('title') : el.textContent.trim();
        if (actual !== value) return;
      }
      el.style.cursor = 'pointer';
      el.addEventListener('click', function(e) {
        e.preventDefault();
        closeModal();
      });
    });
  });
})();
</script>
```

## CLI

```
connect <designs-dir> [options]

Options:
  --out-dir <path>   コピーを作成してそちらに注入
  --clean            注入済みスクリプトを除去
  --config <path>    設定ファイルのパス
  --dry-run          変更内容を表示のみ
  --verbose          詳細表示
  --help             ヘルプ
```

## ファイル構成

```
connect/
├── SPEC.md
├── go.mod
├── main.go
└── connect/
    ├── config.go    # 設定読み込み
    ├── detect.go    # ページ検出
    ├── script.go    # スクリプト生成
    ├── inject.go    # 注入
    └── clean.go     # 除去
```
