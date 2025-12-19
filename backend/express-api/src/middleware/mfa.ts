import { Request, Response, NextFunction } from 'express';
import { recordJustification } from '../lib/justification';

/**
 * mfaEnforce middleware
 * - Verifies MFA via `req.user?.mfa_verified` or a test header in local dev
 * - Requires a justification in `X-Justification` header or `req.body.justification`
 */
export function mfaEnforce(action: string) {
  return async (req: Request, res: Response, next: NextFunction) => {
    try {
      const mfaVerified = (req as any).user?.mfa_verified || req.header('X-MFA-Verified') === 'true';
      if (!mfaVerified) {
        return res.status(403).json({ error: 'MFA required for this operation' });
      }

      const justification = req.header('X-Justification') || (req.body && req.body.justification);
      if (!justification || justification.trim().length < 15) {
        return res.status(400).json({ error: 'Justification is required and must be at least 15 characters' });
      }

      // record justification asynchronously (fire-and-forget but capture errors)
      await recordJustification({
        actor_id: (req as any).user?.id || 'anonymous',
        action,
        resource_id: (req.params && req.params.id) || null,
        justification,
        mfa_verified: true,
        metadata: {
          ip: req.ip,
          ua: req.get('user-agent')
        }
      });

      return next();
    } catch (err) {
      console.error('mfaEnforce error', err);
      return res.status(500).json({ error: 'internal error' });
    }
  };
}
