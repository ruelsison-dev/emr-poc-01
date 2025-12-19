import { Router } from 'express';
import { createAppointment, listAppointments, getAppointment } from '../controllers/appointmentsController';

const router = Router();
router.post('/', createAppointment);
router.get('/', listAppointments);
router.get('/:id', getAppointment);

export default router;
