import fs from 'fs';
import path from 'path';

export interface JustificationRecord {
  id?: string;
  actor_id: string;
  action: string;
  resource_id?: string | null;
  justification: string;
  mfa_verified: boolean;
  metadata?: Record<string, any>;
  timestamp?: string;
}

const EVIDENCE_DIR = process.env.EVIDENCE_DIR || path.join(process.cwd(), 'specs', '001-patient-admin', 'docs', 'compliance_evidence', 'mfa');
// Use process.cwd() so tests and runtime agree on the evidence path (ts-jest can change __dirname resolution)

export async function recordJustification(rec: JustificationRecord) {
  try {
    const payload = {
      ...rec,
      timestamp: new Date().toISOString()
    };
    // ensure evidence dir
    await fs.promises.mkdir(EVIDENCE_DIR, { recursive: true });
    const filename = `justification-${Date.now()}.json`;
    const filepath = path.join(EVIDENCE_DIR, filename);
    await fs.promises.writeFile(filepath, JSON.stringify(payload, null, 2), { encoding: 'utf-8' });
    // Additionally, log to stdout for aggregator ingestion
    console.log('JUSTIFICATION_RECORD', JSON.stringify(payload));
  } catch (err) {
    console.error('Failed to record justification', err);
  }
}
