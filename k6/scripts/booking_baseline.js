import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 20,
  duration: '1m',
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01']
  }
};

export default function () {
  const payload = JSON.stringify({
    patient_id: '00000000-0000-0000-0000-000000000000',
    provider_id: '00000000-0000-0000-0000-000000000001',
    start_time: new Date(Date.now() + 3600000).toISOString(),
    end_time: new Date(Date.now() + 3600000 + 1800000).toISOString()
  });

  const params = { headers: { 'Content-Type': 'application/json' } };
  const res = http.post(`${__ENV.BASE_URL || 'http://localhost:3000'}/appointments`, payload, params);
  check(res, {
    'status is 201': (r) => r.status === 201
  });
  sleep(1);
}
