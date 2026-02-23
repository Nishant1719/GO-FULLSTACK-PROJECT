import { Router } from 'express';
import { createProxyMiddleware } from 'http-proxy-middleware';

/**
 * Users router - proxies /api/v1/users/* to Go domain service.
 * Per ADR-001: BFF routes to Go, does not implement business logic.
 */
export function usersRouter(goApiUrl) {
  const router = Router();

  const proxy = createProxyMiddleware({
    target: goApiUrl,
    changeOrigin: true,
    pathRewrite: (path) => {
      const base = '/api/v1/users'
      return path === '/' || path === '' ? base : `${base}${path}`
    },
    onProxyReq: (proxyReq, req) => {
      // Forward identity context header when JWT auth is added
      if (req.headers['x-user-id']) {
        proxyReq.setHeader('X-User-Id', req.headers['x-user-id']);
      }
    },
  });

  router.use('/', proxy);
  return router;
}
