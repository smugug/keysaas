                import express from 'express';
import bodyParser from 'body-parser';
import { routes } from './routes';

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(bodyParser.json());

// Define routes
app.use('/api', routes);

// Start server
app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
