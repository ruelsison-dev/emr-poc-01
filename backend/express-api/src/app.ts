import express from 'express';
import helmet from 'helmet';
import rateLimit from 'express-rate-limit';

import patientsRouter from './routes/patients';
import appointmentsRouter from './routes/appointments';
import adminRouter from './routes/admin';

const app = express();
app.use(helmet());
app.use(express.json());
app.use(
  rateLimit({
    windowMs: 60 * 1000,
    max: 200
  })
);

app.use('/patients', patientsRouter);
app.use('/appointments', appointmentsRouter);
app.use('/admin', adminRouter);

app.get('/health', (_req, res) => res.json({ status: 'ok' }));

export default app;
