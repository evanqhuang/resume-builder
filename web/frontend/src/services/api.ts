import axios from 'axios';
import type { Resume, JobAnalysisResponse, PartialSectionOrder } from '../types/resume';

const api = axios.create({
  baseURL: '',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const fetchResume = async (): Promise<Resume> => {
  const response = await api.get<Resume>('/api/resume');
  return response.data;
};

export const analyzeJob = async (
  jobTitle: string,
  company: string,
  description: string
): Promise<JobAnalysisResponse> => {
  const response = await api.post<JobAnalysisResponse>('/api/job/analyze', {
    job_title: jobTitle,
    company,
    description,
  });
  return response.data;
};

export const generatePdf = async (selections: Record<string, unknown>): Promise<Blob> => {
  const response = await api.post('/api/generate', selections, {
    responseType: 'blob',
  });
  return response.data;
};

export const saveOrder = async (order: PartialSectionOrder): Promise<void> => {
  await api.put('/api/order', order);
};
