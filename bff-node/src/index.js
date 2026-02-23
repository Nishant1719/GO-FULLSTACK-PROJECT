import express from 'express';
import { usersRouter } from './routes/users.js';
import { healthRouter } from './routes/health.js';

const app = express();
const PORT = process.env.PORT || 3000;
const GO_API_URL = process.env.GO_API_URL || 'http://localhost:8080';

app.use(express.json());

// Liveness probe - early check
app.get('/ping', (req, res) => {
  res.json({ message: 'pong', status: 'healthy' });
});

// Routers
app.use('/api/v1/users', usersRouter(GO_API_URL));
app.use('/health', healthRouter(GO_API_URL));

app.listen(PORT, () => {
  console.log(`BFF listening on http://localhost:${PORT}`);
  console.log(`Proxying to Go API at ${GO_API_URL}`);
});
