package connect

import (
	"encoding/json"
	"fmt"
	"strings"
)

const scriptTemplate = `<script data-connect>
(function() {
  var pages = %s;
  var sel = %s;
  var customMapping = %s;
  var toolbarMapping = %s;

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
</script>`

func GenerateScript(pages []string, selectors Selectors, mapping, toolbar map[string]string) string {
	pagesJSON := toJSONArray(pages)
	selectorsJSON := toSelectorsJSON(selectors)
	mappingJSON := toJSONObject(mapping)
	toolbarJSON := toJSONObject(toolbar)

	return fmt.Sprintf(scriptTemplate, pagesJSON, selectorsJSON, mappingJSON, toolbarJSON)
}

func toJSONArray(items []string) string {
	data, _ := json.Marshal(items)
	return string(data)
}

func toJSONObject(m map[string]string) string {
	if len(m) == 0 {
		return "{}"
	}
	data, _ := json.Marshal(m)
	return string(data)
}

func toSelectorsJSON(s Selectors) string {
	parts := []string{
		fmt.Sprintf(`"nav":"%s"`, s.Nav),
		fmt.Sprintf(`"toolbar":"%s"`, s.Toolbar),
		fmt.Sprintf(`"activeClass":"%s"`, s.ActiveClass),
		fmt.Sprintf(`"modalClose":"%s"`, s.ModalClose),
		fmt.Sprintf(`"modalCloseText":"%s"`, s.ModalCloseText),
	}
	return "{" + strings.Join(parts, ",") + "}"
}
