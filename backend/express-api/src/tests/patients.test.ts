import request from 'supertest';
import app from '../app';

describe('Patients API (scaffold)', () => {
  it('creates and retrieves a patient', async () => {
    const createRes = await request(app).post('/patients').send({ primary_person_id: 'person_1', mrn: 'MRN123' });
    expect(createRes.status).toBe(201);
    expect(createRes.body).toHaveProperty('id');

    const getRes = await request(app).get(`/patients/${createRes.body.id}`);
    expect(getRes.status).toBe(200);
    expect(getRes.body.mrn).toBe('MRN123');
  });

  it('returns 404 for missing patient', async () => {
    const res = await request(app).get('/patients/not-exist');
    expect(res.status).toBe(404);
  });
});
