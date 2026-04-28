package web

const indexHTML = `<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
  <style>
    :root {
      --bg: #f5f7fb;
      --panel: #ffffff;
      --line: #d9e0ea;
      --text: #172033;
      --muted: #667085;
      --accent: #2563eb;
      --accent-hover: #1d4ed8;
      --ok: #16803c;
      --warn: #9a5b00;
      --bad: #b93815;
      --term: #101828;
      --term-line: #283548;
      --term-text: #e4e7ec;
    }
    * { box-sizing: border-box; }
    html, body { min-height: 100%; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font: 14px/1.45 system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    .wrap {
      width: min(1180px, calc(100vw - 32px));
      margin: 0 auto;
    }
    header {
      position: sticky;
      top: 0;
      z-index: 10;
      border-bottom: 1px solid var(--line);
      background: rgba(255, 255, 255, .96);
    }
    .top {
      min-height: 68px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 12px;
      min-width: 0;
      flex-wrap: wrap;
    }
    h1 {
      margin: 0;
      font-size: 21px;
      line-height: 1.15;
      letter-spacing: 0;
    }
    .subtitle {
      margin-top: 2px;
      color: var(--muted);
      font-size: 12px;
    }
    .actions, .row-actions, .page-actions {
      display: flex;
      align-items: center;
      gap: 8px;
      flex-wrap: wrap;
    }
    button {
      min-height: 34px;
      border: 1px solid var(--line);
      border-radius: 6px;
      padding: 6px 12px;
      background: #fff;
      color: var(--text);
      cursor: pointer;
      font: inherit;
    }
    button:hover { background: #f8fafc; border-color: #b9c4d3; }
    button.primary { border-color: var(--accent); background: var(--accent); color: #fff; }
    button.primary:hover { border-color: var(--accent-hover); background: var(--accent-hover); }
    button.danger { color: var(--bad); }
    button.small {
      min-height: 28px;
      padding: 3px 8px;
      font-size: 12px;
    }
    button:disabled { opacity: .62; cursor: wait; }
    main { padding: 16px 0 28px; }
    .status {
      display: grid;
      grid-template-columns: repeat(3, minmax(0, 1fr));
      gap: 12px;
      margin-bottom: 14px;
    }
    .metric, .card, .log-box {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
    }
    .metric {
      min-width: 0;
      min-height: 78px;
      padding: 12px;
    }
    .label {
      margin-bottom: 5px;
      color: var(--muted);
      font-size: 12px;
    }
    .value {
      min-width: 0;
      font-size: 17px;
      font-weight: 650;
      word-break: break-word;
    }
    .subvalue {
      margin-top: 3px;
      color: var(--muted);
      font-size: 12px;
      word-break: break-word;
    }
    .ok { color: var(--ok); }
    .warn { color: var(--warn); }
    .bad { color: var(--bad); }
    .cards {
      display: grid;
      grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
      gap: 14px;
      align-items: start;
      min-width: 0;
    }
    .card {
      min-width: 0;
      height: 560px;
      display: grid;
      grid-template-rows: auto 1fr;
      overflow: hidden;
    }
    .card-head, .block-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
      min-height: 44px;
      padding: 12px 14px;
      border-bottom: 1px solid var(--line);
    }
    h2, h3 {
      margin: 0;
      letter-spacing: 0;
      line-height: 1.2;
    }
    h2 { font-size: 16px; }
    h3 { font-size: 13px; }
    .pill {
      display: inline-flex;
      align-items: center;
      min-height: 23px;
      border: 1px solid var(--line);
      border-radius: 999px;
      padding: 2px 9px;
      background: #fff;
      color: var(--muted);
      font-size: 12px;
      white-space: nowrap;
    }
    .pill.ok { border-color: #b8dec2; background: #ecfdf3; color: var(--ok); }
    .pill.warn { border-color: #ead08a; background: #fff8e6; color: var(--warn); }
    .pill.bad { border-color: #f0b8aa; background: #fff3ef; color: var(--bad); }
    .card-body {
      min-height: 0;
      padding: 14px;
      overflow: hidden;
      display: flex;
      flex-direction: column;
      gap: 12px;
    }
    .form-grid {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 11px;
      min-width: 0;
    }
    .field {
      display: grid;
      gap: 5px;
      min-width: 0;
    }
    .field.full { grid-column: 1 / -1; }
    .field label {
      color: var(--muted);
      font-size: 12px;
    }
    input, select {
      width: 100%;
      min-width: 0;
      height: 36px;
      border: 1px solid var(--line);
      border-radius: 6px;
      padding: 7px 9px;
      background: #fff;
      color: var(--text);
      font: inherit;
    }
    select {
      appearance: none;
      background-image: linear-gradient(45deg, transparent 50%, #667085 50%), linear-gradient(135deg, #667085 50%, transparent 50%);
      background-position: calc(100% - 16px) 50%, calc(100% - 11px) 50%;
      background-size: 5px 5px, 5px 5px;
      background-repeat: no-repeat;
      padding-right: 34px;
    }
    input:focus, select:focus {
      border-color: #8bb7ff;
      box-shadow: 0 0 0 3px rgba(37, 99, 235, .12);
      outline: 0;
    }
    .secret-hidden { -webkit-text-security: disc; }
    .input-row {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 6px;
      min-width: 0;
    }
    .toggle {
      display: flex;
      align-items: center;
      gap: 8px;
      min-height: 34px;
      color: var(--text);
    }
    .toggle input {
      width: 16px;
      height: 16px;
    }
    .message {
      min-height: 20px;
      color: var(--muted);
      font-size: 13px;
    }
    .message.ok { color: var(--ok); }
    .message.bad { color: var(--bad); }
    .issues {
      display: none;
      border: 1px solid #f0b8aa;
      border-radius: 6px;
      padding: 8px 10px;
      background: #fff3ef;
      color: var(--bad);
      font-size: 13px;
    }
    .issues.show { display: block; }
    .issues ul {
      margin: 5px 0 0;
      padding-left: 18px;
    }
    .source-field.hidden { display: none; }
    .list-wrap {
      display: grid;
      gap: 7px;
      min-width: 0;
    }
    .list-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 10px;
    }
    .list-head label {
      color: var(--muted);
      font-size: 12px;
    }
    .listbox {
      height: 104px;
      min-height: 104px;
      max-height: 104px;
      overflow: auto;
      border: 1px solid var(--line);
      border-radius: 6px;
      background: #fff;
      padding: 6px;
    }
    .list-row {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto;
      gap: 6px;
      margin-bottom: 6px;
    }
    .list-row:last-child { margin-bottom: 0; }
    .list-row input { height: 30px; font-size: 13px; }
    .empty-list {
      padding: 8px;
      color: #98a2b3;
      font-size: 13px;
    }
    .hint {
      color: var(--muted);
      font-size: 12px;
    }
    .results-block {
      min-height: 0;
      border: 1px solid var(--line);
      border-radius: 8px;
      overflow: hidden;
      display: grid;
      grid-template-rows: auto 1fr;
    }
    .results-wrap {
      min-height: 0;
      overflow: auto;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      table-layout: fixed;
    }
    th, td {
      border-bottom: 1px solid var(--line);
      padding: 8px 9px;
      text-align: left;
      vertical-align: top;
      font-size: 12px;
      overflow-wrap: anywhere;
    }
    th {
      background: #f8fafc;
      color: var(--muted);
      font-weight: 650;
    }
    .status-tag {
      display: inline-flex;
      align-items: center;
      min-height: 22px;
      border-radius: 999px;
      padding: 1px 8px;
      background: #eef4ff;
      color: #34536f;
      font-size: 12px;
    }
    .status-tag.updated { background: #ecfdf3; color: var(--ok); }
    .status-tag.failed { background: #fff3ef; color: var(--bad); }
    .log-button-mark {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 22px;
      height: 22px;
      border: 1px solid var(--line);
      border-radius: 999px;
      color: var(--muted);
      font-size: 13px;
    }
    .log-page {
      position: fixed;
      inset: 0;
      z-index: 50;
      display: grid;
      grid-template-rows: auto 1fr;
      background: var(--bg);
    }
    .log-page.hidden { display: none; }
    .log-page-head {
      min-height: 68px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      border-bottom: 1px solid var(--line);
      background: #fff;
      padding: 0 max(16px, calc((100vw - 1180px) / 2));
    }
    .log-page-title strong { font-size: 19px; }
    .log-box {
      width: min(1180px, calc(100vw - 32px));
      margin: 16px auto 28px;
      min-height: 0;
      overflow: hidden;
    }
    .term {
      height: 100%;
      min-height: 320px;
      overflow: auto;
      background: var(--term);
    }
    .log-row {
      display: grid;
      grid-template-columns: 156px 64px 1fr;
      gap: 10px;
      padding: 7px 12px;
      border-bottom: 1px solid var(--term-line);
      color: var(--term-text);
      font: 12px/1.45 ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    .log-time { color: #8d9aaa; }
    .log-level { color: #80bfff; }
    .log-level.ERROR { color: #ff9b93; }
    .log-level.WARN { color: #f2cc60; }
    .log-msg { word-break: break-word; }
    .empty {
      padding: 10px;
      color: #98a2b3;
      font-size: 13px;
    }
    @media (max-width: 820px) {
      header { position: static; }
      .top { align-items: flex-start; flex-direction: column; padding: 12px 0; }
      .actions { justify-content: flex-start; }
      .status, .cards, .form-grid { grid-template-columns: 1fr; }
      .card { height: auto; min-height: 0; }
      .log-row { grid-template-columns: 1fr; gap: 2px; }
    }
  </style>
</head>
<body>
  <header>
    <div class="wrap top">
      <div class="brand">
        <div>
          <h1>dnshe-go</h1>
          <div class="subtitle">DNSHE IPv6 DDNS / AAAA</div>
        </div>
        <button id="open-logs" type="button">日志 <span class="log-button-mark">›</span></button>
      </div>
      <div class="actions">
        <button id="refresh" type="button">刷新</button>
        <button id="run" class="primary" type="button">立即同步</button>
        <button id="save" class="primary" type="button">保存配置</button>
      </div>
    </div>
  </header>

  <main class="wrap">
    <div class="status">
      <div class="metric">
        <div class="label">状态</div>
        <div class="value" id="state">--</div>
        <div class="subvalue" id="state-detail">--</div>
      </div>
      <div class="metric">
        <div class="label">上次同步</div>
        <div class="value" id="last-run">--</div>
        <div class="subvalue" id="duration">--</div>
      </div>
      <div class="metric">
        <div class="label">IPv6</div>
        <div class="value" id="ipv6">--</div>
        <div class="subvalue" id="ipv6-domains-count">--</div>
      </div>
    </div>

    <div class="cards">
      <section class="card">
        <div class="card-head">
          <h2>DNSHE</h2>
          <span class="pill" id="credential-state">未配置</span>
        </div>
        <div class="card-body">
          <div class="form-grid">
            <div class="field">
              <label for="api-key">API Key</label>
              <input id="api-key" autocomplete="off" placeholder="留空沿用已保存的 Key">
            </div>
            <div class="field">
              <label for="api-secret">API Secret</label>
              <div class="input-row">
                <input id="api-secret" class="secret-hidden" type="text" autocomplete="off" spellcheck="false" placeholder="留空沿用已保存的 Secret">
                <button id="toggle-secret" class="small" type="button">显示</button>
              </div>
            </div>
            <div class="field full">
              <label for="api-base-url">API Base URL</label>
              <input id="api-base-url">
            </div>
            <div class="field">
              <label for="interval">同步间隔（秒）</label>
              <input id="interval" type="number" min="10" step="1">
            </div>
            <div class="field">
              <label for="ttl">TTL</label>
              <input id="ttl" type="number" min="1" step="1">
            </div>
          </div>
          <div class="message" id="message"></div>
          <div class="issues" id="issues"></div>
          <div class="results-block">
            <div class="block-head">
              <h3>最近同步结果</h3>
              <span class="pill" id="result-count">0 条</span>
            </div>
            <div class="results-wrap" id="results"></div>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="card-head">
          <h2>IPv6</h2>
          <span class="pill" id="ipv6-state">未启用</span>
        </div>
        <div class="card-body">
          <label class="toggle"><input id="ipv6-enable" type="checkbox">启用 AAAA 记录</label>
          <div class="field">
            <label for="ipv6-source">IP 来源</label>
            <select id="ipv6-source">
              <option value="url">URL</option>
              <option value="interface">网卡</option>
              <option value="cmd">命令</option>
            </select>
          </div>
          <div class="field source-field" data-source-field="interface">
            <label for="ipv6-interface">网卡</label>
            <select id="ipv6-interface"></select>
          </div>
          <div class="field source-field" data-source-field="cmd">
            <label for="ipv6-command">命令</label>
            <input id="ipv6-command" placeholder="命令输出中需要包含 IPv6">
          </div>
          <div class="list-wrap source-field" data-source-field="url">
            <div class="list-head">
              <label>URL</label>
              <button class="small" id="add-url" type="button">添加</button>
            </div>
            <div class="listbox" id="ipv6-urls"></div>
          </div>
          <div class="list-wrap">
            <div class="list-head">
              <label>域名</label>
              <button class="small" id="add-domain" type="button">添加</button>
            </div>
            <div class="listbox" id="ipv6-domains"></div>
            <div class="hint">DNSHE 中需要已有对应子域名和 AAAA 记录。</div>
          </div>
        </div>
      </section>
    </div>
  </main>

  <div class="log-page hidden" id="logs-page" aria-hidden="true">
    <div class="log-page-head">
      <div>
        <strong>日志</strong>
        <div class="subtitle">运行与同步记录</div>
      </div>
      <div class="page-actions">
        <button id="clear-logs" type="button">清空</button>
        <button id="close-logs" class="primary" type="button">返回</button>
      </div>
    </div>
    <div class="log-box">
      <div class="term" id="logs"></div>
    </div>
  </div>

  <script>
    const $ = (id) => document.getElementById(id);
    const apiKey = $("api-key");
    const apiSecret = $("api-secret");
    const message = $("message");
    const issues = $("issues");
    let formDirty = false;
    let loadedConfig = false;
    let interfacesLoaded = false;

    function esc(value) {
      return String(value || "")
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;")
        .replaceAll("'", "&#39;");
    }

    function cleanLines(values) {
      const seen = new Set();
      const out = [];
      for (const value of values) {
        const s = String(value || "").trim();
        if (!s || seen.has(s)) continue;
        seen.add(s);
        out.push(s);
      }
      return out;
    }

    function listValues(id) {
      return cleanLines(Array.from($(id).querySelectorAll("input")).map((input) => input.value));
    }

    function markDirty() {
      formDirty = true;
      updateSourceFields();
      updateBadges();
      showProblems([]);
    }

    function addListRow(id, value, placeholder) {
      const box = $(id);
      const empty = box.querySelector(".empty-list");
      if (empty) empty.remove();
      const row = document.createElement("div");
      row.className = "list-row";
      const input = document.createElement("input");
      input.value = value || "";
      input.placeholder = placeholder || "";
      input.autocomplete = "off";
      const remove = document.createElement("button");
      remove.type = "button";
      remove.className = "small danger";
      remove.textContent = "删除";
      row.append(input, remove);
      box.append(row);
      input.addEventListener("input", markDirty);
      input.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
          event.preventDefault();
          const next = addListRow(id, "", placeholder);
          next.focus();
        }
      });
      input.addEventListener("paste", (event) => {
        const text = event.clipboardData.getData("text");
        if (!text.includes("\n")) return;
        event.preventDefault();
        const parts = cleanLines(text.split(/\r?\n/));
        if (parts.length === 0) return;
        input.value = parts.shift();
        for (const part of parts) addListRow(id, part, placeholder);
        markDirty();
      });
      remove.addEventListener("click", () => {
        row.remove();
        if (box.querySelectorAll(".list-row").length === 0) setList(id, [], placeholder);
        markDirty();
      });
      return input;
    }

    function setList(id, values, placeholder) {
      const box = $(id);
      box.innerHTML = "";
      const cleaned = cleanLines(values || []);
      if (cleaned.length === 0) {
        const empty = document.createElement("div");
        empty.className = "empty-list";
        empty.textContent = placeholder || "暂无条目";
        box.append(empty);
        return;
      }
      for (const value of cleaned) addListRow(id, value, placeholder);
    }

    function ipConfig() {
      return {
        enable: $("ipv6-enable").checked,
        source: $("ipv6-source").value,
        urls: listValues("ipv6-urls"),
        interface: $("ipv6-interface").value,
        command: $("ipv6-command").value.trim(),
        domains: listValues("ipv6-domains"),
      };
    }

    function readConfig() {
      return {
        interval_seconds: Number($("interval").value || 300),
        ttl: Number($("ttl").value || 600),
        dnshe: {
          api_key: apiKey.value.trim(),
          api_secret: apiSecret.value.trim(),
          api_base_url: $("api-base-url").value.trim(),
        },
        ipv6: ipConfig(),
      };
    }

    function fillConfig(cfg) {
      $("interval").value = cfg.interval_seconds || 300;
      $("ttl").value = cfg.ttl || 600;
      $("api-base-url").value = cfg.dnshe.api_base_url || "";
      apiKey.placeholder = cfg.dnshe.api_key_masked || "留空沿用已保存的 Key";
      apiSecret.placeholder = cfg.dnshe.api_secret_masked || "留空沿用已保存的 Secret";
      $("credential-state").className = cfg.api_key_configured && cfg.api_secret_configured ? "pill ok" : "pill bad";
      $("credential-state").textContent = cfg.api_key_configured && cfg.api_secret_configured ? "已配置" : "未配置";
      const item = cfg.ipv6 || {};
      $("ipv6-enable").checked = item.enable !== false;
      $("ipv6-source").value = item.source || "url";
      loadInterfaces(item.interface || "");
      $("ipv6-command").value = item.command || "";
      setList("ipv6-urls", item.urls || [], "每行一个接口");
      setList("ipv6-domains", item.domains || [], "每行一个域名");
      updateSourceFields();
      updateBadges();
    }

    async function loadInterfaces(selected) {
      if (interfacesLoaded && selected == null) return;
      const select = $("ipv6-interface");
      const current = selected != null ? selected : select.value;
      try {
        const res = await fetch("/api/interfaces", { cache: "no-store" });
        const data = await res.json();
        select.innerHTML = "";
        const auto = document.createElement("option");
        auto.value = "";
        auto.textContent = "自动选择有效网卡";
        select.append(auto);
        const interfaces = Array.isArray(data.interfaces) ? data.interfaces : [];
        for (const item of interfaces) {
          const option = document.createElement("option");
          option.value = item.name || "";
          option.textContent = item.label || item.name || "";
          select.append(option);
        }
        if (current && !Array.from(select.options).some((option) => option.value === current)) {
          const option = document.createElement("option");
          option.value = current;
          option.textContent = current + "（当前配置）";
          select.append(option);
        }
        select.value = current || "";
        interfacesLoaded = true;
      } catch {
        if (select.options.length === 0) {
          const option = document.createElement("option");
          option.value = "";
          option.textContent = "无法读取网卡，自动选择";
          select.append(option);
        }
      }
    }

    function updateSourceFields() {
      const source = $("ipv6-source").value || "url";
      document.querySelectorAll("[data-source-field]").forEach((field) => {
        field.classList.toggle("hidden", field.dataset.sourceField !== source);
      });
    }

    function updateBadges() {
      const enabled = $("ipv6-enable").checked;
      const count = listValues("ipv6-domains").length;
      const badge = $("ipv6-state");
      if (!enabled) {
        badge.className = "pill warn";
        badge.textContent = "未启用";
      } else if (count === 0) {
        badge.className = "pill warn";
        badge.textContent = "缺少域名";
      } else {
        badge.className = "pill ok";
        badge.textContent = count + " 个域名";
      }
      const source = $("ipv6-source").value || "url";
      $("ipv6-domains-count").textContent = enabled ? source.toUpperCase() + " / " + count + " 个域名" : "未启用";
    }

    function validateBeforeSave(cfg) {
      const problems = [];
      const keyKnown = apiKey.value.trim() || apiKey.placeholder.trim();
      const secretKnown = apiSecret.value.trim() || apiSecret.placeholder.trim();
      if (!keyKnown) problems.push("DNSHE API Key 不能为空");
      if (!secretKnown) problems.push("DNSHE API Secret 不能为空");
      if (!cfg.ipv6.enable) problems.push("IPv6 同步必须启用");
      if (cfg.ipv6.enable && cfg.ipv6.domains.length === 0) problems.push("IPv6 域名不能为空");
      if (cfg.ipv6.enable && cfg.ipv6.source === "url" && cfg.ipv6.urls.length === 0) problems.push("IPv6 使用 URL 来源时 URL 不能为空");
      if (cfg.ipv6.enable && cfg.ipv6.source === "cmd" && !cfg.ipv6.command) problems.push("IPv6 使用命令来源时命令不能为空");
      return problems;
    }

    function showProblems(problems) {
      if (!Array.isArray(problems) || problems.length === 0) {
        issues.classList.remove("show");
        issues.innerHTML = "";
        return;
      }
      issues.classList.add("show");
      issues.innerHTML = "<strong>请先处理这些问题</strong><ul>" + problems.map((item) => "<li>" + esc(item) + "</li>").join("") + "</ul>";
    }

    function setStatus(sync) {
      $("state").innerHTML = sync.running ? '<span class="warn">同步中</span>' : '<span class="ok">空闲</span>';
      $("last-run").textContent = sync.last_run_ended_at ? new Date(sync.last_run_ended_at).toLocaleString() : "--";
      $("duration").textContent = sync.last_duration_ms ? sync.last_duration_ms + " ms" : "--";
      $("ipv6").textContent = sync.current_ipv6 || "--";
      $("state-detail").textContent = sync.last_error || (sync.running ? "正在同步" : "等待下一次同步");
      if (sync.last_error) $("state").innerHTML = '<span class="bad">异常</span>';
      renderResults(sync.results || []);
    }

    function renderResults(results) {
      $("result-count").textContent = results.length + " 条";
      const box = $("results");
      if (!Array.isArray(results) || results.length === 0) {
        box.innerHTML = '<div class="empty">还没有同步结果</div>';
        return;
      }
      box.innerHTML = '<table><thead><tr><th style="width:70px">时间</th><th>域名</th><th>IP</th><th style="width:76px">结果</th></tr></thead><tbody>' +
        results.map((item) => {
          const status = item.status || "unknown";
          return '<tr><td>' + esc(item.time ? new Date(item.time).toLocaleTimeString() : "--") + '</td>' +
            '<td>' + esc(item.domain) + '</td>' +
            '<td>' + esc(item.ip) + '</td>' +
            '<td><span class="status-tag ' + esc(status) + '">' + resultText(status) + '</span></td></tr>';
        }).join("") + '</tbody></table>';
    }

    function resultText(status) {
      if (status === "updated") return "已更新";
      if (status === "unchanged") return "未变化";
      if (status === "failed") return "失败";
      return status || "--";
    }

    function renderLogs(entries) {
      const box = $("logs");
      if (!Array.isArray(entries) || entries.length === 0) {
        box.innerHTML = '<div class="empty">暂无日志</div>';
        return;
      }
      const nearBottom = box.scrollHeight - box.scrollTop - box.clientHeight < 80;
      box.innerHTML = entries.map((entry) => (
        '<div class="log-row"><div class="log-time">' + esc(entry.time) + '</div>' +
        '<div class="log-level ' + esc(entry.level) + '">' + esc(entry.level) + '</div>' +
        '<div class="log-msg">' + esc(entry.message) + '</div></div>'
      )).join("");
      if (nearBottom) box.scrollTop = box.scrollHeight;
    }

    async function refresh(options = {}) {
      const res = await fetch("/api/status", { cache: "no-store" });
      const data = await res.json();
      if (options.populate || !loadedConfig || !formDirty) {
        fillConfig(data.config);
        loadedConfig = true;
        formDirty = false;
      }
      setStatus(data.sync || {});
      renderLogs(data.logs || []);
    }

    async function save() {
      const cfg = readConfig();
      const problems = validateBeforeSave(cfg);
      showProblems(problems);
      if (problems.length > 0) {
        message.className = "message bad";
        message.textContent = "配置还不完整";
        return;
      }
      $("save").disabled = true;
      message.className = "message";
      message.textContent = "保存中...";
      try {
        const res = await fetch("/api/config", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(cfg),
        });
        const data = await res.json();
        if (!res.ok || !data.ok) throw new Error(data.msg || "保存失败");
        apiKey.value = "";
        apiSecret.value = "";
        showProblems([]);
        message.className = "message ok";
        message.textContent = "已保存";
        formDirty = false;
        await refresh({ populate: true });
      } catch (err) {
        message.className = "message bad";
        message.textContent = "失败：" + (err && err.message ? err.message : err);
      } finally {
        $("save").disabled = false;
      }
    }

    async function runNow() {
      $("run").disabled = true;
      try {
        await fetch("/api/run", { method: "POST" });
        message.className = "message";
        message.textContent = "已触发同步";
        setTimeout(refresh, 800);
      } finally {
        setTimeout(() => { $("run").disabled = false; }, 1000);
      }
    }

    async function clearLogs() {
      await fetch("/api/logs/clear", { method: "POST" });
      await refresh();
    }

    document.querySelector("main").addEventListener("input", (event) => {
      if (event.target.matches("input, select")) markDirty();
    });
    $("add-url").addEventListener("click", () => {
      addListRow("ipv6-urls", "", "每行一个接口").focus();
      markDirty();
    });
    $("add-domain").addEventListener("click", () => {
      addListRow("ipv6-domains", "", "每行一个域名").focus();
      markDirty();
    });
    $("toggle-secret").addEventListener("click", () => {
      const hidden = apiSecret.classList.toggle("secret-hidden");
      $("toggle-secret").textContent = hidden ? "显示" : "隐藏";
    });
    $("save").addEventListener("click", save);
    $("run").addEventListener("click", runNow);
    $("refresh").addEventListener("click", () => refresh({ populate: true }));
    $("open-logs").addEventListener("click", () => {
      $("logs-page").classList.remove("hidden");
      $("logs-page").setAttribute("aria-hidden", "false");
      refresh();
    });
    $("close-logs").addEventListener("click", () => {
      $("logs-page").classList.add("hidden");
      $("logs-page").setAttribute("aria-hidden", "true");
    });
    $("clear-logs").addEventListener("click", clearLogs);
    refresh({ populate: true });
    setInterval(refresh, 3000);
  </script>
</body>
</html>`
