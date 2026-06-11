package service

var adminConfigHTML = []byte(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>chat2api admin</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f5f7fb;
      --panel: #ffffff;
      --panel-2: #eef3f8;
      --line: #d7dee8;
      --line-strong: #b8c4d3;
      --text: #17202b;
      --muted: #5d6978;
      --accent: #0f766e;
      --accent-strong: #0b5f59;
      --danger: #b42318;
      --warn: #b54708;
      --ok: #067647;
      --shadow: 0 12px 32px rgba(19, 33, 53, .08);
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font-size: 14px;
    }
    button, input, select {
      font: inherit;
    }
    button {
      border: 1px solid var(--line-strong);
      background: var(--panel);
      color: var(--text);
      height: 34px;
      padding: 0 12px;
      border-radius: 6px;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 6px;
      white-space: nowrap;
    }
    button:hover { border-color: var(--accent); }
    button.primary {
      background: var(--accent);
      border-color: var(--accent);
      color: #fff;
    }
    button.primary:hover {
      background: var(--accent-strong);
      border-color: var(--accent-strong);
    }
    button.danger {
      color: var(--danger);
      border-color: #f0b8b2;
    }
    button.icon {
      width: 34px;
      padding: 0;
    }
    input, select {
      width: 100%;
      height: 34px;
      border: 1px solid var(--line);
      background: #fff;
      color: var(--text);
      border-radius: 6px;
      padding: 0 10px;
      outline: none;
    }
    input:focus, select:focus {
      border-color: var(--accent);
      box-shadow: 0 0 0 3px rgba(15, 118, 110, .12);
    }
    label {
      display: grid;
      gap: 6px;
      color: var(--muted);
      font-size: 12px;
      font-weight: 600;
    }
    .app {
      min-height: 100vh;
      display: grid;
      grid-template-rows: auto 1fr;
    }
    .topbar {
      position: sticky;
      top: 0;
      z-index: 5;
      display: grid;
      grid-template-columns: minmax(180px, 1fr) minmax(260px, 520px) auto;
      gap: 12px;
      align-items: center;
      padding: 14px 18px;
      border-bottom: 1px solid var(--line);
      background: rgba(255, 255, 255, .92);
      backdrop-filter: blur(14px);
    }
    .brand {
      display: flex;
      flex-direction: column;
      min-width: 0;
    }
    h1 {
      margin: 0;
      font-size: 18px;
      line-height: 1.2;
      letter-spacing: 0;
    }
    .sub {
      color: var(--muted);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      margin-top: 3px;
      font-size: 12px;
    }
    .auth {
      display: grid;
      grid-template-columns: 1fr auto;
      gap: 8px;
      align-items: end;
    }
    .toolbar {
      display: flex;
      gap: 8px;
      justify-content: flex-end;
      align-items: center;
    }
    .status {
      min-height: 34px;
      display: flex;
      align-items: center;
      color: var(--muted);
      font-size: 13px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .status.ok { color: var(--ok); }
    .status.err { color: var(--danger); }
    .layout {
      display: grid;
      grid-template-columns: minmax(300px, 380px) minmax(0, 1fr);
      gap: 16px;
      padding: 16px 18px 22px;
      max-width: 1480px;
      width: 100%;
      margin: 0 auto;
    }
    section {
      min-width: 0;
      display: grid;
      gap: 16px;
      align-content: start;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      box-shadow: var(--shadow);
      min-width: 0;
    }
    .panel-head {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 12px;
      padding: 13px 14px;
      border-bottom: 1px solid var(--line);
    }
    h2 {
      margin: 0;
      font-size: 14px;
      letter-spacing: 0;
    }
    .panel-body {
      display: grid;
      gap: 12px;
      padding: 14px;
    }
    .row {
      display: grid;
      grid-template-columns: 1fr auto;
      gap: 8px;
      align-items: end;
    }
    .secret-row {
      display: grid;
      grid-template-columns: 1fr 34px;
      gap: 8px;
      align-items: end;
    }
    .secret-input {
      position: relative;
    }
    .mask {
      color: var(--muted);
      font-size: 11px;
      margin-top: 4px;
      min-height: 14px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .grid-2 {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 10px;
    }
    .accounts {
      display: grid;
      gap: 10px;
    }
    .account {
      background: var(--panel);
      border: 1px solid var(--line);
      border-radius: 8px;
      overflow: hidden;
    }
    .account-head {
      display: grid;
      grid-template-columns: minmax(0, 1fr) auto auto;
      gap: 10px;
      align-items: center;
      padding: 10px 12px;
      background: var(--panel-2);
      border-bottom: 1px solid var(--line);
    }
    .account-title {
      display: flex;
      align-items: center;
      gap: 8px;
      min-width: 0;
      font-weight: 700;
    }
    .badge {
      display: inline-flex;
      align-items: center;
      min-width: 0;
      height: 22px;
      border-radius: 999px;
      padding: 0 8px;
      background: #e8f5f2;
      color: var(--accent-strong);
      border: 1px solid #b9ded7;
      font-size: 12px;
      font-weight: 700;
    }
    .badge.warn {
      background: #fff3e5;
      color: var(--warn);
      border-color: #ffd6a8;
    }
    .account-body {
      display: grid;
      gap: 10px;
      padding: 12px;
    }
    .empty {
      border: 1px dashed var(--line-strong);
      border-radius: 8px;
      color: var(--muted);
      padding: 26px 16px;
      text-align: center;
      background: rgba(255, 255, 255, .6);
    }
    .hide {
      display: none !important;
    }
    @media (max-width: 980px) {
      .topbar {
        grid-template-columns: 1fr;
      }
      .toolbar {
        justify-content: flex-start;
        flex-wrap: wrap;
      }
      .layout {
        grid-template-columns: 1fr;
      }
    }
    @media (max-width: 640px) {
      .topbar, .layout {
        padding-left: 10px;
        padding-right: 10px;
      }
      .grid-2, .account-head {
        grid-template-columns: 1fr;
      }
      .toolbar, .auth {
        grid-template-columns: 1fr;
        display: grid;
      }
      button {
        width: 100%;
      }
      button.icon {
        width: 34px;
      }
    }
  </style>
</head>
<body>
  <div class="app">
    <header class="topbar">
      <div class="brand">
        <h1>chat2api admin</h1>
        <div class="sub" id="configPath">Config file not loaded</div>
      </div>
      <div class="auth">
        <label>Local API key
          <input id="apiKey" type="password" autocomplete="off" placeholder="sk-...">
        </label>
        <button id="loadBtn" class="primary" type="button">Load</button>
      </div>
      <div class="toolbar">
        <div class="status" id="status">Enter a local key to manage the account pool.</div>
        <button id="saveBtn" class="primary" type="button" disabled>Save</button>
      </div>
    </header>

    <main class="layout">
      <section>
        <div class="panel">
          <div class="panel-head">
            <h2>Runtime</h2>
          </div>
          <div class="panel-body">
            <label>Global proxy
              <input id="proxy" placeholder="http://127.0.0.1:7890">
            </label>
            <label>ChatGPT base URL
              <input id="baseUrl" placeholder="https://chatgpt.com">
            </label>
          </div>
        </div>

        <div class="panel">
          <div class="panel-head">
            <h2>Local API keys</h2>
            <button id="addAuthBtn" type="button">Add</button>
          </div>
          <div class="panel-body" id="authTokens"></div>
        </div>

        <div class="panel">
          <div class="panel-head">
            <h2>Direct token prefixes</h2>
            <button id="addPrefixBtn" type="button">Add</button>
          </div>
          <div class="panel-body" id="prefixes"></div>
        </div>
      </section>

      <section>
        <div class="panel">
          <div class="panel-head">
            <h2>Account pool</h2>
            <button id="addAccountBtn" type="button">Add account</button>
          </div>
          <div class="panel-body">
            <div class="accounts" id="accounts"></div>
          </div>
        </div>
      </section>
    </main>
  </div>

  <script>
    const state = { config: null };
    const els = {
      apiKey: document.getElementById('apiKey'),
      loadBtn: document.getElementById('loadBtn'),
      saveBtn: document.getElementById('saveBtn'),
      status: document.getElementById('status'),
      configPath: document.getElementById('configPath'),
      proxy: document.getElementById('proxy'),
      baseUrl: document.getElementById('baseUrl'),
      authTokens: document.getElementById('authTokens'),
      prefixes: document.getElementById('prefixes'),
      accounts: document.getElementById('accounts'),
      addAuthBtn: document.getElementById('addAuthBtn'),
      addPrefixBtn: document.getElementById('addPrefixBtn'),
      addAccountBtn: document.getElementById('addAccountBtn')
    };

    function setStatus(text, kind) {
      els.status.textContent = text;
      els.status.className = 'status' + (kind ? ' ' + kind : '');
    }

    function apiHeaders() {
      return {
        'Authorization': 'Bearer ' + els.apiKey.value.trim(),
        'Content-Type': 'application/json'
      };
    }

    async function requestConfig(method, body) {
      const res = await fetch('/admin/api/config', {
        method,
        headers: apiHeaders(),
        body: body ? JSON.stringify(body) : undefined
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        throw new Error(data.error || (data.detail && data.detail.msg) || 'Request failed');
      }
      return data;
    }

    function secretItem(item) {
      return {
        index: Number.isInteger(item && item.index) ? item.index : -1,
        set: Boolean(item && item.set),
        masked: item && item.masked ? item.masked : '',
        value: ''
      };
    }

    function emptySecret() {
      return { index: -1, set: false, masked: '', value: '' };
    }

    function newAccount() {
      return {
        index: -1,
        id_token: emptySecret(),
        access_token: emptySecret(),
        refresh_token: emptySecret(),
        account_id: '',
        last_refresh: '',
        email: '',
        type: 'codex',
        expired: '',
        proxy: ''
      };
    }

    function normalizeLoadedConfig(cfg) {
      cfg.auth_tokens = (cfg.auth_tokens || []).map(secretItem);
      cfg.access_token_prefixes = (cfg.access_token_prefixes || []).map(secretItem);
      cfg.chatgpts = (cfg.chatgpts || []).map((account) => ({
        index: Number.isInteger(account.index) ? account.index : -1,
        id_token: secretItem(account.id_token),
        access_token: secretItem(account.access_token),
        refresh_token: secretItem(account.refresh_token),
        account_id: account.account_id || '',
        last_refresh: account.last_refresh || '',
        email: account.email || '',
        type: account.type || '',
        expired: account.expired || '',
        proxy: account.proxy || ''
      }));
      return cfg;
    }

    function renderSecretList(root, list, label, placeholder) {
      root.innerHTML = '';
      if (list.length === 0) {
        const empty = document.createElement('div');
        empty.className = 'empty';
        empty.textContent = 'No entries';
        root.appendChild(empty);
        return;
      }
      list.forEach((item, index) => {
        const row = document.createElement('div');
        row.className = 'secret-row';
        row.innerHTML =
          '<label>' + label + ' ' + (index + 1) +
          '<input type="password" autocomplete="off" placeholder="' + escapeAttr(placeholder) + '" value="' + escapeAttr(item.value || '') + '">' +
          '<span class="mask">' + (item.set ? 'Saved: ' + escapeText(item.masked) : 'New value required') + '</span>' +
          '</label>' +
          '<button class="icon danger" title="Remove" type="button">x</button>';
        row.querySelector('input').addEventListener('input', (event) => {
          item.value = event.target.value;
        });
        row.querySelector('button').addEventListener('click', () => {
          list.splice(index, 1);
          render();
        });
        root.appendChild(row);
      });
    }

    function renderAccounts() {
      els.accounts.innerHTML = '';
      const accounts = state.config.chatgpts;
      if (accounts.length === 0) {
        const empty = document.createElement('div');
        empty.className = 'empty';
        empty.textContent = 'No upstream accounts configured';
        els.accounts.appendChild(empty);
        return;
      }
      accounts.forEach((account, index) => {
        const card = document.createElement('div');
        card.className = 'account';
        const title = account.email || account.account_id || 'Account ' + (index + 1);
        const tokenState = account.access_token.set || account.access_token.value ? 'token set' : 'missing token';
        const tokenBadge = account.access_token.set || account.access_token.value ? 'badge' : 'badge warn';
        card.innerHTML =
          '<div class="account-head">' +
            '<div class="account-title"><span>' + escapeText(title) + '</span><span class="' + tokenBadge + '">' + tokenState + '</span></div>' +
            '<select data-field="type">' +
              option('codex', account.type) +
              option('plus', account.type) +
              option('pro', account.type) +
              option('free', account.type) +
              option('', account.type || '') +
            '</select>' +
            '<button class="danger" type="button">Delete</button>' +
          '</div>' +
          '<div class="account-body">' +
            '<div class="grid-2">' +
              input('Email', 'email', account.email, 'user@example.com') +
              input('Account ID', 'account_id', account.account_id, '') +
            '</div>' +
            '<div class="grid-2">' +
              secretInput('Access token', 'access_token', account.access_token, 'new access token') +
              secretInput('Refresh token', 'refresh_token', account.refresh_token, 'optional refresh token') +
            '</div>' +
            '<div class="grid-2">' +
              secretInput('ID token', 'id_token', account.id_token, 'optional id token') +
              input('Proxy', 'proxy', account.proxy, 'account proxy overrides global proxy') +
            '</div>' +
            '<div class="grid-2">' +
              input('Last refresh', 'last_refresh', account.last_refresh, '') +
              input('Expired', 'expired', account.expired, '') +
            '</div>' +
          '</div>';
        card.querySelector('select[data-field="type"]').addEventListener('change', (event) => {
          account.type = event.target.value;
        });
        card.querySelector('button.danger').addEventListener('click', () => {
          accounts.splice(index, 1);
          render();
        });
        card.querySelectorAll('input[data-field]').forEach((inputEl) => {
          const field = inputEl.dataset.field;
          inputEl.addEventListener('input', (event) => {
            if (account[field] && typeof account[field] === 'object') {
              account[field].value = event.target.value;
            } else {
              account[field] = event.target.value;
            }
          });
        });
        els.accounts.appendChild(card);
      });
    }

    function input(label, field, value, placeholder) {
      return '<label>' + label +
        '<input data-field="' + field + '" value="' + escapeAttr(value || '') + '" placeholder="' + escapeAttr(placeholder || '') + '">' +
        '</label>';
    }

    function secretInput(label, field, item, placeholder) {
      return '<label>' + label +
        '<input type="password" autocomplete="off" data-field="' + field + '" value="' + escapeAttr(item.value || '') + '" placeholder="' + escapeAttr(placeholder || '') + '">' +
        '<span class="mask">' + (item.set ? 'Saved: ' + escapeText(item.masked) : 'Not set') + '</span>' +
        '</label>';
    }

    function option(value, selected) {
      const label = value || 'custom';
      return '<option value="' + escapeAttr(value) + '"' + (value === selected ? ' selected' : '') + '>' + escapeText(label) + '</option>';
    }

    function escapeText(value) {
      return String(value == null ? '' : value)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
    }

    function escapeAttr(value) {
      return escapeText(value).replace(/"/g, '&quot;');
    }

    function render() {
      if (!state.config) {
        return;
      }
      els.saveBtn.disabled = false;
      els.configPath.textContent = state.config.config_path || 'Config file loaded';
      els.proxy.value = state.config.proxy || '';
      els.baseUrl.value = state.config.chatgpt_base_url || '';
      renderSecretList(els.authTokens, state.config.auth_tokens, 'Key', 'new local API key');
      renderSecretList(els.prefixes, state.config.access_token_prefixes, 'Prefix', 'private-prefix-');
      renderAccounts();
    }

    async function loadConfig() {
      const key = els.apiKey.value.trim();
      if (!key) {
        setStatus('Local API key is required.', 'err');
        return;
      }
      localStorage.setItem('chat2apiAdminKey', key);
      setStatus('Loading...', '');
      try {
        state.config = normalizeLoadedConfig(await requestConfig('GET'));
        render();
        setStatus('Loaded.', 'ok');
      } catch (err) {
        setStatus(err.message, 'err');
      }
    }

    async function saveConfig() {
      if (!state.config) {
        return;
      }
      state.config.proxy = els.proxy.value;
      state.config.chatgpt_base_url = els.baseUrl.value;
      setStatus('Saving...', '');
      try {
        state.config = normalizeLoadedConfig(await requestConfig('PUT', state.config));
        render();
        setStatus('Saved. The running pool is refreshed.', 'ok');
      } catch (err) {
        setStatus(err.message, 'err');
      }
    }

    els.loadBtn.addEventListener('click', loadConfig);
    els.saveBtn.addEventListener('click', saveConfig);
    els.addAuthBtn.addEventListener('click', () => {
      if (!state.config) return;
      state.config.auth_tokens.push(emptySecret());
      render();
    });
    els.addPrefixBtn.addEventListener('click', () => {
      if (!state.config) return;
      state.config.access_token_prefixes.push(emptySecret());
      render();
    });
    els.addAccountBtn.addEventListener('click', () => {
      if (!state.config) return;
      state.config.chatgpts.push(newAccount());
      render();
    });
    els.apiKey.addEventListener('keydown', (event) => {
      if (event.key === 'Enter') loadConfig();
    });
    els.apiKey.value = localStorage.getItem('chat2apiAdminKey') || '';
  </script>
</body>
</html>`)
