import { Router } from 'express';
import request from 'request';
import { token,cert } from '../server';

// curl --post --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" https://kubernetes.default.svc/apis/keysaascontroller.keysaas/v1/namespaces/customer2/keysaases
const router = Router();
const BASE_URL = 'https://kubernetes.default.svc/apis/keysaascontroller.keysaas/v1/namespaces/customer2/keysaases';

// GET all KeySaaS instances
router.get('/', (req, res) => {
  const options = {
    url: BASE_URL,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    agentOptions: {
      ca: cert,
      rejectUnauthorized: true,
    },
  };

  request.get(options, (error, response, body) => {
    if (error) {
      res.status(500).json({ message: 'Error fetching KeySaaS instances', error });
    } else {
      res.json(JSON.parse(body));
    }
  });
});

// POST (Create) a new KeySaaS instance
router.post('/', (req, res) => {
  const options = {
    url: BASE_URL,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    agentOptions: {
      ca: cert,
      rejectUnauthorized: true,
    },
    json: true,
    body: req.body,
  };

  request.post(options, (error, response, body) => {
    if (error) {
      res.status(500).json({ message: 'Error creating KeySaaS instance', error });
    } else {
      res.json(body);
    }
  });
});

// PUT (Update) a KeySaaS instance
router.put('/:name', (req, res) => {
  const options = {
    url: `${BASE_URL}/${req.params.name}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    agentOptions: {
      ca: cert,
      rejectUnauthorized: true,
    },
    json: true,
    body: req.body,
  };

  request.put(options, (error, response, body) => {
    if (error) {
      res.status(500).json({ message: `Error updating KeySaaS instance ${req.params.name}`, error });
    } else {
      res.json(body);
    }
  });
});

// DELETE a KeySaaS instance
router.delete('/:name', (req, res) => {
  const options = {
    url: `${BASE_URL}/${req.params.name}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    agentOptions: {
      ca: cert,
      rejectUnauthorized: true,
    },
  };

  request.delete(options, (error, response, body) => {
    if (error) {
      res.status(500).json({ message: `Error deleting KeySaaS instance ${req.params.name}`, error });
    } else {
      res.json(JSON.parse(body));
    }
  });
});

export default router