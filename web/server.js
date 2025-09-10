import http from 'http';
import fs from 'fs';
import path from 'path';
import url from 'url';
import httpProxy from 'http-proxy';
const { createProxyServer } = httpProxy;

const __dirname = path.dirname(url.fileURLToPath(import.meta.url));
const PORT = Number(process.env.PORT || 5173);
const BASE_URL = process.env.BASE_URL || 'https://danildanil.duckdns.org';

const proxy = createProxyServer({ target: BASE_URL, changeOrigin: true });

function setCors(res) {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET,POST,PUT,DELETE,OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
}

const server = http.createServer((req, res) => {
  const parsed = url.parse(req.url);

  // CORS preflight
  if (req.method === 'OPTIONS') {
    setCors(res);
    res.writeHead(204);
    res.end();
    return;
  }

  // Proxy API calls
  if (parsed.pathname.startsWith('/api') || parsed.pathname === '/login' || parsed.pathname === '/new_user' || parsed.pathname === '/refresh') {
    setCors(res);
    proxy.web(req, res);
    return;
  }

  // Serve static files
  let filePath = path.join(__dirname, parsed.pathname === '/' ? '/index.html' : parsed.pathname);
  if (!filePath.startsWith(__dirname)) {
    res.writeHead(403);
    res.end('Forbidden');
    return;
  }
  fs.readFile(filePath, (err, data) => {
    if (err) {
      res.writeHead(404);
      res.end('Not found');
      return;
    }
    const ext = path.extname(filePath);
    const type = ext === '.html' ? 'text/html' : ext === '.css' ? 'text/css' : ext === '.js' ? 'text/javascript' : 'application/octet-stream';
    res.writeHead(200, { 'Content-Type': type });
    res.end(data);
  });
});

server.listen(PORT, () => {
  console.log(`Web server running at http://localhost:${PORT} (proxy -> ${BASE_URL})`);
}); 