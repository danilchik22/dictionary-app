# DictionaryApp Web

A minimal web client (HTML + JS) for the same backend used by the Android app.

How to run (no server):
- Open `index.html` directly in a browser, or serve the `web/` folder with any static server.
- Backend must have CORS enabled.

Recommended dev server (with proxy to avoid CORS):
- Requires Node.js 18+
- In `web/` run:
```bash
npm install http-proxy@1 --no-save
npm start
```
- By default the proxy target is `http://192.168.3.60:8081`. Override via env:
```bash
BASE_URL=http://HOST:PORT PORT=5173 npm start
```
- Open `http://localhost:5173`.

Features:
- Login / Registration
- Token refresh
- Search with debounce (300ms)
- Total words
- Add new word

Notes:
- Tokens are stored in `localStorage`. 