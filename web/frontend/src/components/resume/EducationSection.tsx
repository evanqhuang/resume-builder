import type { EducationEntry } from '../../types/resume';

interface EducationSectionProps {
  education: EducationEntry;
}

export const EducationSection = ({ education }: EducationSectionProps) => {
  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-xl font-bold text-gray-900 mb-4">Education</h2>
      <div className="space-y-2">
        <div className="flex justify-between items-start">
          <div>
            <h3 className="font-semibold text-gray-900">{education.institution}</h3>
            <p className="text-gray-700">{education.degree}</p>
          </div>
          <p className="text-gray-600 text-sm">{education.location}</p>
        </div>
        {education.minor && (
          <p className="text-gray-600 text-sm">Minor: {education.minor}</p>
        )}
        {education.gpa && (
          <p className="text-gray-600 text-sm">GPA: {education.gpa}</p>
        )}
        {education.honors && (
          <p className="text-gray-600 text-sm">{education.honors}</p>
        )}
      </div>
    </div>
  );
};
