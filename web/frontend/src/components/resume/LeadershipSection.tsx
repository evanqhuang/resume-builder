import { useState } from 'react';
import { Checkbox } from '../common/Checkbox';
import { useResume } from '../../hooks/useResume';
import type { LeadershipEntry } from '../../types/resume';

interface LeadershipSectionProps {
  leadership: LeadershipEntry[];
}

export const LeadershipSection = ({ leadership }: LeadershipSectionProps) => {
  const { dispatch } = useResume();
  const [isExpanded, setIsExpanded] = useState(true);

  const handleToggle = (id: string) => {
    dispatch({ type: 'TOGGLE_LEADERSHIP', payload: id });
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex justify-between items-center mb-4"
      >
        <h2 className="text-xl font-bold text-gray-900">Leadership & Activities</h2>
        <span className="text-gray-500">{isExpanded ? '−' : '+'}</span>
      </button>

      {isExpanded && (
        <ul className="space-y-3">
          {leadership.map((entry) => (
            <li key={entry.id} className="flex items-start gap-2">
              <Checkbox
                checked={entry.selected}
                onChange={() => handleToggle(entry.id)}
                className="mt-1"
              />
              <p className={`text-sm flex-1 ${entry.selected ? 'text-gray-700' : 'text-gray-500 line-through opacity-50'}`}>
                {entry.text}
              </p>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
