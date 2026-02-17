import { useState } from 'react';
import { Checkbox } from '../common/Checkbox';
import { RelevanceBar } from '../job/RelevanceBar';
import { useResume } from '../../hooks/useResume';
import type { ExperienceEntry } from '../../types/resume';

interface ExperienceSectionProps {
  experiences: ExperienceEntry[];
}

export const ExperienceSection = ({ experiences }: ExperienceSectionProps) => {
  const { dispatch, state } = useResume();
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set(experiences.map(e => e.id)));
  const { jobAnalysis } = state;

  const toggleExpanded = (id: string) => {
    const newExpanded = new Set(expandedIds);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedIds(newExpanded);
  };

  const handleToggleExperience = (id: string) => {
    dispatch({ type: 'TOGGLE_EXPERIENCE', payload: id });
  };

  const handleToggleBullet = (entryId: string, bulletId: string) => {
    dispatch({ type: 'TOGGLE_BULLET', payload: { entryId, bulletId, entryType: 'experience' } });
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Experience</h2>
      <div className="space-y-4">
        {experiences.map((exp) => {
          const isExpanded = expandedIds.has(exp.id);
          const dateRange = `${exp.start_date} - ${exp.end_date || 'Present'}`;

          return (
            <div
              key={exp.id}
              className={`border rounded-lg p-4 transition-all ${
                exp.selected ? 'border-indigo-600 bg-indigo-50' : 'border-gray-300 bg-white opacity-50'
              }`}
            >
              <div className="flex items-start gap-3 mb-2">
                <Checkbox
                  checked={exp.selected}
                  onChange={() => handleToggleExperience(exp.id)}
                />
                <div className="flex-1">
                  <button
                    onClick={() => toggleExpanded(exp.id)}
                    className="w-full text-left"
                  >
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className={`font-semibold ${exp.selected ? 'text-gray-900' : 'text-gray-500'}`}>
                          {exp.title}
                        </h3>
                        <p className={`text-sm ${exp.selected ? 'text-gray-700' : 'text-gray-500'}`}>
                          {exp.company}
                        </p>
                      </div>
                      <div className="text-right flex items-center gap-2">
                        <p className={`text-sm ${exp.selected ? 'text-gray-600' : 'text-gray-500'}`}>
                          {exp.location} • {dateRange}
                        </p>
                        <span className="text-gray-500 text-sm">{isExpanded ? '−' : '+'}</span>
                      </div>
                    </div>
                  </button>

                  {isExpanded && (
                    <ul className="mt-3 space-y-2">
                      {exp.bullets.map((bullet) => (
                        <li key={bullet.id} className="flex items-start gap-2">
                          <Checkbox
                            checked={bullet.selected}
                            onChange={() => handleToggleBullet(exp.id, bullet.id)}
                            className="mt-1"
                          />
                          <div className="flex-1">
                            <p className={`text-sm ${bullet.selected ? 'text-gray-700' : 'text-gray-500 line-through'}`}>
                              {bullet.text}
                            </p>
                          </div>
                          {jobAnalysis && bullet.relevanceScore !== undefined && (
                            <div className="ml-2">
                              <RelevanceBar score={bullet.relevanceScore} />
                            </div>
                          )}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
