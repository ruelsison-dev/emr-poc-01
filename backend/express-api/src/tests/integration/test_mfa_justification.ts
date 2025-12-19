import request from 'supertest';
import app from '../../app';
import fs from 'fs';
import path from 'path';

const EVIDENCE_DIR = path.join(__dirname, '..', '..', '..', 'specs', '001-patient-admin', 'docs', 'compliance_evidence', 'mfa');

describe('MFA & Justification enforcement', () => {
  it('rejects when MFA not present', async () => {
    const res = await request(app).post('/admin/revoke-consent/123').send({ justification: 'This action is needed' });
    expect(res.status).toBe(403);
    expect(res.body.error).toMatch(/MFA required/);
  });

  it('rejects when justification is missing or too short', async () => {
    const res = await request(app).post('/admin/revoke-consent/123').set('X-MFA-Verified', 'true').send({ justification: 'short' });
    expect(res.status).toBe(400);
  });

  it('accepts when MFA and valid justification present and writes evidence', async () => {
    // cleanup folder
    await fs.promises.mkdir(EVIDENCE_DIR, { recursive: true });
    const before = await fs.promises.readdir(EVIDENCE_DIR).catch(() => []);

    const res = await request(app)
      .post('/admin/revoke-consent/abc')
      .set('X-MFA-Verified', 'true')
      .set('X-Justification', 'Valid justification for emergency override reasons')
      .send();

    expect(res.status).toBe(200);
    const after = await fs.promises.readdir(EVIDENCE_DIR).catch(() => []);
    expect(after.length).toBeGreaterThan(before.length);
  });
});
