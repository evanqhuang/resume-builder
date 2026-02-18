import { useState } from 'react';
import { DndContext, closestCenter } from '@dnd-kit/core';
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable';
import { Checkbox } from '../common/Checkbox';
import { RelevanceBar } from '../job/RelevanceBar';
import { useResume } from '../../hooks/useResume';
import { useSectionReorder } from '../../hooks/useSectionReorder';
import { SortableItem } from '../dnd/SortableItem';
import { DragHandle } from '../dnd/DragHandle';
import type { ProjectEntry } from '../../types/resume';

interface ProjectsSectionProps {
  projects: ProjectEntry[];
}

export const ProjectsSection = ({ projects }: ProjectsSectionProps) => {
  const { dispatch, state } = useResume();
  const { handleDragEnd } = useSectionReorder('projects', projects);
  const [expandedIds, setExpandedIds] = useState<Set<string>>(new Set(projects.map(p => p.id)));
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

  const handleToggleProject = (id: string) => {
    dispatch({ type: 'TOGGLE_PROJECT', payload: id });
  };

  const handleToggleBullet = (entryId: string, bulletId: string) => {
    dispatch({ type: 'TOGGLE_BULLET', payload: { entryId, bulletId, entryType: 'project' } });
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Projects</h2>
      <DndContext collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
        <SortableContext items={projects.map(p => p.id)} strategy={verticalListSortingStrategy}>
          <div className="space-y-4">
            {projects.map((project) => {
              const isExpanded = expandedIds.has(project.id);

              return (
                <SortableItem key={project.id} id={project.id}>
                  {({ dragHandleProps }) => (
                    <div
                      className={`border rounded-lg p-4 transition-all ${
                        project.selected ? 'border-indigo-600 bg-indigo-50' : 'border-gray-300 bg-white opacity-50'
                      }`}
                    >
                      <div className="flex items-start gap-2 mb-2">
                        <DragHandle {...dragHandleProps} />
                        <Checkbox
                          checked={project.selected}
                          onChange={() => handleToggleProject(project.id)}
                        />
                        <div className="flex-1">
                          <button
                            onClick={() => toggleExpanded(project.id)}
                            className="w-full text-left"
                          >
                            <div className="flex justify-between items-start">
                              <div>
                                <h3 className={`font-semibold ${project.selected ? 'text-gray-900' : 'text-gray-500'}`}>
                                  {project.title}
                                </h3>
                                <p className={`text-sm ${project.selected ? 'text-gray-700' : 'text-gray-500'}`}>
                                  {project.technologies}
                                </p>
                                {project.github && (
                                  <a
                                    href={`https://${project.github}`}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="text-xs text-indigo-600 hover:text-indigo-800 hover:underline"
                                    onClick={(e) => e.stopPropagation()}
                                  >
                                    {project.github}
                                  </a>
                                )}
                              </div>
                              <span className="text-gray-500 text-sm">{isExpanded ? '−' : '+'}</span>
                            </div>
                          </button>

                          {isExpanded && (
                            <ul className="mt-3 space-y-2">
                              {project.bullets.map((bullet) => (
                                <li key={bullet.id} className="flex items-start gap-2">
                                  <Checkbox
                                    checked={bullet.selected}
                                    onChange={() => handleToggleBullet(project.id, bullet.id)}
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
                  )}
                </SortableItem>
              );
            })}
          </div>
        </SortableContext>
      </DndContext>
    </div>
  );
};
