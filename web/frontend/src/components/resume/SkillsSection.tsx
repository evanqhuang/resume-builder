import { useState } from 'react';
import { Checkbox } from '../common/Checkbox';
import { useResume } from '../../hooks/useResume';
import type { SkillCategory } from '../../types/resume';

interface SkillsSectionProps {
  skills: SkillCategory[];
}

export const SkillsSection = ({ skills }: SkillsSectionProps) => {
  const { dispatch, state } = useResume();
  const [isExpanded, setIsExpanded] = useState(true);
  const { jobAnalysis } = state;

  const getRelevanceColor = (skillName: string): string => {
    if (!jobAnalysis) return '';
    const score = jobAnalysis.scores[skillName] || 0;
    if (score >= 70) return 'border-green-500 bg-green-50';
    if (score >= 40) return 'border-yellow-500 bg-yellow-50';
    if (score > 0) return 'border-red-500 bg-red-50';
    return '';
  };

  const handleToggle = (category: string, skillName: string) => {
    dispatch({ type: 'TOGGLE_SKILL', payload: { category, skillName } });
  };

  const handleToggleCategory = (category: string, items: { selected: boolean }[]) => {
    const allSelected = items.every((item) => item.selected);
    dispatch({ type: 'TOGGLE_SKILL_CATEGORY', payload: { category, selected: !allSelected } });
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <button
        onClick={() => setIsExpanded(!isExpanded)}
        className="w-full flex justify-between items-center mb-4"
      >
        <h2 className="text-xl font-bold text-gray-900">Skills</h2>
        <span className="text-gray-500">{isExpanded ? '−' : '+'}</span>
      </button>

      {isExpanded && (
        <div className="space-y-6">
          {skills.map((category) => (
            <div key={category.category}>
              <div className="flex items-center gap-2 mb-2">
                <Checkbox
                  checked={category.items.every((item) => item.selected)}
                  onChange={() => handleToggleCategory(category.category, category.items)}
                />
                <h3 className="text-sm font-semibold text-gray-700 capitalize">
                  {category.category}
                </h3>
              </div>
              <div className="flex flex-wrap gap-2">
                {category.items.map((skill) => (
                  <div
                    key={skill.name}
                    className={`inline-flex items-center border rounded-full px-3 py-1 transition-all ${
                      skill.selected
                        ? `border-indigo-600 bg-indigo-50 ${getRelevanceColor(skill.name)}`
                        : 'border-gray-300 bg-white opacity-50'
                    }`}
                  >
                    <Checkbox
                      checked={skill.selected}
                      onChange={() => handleToggle(category.category, skill.name)}
                    />
                    <span className={`ml-2 text-sm ${skill.selected ? 'text-gray-900' : 'text-gray-500 line-through'}`}>
                      {skill.name}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
