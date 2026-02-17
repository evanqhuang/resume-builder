import { useState } from 'react';
import { Button } from '../common/Button';
import { analyzeJob } from '../../services/api';
import { useResume } from '../../hooks/useResume';

export const JobInput = () => {
  const { dispatch } = useResume();
  const [jobTitle, setJobTitle] = useState('');
  const [company, setCompany] = useState('');
  const [description, setDescription] = useState('');
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  const handleAnalyze = async () => {
    if (!jobTitle.trim() || !description.trim()) {
      return;
    }

    setIsAnalyzing(true);
    dispatch({ type: 'SET_LOADING', payload: true });

    try {
      const analysis = await analyzeJob(jobTitle, company, description);
      dispatch({ type: 'SET_JOB_ANALYSIS', payload: analysis });
      dispatch({ type: 'SET_ERROR', payload: null });
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to analyze job';
      dispatch({ type: 'SET_ERROR', payload: message });
    } finally {
      setIsAnalyzing(false);
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  return (
    <div className="space-y-4">
      <h2 className="text-lg font-semibold text-gray-900">Job Description</h2>

      <div>
        <label htmlFor="jobTitle" className="block text-sm font-medium text-gray-700 mb-1">
          Job Title <span className="text-red-500">*</span>
        </label>
        <input
          id="jobTitle"
          type="text"
          value={jobTitle}
          onChange={(e) => setJobTitle(e.target.value)}
          placeholder="e.g., Senior Software Engineer"
          className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>

      <div>
        <label htmlFor="company" className="block text-sm font-medium text-gray-700 mb-1">
          Company
        </label>
        <input
          id="company"
          type="text"
          value={company}
          onChange={(e) => setCompany(e.target.value)}
          placeholder="e.g., Google"
          className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>

      <div>
        <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
          Job Description <span className="text-red-500">*</span>
        </label>
        <textarea
          id="description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Paste the full job description here..."
          rows={8}
          className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>

      <Button
        onClick={handleAnalyze}
        disabled={isAnalyzing || !jobTitle.trim() || !description.trim()}
        className="w-full"
      >
        {isAnalyzing ? 'Analyzing...' : 'Analyze Job'}
      </Button>
    </div>
  );
};
