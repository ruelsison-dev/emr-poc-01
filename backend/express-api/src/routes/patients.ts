import { Router } from 'express';
import { createPatient, listPatients, getPatient } from '../controllers/patientsController';

const router = Router();
router.post('/', createPatient);
router.get('/', listPatients);
router.get('/:id', getPatient);

export default router;
