import { Request, Response } from 'express';

const appointments: Record<string, any> = {};

export const createAppointment = (req: Request, res: Response) => {
  const id = `appt_${Date.now()}`;
  const payload = { id, ...req.body };
  appointments[id] = payload;
  res.status(201).json(payload);
};

export const listAppointments = (_req: Request, res: Response) => {
  res.json(Object.values(appointments));
};

export const getAppointment = (req: Request, res: Response) => {
  const { id } = req.params;
  const a = appointments[id];
  if (!a) return res.status(404).json({ error: 'not_found' });
  return res.json(a);
};
