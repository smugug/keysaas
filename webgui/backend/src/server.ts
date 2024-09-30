import express from 'express';
import keySaaSRoutes from './routes/keysaas';
import prometheusRoutes from './routes/prometheus';
import fs from 'fs';
import cors from 'cors';

const app = express();
const port = 4000;
const token_path =  '/var/run/secrets/kubernetes.io/serviceaccount/token';
const cert_path = '/var/run/secrets/kubernetes.io/serviceaccount/ca.crt'
const readFile = (path: string): string => {
  try {
    const content = fs.readFileSync(path, 'utf8');
    return content;
  } catch (error) {
    console.error('Error reading service account', error);
    return '';
  }
};
export const token = readFile(token_path);
export const cert = readFile(cert_path);

console.log('0.0.12')

app.use(cors());
app.use(express.json());

// Routes
app.use('/api/keysaas', keySaaSRoutes);
app.use('/api/prometheus', prometheusRoutes);


app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});


