import { Router } from 'express';
import { mfaEnforce } from '../middleware/mfa';

const router = Router();

// Example privileged endpoint: revoke consent for a patient (simulated)
router.post('/revoke-consent/:id', mfaEnforce('revoke_consent'), async (req, res) => {
  const id = req.params.id;
  // Simulate action
  return res.status(200).json({ id, status: 'consent_revoked' });
});

export default router;
