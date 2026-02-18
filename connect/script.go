package connect

import (
	"encoding/json"
	"fmt"
)

const scriptTemplate = `<script data-connect>
(function() {
  var pages = %s;
  var modals = %s;
  var links = %s;
  var close = %s;

  // 現在のページを取得
  var currentPage = location.pathname.split('/').slice(-2)[0];

  // 非モーダルページなら記録
  if (modals.indexOf(currentPage) === -1) {
    sessionStorage.setItem('connect-last-page', currentPage);
  }

  var TEXT_NODE = 3;

  function toPageId(text) {
    return text.trim().toLowerCase().replace(/\s+/g, '-');
  }

  function getDirectText(el) {
    return Array.from(el.childNodes)
      .filter(function(n) { return n.nodeType === TEXT_NODE; })
      .map(function(n) { return n.textContent; })
      .join('').trim();
  }

  function getText(el, match) {
    if (match === 'title') return el.getAttribute('title') || '';
    if (match === 'auto') return getDirectText(el) || el.getAttribute('title') || '';
    return getDirectText(el);
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
    document.querySelectorAll(rule.selector).forEach(function(el) {
      if (rule.value) {
        var actual = rule.match === 'title' ? el.getAttribute('title') : el.textContent.trim();
        if (actual !== rule.value) return;
      }
      el.style.cursor = 'pointer';
      el.addEventListener('click', function(e) {
        e.preventDefault();
        closeModal();
      });
    });
  });
})();
</script>`

type jsLinkRule struct {
	Selector string            `json:"selector"`
	Match    string            `json:"match"`
	Mapping  map[string]string `json:"mapping"`
	Target   *string           `json:"target"`
}

type jsCloseRule struct {
	Selector string  `json:"selector"`
	Match    *string `json:"match"`
	Value    *string `json:"value"`
}

func GenerateScript(pages []string, cfg Config) string {
	pagesJSON := toJSON(pages)
	modalsJSON := toJSON(cfg.Modals)
	linksJSON := linksToJSON(cfg.Links)
	closeJSON := closeToJSON(cfg.Close)

	return fmt.Sprintf(scriptTemplate, pagesJSON, modalsJSON, linksJSON, closeJSON)
}

func toJSON(v any) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func linksToJSON(links []LinkRule) string {
	jsRules := make([]jsLinkRule, len(links))
	for i, link := range links {
		jsRules[i] = jsLinkRule{
			Selector: link.Selector,
			Match:    link.Match,
			Mapping:  link.Mapping,
		}
		if link.Target != "" {
			target := link.Target
			jsRules[i].Target = &target
		}
	}
	return toJSON(jsRules)
}

func closeToJSON(closes []CloseRule) string {
	jsRules := make([]jsCloseRule, len(closes))
	for i, c := range closes {
		jsRules[i] = jsCloseRule{
			Selector: c.Selector,
		}
		if c.Match != "" {
			match := c.Match
			jsRules[i].Match = &match
		}
		if c.Value != "" {
			value := c.Value
			jsRules[i].Value = &value
		}
	}
	return toJSON(jsRules)
}
