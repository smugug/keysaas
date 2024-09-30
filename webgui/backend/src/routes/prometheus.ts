import { Router } from 'express';
import axios from 'axios';

const router = Router();
const PROMETHEUS_URL = 'http://prometheus-operated.default.svc:9090/api/v1';

// Query Prometheus metrics
router.get('/', async (req, res) => {
  try {
    const query = req.query.query as string;
    if (!query) {
      res.status(400).json({ message: 'Missing Prometheus query parameter' });
    } else {
      const url = `${PROMETHEUS_URL}/query`;
      const response = await axios.get(url, {
        params: { query },
      });
      res.json(response.data);
    }
  } catch (error) {
    res.status(500).json({ message: 'Error fetching Prometheus metrics', error });
  }
});

export default router;
