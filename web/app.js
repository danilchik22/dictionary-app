(() => {
  'use strict';

  // --- Configuration ---
  // You can override base URL via ?baseUrl=... query param
  const urlParams = new URLSearchParams(window.location.search);
  const DEFAULT_BACKEND = 'http://81.177.48.223:8081';
  const isDevProxy = (location.hostname === 'localhost' || location.hostname === '127.0.0.1') && location.port === '5173' && !urlParams.get('baseUrl');
  const API_BASE = isDevProxy ? '' : (urlParams.get('baseUrl') || DEFAULT_BACKEND);

  // --- DOM helpers ---
  const $ = (sel) => document.querySelector(sel);
  const notify = (text, isError = false) => {
    const el = $('#notification');
    if (!el) return;
    el.textContent = text || '';
    el.className = 'notification' + (text ? ' show' : '') + (isError ? ' error' : '');
  };

  const setView = (id) => {
    document.querySelectorAll('.view').forEach(v => v.classList.add('hidden'));
    const el = document.getElementById(id);
    if (el) el.classList.remove('hidden');
  };

  const storeTokens = (accessToken, refreshToken) => {
    if (accessToken) localStorage.setItem('access_token', accessToken);
    if (refreshToken) localStorage.setItem('refresh_token', refreshToken);
    updateTokenPreview();
  };
  const getAccessToken = () => localStorage.getItem('access_token');
  const getRefreshToken = () => localStorage.getItem('refresh_token');
  const clearTokens = () => { localStorage.removeItem('access_token'); localStorage.removeItem('refresh_token'); updateTokenPreview(); };
  const updateTokenPreview = () => {
    const token = getAccessToken();
    const preview = token ? token.slice(0, 12) + '...' : '(no token)';
    const el = $('#accessTokenPreview');
    if (el) el.textContent = `Access: ${preview}`;
  };

  // --- API layer ---
  const jsonHeaders = () => ({ 'Content-Type': 'application/json' });
  const authHeaders = () => ({ ...jsonHeaders(), 'Authorization': `Bearer ${getAccessToken() || ''}` });

  async function apiRegister(username, password, age, sex) {
    const res = await fetch(`${API_BASE}/new_user`, { mode: 'cors', credentials: 'omit',
      method: 'POST',
      headers: jsonHeaders(),
      body: JSON.stringify({ username, password, age: Number(age), sex: sex === 'Male' })
    });
    if (!res.ok) throw await parseApiError(res);
    return true;
  }

  async function apiLogin(username, password) {
    const res = await fetch(`${API_BASE}/login`, { mode: 'cors', credentials: 'omit',
      method: 'POST',
      headers: jsonHeaders(),
      body: JSON.stringify({ username, password })
    });
    if (!res.ok) throw await parseApiError(res);
    const data = await res.json();
    return { accessToken: data.access_token, refreshToken: data.refresh_token };
  }

  async function apiRefresh() {
    const token = getRefreshToken();
    if (!token) throw new Error('No refresh token');
    const res = await fetch(`${API_BASE}/refresh`, { mode: 'cors', credentials: 'omit',
      method: 'POST',
      headers: jsonHeaders(),
      body: JSON.stringify({ refresh_token: token })
    });
    if (!res.ok) throw await parseApiError(res);
    const data = await res.json();
    return { accessToken: data.access_token, refreshToken: data.refresh_token };
  }

  async function apiSearch(query) {
    const res = await fetch(`${API_BASE}/api/search`, { mode: 'cors', credentials: 'omit',
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ query, limit: 10, one: true })
    });
    if (res.status === 401) return { unauthorized: true };
    if (!res.ok) throw await parseApiError(res);
    return await res.json();
  }

  async function apiCount() {
    const res = await fetch(`${API_BASE}/api/total_words`, { mode: 'cors', credentials: 'omit',
      method: 'GET',
      headers: authHeaders()
    });
    if (res.status === 401) return { unauthorized: true };
    if (!res.ok) throw await parseApiError(res);
    return await res.json();
  }

  async function apiNewWord(word, definition) {
    const res = await fetch(`${API_BASE}/api/new_word`, { mode: 'cors', credentials: 'omit',
      method: 'POST',
      headers: authHeaders(),
      body: JSON.stringify({ word, definition })
    });
    if (res.status === 401) return { unauthorized: true };
    if (!res.ok) throw await parseApiError(res);
    return await res.json();
  }

  async function parseApiError(res) {
    try {
      const text = await res.text();
      const data = JSON.parse(text);
      const message = data?.message || text || `HTTP ${res.status}`;
      return new Error(message);
    } catch {
      return new Error(`HTTP ${res.status}`);
    }
  }

  // --- UI wiring ---
  function setupAuthFlows() {
    $('#toRegisterBtn').addEventListener('click', () => { setView('registerView'); notify(''); });
    $('#toLoginFromRegisterBtn').addEventListener('click', () => { setView('loginView'); notify(''); });

    $('#loginForm').addEventListener('submit', async (e) => {
      e.preventDefault();
      const username = $('#loginUsername').value.trim();
      const password = $('#loginPassword').value;
      if (!username || !password) return notify('Please fill all fields', true);
      notify('Logging in...');
      try {
        const { accessToken, refreshToken } = await apiLogin(username, password);
        storeTokens(accessToken, refreshToken);
        setView('dictionaryView');
        notify('Login successful');
      } catch (err) {
        notify(err.message || 'Login failed', true);
      }
    });

    $('#registerForm').addEventListener('submit', async (e) => {
      e.preventDefault();
      const username = $('#regUsername').value.trim();
      const password = $('#regPassword').value;
      const age = $('#regAge').value;
      const sex = $('#regSex').value;
      if (!username || !password || !age) return notify('Please fill all fields', true);
      notify('Creating account...');
      try {
        await apiRegister(username, password, age, sex);
        notify('Registration successful, now login');
        setView('loginView');
      } catch (err) {
        notify(err.message || 'Registration failed', true);
      }
    });

    $('#logoutBtn').addEventListener('click', () => {
      clearTokens();
      setView('loginView');
      notify('Logged out');
    });

    $('#refreshTokenBtn').addEventListener('click', async () => {
      notify('Refreshing token...');
      try {
        const { accessToken, refreshToken } = await apiRefresh();
        storeTokens(accessToken, refreshToken);
        notify('Token refreshed');
      } catch (err) {
        notify(err.message || 'Refresh failed', true);
      }
    });
  }

  function setupDictionaryFlows() {
    // Debounced search
    let searchTimer = null;
    $('#searchInput').addEventListener('input', (e) => {
      const q = e.target.value.trim();
      clearTimeout(searchTimer);
      if (!q) { $('#searchResult').textContent = ''; return; }
      searchTimer = setTimeout(async () => {
        $('#searchResult').textContent = 'Searching...';
        try {
          const res = await apiSearch(q);
          if (res?.unauthorized) {
            $('#searchResult').textContent = 'Unauthorized. Please refresh token or login again.';
            return;
          }
          $('#searchResult').textContent = formatSearchResponse(res);
        } catch (err) {
          $('#searchResult').textContent = `Error: ${err.message || 'Search failed'}`;
        }
      }, 300);
    });

    $('#countBtn').addEventListener('click', async () => {
      $('#countResult').textContent = 'Loading...';
      try {
        const res = await apiCount();
        if (res?.unauthorized) {
          $('#countResult').textContent = 'Unauthorized. Please refresh token or login again.';
          return;
        }
        $('#countResult').textContent = `Total words: ${res.count_words ?? res.count ?? 'N/A'}`;
      } catch (err) {
        $('#countResult').textContent = `Error: ${err.message || 'Failed to get total'}`;
      }
    });

    $('#newWordForm').addEventListener('submit', async (e) => {
      e.preventDefault();
      const word = $('#newWord').value.trim();
      const definition = $('#newDefinition').value.trim();
      if (!word || !definition) return notify('Fill both word and definition', true);
      $('#newWordResult').textContent = 'Submitting...';
      try {
        const res = await apiNewWord(word, definition);
        if (res?.unauthorized) {
          $('#newWordResult').textContent = 'Unauthorized. Please refresh token or login again.';
          return;
        }
        $('#newWordResult').textContent = `Added. Word ID: ${res.user_id ?? 'N/A'}`;
      } catch (err) {
        $('#newWordResult').textContent = `Error: ${err.message || 'Failed to add word'}`;
      }
    });
  }

  function formatSearchResponse(data) {
    if (!data) return 'No data';
    const one = data.one_word || data.oneWord;
    const list = data.all_words || data.allWords;
    const similar = data.similar_words || data.similarWords;

    const sections = [];

    if (one) {
      sections.push(`Main word:\n- ${one.word} — ${one.definition}`);
    }

    if (Array.isArray(list) && list.length) {
      sections.push('All words:\n' + list.map((w, i) => `${i + 1}. ${w.word} — ${w.definition}`).join('\n'));
    }

    if (Array.isArray(similar) && similar.length) {
      sections.push('Similar words:\n' + similar.map((w, i) => {
        if (typeof w === 'string') return `${i + 1}. ${w}`;
        if (w && typeof w === 'object' && 'word' in w) return `${i + 1}. ${w.word}${w.definition ? ' — ' + w.definition : ''}`;
        return `${i + 1}. ${String(w)}`;
      }).join('\n'));
    }

    if (sections.length) return sections.join('\n\n');

    return JSON.stringify(data, null, 2);
  }

  function boot() {
    setupAuthFlows();
    setupDictionaryFlows();
    updateTokenPreview();

    if (getAccessToken() && getRefreshToken()) {
      setView('dictionaryView');
    } else {
      setView('loginView');
    }
  }

  window.addEventListener('DOMContentLoaded', boot);
})(); 