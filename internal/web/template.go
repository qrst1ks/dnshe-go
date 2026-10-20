package web

const indexHTML = `<!doctype html>
<html lang="zh-CN" data-theme="dark">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, viewport-fit=cover">
  <meta name="color-scheme" content="dark light">
  <meta name="description" content="dnshe-go IPv6 DDNS 控制台，同时管理 DNSHE 与 Cloudflare 的 AAAA 记录。">
  <title>{{.Title}}</title>
  <script>
    (function () {
      try {
        var stored = localStorage.getItem("dnshe-go-theme");
        var dark = !window.matchMedia || !window.matchMedia("(prefers-color-scheme: light)").matches;
        document.documentElement.setAttribute("data-theme", stored === "light" || stored === "dark" ? stored : (dark ? "dark" : "light"));
      } catch (err) {
        document.documentElement.setAttribute("data-theme", "dark");
      }
    })();
  </script>
  <style>
    :root {
      --mono: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace;
      --sans: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
      --radius: 20px;
    }
    [data-theme="dark"] {
      color-scheme: dark;
      --bg: #070a12;
      --bg-soft: #0a0f1c;
      --surface: rgba(16, 24, 40, .66);
      --surface-solid: #111a2c;
      --surface-2: rgba(255, 255, 255, .04);
      --input: rgba(4, 8, 16, .5);
      --input-focus: rgba(7, 12, 23, .75);
      --line: rgba(148, 163, 184, .15);
      --line-strong: rgba(148, 163, 184, .3);
      --text: #e8eefb;
      --text-2: #c3cddf;
      --muted: #8b9ab4;
      --faint: #64748b;
      --accent: #818cf8;
      --accent-strong: #6366f1;
      --accent-soft: rgba(129, 140, 248, .13);
      --focus: rgba(129, 140, 248, .22);
      --ok: #4ade80;
      --ok-soft: rgba(74, 222, 128, .12);
      --warn: #fbbf24;
      --warn-soft: rgba(251, 191, 36, .12);
      --bad: #fb7185;
      --bad-soft: rgba(251, 113, 133, .11);
      --shadow: 0 28px 70px -40px rgba(0, 0, 0, .95);
      --glow-a: rgba(99, 102, 241, .14);
      --glow-b: rgba(34, 211, 238, .08);
      --topbar: rgba(7, 10, 18, .72);
      --input-bg: #0b1220;
    }
    [data-theme="light"] {
      color-scheme: light;
      --bg: #f4f6fb;
      --bg-soft: #eef2f9;
      --surface: rgba(255, 255, 255, .82);
      --surface-solid: #ffffff;
      --surface-2: rgba(15, 23, 42, .035);
      --input: #f8fafc;
      --input-focus: #ffffff;
      --line: #e2e8f1;
      --line-strong: #cbd5e1;
      --text: #0f172a;
      --text-2: #334155;
      --muted: #64748b;
      --faint: #94a3b8;
      --accent: #4f46e5;
      --accent-strong: #4338ca;
      --accent-soft: rgba(79, 70, 229, .08);
      --focus: rgba(79, 70, 229, .16);
      --ok: #059669;
      --ok-soft: rgba(5, 150, 105, .09);
      --warn: #b45309;
      --warn-soft: rgba(180, 83, 9, .08);
      --bad: #dc2626;
      --bad-soft: rgba(220, 38, 38, .07);
      --shadow: 0 26px 60px -44px rgba(15, 23, 42, .55);
      --glow-a: rgba(99, 102, 241, .1);
      --glow-b: rgba(14, 165, 233, .08);
      --topbar: rgba(244, 246, 251, .8);
      --input-bg: #f8fafc;
    }
    *, *::before, *::after { box-sizing: border-box; }
    html { -webkit-text-size-adjust: 100%; scroll-behavior: smooth; }
    body {
      margin: 0;
      min-height: 100vh;
      background: var(--bg);
      color: var(--text);
      font: 14px/1.55 var(--sans);
      -webkit-font-smoothing: antialiased;
    }
    body::before {
      content: "";
      position: fixed;
      inset: 0;
      z-index: -1;
      pointer-events: none;
      background:
        radial-gradient(820px 480px at 8% -12%, var(--glow-a), transparent 62%),
        radial-gradient(720px 460px at 102% 4%, var(--glow-b), transparent 60%);
    }
    body.drawer-open { overflow: hidden; }
    .svg-sprite { position: absolute; width: 0; height: 0; overflow: hidden; }
    .icon { width: 18px; height: 18px; display: inline-block; vertical-align: middle; flex: none; }
    .mono { font-family: var(--mono); }
    button, input, select, textarea { font: inherit; }
    button { cursor: pointer; }
    a { color: var(--accent); text-decoration: none; }
    a:hover { color: var(--accent-strong); }
    [hidden] { display: none !important; }
    .shell { width: min(1200px, calc(100% - 40px)); margin: 0 auto; }
    .sr-only {
      position: absolute;
      width: 1px;
      height: 1px;
      padding: 0;
      margin: -1px;
      overflow: hidden;
      clip: rect(0, 0, 0, 0);
      white-space: nowrap;
      border: 0;
    }

    .topbar {
      position: sticky;
      top: 0;
      z-index: 50;
      border-bottom: 1px solid var(--line);
      background: var(--topbar);
      -webkit-backdrop-filter: blur(18px) saturate(1.4);
      backdrop-filter: blur(18px) saturate(1.4);
    }
    .topbar-inner {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      min-height: 68px;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 11px;
      min-width: 0;
      color: var(--text);
    }
    .brand:hover { color: var(--text); }
    .brand-mark {
      width: 38px;
      height: 38px;
      flex: none;
      display: grid;
      place-items: center;
      border-radius: 12px;
      color: #eef2ff;
      background: linear-gradient(135deg, #6366f1, #0ea5e9);
      box-shadow: 0 10px 26px -12px rgba(99, 102, 241, .8);
    }
    .brand-mark .icon { width: 22px; height: 22px; }
    .brand-text { display: flex; flex-direction: column; line-height: 1.2; min-width: 0; }
    .brand-text strong { font-size: 16px; letter-spacing: -.01em; }
    .brand-text span { color: var(--muted); font-size: 11px; margin-top: 2px; }
    .topbar-right { display: flex; align-items: center; gap: 12px; min-width: 0; }
    .topbar-actions { display: flex; align-items: center; gap: 7px; }

    .sync-pill {
      display: inline-flex;
      align-items: center;
      gap: 7px;
      height: 30px;
      padding: 0 11px;
      border-radius: 999px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      color: var(--muted);
      font-size: 11.5px;
      font-weight: 600;
      white-space: nowrap;
    }
    .sync-pill .dot { width: 6px; height: 6px; border-radius: 50%; background: currentColor; }
    .sync-pill.is-on { color: var(--ok); border-color: transparent; background: var(--ok-soft); }
    .sync-pill.is-running { color: var(--accent); border-color: transparent; background: var(--accent-soft); }
    .sync-pill.is-running .dot { animation: pulse 1.3s ease-out infinite; }
    .sync-pill.is-bad { color: var(--bad); border-color: transparent; background: var(--bad-soft); }
    .sync-pill.is-idle { color: var(--warn); border-color: transparent; background: var(--warn-soft); }

    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 7px;
      height: 38px;
      padding: 0 14px;
      border-radius: 11px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      color: var(--text-2);
      font-size: 12.5px;
      font-weight: 600;
      transition: background .16s, border-color .16s, color .16s, transform .16s;
      white-space: nowrap;
    }
    .btn:hover { color: var(--text); border-color: var(--line-strong); background: var(--surface-2); }
    .btn:active { transform: translateY(1px); }
    .btn:focus-visible, .icon-btn:focus-visible, .reveal:focus-visible, .seg-btn:focus-visible, .preset:focus-visible, .filter:focus-visible, .chip-remove:focus-visible, .copy-btn:focus-visible {
      outline: 2px solid var(--accent);
      outline-offset: 2px;
    }
    .btn.primary {
      border-color: transparent;
      color: #fff;
      background: linear-gradient(135deg, #6366f1, #4f6ef7);
      box-shadow: 0 10px 24px -14px rgba(99, 102, 241, .95);
    }
    .btn.primary:hover { background: linear-gradient(135deg, #7477f5, #5b7bff); color: #fff; }
    .btn:disabled { opacity: .6; cursor: wait; }
    .btn .icon { width: 16px; height: 16px; }
    .btn-dot { display: none; width: 6px; height: 6px; border-radius: 50%; background: #c7d2fe; box-shadow: 0 0 0 3px rgba(199, 210, 254, .25); }
    .btn.is-dirty .btn-dot { display: inline-block; }
    .btn.loading .icon { animation: spin .9s linear infinite; }
    .icon-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 38px;
      height: 38px;
      padding: 0;
      border-radius: 11px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      color: var(--muted);
      transition: background .16s, color .16s, border-color .16s;
    }
    .icon-btn:hover { color: var(--text); border-color: var(--line-strong); }
    .icon-btn .icon { width: 18px; height: 18px; }
    .icon-btn.small { width: 30px; height: 30px; border-radius: 9px; }
    .icon-btn.small .icon { width: 15px; height: 15px; }
    .icon-sun { display: none; }
    [data-theme="dark"] .icon-sun { display: block; }
    [data-theme="dark"] .icon-moon { display: none; }
    .spinning .icon { animation: spin .8s linear infinite; }

    main { padding: 26px 0 8px; }
    .alert {
      display: flex;
      align-items: flex-start;
      gap: 12px;
      margin-bottom: 18px;
      padding: 14px 16px;
      border-radius: 16px;
      border: 1px solid rgba(251, 113, 133, .3);
      background: var(--bad-soft);
      color: var(--text);
    }
    .alert > .icon { color: var(--bad); margin-top: 1px; }
    .alert-body { flex: 1; min-width: 0; }
    .alert-body strong { display: block; font-size: 13px; }
    .alert-body ul { margin: 6px 0 0; padding-left: 18px; color: var(--bad); font-size: 12px; line-height: 1.65; }
    .alert .icon-btn { border: 0; background: transparent; }

    .hero {
      position: relative;
      overflow: hidden;
      padding: 26px 28px 24px;
      border-radius: 24px;
      border: 1px solid rgba(148, 163, 255, .16);
      color: #eef2ff;
      background:
        radial-gradient(900px 380px at 0% 0%, rgba(99, 102, 241, .38), transparent 60%),
        radial-gradient(680px 460px at 100% 100%, rgba(34, 211, 238, .16), transparent 58%),
        linear-gradient(135deg, #131c38 0%, #0a1020 68%);
      box-shadow: 0 34px 90px -50px rgba(30, 41, 90, .95);
    }
    .hero::after {
      content: "";
      position: absolute;
      inset: 0;
      pointer-events: none;
      opacity: .5;
      background-image:
        linear-gradient(rgba(255, 255, 255, .035) 1px, transparent 1px),
        linear-gradient(90deg, rgba(255, 255, 255, .035) 1px, transparent 1px);
      background-size: 46px 46px;
      -webkit-mask-image: radial-gradient(680px 340px at 78% 0%, #000, transparent 75%);
      mask-image: radial-gradient(680px 340px at 78% 0%, #000, transparent 75%);
    }
    .hero > * { position: relative; z-index: 1; }
    .hero-top { display: flex; align-items: center; justify-content: space-between; gap: 14px; flex-wrap: wrap; }
    .eyebrow { display: inline-flex; align-items: center; gap: 9px; color: #c7d2fe; font-size: 12.5px; font-weight: 650; letter-spacing: .01em; }
    .live-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: #4ade80;
      box-shadow: 0 0 0 4px rgba(74, 222, 128, .16);
      animation: breathe 2.6s ease-in-out infinite;
    }
    .hero.is-running .live-dot { background: #a5b4fc; box-shadow: 0 0 0 4px rgba(165, 180, 252, .2); animation: pulse 1.3s ease-out infinite; }
    .hero.is-bad .live-dot { background: #fb7185; box-shadow: 0 0 0 4px rgba(251, 113, 133, .18); animation: none; }
    .hero.is-idle .live-dot { background: #fbbf24; box-shadow: 0 0 0 4px rgba(251, 191, 36, .16); animation: none; }
    .hero-badges { display: flex; align-items: center; gap: 7px; flex-wrap: wrap; }
    .hero-provider {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 5px 10px;
      border-radius: 999px;
      border: 1px solid rgba(255, 255, 255, .12);
      background: rgba(255, 255, 255, .06);
      color: #dbe4ff;
      font-size: 11.5px;
      font-weight: 600;
    }
    .hero-provider .icon { width: 13px; height: 13px; opacity: .9; }
    .hero-provider.is-empty { color: #94a3c4; border-style: dashed; }
    .ip-block { margin-top: 24px; }
    .ip-label { color: #8fa1c7; font-size: 10.5px; font-weight: 700; letter-spacing: .16em; text-transform: uppercase; }
    .ip-row { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; margin-top: 7px; }
    .ip-value {
      font-family: var(--mono);
      font-size: clamp(20px, 3.5vw, 30px);
      font-weight: 600;
      letter-spacing: -.02em;
      line-height: 1.25;
      color: #ffffff;
      word-break: break-all;
    }
    .copy-btn {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      height: 32px;
      padding: 0 11px;
      border-radius: 9px;
      border: 1px solid rgba(255, 255, 255, .14);
      background: rgba(255, 255, 255, .06);
      color: #c7d2fe;
      font-size: 11.5px;
      font-weight: 600;
      transition: background .16s, color .16s;
    }
    .copy-btn:hover { background: rgba(255, 255, 255, .12); color: #fff; }
    .copy-btn:disabled { opacity: .4; cursor: default; }
    .copy-btn .icon { width: 14px; height: 14px; }
    .copy-done { display: none; }
    .copy-btn.copied { color: #4ade80; border-color: rgba(74, 222, 128, .3); background: rgba(74, 222, 128, .1); }
    .copy-btn.copied .copy-text { display: none; }
    .copy-btn.copied .copy-done { display: inline; }
    .hero-detail { margin-top: 10px; max-width: 820px; color: #93a4c8; font-size: 12.5px; line-height: 1.6; word-break: break-word; }
    .hero-stats {
      display: grid;
      grid-template-columns: repeat(4, minmax(0, 1fr));
      gap: 18px;
      margin-top: 24px;
      padding-top: 20px;
      border-top: 1px solid rgba(255, 255, 255, .09);
    }
    .stat { min-width: 0; display: flex; flex-direction: column; gap: 3px; }
    .stat-label { color: #8494b8; font-size: 10.5px; font-weight: 700; letter-spacing: .1em; text-transform: uppercase; }
    .stat strong { color: #f1f5ff; font-size: 17px; font-weight: 650; letter-spacing: -.01em; }
    .stat-sub { color: #7f90b4; font-size: 11.5px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

    .section-head {
      display: flex;
      align-items: flex-end;
      justify-content: space-between;
      gap: 16px;
      margin: 30px 2px 14px;
    }
    .section-head h2 { margin: 0; font-size: 16px; letter-spacing: -.01em; }
    .section-head p { margin: 3px 0 0; color: var(--muted); font-size: 12px; }
    .section-meta { color: var(--muted); font-size: 11.5px; font-weight: 600; white-space: nowrap; }
    .provider-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; }
    .card {
      min-width: 0;
      border-radius: var(--radius);
      border: 1px solid var(--line);
      background: var(--surface);
      box-shadow: var(--shadow);
      -webkit-backdrop-filter: blur(14px);
      backdrop-filter: blur(14px);
      overflow: hidden;
    }
    .provider-card { position: relative; --brand-a: #22d3ee; --brand-b: #3b82f6; }
    .provider-card.dnshe { --brand-a: #22d3ee; --brand-b: #2563eb; }
    .provider-card.cloudflare { --brand-a: #fbbf24; --brand-b: #f97316; }
    .provider-card::before {
      content: "";
      position: absolute;
      inset: 0 0 auto;
      height: 2px;
      background: linear-gradient(90deg, var(--brand-a), var(--brand-b), transparent 92%);
      opacity: .9;
    }
    .provider-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 14px; padding: 20px 22px 15px; }
    .provider-id { display: flex; align-items: center; gap: 12px; min-width: 0; }
    .provider-logo {
      width: 46px;
      height: 46px;
      flex: none;
      display: grid;
      place-items: center;
      border-radius: 14px;
      background: #fff;
      border: 1px solid rgba(148, 163, 184, .22);
      box-shadow: 0 10px 24px -18px rgba(2, 6, 23, .8);
    }
    .provider-logo img { display: block; width: 28px; height: 28px; object-fit: contain; }
    .provider-card.cloudflare .provider-logo img { width: 32px; }
    .provider-id h2 { margin: 0; font-size: 16px; letter-spacing: -.01em; }
    .provider-id p { margin: 2px 0 0; color: var(--muted); font-size: 11.5px; }
    .provider-head-right { display: flex; align-items: center; gap: 10px; flex: none; }
    .status-pill {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      padding: 5px 9px;
      border-radius: 999px;
      font-size: 11px;
      font-weight: 650;
      white-space: nowrap;
      border: 1px solid transparent;
    }
    .status-pill .dot { width: 6px; height: 6px; border-radius: 50%; background: currentColor; }
    .status-pill.is-on { color: var(--ok); background: var(--ok-soft); }
    .status-pill.is-off { color: var(--muted); background: var(--surface-2); border-color: var(--line); }
    .status-pill.is-warn { color: var(--warn); background: var(--warn-soft); }
    .switch { position: relative; display: inline-flex; align-items: center; flex: none; cursor: pointer; }
    .switch input { position: absolute; opacity: 0; width: 1px; height: 1px; }
    .switch-track {
      position: relative;
      display: block;
      width: 46px;
      height: 26px;
      border-radius: 999px;
      background: var(--line-strong);
      transition: background .2s;
    }
    .switch-thumb {
      position: absolute;
      top: 3px;
      left: 3px;
      width: 20px;
      height: 20px;
      border-radius: 50%;
      background: #fff;
      box-shadow: 0 2px 6px rgba(0, 0, 0, .25);
      transition: transform .22s cubic-bezier(.22, 1, .36, 1);
    }
    .switch input:checked + .switch-track { background: linear-gradient(135deg, var(--brand-a), var(--brand-b)); }
    .switch input:checked + .switch-track .switch-thumb { transform: translateX(20px); }
    .switch input:focus-visible + .switch-track { outline: 2px solid var(--accent); outline-offset: 3px; }
    .provider-body {
      display: flex;
      flex-direction: column;
      gap: 16px;
      padding: 0 22px 22px;
      transition: opacity .2s, filter .2s;
    }
    .provider-card.is-off .provider-body { opacity: .55; filter: saturate(.65); }
    .provider-card.is-off .provider-logo { filter: saturate(.35); }
    .form-row { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 13px; }
    .field { display: flex; flex-direction: column; gap: 7px; min-width: 0; }
    .field-head { display: flex; align-items: center; justify-content: space-between; gap: 10px; min-height: 18px; }
    .field label, .field-label { color: var(--text-2); font-size: 12px; font-weight: 650; }
    .mini-tag {
      padding: 1px 7px;
      border-radius: 999px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      color: var(--faint);
      font-size: 10px;
      font-weight: 650;
    }
    .mini-tag.ok { border-color: transparent; background: var(--ok-soft); color: var(--ok); }
    .input-wrap { position: relative; display: flex; align-items: center; }
    input[type="text"], input[type="password"], input[type="number"], select, textarea {
      width: 100%;
      height: 42px;
      min-width: 0;
      padding: 0 13px;
      border-radius: 11px;
      border: 1px solid var(--line);
      background: var(--input);
      color: var(--text);
      transition: border-color .16s, box-shadow .16s, background .16s;
    }
    input::placeholder, textarea::placeholder { color: var(--faint); }
    input:hover, select:hover { border-color: var(--line-strong); }
    input:focus, select:focus, textarea:focus {
      outline: 0;
      border-color: var(--accent);
      background: var(--input-focus);
      box-shadow: 0 0 0 3px var(--focus);
    }
    input[type="password"] { letter-spacing: .08em; }
    .input-wrap input { padding-right: 44px; }
    .reveal {
      position: absolute;
      right: 4px;
      display: grid;
      place-items: center;
      width: 34px;
      height: 34px;
      padding: 0;
      border: 0;
      border-radius: 9px;
      background: transparent;
      color: var(--muted);
    }
    .reveal:hover { background: var(--surface-2); color: var(--text); }
    .reveal:disabled { opacity: .45; cursor: default; }
    .reveal .icon { width: 17px; height: 17px; }
    .select-wrap { position: relative; }
    .select-wrap select { appearance: none; -webkit-appearance: none; padding-right: 40px; }
    .select-wrap .icon { position: absolute; right: 13px; top: 50%; transform: translateY(-50%) rotate(90deg); width: 15px; height: 15px; color: var(--muted); pointer-events: none; }
    .hint { display: flex; align-items: flex-start; gap: 6px; color: var(--faint); font-size: 11.5px; line-height: 1.6; }
    .hint .icon { width: 13px; height: 13px; margin-top: 2px; color: var(--faint); }
    .advanced { padding-top: 13px; border-top: 1px dashed var(--line); }
    .advanced summary {
      display: flex;
      align-items: center;
      gap: 6px;
      width: fit-content;
      list-style: none;
      color: var(--muted);
      font-size: 11.5px;
      font-weight: 650;
      cursor: pointer;
      user-select: none;
    }
    .advanced summary::-webkit-details-marker { display: none; }
    .advanced summary .icon { width: 13px; height: 13px; transition: transform .2s; }
    .advanced[open] summary .icon { transform: rotate(90deg); }
    .advanced .field { margin-top: 11px; }
    .chips {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 7px;
      min-height: 46px;
      padding: 8px;
      border-radius: 11px;
      border: 1px solid var(--line);
      background: var(--input);
      cursor: text;
      transition: border-color .16s, box-shadow .16s, background .16s;
    }
    .chips:hover { border-color: var(--line-strong); }
    .chips:focus-within { border-color: var(--accent); background: var(--input-focus); box-shadow: 0 0 0 3px var(--focus); }
    .chip {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      max-width: 100%;
      padding: 4px 5px 4px 9px;
      border-radius: 8px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      font-family: var(--mono);
      font-size: 11.5px;
      color: var(--text-2);
    }
    .chip-text { max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    .chip-remove {
      display: grid;
      place-items: center;
      width: 18px;
      height: 18px;
      padding: 0;
      border: 0;
      border-radius: 5px;
      background: transparent;
      color: var(--faint);
    }
    .chip-remove:hover { color: var(--bad); background: var(--bad-soft); }
    .chip-remove .icon { width: 12px; height: 12px; }
    .chip-entry {
      flex: 1 0 100%;
      min-width: 0;
      height: 28px !important;
      padding: 0 4px !important;
      border: 0 !important;
      border-radius: 0 !important;
      background: transparent !important;
      box-shadow: none !important;
      font-family: var(--mono);
      font-size: 12px;
    }
    .chip-entry:focus { box-shadow: none !important; }
    .counter { color: var(--faint); font-size: 10.5px; font-variant-numeric: tabular-nums; }

    .lower-grid { display: grid; grid-template-columns: minmax(0, 1.32fr) minmax(0, 1fr); gap: 18px; margin-top: 18px; }
    .panel-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 14px; padding: 20px 22px 0; }
    .panel-head h2 { margin: 0; font-size: 15px; letter-spacing: -.01em; }
    .panel-head p { margin: 3px 0 0; color: var(--muted); font-size: 11.5px; }
    .panel-body { padding: 17px 22px 22px; display: flex; flex-direction: column; gap: 16px; }
    .segmented {
      display: inline-flex;
      gap: 3px;
      padding: 4px;
      border-radius: 12px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      flex-wrap: wrap;
    }
    .seg-btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      height: 34px;
      padding: 0 13px;
      border: 0;
      border-radius: 9px;
      background: transparent;
      color: var(--muted);
      font-size: 12px;
      font-weight: 650;
      transition: background .16s, color .16s;
    }
    .seg-btn .icon { width: 15px; height: 15px; }
    .seg-btn:hover { color: var(--text); }
    .seg-btn.active { color: var(--text); background: var(--surface-solid); box-shadow: 0 2px 8px rgba(0, 0, 0, .12); }
    .source-panel { display: flex; flex-direction: column; gap: 13px; }
    .interface-row { display: flex; align-items: center; gap: 8px; }
    .interface-row .select-wrap { flex: 1; min-width: 0; }
    .mono-input { font-family: var(--mono); font-size: 12px; }
    .setting-block + .setting-block { margin-top: 15px; padding-top: 15px; border-top: 1px solid var(--line); }
    .preset-row { display: flex; flex-wrap: wrap; gap: 6px; }
    .preset {
      height: 27px;
      padding: 0 9px;
      border-radius: 7px;
      border: 1px solid var(--line);
      background: transparent;
      color: var(--muted);
      font-size: 11px;
      font-weight: 600;
      transition: color .16s, border-color .16s, background .16s;
    }
    .preset:hover { color: var(--text); border-color: var(--line-strong); }
    .preset.active { color: var(--accent); border-color: var(--accent); background: var(--accent-soft); }
    .input-suffix { position: relative; display: flex; align-items: center; }
    .input-suffix input { padding-right: 48px; }
    .input-suffix span { position: absolute; right: 13px; color: var(--faint); font-size: 11px; pointer-events: none; }
    .source-current {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 11px 12px;
      border-radius: 13px;
      border: 1px solid var(--line);
      background: var(--surface-2);
    }
    .source-current-icon {
      width: 34px;
      height: 34px;
      flex: none;
      display: grid;
      place-items: center;
      border-radius: 10px;
      border: 1px solid rgba(34, 211, 238, .18);
      background: rgba(34, 211, 238, .09);
      color: #22d3ee;
    }
    .source-current-icon .icon { width: 17px; height: 17px; }
    .source-current-copy { display: flex; flex-direction: column; gap: 1px; min-width: 0; flex: 1; }
    .source-current-copy span { color: var(--muted); font-size: 10px; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; }
    .source-current-copy strong { color: var(--text); font-family: var(--mono); font-size: 12.5px; font-weight: 600; word-break: break-all; }
    .source-current-mode {
      flex: none;
      padding: 4px 9px;
      border-radius: 999px;
      border: 1px solid var(--line);
      background: var(--surface);
      color: var(--muted);
      font-size: 10.5px;
      font-weight: 650;
      white-space: nowrap;
    }
    .source-facts { display: flex; flex-wrap: wrap; gap: 14px; color: var(--faint); font-size: 11px; }
    .source-facts span { display: inline-flex; align-items: center; gap: 6px; }
    .source-facts .dot { width: 5px; height: 5px; border-radius: 50%; background: var(--accent); }
    .sync-note {
      display: flex;
      align-items: flex-start;
      gap: 8px;
      padding: 10px 11px;
      border-radius: 11px;
      background: var(--surface-2);
      border: 1px solid var(--line);
      color: var(--muted);
      font-size: 11px;
      line-height: 1.55;
    }
    .sync-note .icon { width: 14px; height: 14px; color: var(--accent); margin-top: 1px; }

    .results-card { margin-top: 18px; }
    .results-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 14px;
      flex-wrap: wrap;
      padding: 19px 22px;
      border-bottom: 1px solid var(--line);
    }
    .results-head h2 { margin: 0; font-size: 15px; letter-spacing: -.01em; }
    .results-head p { margin: 3px 0 0; color: var(--muted); font-size: 11.5px; }
    .filters {
      display: inline-flex;
      gap: 3px;
      padding: 4px;
      border-radius: 11px;
      border: 1px solid var(--line);
      background: var(--surface-2);
    }
    .filter {
      display: inline-flex;
      align-items: center;
      gap: 6px;
      height: 29px;
      padding: 0 10px;
      border: 0;
      border-radius: 8px;
      background: transparent;
      color: var(--muted);
      font-size: 11.5px;
      font-weight: 650;
      transition: background .16s, color .16s;
    }
    .filter span { color: var(--faint); font-size: 10.5px; font-variant-numeric: tabular-nums; }
    .filter:hover { color: var(--text); }
    .filter.active { color: var(--text); background: var(--surface-solid); box-shadow: 0 2px 8px rgba(0, 0, 0, .12); }
    .filter.active span { color: var(--accent); }
    .results-list { padding: 6px 8px 8px; }
    .result-row {
      display: flex;
      align-items: flex-start;
      gap: 13px;
      padding: 13px 14px;
      border-radius: 13px;
      transition: background .15s;
    }
    .result-row + .result-row { border-top: 1px solid var(--line); border-radius: 0; }
    .result-row:hover { background: var(--surface-2); }
    .result-logo {
      width: 34px;
      height: 34px;
      flex: none;
      display: grid;
      place-items: center;
      border-radius: 10px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      color: var(--muted);
    }
    .result-logo .icon { width: 17px; height: 17px; }
    .result-logo.dnshe { color: #22d3ee; background: rgba(34, 211, 238, .09); border-color: rgba(34, 211, 238, .18); }
    .result-logo.cloudflare { color: #fb923c; background: rgba(251, 146, 60, .09); border-color: rgba(251, 146, 60, .18); }
    .result-main { flex: 1; min-width: 0; }
    .result-title { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }
    .result-domain { font-family: var(--mono); font-size: 12.5px; font-weight: 650; word-break: break-all; }
    .provider-name { color: var(--faint); font-size: 10px; font-weight: 700; letter-spacing: .08em; text-transform: uppercase; }
    .tag {
      padding: 2px 7px;
      border-radius: 999px;
      background: var(--surface-2);
      color: var(--muted);
      font-size: 10px;
      font-weight: 700;
    }
    .tag.updated, .tag.created { background: var(--ok-soft); color: var(--ok); }
    .tag.failed { background: var(--bad-soft); color: var(--bad); }
    .result-meta { display: flex; align-items: center; gap: 7px; flex-wrap: wrap; margin-top: 5px; color: var(--muted); font-size: 11px; }
    .result-meta .mono { color: var(--text-2); font-size: 10.5px; word-break: break-all; }
    .result-meta .arrow { color: var(--faint); }
    .result-error { margin-top: 5px; color: var(--bad); font-size: 11px; line-height: 1.55; word-break: break-word; }
    .result-side { display: flex; flex-direction: column; align-items: flex-end; gap: 5px; flex: none; color: var(--faint); font-size: 10.5px; white-space: nowrap; }
    .record-badge {
      padding: 2px 6px;
      border-radius: 6px;
      border: 1px solid var(--line);
      background: var(--surface-2);
      color: var(--muted);
      font-family: var(--mono);
      font-size: 10px;
    }
    .empty-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 7px;
      padding: 44px 20px;
      text-align: center;
      color: var(--faint);
    }
    .empty-state .icon { width: 26px; height: 26px; opacity: .55; margin-bottom: 2px; }
    .empty-state strong { color: var(--muted); font-size: 12.5px; font-weight: 650; }
    .empty-state p { margin: 0; max-width: 340px; font-size: 11.5px; line-height: 1.6; }

    .footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 14px;
      flex-wrap: wrap;
      padding: 22px 3px 30px;
      color: var(--faint);
      font-size: 11px;
    }
    .footer-left, .footer-right { display: flex; align-items: center; gap: 10px; min-width: 0; }
    .footer-left .icon { width: 14px; height: 14px; color: var(--muted); }
    .config-path {
      max-width: min(420px, 60vw);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      color: var(--muted);
      font-family: var(--mono);
      font-size: 10.5px;
    }
    .footer-link { display: inline-flex; align-items: center; gap: 5px; color: var(--muted); font-weight: 600; }
    .footer-link:hover { color: var(--text); }
    .footer-link .icon { width: 14px; height: 14px; }

    .drawer {
      position: fixed;
      top: 0;
      right: 0;
      bottom: 0;
      z-index: 100;
      display: flex;
      flex-direction: column;
      width: min(590px, 100vw);
      border-left: 1px solid rgba(148, 163, 184, .16);
      background: #080d17;
      color: #dce4f5;
      box-shadow: -40px 0 100px rgba(0, 0, 0, .55);
      transform: translateX(100%);
      visibility: hidden;
      transition: transform .28s cubic-bezier(.22, 1, .36, 1), visibility 0s linear .3s;
    }
    .drawer.open { transform: none; visibility: visible; transition: transform .28s cubic-bezier(.22, 1, .36, 1), visibility 0s; }
    .drawer-backdrop {
      position: fixed;
      inset: 0;
      z-index: 90;
      background: rgba(2, 6, 15, .55);
      -webkit-backdrop-filter: blur(3px);
      backdrop-filter: blur(3px);
      opacity: 0;
      pointer-events: none;
      transition: opacity .25s;
    }
    .drawer-backdrop.open { opacity: 1; pointer-events: auto; }
    .drawer-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 14px;
      padding: 18px 20px;
      border-bottom: 1px solid rgba(148, 163, 184, .12);
      background: rgba(255, 255, 255, .015);
    }
    .drawer-head h2 { margin: 0; color: #f1f5ff; font-size: 15px; }
    .drawer-head p { margin: 3px 0 0; color: #6f80a0; font-size: 11px; }
    .drawer-actions { display: flex; align-items: center; gap: 7px; }
    .drawer-actions .icon-btn { border-color: rgba(148, 163, 184, .16); background: rgba(255, 255, 255, .04); color: #a9b7d0; }
    .drawer-actions .icon-btn:hover { color: #fff; background: rgba(255, 255, 255, .08); }
    .drawer-actions .btn { height: 34px; border-color: rgba(148, 163, 184, .16); background: rgba(255, 255, 255, .04); color: #a9b7d0; font-size: 11.5px; }
    .drawer-actions .btn:hover { color: #fff; background: rgba(255, 255, 255, .08); }
    .term { flex: 1; min-height: 0; overflow: auto; overscroll-behavior: contain; }
    .log-line {
      display: grid;
      grid-template-columns: 138px 46px minmax(0, 1fr);
      gap: 9px;
      padding: 8px 20px;
      border-bottom: 1px solid rgba(148, 163, 184, .06);
      font-family: var(--mono);
      font-size: 11px;
      line-height: 1.55;
    }
    .log-time { color: #5a6b89; }
    .log-level { font-size: 9.5px; font-weight: 700; letter-spacing: .06em; }
    .log-level.INFO { color: #7dd3fc; }
    .log-level.WARN { color: #fbbf24; }
    .log-level.ERROR { color: #fb7185; }
    .log-msg { color: #c3cfe4; word-break: break-word; }
    .logs-empty { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 48px 20px; color: #5a6b89; font-size: 11.5px; }
    .logs-empty .icon { width: 24px; height: 24px; opacity: .6; }

    .toasts {
      position: fixed;
      top: 84px;
      right: 20px;
      z-index: 130;
      display: flex;
      flex-direction: column;
      gap: 9px;
      width: min(360px, calc(100vw - 32px));
      pointer-events: none;
    }
    .toast {
      display: flex;
      align-items: flex-start;
      gap: 10px;
      padding: 12px 12px 12px 14px;
      border-radius: 14px;
      border: 1px solid var(--line);
      background: var(--surface-solid);
      color: var(--text);
      box-shadow: 0 20px 50px -24px rgba(0, 0, 0, .55);
      pointer-events: auto;
      animation: toast-in .22s ease both;
    }
    .toast.leaving { animation: toast-out .18s ease both; }
    .toast-icon { display: flex; margin-top: 1px; color: var(--accent); }
    .toast-icon .icon { width: 16px; height: 16px; }
    .toast.success .toast-icon { color: var(--ok); }
    .toast.error .toast-icon { color: var(--bad); }
    .toast-msg { flex: 1; min-width: 0; font-size: 12px; line-height: 1.55; word-break: break-word; }
    .toast-close {
      display: grid;
      place-items: center;
      width: 20px;
      height: 20px;
      padding: 0;
      border: 0;
      border-radius: 6px;
      background: transparent;
      color: var(--faint);
      flex: none;
    }
    .toast-close:hover { color: var(--text); background: var(--surface-2); }
    .toast-close .icon { width: 12px; height: 12px; }

    @keyframes spin { to { transform: rotate(360deg); } }
    @keyframes pulse { 0% { box-shadow: 0 0 0 0 rgba(165, 180, 252, .35); } 70% { box-shadow: 0 0 0 7px rgba(165, 180, 252, 0); } 100% { box-shadow: 0 0 0 0 rgba(165, 180, 252, 0); } }
    @keyframes breathe { 0%, 100% { box-shadow: 0 0 0 4px rgba(74, 222, 128, .14); } 50% { box-shadow: 0 0 0 6px rgba(74, 222, 128, .05); } }
    @keyframes toast-in { from { opacity: 0; transform: translateY(-6px) scale(.98); } to { opacity: 1; transform: none; } }
    @keyframes toast-out { to { opacity: 0; transform: translateY(-4px) scale(.98); } }
    @media (prefers-reduced-motion: reduce) {
      *, *::before, *::after { animation-duration: .01ms !important; animation-iteration-count: 1 !important; transition-duration: .01ms !important; scroll-behavior: auto !important; }
    }

    @media (max-width: 1020px) {
      .hero-stats { grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 16px; }
      .lower-grid { grid-template-columns: minmax(0, 1fr); }
    }
    @media (max-width: 860px) {
      .provider-grid { grid-template-columns: minmax(0, 1fr); }
    }
    @media (max-width: 760px) {
      .shell { width: calc(100% - 24px); }
      .topbar-inner { flex-wrap: wrap; gap: 9px; padding: 10px 0; }
      .topbar-right { width: 100%; justify-content: space-between; }
      .topbar-actions { flex: 1; justify-content: flex-end; }
      .btn.optional .btn-label { display: none; }
      .btn.optional { width: 38px; padding: 0; }
      .github-top { display: none; }
      main { padding-top: 18px; }
      .hero { padding: 22px 20px; border-radius: 20px; }
      .ip-block { margin-top: 20px; }
      .provider-head { padding: 18px 17px 14px; }
      .provider-body { padding: 0 17px 18px; }
      .panel-head { padding: 18px 17px 0; }
      .panel-body { padding: 16px 17px 18px; }
      .results-head { padding: 16px 17px; }
      .result-side { display: none; }
      .section-head { margin-top: 24px; }
      .toasts { top: auto; bottom: 18px; right: 12px; }
      .config-path { max-width: 55vw; }
    }
    @media (max-width: 560px) {
      .brand-text span { display: none; }
      .sync-pill { padding: 0 9px; }
      .form-row { grid-template-columns: minmax(0, 1fr); }
      .hero-stats { gap: 14px; }
      .ip-value { font-size: 18px; }
      .provider-head { align-items: flex-start; }
      .provider-head-right { flex-direction: column-reverse; align-items: flex-end; gap: 7px; }
      .status-pill { font-size: 10px; }
      .segmented { width: 100%; }
      .seg-btn { flex: 1; padding: 0 8px; }
      .log-line { grid-template-columns: minmax(0, 1fr) auto; }
      .log-time { grid-column: 1 / -1; font-size: 10px; }
      .log-msg { grid-column: 1 / -1; }
      .result-row { padding: 12px 10px; }
      .result-logo { width: 30px; height: 30px; border-radius: 9px; }
      .filters { width: 100%; }
      .filter { flex: 1; justify-content: center; }
      .footer { flex-direction: column; align-items: flex-start; gap: 8px; }
      .footer-left { flex-wrap: wrap; row-gap: 6px; }
      .config-path { flex-basis: 100%; max-width: 100%; }
    }
    @media (max-width: 420px) {
      .btn#save .btn-label { display: none; }
      .btn#save { width: 38px; padding: 0; }
      .topbar-actions { gap: 6px; }
      .topbar-right { flex-wrap: wrap; row-gap: 8px; }
      .sync-pill { font-size: 10.5px; }
    }
  </style>
</head>
<body>
  <svg xmlns="http://www.w3.org/2000/svg" class="svg-sprite" aria-hidden="true" focusable="false">
    <symbol id="i-logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 3 4.5 7.5v9L12 21l7.5-4.5v-9L12 3Z"></path>
      <path d="M12 15.2V9.2"></path>
      <path d="m9.4 11.6 2.6-2.5 2.6 2.5"></path>
    </symbol>
    <symbol id="i-globe" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="9"></circle>
      <path d="M3 12h18"></path>
      <path d="M12 3a14.5 14.5 0 0 1 0 18a14.5 14.5 0 0 1 0-18Z"></path>
    </symbol>
    <symbol id="i-cloud" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M17.4 19a4.4 4.4 0 0 0 .3-8.78A6.5 6.5 0 0 0 5.3 12.1 3.6 3.6 0 0 0 6.6 19h10.8Z"></path>
    </symbol>
    <symbol id="i-dns" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="3" y="4" width="18" height="7" rx="2"></rect>
      <rect x="3" y="13" width="18" height="7" rx="2"></rect>
      <path d="M7 7.5h.01M7 16.5h.01"></path>
    </symbol>
    <symbol id="i-link" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="m10.5 13.5 3-3"></path>
      <path d="M8.5 16.5 7 18a3.54 3.54 0 0 1-5-5l3-3a3.54 3.54 0 0 1 5 0"></path>
      <path d="M15.5 7.5 17 6a3.54 3.54 0 0 1 5 5l-3 3a3.54 3.54 0 0 1-5 0"></path>
    </symbol>
    <symbol id="i-nic" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="3" y="4" width="18" height="13" rx="2"></rect>
      <path d="M8 21h8M12 17v4"></path>
    </symbol>
    <symbol id="i-terminal" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="3" y="4" width="18" height="16" rx="2"></rect>
      <path d="m7 9 3 3-3 3"></path>
      <path d="M13 15h4"></path>
    </symbol>
    <symbol id="i-refresh" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M20 11a8 8 0 0 0-13.6-5.4L4 8"></path>
      <path d="M4 4v4h4"></path>
      <path d="M4 13a8 8 0 0 0 13.6 5.4L20 16"></path>
      <path d="M20 20v-4h-4"></path>
    </symbol>
    <symbol id="i-play" viewBox="0 0 24 24" fill="currentColor">
      <path d="M8 5.5v13a1 1 0 0 0 1.54.84l10-6.5a1 1 0 0 0 0-1.68l-10-6.5A1 1 0 0 0 8 5.5Z"></path>
    </symbol>
    <symbol id="i-save" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2Z"></path>
      <path d="M17 21v-8H7v8M7 3v5h8"></path>
    </symbol>
    <symbol id="i-sun" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="4"></circle>
      <path d="M12 2v2M12 20v2M4.93 4.93l1.42 1.42M17.65 17.65l1.42 1.42M2 12h2M20 12h2M4.93 19.07l1.42-1.42M17.65 6.35l1.42-1.42"></path>
    </symbol>
    <symbol id="i-moon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8Z"></path>
    </symbol>
    <symbol id="i-github" viewBox="0 0 16 16" fill="currentColor">
      <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.5-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82a7.6 7.6 0 0 1 2-.27c.68 0 1.37.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path>
    </symbol>
    <symbol id="i-copy" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <rect x="9" y="9" width="12" height="12" rx="2"></rect>
      <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
    </symbol>
    <symbol id="i-x" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
      <path d="m6 6 12 12M18 6 6 18"></path>
    </symbol>
    <symbol id="i-eye" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M2 12s3.5-6.5 10-6.5S22 12 22 12s-3.5 6.5-10 6.5S2 12 2 12Z"></path>
      <circle cx="12" cy="12" r="2.5"></circle>
    </symbol>
    <symbol id="i-eye-off" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M3 3l18 18"></path>
      <path d="M10.6 6.1A10.8 10.8 0 0 1 12 6c6.5 0 10 6 10 6a17.4 17.4 0 0 1-3.2 3.8"></path>
      <path d="M6.1 7.3A17 17 0 0 0 2 12s3.5 6 10 6a10.9 10.9 0 0 0 4.2-.8"></path>
      <path d="M9.9 9.9a3 3 0 0 0 4.2 4.2"></path>
    </symbol>
    <symbol id="i-alert" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <path d="M10.3 3.9 1.9 18a2 2 0 0 0 1.7 3h16.8a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z"></path>
      <path d="M12 9v4M12 17h.01"></path>
    </symbol>
    <symbol id="i-info" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="9"></circle>
      <path d="M12 11v5M12 8h.01"></path>
    </symbol>
    <symbol id="i-check" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="m4.5 12.5 5 5 10-11"></path>
    </symbol>
    <symbol id="i-chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="m9 6 6 6-6 6"></path>
    </symbol>
    <symbol id="i-list" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round">
      <path d="M8 6h13M8 12h13M8 18h13"></path>
      <path d="M3.5 6h.01M3.5 12h.01M3.5 18h.01"></path>
    </symbol>
  </svg>

  <header class="topbar">
    <div class="shell topbar-inner">
      <a class="brand" href="/" aria-label="dnshe-go 控制台首页">
        <span class="brand-mark"><svg class="icon" aria-hidden="true"><use href="#i-logo"></use></svg></span>
        <span class="brand-text">
          <strong>dnshe-go</strong>
          <span>IPv6 DDNS 控制台</span>
        </span>
      </a>
      <div class="topbar-right">
        <span class="sync-pill is-idle" id="sync-pill"><span class="dot"></span><span id="sync-pill-text">连接中</span></span>
        <div class="topbar-actions">
          <a class="icon-btn github-top" href="https://github.com/qrst1ks/dnshe-go" target="_blank" rel="noopener noreferrer" aria-label="GitHub 仓库" title="GitHub">
            <svg class="icon" aria-hidden="true"><use href="#i-github"></use></svg>
          </a>
          <button class="icon-btn" id="theme-toggle" type="button" aria-label="切换主题" title="切换主题">
            <svg class="icon icon-sun" aria-hidden="true"><use href="#i-sun"></use></svg>
            <svg class="icon icon-moon" aria-hidden="true"><use href="#i-moon"></use></svg>
          </button>
          <button class="icon-btn" id="open-logs" type="button" aria-label="打开运行日志" title="运行日志">
            <svg class="icon" aria-hidden="true"><use href="#i-terminal"></use></svg>
          </button>
          <button class="btn optional" id="refresh" type="button" title="刷新状态">
            <svg class="icon" aria-hidden="true"><use href="#i-refresh"></use></svg>
            <span class="btn-label">刷新</span>
          </button>
          <button class="btn optional" id="run" type="button" title="立即同步">
            <svg class="icon" aria-hidden="true"><use href="#i-play"></use></svg>
            <span class="btn-label">立即同步</span>
          </button>
          <button class="btn primary" id="save" type="button" title="保存配置（Ctrl / Cmd + S）">
            <svg class="icon" aria-hidden="true"><use href="#i-save"></use></svg>
            <span class="btn-label">保存配置</span>
            <span class="btn-dot" aria-hidden="true"></span>
          </button>
        </div>
      </div>
    </div>
  </header>

  <main class="shell">
    <div class="alert" id="alert" role="alert" hidden>
      <svg class="icon" aria-hidden="true"><use href="#i-alert"></use></svg>
      <div class="alert-body">
        <strong>配置还不完整</strong>
        <ul id="alert-list"></ul>
      </div>
      <button class="icon-btn small" id="alert-close" type="button" aria-label="关闭提示">
        <svg class="icon" aria-hidden="true"><use href="#i-x"></use></svg>
      </button>
    </div>

    <section class="hero" id="hero" aria-label="运行状态">
      <div class="hero-top">
        <span class="eyebrow"><span class="live-dot" id="live-dot"></span><span id="state">正在连接</span></span>
        <div class="hero-badges" id="provider-summary"></div>
      </div>
      <div class="ip-block">
        <div class="ip-label">当前公网 IPv6</div>
        <div class="ip-row">
          <div class="ip-value" id="ipv6">--</div>
          <button class="copy-btn" id="copy-ip" type="button" title="复制 IPv6 地址" disabled>
            <svg class="icon" aria-hidden="true"><use href="#i-copy"></use></svg>
            <span class="copy-text">复制</span>
            <span class="copy-done">已复制</span>
          </button>
        </div>
        <div class="hero-detail" id="state-detail">正在读取运行状态…</div>
      </div>
      <div class="hero-stats">
        <div class="stat">
          <span class="stat-label">上次同步</span>
          <strong id="last-run">--</strong>
          <span class="stat-sub" id="last-run-abs">尚未同步</span>
        </div>
        <div class="stat">
          <span class="stat-label">同步耗时</span>
          <strong id="duration">--</strong>
          <span class="stat-sub">最近一次</span>
        </div>
        <div class="stat">
          <span class="stat-label">启用域名</span>
          <strong id="domain-count">0</strong>
          <span class="stat-sub" id="domain-summary">未启用服务商</span>
        </div>
        <div class="stat">
          <span class="stat-label">同步间隔</span>
          <strong id="interval-display">--</strong>
          <span class="stat-sub" id="ttl-display">--</span>
        </div>
      </div>
    </section>

    <div class="section-head">
      <div>
        <h2>DNS 服务商</h2>
        <p>两个服务商相互独立，可以单独启用，也可以同时运行</p>
      </div>
      <span class="section-meta" id="provider-count">未启用</span>
    </div>

    <div class="provider-grid">
      <section class="card provider-card dnshe" id="dnshe-card">
        <div class="provider-head">
          <div class="provider-id">
            <span class="provider-logo"><img src="/assets/dnshe.png?v=1" alt="" decoding="async"></span>
            <div>
              <h2>DNSHE</h2>
              <p>更新已有子域名的 AAAA 记录</p>
            </div>
          </div>
          <div class="provider-head-right">
            <span class="status-pill is-off" id="dnshe-state"><span class="dot"></span>已停用</span>
            <label class="switch" title="启用 DNSHE">
              <input id="dnshe-enable" type="checkbox">
              <span class="switch-track"><span class="switch-thumb"></span></span>
              <span class="sr-only">启用 DNSHE</span>
            </label>
          </div>
        </div>
        <div class="provider-body">
          <div class="form-row">
            <div class="field">
              <div class="field-head">
                <label for="dnshe-api-key">API Key</label>
                <span class="mini-tag" data-secret-tag="dnshe-api-key">未设置</span>
              </div>
              <div class="input-wrap">
                <input id="dnshe-api-key" type="password" autocomplete="off" spellcheck="false" placeholder="请输入 DNSHE API Key">
                <button class="reveal" id="toggle-dnshe-key" data-reveal-for="dnshe-api-key" type="button" aria-label="显示完整密钥" title="显示完整密钥">
                  <svg class="icon" aria-hidden="true"><use href="#i-eye"></use></svg>
                </button>
              </div>
            </div>
            <div class="field">
              <div class="field-head">
                <label for="dnshe-api-secret">API Secret</label>
                <span class="mini-tag" data-secret-tag="dnshe-api-secret">未设置</span>
              </div>
              <div class="input-wrap">
                <input id="dnshe-api-secret" type="password" autocomplete="off" spellcheck="false" placeholder="请输入 DNSHE API Secret">
                <button class="reveal" id="toggle-dnshe-secret" data-reveal-for="dnshe-api-secret" type="button" aria-label="显示完整密钥" title="显示完整密钥">
                  <svg class="icon" aria-hidden="true"><use href="#i-eye"></use></svg>
                </button>
              </div>
            </div>
          </div>
          <details class="advanced">
            <summary>
              <svg class="icon" aria-hidden="true"><use href="#i-chevron"></use></svg>
              高级设置 · API 地址
            </summary>
            <div class="field">
              <label for="dnshe-api-base-url">API Base URL</label>
              <input id="dnshe-api-base-url" type="text" autocomplete="off" spellcheck="false" placeholder="https://api005.dnshe.com/index.php?m=domain_hub">
            </div>
          </details>
          <div class="field">
            <div class="field-head">
              <label>域名</label>
              <span class="counter" id="dnshe-domain-count">0 个</span>
            </div>
            <div class="chips" id="dnshe-domains" data-counter="dnshe-domain-count" data-counter-unit="个">
              <input class="chip-entry" type="text" autocomplete="off" spellcheck="false" placeholder="输入域名，回车添加">
            </div>
            <div class="hint">
              <svg class="icon" aria-hidden="true"><use href="#i-info"></use></svg>
              <span>需要先在 DNSHE 创建子域名与 AAAA 记录，本程序只更新已有记录。</span>
            </div>
          </div>
        </div>
      </section>

      <section class="card provider-card cloudflare" id="cloudflare-card">
        <div class="provider-head">
          <div class="provider-id">
            <span class="provider-logo"><img src="/assets/cloudflare.png?v=1" alt="" decoding="async"></span>
            <div>
              <h2>Cloudflare</h2>
              <p>自动更新或创建 AAAA 记录</p>
            </div>
          </div>
          <div class="provider-head-right">
            <span class="status-pill is-off" id="cloudflare-state"><span class="dot"></span>已停用</span>
            <label class="switch" title="启用 Cloudflare">
              <input id="cloudflare-enable" type="checkbox">
              <span class="switch-track"><span class="switch-thumb"></span></span>
              <span class="sr-only">启用 Cloudflare</span>
            </label>
          </div>
        </div>
        <div class="provider-body">
          <div class="field">
            <div class="field-head">
              <label for="cloudflare-api-token">API Token</label>
              <span class="mini-tag" data-secret-tag="cloudflare-api-token">未设置</span>
            </div>
            <div class="input-wrap">
              <input id="cloudflare-api-token" type="password" autocomplete="off" spellcheck="false" placeholder="请输入 Cloudflare API Token">
              <button class="reveal" id="toggle-cloudflare-token" data-reveal-for="cloudflare-api-token" type="button" aria-label="显示完整密钥" title="显示完整密钥">
                <svg class="icon" aria-hidden="true"><use href="#i-eye"></use></svg>
              </button>
            </div>
          </div>
          <details class="advanced">
            <summary>
              <svg class="icon" aria-hidden="true"><use href="#i-chevron"></use></svg>
              高级设置 · API 地址
            </summary>
            <div class="field">
              <label for="cloudflare-api-base-url">API Base URL</label>
              <input id="cloudflare-api-base-url" type="text" autocomplete="off" spellcheck="false" placeholder="https://api.cloudflare.com/client/v4">
            </div>
          </details>
          <div class="field">
            <div class="field-head">
              <label>域名</label>
              <span class="counter" id="cloudflare-domain-count">0 个</span>
            </div>
            <div class="chips" id="cloudflare-domains" data-counter="cloudflare-domain-count" data-counter-unit="个">
              <input class="chip-entry" type="text" autocomplete="off" spellcheck="false" placeholder="输入域名，回车添加">
            </div>
            <div class="hint">
              <svg class="icon" aria-hidden="true"><use href="#i-info"></use></svg>
              <span>Token 需要 Zone:Read 与 DNS:Edit 权限；记录不存在时会自动创建。</span>
            </div>
          </div>
        </div>
      </section>
    </div>

    <div class="lower-grid">
      <section class="card">
        <div class="panel-head">
          <div>
            <h2>IPv6 地址来源</h2>
            <p>每轮同步会读取一次公网 IPv6，再更新所有启用域名</p>
          </div>
        </div>
        <div class="panel-body">
          <div class="segmented" role="tablist" aria-label="IPv6 来源">
            <button class="seg-btn active" type="button" role="tab" aria-selected="true" data-source="url">
              <svg class="icon" aria-hidden="true"><use href="#i-link"></use></svg>URL 接口
            </button>
            <button class="seg-btn" type="button" role="tab" aria-selected="false" data-source="interface">
              <svg class="icon" aria-hidden="true"><use href="#i-nic"></use></svg>本机网卡
            </button>
            <button class="seg-btn" type="button" role="tab" aria-selected="false" data-source="cmd">
              <svg class="icon" aria-hidden="true"><use href="#i-terminal"></use></svg>自定义命令
            </button>
          </div>
          <input id="ipv6-source" type="hidden" value="url">

          <div class="source-panel" data-panel="url">
            <div class="field">
              <div class="field-head">
                <label>查询接口</label>
                <span class="counter" id="url-count">0 条</span>
              </div>
              <div class="chips" id="ipv6-urls" data-counter="url-count" data-counter-unit="条">
                <input class="chip-entry" type="text" autocomplete="off" spellcheck="false" placeholder="粘贴接口地址，回车添加">
              </div>
              <div class="hint">
                <svg class="icon" aria-hidden="true"><use href="#i-info"></use></svg>
                <span>按顺序尝试，返回内容中需要包含公网 IPv6 地址。</span>
              </div>
            </div>
          </div>

          <div class="source-panel" data-panel="interface" hidden>
            <div class="field">
              <label for="ipv6-interface">选择网卡</label>
              <div class="interface-row">
                <div class="select-wrap">
                  <select id="ipv6-interface">
                    <option value="">自动选择有公网 IPv6 的网卡</option>
                  </select>
                  <svg class="icon" aria-hidden="true"><use href="#i-chevron"></use></svg>
                </div>
                <button class="icon-btn" id="refresh-interfaces" type="button" aria-label="重新检测网卡" title="重新检测网卡">
                  <svg class="icon" aria-hidden="true"><use href="#i-refresh"></use></svg>
                </button>
              </div>
              <div class="hint">
                <svg class="icon" aria-hidden="true"><use href="#i-info"></use></svg>
                <span>自动模式会选择第一个可用的全局 IPv6 地址。</span>
              </div>
            </div>
          </div>

          <div class="source-panel" data-panel="cmd" hidden>
            <div class="field">
              <label for="ipv6-command">执行命令</label>
              <input class="mono-input" id="ipv6-command" type="text" autocomplete="off" spellcheck="false" placeholder="ip -6 addr show scope global">
              <div class="hint">
                <svg class="icon" aria-hidden="true"><use href="#i-info"></use></svg>
                <span>程序会执行命令，并从输出中提取第一个公网 IPv6 地址。</span>
              </div>
            </div>
          </div>

          <div class="source-current">
            <span class="source-current-icon"><svg class="icon" aria-hidden="true"><use href="#i-globe"></use></svg></span>
            <div class="source-current-copy">
              <span>当前解析地址</span>
              <strong id="source-current-ip">--</strong>
            </div>
            <span class="source-current-mode" id="source-current-mode">URL 接口</span>
          </div>
          <div class="source-facts">
            <span><span class="dot"></span>每轮同步重新读取</span>
            <span><span class="dot"></span>公网 IPv6 自动校验</span>
          </div>
        </div>
      </section>

      <section class="card">
        <div class="panel-head">
          <div>
            <h2>同步设置</h2>
            <p>控制检查频率与记录 TTL</p>
          </div>
        </div>
        <div class="panel-body">
          <div class="setting-block">
            <div class="field">
              <div class="field-head">
                <label for="interval">同步间隔</label>
                <span class="counter">建议不低于 1 分钟</span>
              </div>
              <div class="input-suffix">
                <input id="interval" type="number" min="10" step="1" inputmode="numeric">
                <span>秒</span>
              </div>
              <div class="preset-row" aria-label="同步间隔快捷选项">
                <button class="preset" type="button" data-target="interval" data-value="60">1 分钟</button>
                <button class="preset" type="button" data-target="interval" data-value="300">5 分钟</button>
                <button class="preset" type="button" data-target="interval" data-value="900">15 分钟</button>
                <button class="preset" type="button" data-target="interval" data-value="3600">1 小时</button>
              </div>
            </div>
          </div>
          <div class="setting-block">
            <div class="field">
              <div class="field-head">
                <label for="ttl">记录 TTL</label>
                <span class="counter">单位：秒</span>
              </div>
              <div class="input-suffix">
                <input id="ttl" type="number" min="1" step="1" inputmode="numeric">
                <span>秒</span>
              </div>
              <div class="preset-row" aria-label="TTL 快捷选项">
                <button class="preset" type="button" data-target="ttl" data-value="600">10 分钟</button>
                <button class="preset" type="button" data-target="ttl" data-value="1800">30 分钟</button>
                <button class="preset" type="button" data-target="ttl" data-value="3600">1 小时</button>
              </div>
            </div>
          </div>
          <div class="sync-note">
            <svg class="icon" aria-hidden="true"><use href="#i-info"></use></svg>
            <span>保存配置后会立即触发一次同步；修改未保存时，“立即同步”仍使用已保存的配置。</span>
          </div>
        </div>
      </section>
    </div>

    <section class="card results-card">
      <div class="results-head">
        <div>
          <h2>最近同步结果</h2>
          <p>每个域名的最新处理状态</p>
        </div>
        <div class="filters" role="tablist" aria-label="筛选同步结果">
          <button class="filter active" type="button" role="tab" aria-selected="true" data-filter="all">全部 <span id="count-all">0</span></button>
          <button class="filter" type="button" role="tab" aria-selected="false" data-filter="success">成功 <span id="count-success">0</span></button>
          <button class="filter" type="button" role="tab" aria-selected="false" data-filter="failed">失败 <span id="count-failed">0</span></button>
        </div>
      </div>
      <div class="results-list" id="results"></div>
    </section>

    <footer class="footer">
      <div class="footer-left">
        <svg class="icon" aria-hidden="true"><use href="#i-logo"></use></svg>
        <span>dnshe-go · IPv6 AAAA DDNS</span>
        <span class="config-path" id="config-path" title="配置文件路径">--</span>
      </div>
      <div class="footer-right">
        <a class="footer-link" href="https://github.com/qrst1ks/dnshe-go" target="_blank" rel="noopener noreferrer">
          <svg class="icon" aria-hidden="true"><use href="#i-github"></use></svg>GitHub
        </a>
      </div>
    </footer>
  </main>

  <div class="drawer-backdrop" id="logs-backdrop"></div>
  <aside class="drawer" id="logs-drawer" aria-hidden="true" aria-label="运行日志">
    <div class="drawer-head">
      <div>
        <h2>运行日志</h2>
        <p><span id="log-count">0</span> 条记录 · 实时刷新</p>
      </div>
      <div class="drawer-actions">
        <button class="btn" id="clear-logs" type="button">清空</button>
        <button class="icon-btn" id="close-logs" type="button" aria-label="关闭日志">
          <svg class="icon" aria-hidden="true"><use href="#i-x"></use></svg>
        </button>
      </div>
    </div>
    <div class="term" id="logs"></div>
  </aside>

  <div class="toasts" id="toasts" aria-live="polite" aria-atomic="false"></div>

  <script>
    (function () {
      "use strict";

      var $ = function (id) { return document.getElementById(id); };
      var state = {
        loaded: false,
        dirty: false,
        filter: "all",
        results: [],
        logs: [],
        source: "url",
        offline: false,
        polling: false,
        saving: false,
        interfacesLoaded: false,
        activeProvider: false
      };

      function esc(value) {
        return String(value == null ? "" : value)
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")
          .replace(/"/g, "&quot;")
          .replace(/'/g, "&#39;");
      }

      function icon(name) {
        return '<svg class="icon" aria-hidden="true"><use href="#i-' + name + '"></use></svg>';
      }

      function cleanList(values) {
        var seen = {};
        var out = [];
        (values || []).forEach(function (value) {
          var item = String(value == null ? "" : value).trim();
          if (!item) return;
          var key = item.toLowerCase();
          if (seen[key]) return;
          seen[key] = true;
          out.push(item);
        });
        return out;
      }

      function humanDuration(raw) {
        var sec = Number(raw) || 0;
        if (sec <= 0) return "--";
        if (sec < 60) return sec + " 秒";
        if (sec < 3600) return Math.round(sec / 60) + " 分钟";
        var hours = Math.floor(sec / 3600);
        var minutes = Math.round((sec % 3600) / 60);
        return minutes === 0 ? hours + " 小时" : hours + " 小时 " + minutes + " 分";
      }

      function timeAgo(value) {
        if (!value) return "--";
        var time = new Date(value).getTime();
        if (!isFinite(time)) return "--";
        var diff = Math.max(0, (Date.now() - time) / 1000);
        if (diff < 10) return "刚刚";
        if (diff < 60) return Math.floor(diff) + " 秒前";
        if (diff < 3600) return Math.floor(diff / 60) + " 分钟前";
        if (diff < 86400) return Math.floor(diff / 3600) + " 小时前";
        if (diff < 604800) return Math.floor(diff / 86400) + " 天前";
        return new Date(time).toLocaleDateString();
      }

      function formatAbsolute(value) {
        if (!value) return "--";
        var time = new Date(value);
        return isNaN(time.getTime()) ? String(value) : time.toLocaleString();
      }

      function toast(type, message) {
        var box = $("toasts");
        var el = document.createElement("div");
        el.className = "toast " + type;
        var iconName = type === "success" ? "check" : type === "error" ? "alert" : "info";
        el.innerHTML = '<span class="toast-icon">' + icon(iconName) + '</span><span class="toast-msg"></span><button class="toast-close" type="button" aria-label="关闭通知">' + icon("x") + '</button>';
        el.querySelector(".toast-msg").textContent = message;
        function dismiss() {
          clearTimeout(timer);
          if (!el.parentNode) return;
          el.classList.add("leaving");
          setTimeout(function () { el.remove(); }, 180);
        }
        var timer = setTimeout(dismiss, 4200);
        el.querySelector(".toast-close").addEventListener("click", dismiss);
        box.appendChild(el);
      }

      function showAlert(problems) {
        var box = $("alert");
        var list = $("alert-list");
        if (!Array.isArray(problems) || problems.length === 0) {
          box.hidden = true;
          list.innerHTML = "";
          return;
        }
        list.innerHTML = problems.map(function (item) { return "<li>" + esc(item) + "</li>"; }).join("");
        box.hidden = false;
        box.scrollIntoView({ behavior: "smooth", block: "nearest" });
      }

      function hideAlert() {
        showAlert([]);
      }

      function chipValues(id) {
        return Array.prototype.map.call($(id).querySelectorAll(".chip"), function (chip) {
          return chip.getAttribute("data-value") || "";
        });
      }

      function makeChip(value) {
        var chip = document.createElement("span");
        chip.className = "chip";
        chip.setAttribute("data-value", value);
        chip.title = value;
        var label = document.createElement("span");
        label.className = "chip-text";
        label.textContent = value;
        var remove = document.createElement("button");
        remove.type = "button";
        remove.className = "chip-remove";
        remove.setAttribute("aria-label", "移除 " + value);
        remove.innerHTML = icon("x");
        chip.appendChild(label);
        chip.appendChild(remove);
        return chip;
      }

      function updateChipCount(id) {
        var box = $(id);
        var counterId = box.getAttribute("data-counter");
        if (!counterId) return;
        var unit = box.getAttribute("data-counter-unit") || "个";
        $(counterId).textContent = box.querySelectorAll(".chip").length + " " + unit;
      }

      function setChips(id, values) {
        var box = $(id);
        var entry = box.querySelector(".chip-entry");
        Array.prototype.forEach.call(box.querySelectorAll(".chip"), function (chip) { chip.remove(); });
        cleanList(values).forEach(function (value) {
          box.insertBefore(makeChip(value), entry);
        });
        updateChipCount(id);
      }

      function addChip(id, raw) {
        var value = String(raw == null ? "" : raw).trim().replace(/,$/, "").trim();
        if (!value) return false;
        var lower = value.toLowerCase();
        var exists = chipValues(id).some(function (item) { return item.toLowerCase() === lower; });
        if (exists) return false;
        var box = $(id);
        var entry = box.querySelector(".chip-entry");
        box.insertBefore(makeChip(value), entry);
        updateChipCount(id);
        markDirty();
        return true;
      }

      function setupChips(id) {
        var box = $(id);
        var entry = box.querySelector(".chip-entry");

        function commit() {
          var parts = entry.value.split(/[\s,]+/);
          var added = false;
          parts.forEach(function (part) { if (addChip(id, part)) added = true; });
          if (added) entry.value = "";
        }

        box.addEventListener("click", function (event) {
          var remove = event.target.closest(".chip-remove");
          if (remove) {
            var chip = remove.closest(".chip");
            if (chip) chip.remove();
            updateChipCount(id);
            markDirty();
            entry.focus();
            return;
          }
          if (event.target === box) entry.focus();
        });

        entry.addEventListener("keydown", function (event) {
          if (event.key === "Enter" || event.key === ",") {
            event.preventDefault();
            commit();
            return;
          }
          if (event.key === "Backspace" && entry.value === "") {
            var chips = box.querySelectorAll(".chip");
            if (chips.length === 0) return;
            chips[chips.length - 1].remove();
            updateChipCount(id);
            markDirty();
          }
        });

        entry.addEventListener("paste", function (event) {
          var clipboard = event.clipboardData || window.clipboardData;
          var text = clipboard ? clipboard.getData("text") : "";
          if (!text || !/[\s,]/.test(text)) return;
          event.preventDefault();
          text.split(/[\s,]+/).forEach(function (part) { addChip(id, part); });
        });

        entry.addEventListener("blur", commit);
      }

      function hasSecret(input) {
        return input.value.trim() !== "" || input.getAttribute("data-saved") === "true";
      }

      function providerInfo() {
        return {
          dnshe: {
            enabled: $("dnshe-enable").checked,
            configured: hasSecret($("dnshe-api-key")) && hasSecret($("dnshe-api-secret")),
            domains: chipValues("dnshe-domains").length
          },
          cloudflare: {
            enabled: $("cloudflare-enable").checked,
            configured: hasSecret($("cloudflare-api-token")),
            domains: chipValues("cloudflare-domains").length
          }
        };
      }

      function paintProviderCard(prefix, info) {
        var card = $(prefix + "-card");
        var pill = $(prefix + "-state");
        card.classList.toggle("is-off", !info.enabled);
        var label;
        var kind;
        if (!info.enabled) {
          label = "已停用";
          kind = "is-off";
        } else if (!info.configured) {
          label = "待填写凭据";
          kind = "is-warn";
        } else if (info.domains === 0) {
          label = "待添加域名";
          kind = "is-warn";
        } else {
          label = "已启用";
          kind = "is-on";
        }
        pill.className = "status-pill " + kind;
        pill.innerHTML = '<span class="dot"></span>' + label;
      }

      function highlightPresets() {
        Array.prototype.forEach.call(document.querySelectorAll(".preset"), function (button) {
          var target = button.getAttribute("data-target");
          var value = Number(button.getAttribute("data-value"));
          button.classList.toggle("active", Number($(target).value) === value);
        });
      }

      function updateProviderUI() {
        var info = providerInfo();
        paintProviderCard("dnshe", info.dnshe);
        paintProviderCard("cloudflare", info.cloudflare);

        var enabled = [];
        var domains = 0;
        if (info.dnshe.enabled) {
          enabled.push({ name: "DNSHE", icon: "globe" });
          domains += info.dnshe.domains;
        }
        if (info.cloudflare.enabled) {
          enabled.push({ name: "Cloudflare", icon: "cloud" });
          domains += info.cloudflare.domains;
        }

        $("provider-summary").innerHTML = enabled.length === 0
          ? '<span class="hero-provider is-empty">未启用服务商</span>'
          : enabled.map(function (item) {
              return '<span class="hero-provider">' + icon(item.icon) + esc(item.name) + '</span>';
            }).join("");

        $("provider-count").textContent = enabled.length ? enabled.length + " 个已启用" : "未启用";
        $("domain-count").textContent = String(domains);
        $("domain-summary").textContent = enabled.length
          ? enabled.map(function (item) { return item.name; }).join(" · ")
          : "未启用服务商";
        $("interval-display").textContent = humanDuration(Number($("interval").value));
        $("ttl-display").textContent = "TTL " + (Number($("ttl").value) || 0) + " 秒";
        highlightPresets();
      }

      var SOURCE_LABELS = { url: "URL 接口", interface: "本机网卡", cmd: "自定义命令" };

      function setSourceMode(source, silent) {
        if (source !== "url" && source !== "interface" && source !== "cmd") source = "url";
        state.source = source;
        $("ipv6-source").value = source;
        $("source-current-mode").textContent = SOURCE_LABELS[source];
        Array.prototype.forEach.call(document.querySelectorAll(".seg-btn[data-source]"), function (button) {
          var active = button.getAttribute("data-source") === source;
          button.classList.toggle("active", active);
          button.setAttribute("aria-selected", active ? "true" : "false");
        });
        Array.prototype.forEach.call(document.querySelectorAll("[data-panel]"), function (panel) {
          panel.hidden = panel.getAttribute("data-panel") !== source;
        });
        if (!silent) markDirty();
      }

      function selectInterface(name) {
        var select = $("ipv6-interface");
        var value = name || "";
        var exists = Array.prototype.some.call(select.options, function (option) { return option.value === value; });
        if (value && !exists) {
          var option = document.createElement("option");
          option.value = value;
          option.textContent = value + "（当前配置）";
          select.appendChild(option);
        }
        select.value = value;
      }

      async function loadInterfaces(selected, force) {
        if (state.interfacesLoaded && !force) {
          if (selected != null) selectInterface(selected);
          return;
        }
        var select = $("ipv6-interface");
        var current = selected != null ? selected : select.value;
        var button = $("refresh-interfaces");
        button.classList.add("spinning");
        try {
          var res = await fetch("/api/interfaces", { cache: "no-store" });
          var data = await res.json();
          var items = Array.isArray(data.interfaces) ? data.interfaces : [];
          select.innerHTML = "";
          var auto = document.createElement("option");
          auto.value = "";
          auto.textContent = "自动选择有公网 IPv6 的网卡";
          select.appendChild(auto);
          items.forEach(function (item) {
            var option = document.createElement("option");
            option.value = item.name || "";
            option.textContent = item.label || item.name || "";
            select.appendChild(option);
          });
          selectInterface(current || "");
          state.interfacesLoaded = true;
        } catch (err) {
          if (select.options.length === 0) {
            var fallback = document.createElement("option");
            fallback.value = current || "";
            fallback.textContent = current ? current + "（当前配置）" : "无法读取网卡，自动选择";
            select.appendChild(fallback);
            select.value = current || "";
          }
          if (force) toast("error", "读取网卡列表失败");
        } finally {
          button.classList.remove("spinning");
        }
      }

      function resetReveal(input) {
        input.type = "password";
        input.removeAttribute("data-revealed");
        var button = document.querySelector('[data-reveal-for="' + input.id + '"]');
        if (!button) return;
        button.innerHTML = icon("eye");
        button.setAttribute("aria-label", "显示完整密钥");
        button.title = "显示完整密钥";
        button.disabled = false;
      }

      function setSecretField(input, masked, placeholder) {
        resetReveal(input);
        input.value = "";
        input.setAttribute("data-saved", masked ? "true" : "false");
        input.placeholder = masked || placeholder;
        input.title = masked ? "已保存，留空保持不变" : "";
        var tag = document.querySelector('[data-secret-tag="' + input.id + '"]');
        if (!tag) return;
        tag.textContent = masked ? "已保存" : "未设置";
        tag.title = masked ? "已保存，留空保持不变" : "";
        tag.className = "mini-tag" + (masked ? " ok" : "");
      }

      function fillConfig(cfg) {
        cfg = cfg || {};
        $("interval").value = cfg.interval_seconds || 300;
        $("ttl").value = cfg.ttl || 600;

        var dnshe = cfg.dnshe || {};
        $("dnshe-enable").checked = dnshe.enable === true;
        setSecretField($("dnshe-api-key"), dnshe.api_key_masked, "请输入 DNSHE API Key");
        setSecretField($("dnshe-api-secret"), dnshe.api_secret_masked, "请输入 DNSHE API Secret");
        $("dnshe-api-base-url").value = dnshe.api_base_url || "";
        setChips("dnshe-domains", dnshe.domains || []);

        var cloudflare = cfg.cloudflare || {};
        $("cloudflare-enable").checked = cloudflare.enable === true;
        setSecretField($("cloudflare-api-token"), cloudflare.api_token_masked, "请输入 Cloudflare API Token");
        $("cloudflare-api-base-url").value = cloudflare.api_base_url || "";
        setChips("cloudflare-domains", cloudflare.domains || []);

        var dnsheReady = dnshe.enable === true && !!dnshe.api_key_masked && (dnshe.domains || []).length > 0;
        var cloudflareReady = cloudflare.enable === true && !!cloudflare.api_token_masked && (cloudflare.domains || []).length > 0;
        state.activeProvider = dnsheReady || cloudflareReady;

        var ipv6 = cfg.ipv6 || {};
        setSourceMode(ipv6.source || "url", true);
        loadInterfaces(ipv6.interface || "");
        $("ipv6-command").value = ipv6.command || "";
        setChips("ipv6-urls", ipv6.urls || []);

        updateProviderUI();
        hideAlert();
        $("save").classList.remove("is-dirty");
      }

      function markDirty() {
        state.dirty = true;
        $("save").classList.add("is-dirty");
        updateProviderUI();
        hideAlert();
      }

      function readConfig() {
        return {
          interval_seconds: Number($("interval").value || 300),
          ttl: Number($("ttl").value || 600),
          dnshe: {
            enable: $("dnshe-enable").checked,
            api_key: $("dnshe-api-key").value.trim(),
            api_secret: $("dnshe-api-secret").value.trim(),
            api_base_url: $("dnshe-api-base-url").value.trim(),
            domains: chipValues("dnshe-domains")
          },
          cloudflare: {
            enable: $("cloudflare-enable").checked,
            api_token: $("cloudflare-api-token").value.trim(),
            api_base_url: $("cloudflare-api-base-url").value.trim(),
            domains: chipValues("cloudflare-domains")
          },
          ipv6: {
            source: $("ipv6-source").value,
            urls: chipValues("ipv6-urls"),
            interface: $("ipv6-interface").value,
            command: $("ipv6-command").value.trim()
          }
        };
      }

      function validateBeforeSave(cfg) {
        var problems = [];
        if (!cfg.dnshe.enable && !cfg.cloudflare.enable) problems.push("请至少启用一个 DNS 服务商");
        if (cfg.dnshe.enable) {
          if (!hasSecret($("dnshe-api-key"))) problems.push("请填写 DNSHE API Key");
          if (!hasSecret($("dnshe-api-secret"))) problems.push("请填写 DNSHE API Secret");
          if (cfg.dnshe.domains.length === 0) problems.push("请为 DNSHE 添加至少一个域名");
        }
        if (cfg.cloudflare.enable) {
          if (!hasSecret($("cloudflare-api-token"))) problems.push("请填写 Cloudflare API Token");
          if (cfg.cloudflare.domains.length === 0) problems.push("请为 Cloudflare 添加至少一个域名");
        }
        if (cfg.ipv6.source === "url" && cfg.ipv6.urls.length === 0) problems.push("URL 来源至少需要一个接口地址");
        if (cfg.ipv6.source === "cmd" && !cfg.ipv6.command) problems.push("命令来源需要填写执行命令");
        return problems;
      }

      function setConnection(ok) {
        if (ok) {
          if (state.offline) toast("success", "已重新连接服务");
          state.offline = false;
          return;
        }
        if (!state.offline) toast("error", "无法连接服务，请确认 dnshe-go 正在运行");
        state.offline = true;
        $("sync-pill").className = "sync-pill is-bad";
        $("sync-pill-text").textContent = "连接断开";
        $("state").textContent = "连接断开";
        $("state-detail").textContent = "无法获取运行状态，正在自动重试…";
        document.getElementById("hero").classList.add("is-bad");
      }

      function setHeroState(kind) {
        var hero = $("hero");
        hero.classList.toggle("is-bad", kind === "bad");
        hero.classList.toggle("is-running", kind === "running");
        hero.classList.toggle("is-idle", kind === "idle");
      }

      function applySync(sync) {
        sync = sync || {};
        var running = sync.running === true;
        var error = String(sync.last_error || "").trim();
        var done = !!sync.last_run_ended_at;
        var kind = "idle";
        var stateText = "等待首次同步";
        var pillText = "等待同步";
        var pillClass = "sync-pill is-idle";
        var detail = "保存配置后会自动开始第一次同步";
        var active = state.activeProvider;
        var showRun = done && active;

        if (running) {
          kind = "running";
          stateText = "同步中";
          pillText = "正在同步";
          pillClass = "sync-pill is-running";
          detail = "正在读取 IPv6 并更新已启用的服务商…";
        } else if (!active) {
          stateText = "待配置";
          pillText = "待配置";
          detail = "填写服务商凭据和域名并保存，即可开始首次同步";
        } else if (error) {
          kind = "bad";
          stateText = "需要关注";
          pillText = "同步异常";
          pillClass = "sync-pill is-bad";
          detail = error;
        } else if (done) {
          kind = "ok";
          stateText = "运行正常";
          pillText = "运行正常";
          pillClass = "sync-pill is-on";
          detail = "所有已启用服务商的记录均已检查完成";
        }

        setHeroState(kind);
        $("state").textContent = stateText;
        $("state-detail").textContent = detail;
        $("state-detail").title = detail;
        $("sync-pill").className = pillClass;
        $("sync-pill-text").textContent = pillText;

        var ip = sync.current_ipv6 || "";
        $("ipv6").textContent = ip || "--";
        $("source-current-ip").textContent = ip || "--";
        $("copy-ip").disabled = !ip;
        $("copy-ip").setAttribute("data-ip", ip);
        $("last-run").textContent = showRun ? timeAgo(sync.last_run_ended_at) : "--";
        $("last-run-abs").textContent = showRun ? formatAbsolute(sync.last_run_ended_at) : "尚未同步";
        $("last-run-abs").title = showRun ? formatAbsolute(sync.last_run_ended_at) : "";
        $("duration").textContent = showRun && typeof sync.last_duration_ms === "number" && sync.last_duration_ms >= 0 ? sync.last_duration_ms + " ms" : "--";

        state.results = Array.isArray(sync.results) ? sync.results : [];
        renderResults();
      }

      var RESULT_TEXT = { unchanged: "未变化", updated: "已更新", created: "已创建", failed: "失败" };

      function providerKey(name) {
        var value = String(name || "").toLowerCase();
        if (value === "cloudflare") return "cloudflare";
        if (value === "dnshe") return "dnshe";
        return "other";
      }

      function statusTag(status) {
        var key = String(status || "unknown");
        return '<span class="tag ' + esc(key) + '">' + esc(RESULT_TEXT[key] || key) + '</span>';
      }

      function renderResults() {
        var results = state.results || [];
        var success = results.filter(function (item) { return item.status !== "failed"; }).length;
        var failed = results.filter(function (item) { return item.status === "failed"; }).length;
        $("count-all").textContent = String(results.length);
        $("count-success").textContent = String(success);
        $("count-failed").textContent = String(failed);

        Array.prototype.forEach.call(document.querySelectorAll(".filter"), function (button) {
          var active = button.getAttribute("data-filter") === state.filter;
          button.classList.toggle("active", active);
          button.setAttribute("aria-selected", active ? "true" : "false");
        });

        var filtered = results.filter(function (item) {
          if (state.filter === "success") return item.status !== "failed";
          if (state.filter === "failed") return item.status === "failed";
          return true;
        });

        var box = $("results");
        if (filtered.length === 0) {
          var title = "还没有同步结果";
          var subtitle = "保存配置或点击“立即同步”后，这里会显示每个域名的处理结果。";
          if (results.length > 0 && state.filter === "failed") {
            title = "没有失败的记录";
            subtitle = "所有已启用服务商目前均同步成功。";
          } else if (results.length > 0 && state.filter === "success") {
            title = "没有成功的记录";
            subtitle = "切换到“全部”查看失败详情。";
          }
          box.innerHTML = '<div class="empty-state">' + icon("list") + '<strong>' + title + '</strong><p>' + subtitle + '</p></div>';
          return;
        }

        box.innerHTML = filtered.map(function (item) {
          var status = item.status || "unknown";
          var provider = item.provider || "未知服务商";
          var key = providerKey(provider);
          var ip = item.ip || "";
          var oldIP = item.old_ip || "";
          var meta;
          if (status === "failed") {
            meta = '<div class="result-error">' + esc(item.message || "同步失败，请查看运行日志") + '</div>';
          } else if (oldIP && oldIP !== ip) {
            meta = '<div class="result-meta"><span class="mono">' + esc(oldIP) + '</span><span class="arrow">→</span><span class="mono">' + esc(ip) + '</span></div>';
          } else {
            meta = '<div class="result-meta"><span class="mono">' + esc(ip || "--") + '</span></div>';
          }
          var logo = key === "cloudflare" ? "cloud" : key === "dnshe" ? "globe" : "dns";
          return '<div class="result-row">' +
            '<span class="result-logo ' + key + '">' + icon(logo) + '</span>' +
            '<div class="result-main">' +
              '<div class="result-title">' +
                '<span class="result-domain">' + esc(item.domain || "--") + '</span>' +
                statusTag(status) +
                '<span class="provider-name">' + esc(provider) + '</span>' +
              '</div>' +
              meta +
            '</div>' +
            '<div class="result-side">' +
              '<span class="record-badge">' + esc(item.record_type || "AAAA") + '</span>' +
              '<span title="' + esc(formatAbsolute(item.time)) + '">' + esc(timeAgo(item.time)) + '</span>' +
            '</div>' +
          '</div>';
        }).join("");
      }

      function renderLogs(entries) {
        state.logs = Array.isArray(entries) ? entries : [];
        $("log-count").textContent = String(state.logs.length);
        var box = $("logs");
        if (state.logs.length === 0) {
          box.innerHTML = '<div class="logs-empty">' + icon("terminal") + '<span>暂无日志记录</span></div>';
          return;
        }
        var nearBottom = box.scrollHeight - box.scrollTop - box.clientHeight < 90;
        box.innerHTML = state.logs.map(function (entry) {
          var level = String(entry.level || "INFO").toUpperCase();
          return '<div class="log-line">' +
            '<span class="log-time">' + esc(entry.time) + '</span>' +
            '<span class="log-level ' + esc(level) + '">' + esc(level) + '</span>' +
            '<span class="log-msg">' + esc(entry.message) + '</span>' +
          '</div>';
        }).join("");
        if (nearBottom) box.scrollTop = box.scrollHeight;
      }

      async function refresh(options) {
        options = options || {};
        if (state.polling) return;
        if (document.hidden && !options.populate) return;
        state.polling = true;
        try {
          var res = await fetch("/api/status", { cache: "no-store" });
          if (!res.ok) throw new Error("HTTP " + res.status);
          var data = await res.json();
          if (data.config && (options.populate || !state.loaded || !state.dirty)) {
            fillConfig(data.config);
            state.loaded = true;
            state.dirty = false;
          }
          applySync(data.sync || {});
          renderLogs(data.logs || []);
          if (data.config_path) {
            $("config-path").textContent = data.config_path;
            $("config-path").title = data.config_path;
          }
          setConnection(true);
        } catch (err) {
          setConnection(false);
        } finally {
          state.polling = false;
        }
      }

      async function save() {
        if (state.saving) return;
        var cfg = readConfig();
        var problems = validateBeforeSave(cfg);
        showAlert(problems);
        if (problems.length > 0) {
          toast("error", "配置还不完整，请检查 " + problems.length + " 项问题");
          return;
        }
        state.saving = true;
        var button = $("save");
        button.disabled = true;
        button.classList.add("loading");
        try {
          var res = await fetch("/api/config", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(cfg)
          });
          var data = await res.json().catch(function () { return {}; });
          if (!res.ok || !data.ok) {
            if (Array.isArray(data.problems) && data.problems.length > 0) showAlert(data.problems);
            throw new Error(data.msg || "保存失败");
          }
          $("dnshe-api-key").value = "";
          $("dnshe-api-secret").value = "";
          $("cloudflare-api-token").value = "";
          state.dirty = false;
          button.classList.remove("is-dirty");
          hideAlert();
          toast("success", "配置已保存，正在同步");
          await refresh({ populate: true });
        } catch (err) {
          toast("error", "保存失败：" + (err && err.message ? err.message : String(err)));
        } finally {
          state.saving = false;
          button.disabled = false;
          button.classList.remove("loading");
        }
      }

      async function runNow() {
        var button = $("run");
        if (button.disabled) return;
        button.disabled = true;
        button.classList.add("spinning");
        try {
          var res = await fetch("/api/run", { method: "POST" });
          if (!res.ok) throw new Error("HTTP " + res.status);
          toast("info", "已触发同步，结果会自动刷新");
          setTimeout(function () { refresh(); }, 700);
          setTimeout(function () { refresh(); }, 1800);
          setTimeout(function () { refresh(); }, 3600);
        } catch (err) {
          toast("error", "触发同步失败");
        } finally {
          setTimeout(function () {
            button.disabled = false;
            button.classList.remove("spinning");
          }, 1200);
        }
      }

      async function clearLogs() {
        try {
          var res = await fetch("/api/logs/clear", { method: "POST" });
          if (!res.ok) throw new Error("HTTP " + res.status);
          state.logs = [];
          renderLogs([]);
          toast("success", "日志已清空");
        } catch (err) {
          toast("error", "清空日志失败");
        }
      }

      function openDrawer() {
        $("logs-drawer").classList.add("open");
        $("logs-backdrop").classList.add("open");
        $("logs-drawer").setAttribute("aria-hidden", "false");
        document.body.classList.add("drawer-open");
        setTimeout(function () { $("close-logs").focus(); }, 90);
        var box = $("logs");
        box.scrollTop = box.scrollHeight;
        refresh();
      }

      function closeDrawer() {
        $("logs-drawer").classList.remove("open");
        $("logs-backdrop").classList.remove("open");
        $("logs-drawer").setAttribute("aria-hidden", "true");
        document.body.classList.remove("drawer-open");
        $("open-logs").focus();
      }

      async function copyIP() {
        var ip = $("copy-ip").getAttribute("data-ip") || "";
        if (!ip) {
          toast("info", "当前还没有可复制的 IPv6 地址");
          return;
        }
        try {
          if (navigator.clipboard && navigator.clipboard.writeText) {
            await navigator.clipboard.writeText(ip);
          } else {
            var area = document.createElement("textarea");
            area.value = ip;
            area.style.position = "fixed";
            area.style.opacity = "0";
            document.body.appendChild(area);
            area.select();
            document.execCommand("copy");
            area.remove();
          }
          var button = $("copy-ip");
          button.classList.add("copied");
          setTimeout(function () { button.classList.remove("copied"); }, 1400);
          toast("success", "IPv6 地址已复制");
        } catch (err) {
          toast("error", "复制失败，请手动选择地址");
        }
      }

      function setupTheme() {
        var button = $("theme-toggle");
        function apply(theme) {
          document.documentElement.setAttribute("data-theme", theme);
          try { localStorage.setItem("dnshe-go-theme", theme); } catch (err) {}
          var label = theme === "dark" ? "切换到浅色主题" : "切换到深色主题";
          button.setAttribute("aria-label", label);
          button.title = label;
        }
        apply(document.documentElement.getAttribute("data-theme") || "dark");
        button.addEventListener("click", function () {
          apply(document.documentElement.getAttribute("data-theme") === "dark" ? "light" : "dark");
        });
      }

      function setupReveal(buttonId, inputId) {
        var button = $(buttonId);
        var input = $(inputId);
        input.addEventListener("input", function () {
          input.removeAttribute("data-revealed");
        });
        button.addEventListener("click", async function () {
          if (input.type === "text") {
            input.type = "password";
            if (input.getAttribute("data-revealed") === "true") {
              input.value = "";
              input.removeAttribute("data-revealed");
            }
            button.innerHTML = icon("eye");
            button.setAttribute("aria-label", "显示完整密钥");
            button.title = "显示完整密钥";
            return;
          }
          if (!input.value && input.getAttribute("data-saved") !== "true") return;
          if (!input.value) {
            button.disabled = true;
            try {
              var res = await fetch("/api/secret?id=" + encodeURIComponent(inputId), { cache: "no-store" });
              var data = await res.json().catch(function () { return {}; });
              if (!res.ok || !data.ok) throw new Error(data.msg || "读取失败");
              if (!data.value) {
                toast("error", "未找到已保存的密钥");
                return;
              }
              input.value = data.value;
              input.setAttribute("data-revealed", "true");
            } catch (err) {
              toast("error", "读取密钥失败：" + (err && err.message ? err.message : String(err)));
              return;
            } finally {
              button.disabled = false;
            }
          }
          input.type = "text";
          button.innerHTML = icon("eye-off");
          button.setAttribute("aria-label", "隐藏密钥");
          button.title = "隐藏密钥";
        });
      }

      function onFormChange(event) {
        if (event.target.matches("input, select, textarea")) markDirty();
      }

      document.querySelector("main").addEventListener("input", onFormChange);
      document.querySelector("main").addEventListener("change", onFormChange);

      ["dnshe-domains", "cloudflare-domains", "ipv6-urls"].forEach(setupChips);
      setupReveal("toggle-dnshe-key", "dnshe-api-key");
      setupReveal("toggle-dnshe-secret", "dnshe-api-secret");
      setupReveal("toggle-cloudflare-token", "cloudflare-api-token");

      Array.prototype.forEach.call(document.querySelectorAll(".seg-btn[data-source]"), function (button) {
        button.addEventListener("click", function () { setSourceMode(button.getAttribute("data-source")); });
      });

      Array.prototype.forEach.call(document.querySelectorAll(".preset"), function (button) {
        button.addEventListener("click", function () {
          $(button.getAttribute("data-target")).value = button.getAttribute("data-value");
          markDirty();
        });
      });

      Array.prototype.forEach.call(document.querySelectorAll(".filter"), function (button) {
        button.addEventListener("click", function () {
          state.filter = button.getAttribute("data-filter") || "all";
          renderResults();
        });
      });

      $("save").addEventListener("click", save);
      $("run").addEventListener("click", runNow);
      $("refresh").addEventListener("click", async function () {
        var button = $("refresh");
        button.classList.add("spinning");
        await refresh({ populate: true });
        button.classList.remove("spinning");
      });
      $("open-logs").addEventListener("click", openDrawer);
      $("close-logs").addEventListener("click", closeDrawer);
      $("logs-backdrop").addEventListener("click", closeDrawer);
      $("clear-logs").addEventListener("click", clearLogs);
      $("alert-close").addEventListener("click", hideAlert);
      $("copy-ip").addEventListener("click", copyIP);
      $("refresh-interfaces").addEventListener("click", function () { loadInterfaces(null, true); });

      document.addEventListener("keydown", function (event) {
        if (event.key === "Escape" && $("logs-drawer").classList.contains("open")) closeDrawer();
        if ((event.metaKey || event.ctrlKey) && String(event.key).toLowerCase() === "s") {
          event.preventDefault();
          save();
        }
      });
      document.addEventListener("visibilitychange", function () {
        if (!document.hidden) refresh();
      });

      setupTheme();
      refresh({ populate: true });
      setInterval(function () { refresh(); }, 3000);
    })();
  </script>
</body>
</html>`
