import { Request, Response } from 'express';

// Simple in-memory store for the scaffold
const patients: Record<string, any> = {};

export const createPatient = (req: Request, res: Response) => {
  const id = `pat_${Date.now()}`;
  const payload = { id, ...req.body };
  patients[id] = payload;
  res.status(201).json(payload);
};

export const listPatients = (_req: Request, res: Response) => {
  res.json(Object.values(patients));
};

export const getPatient = (req: Request, res: Response) => {
  const { id } = req.params;
  const p = patients[id];
  if (!p) return res.status(404).json({ error: 'not_found' });
  return res.json(p);
};
