import { useResume } from '../../hooks/useResume';
import { SkillsSection } from './SkillsSection';
import { ExperienceSection } from './ExperienceSection';
import { ProjectsSection } from './ProjectsSection';
import { EducationSection } from './EducationSection';
import { LeadershipSection } from './LeadershipSection';

export const ResumeEditor = () => {
  const { state } = useResume();
  const { resume } = state;

  if (!resume) {
    return (
      <div className="bg-white rounded-lg shadow-sm p-8 text-center">
        <p className="text-gray-600">Loading resume...</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-lg shadow-sm p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-2">{resume.contact.name}</h2>
        <div className="flex flex-wrap gap-4 text-sm text-gray-600">
          <span>{resume.contact.email}</span>
          <span>{resume.contact.phone}</span>
          {resume.contact.linkedin && <span>{resume.contact.linkedin}</span>}
          {resume.contact.github && <span>{resume.contact.github}</span>}
        </div>
      </div>

      <EducationSection education={resume.education} />
      <SkillsSection skills={resume.skills} />
      <ExperienceSection experiences={resume.experience} />
      <ProjectsSection projects={resume.projects} />
      <LeadershipSection leadership={resume.leadership} />
    </div>
  );
};
