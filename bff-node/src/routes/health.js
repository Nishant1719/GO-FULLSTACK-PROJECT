import { Router } from 'express';

/**
 * Health router - aggregates health from BFF and Go domain.
 * Early checks for liveness/readiness.
 */
export function healthRouter(goApiUrl) {
  const router = Router();

  router.get('/', async (req, res) => {
    const result = {
      status: 'healthy',
      bff: 'ok',
    };

    try {
      const resp = await fetch(`${goApiUrl}/ping`);
      const data = await resp.json();
      result.goDomain = resp.ok ? 'ok' : `error: ${data?.message || resp.status}`;
      if (!resp.ok) result.status = 'degraded';
    } catch (err) {
      result.goDomain = `unreachable: ${err.message}`;
      result.status = 'degraded';
    }

    res.json(result);
  });

  return router;
}
