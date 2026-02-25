import { Router } from 'express';

/**
 * Users router - proxies /api/v1/users/* to Go domain service.
 * Per ADR-001: BFF routes to Go, does not implement business logic.
 * Uses fetch-based forward to avoid 301 redirect loops from http-proxy-middleware.
 */
export function usersRouter(goApiUrl) {
  const router = Router();

  router.use('/', async (req, res) => {
    const url = `${goApiUrl.replace(/\/$/, '')}${req.originalUrl}`;

    try {
      const headers = {};
      for (const [k, v] of Object.entries(req.headers)) {
        if (k.toLowerCase() !== 'host' && v) headers[k] = v;
      }
      const body = req.method !== 'GET' && req.method !== 'HEAD' && req.body
        ? JSON.stringify(req.body) : undefined;
      const upstream = await fetch(url, {
        method: req.method,
        headers,
        body,
      });
      res.status(upstream.status);
      upstream.headers.forEach((v, k) => res.setHeader(k, v));
      const text = await upstream.text();
      res.send(text);
    } catch (err) {
      res.status(502).json({ error: 'Upstream unavailable' });
    }
  });
  return router;
}
