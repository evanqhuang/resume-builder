import { Button } from '../common/Button';
import { JobInput } from '../job/JobInput';
import { KeywordBadges } from '../job/KeywordBadges';
import { useResume } from '../../hooks/useResume';
import { generatePdf } from '../../services/api';

export const Sidebar = () => {
  const { state, dispatch } = useResume();
  const { jobAnalysis, resume, error } = state;

  const handleSelectAll = () => {
    dispatch({ type: 'SELECT_ALL' });
  };

  const handleDeselectAll = () => {
    dispatch({ type: 'DESELECT_ALL' });
  };

  const handleApplySuggestions = () => {
    dispatch({ type: 'APPLY_SUGGESTIONS', payload: { threshold: 70 } });
  };

  const handleGeneratePdf = async () => {
    if (!resume) return;

    try {
      const selectedSkills = resume.skills
        .flatMap((cat) => cat.items.filter((s) => s.selected).map((s) => s.name));
      const selectedBullets: string[] = [];
      const selectedExperience: string[] = [];
      const selectedProjects: string[] = [];

      resume.experience.forEach((entry) => {
        if (entry.selected) {
          selectedExperience.push(entry.id);
          entry.bullets.filter((b) => b.selected).forEach((b) => selectedBullets.push(b.id));
        }
      });

      resume.projects.forEach((entry) => {
        if (entry.selected) {
          selectedProjects.push(entry.id);
          entry.bullets.filter((b) => b.selected).forEach((b) => selectedBullets.push(b.id));
        }
      });

      const selectedLeadership = resume.leadership
        .filter((l) => l.selected)
        .map((l) => l.id);

      const selections = {
        skill_ids: selectedSkills,
        experience_ids: selectedExperience,
        bullet_ids: selectedBullets,
        project_ids: selectedProjects,
        leadership_ids: selectedLeadership,
      };

      const blob = await generatePdf({ selections });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${resume.contact.name.replace(/\s+/g, '_')}_Resume.pdf`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to generate PDF';
      dispatch({ type: 'SET_ERROR', payload: message });
    }
  };

  return (
    <div className="bg-white border-r h-screen sticky top-0 overflow-y-auto">
      <div className="p-6 space-y-6">
        <JobInput />

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-sm text-red-700">{error}</p>
          </div>
        )}

        {jobAnalysis && (
          <div className="space-y-4">
            <KeywordBadges keywords={jobAnalysis.keywords} />

            <div className="pt-4 border-t space-y-2">
              <Button
                onClick={handleApplySuggestions}
                variant="primary"
                className="w-full"
              >
                Apply Suggestions
              </Button>
            </div>
          </div>
        )}

        <div className="pt-4 border-t space-y-2">
          <h3 className="text-sm font-semibold text-gray-700 mb-2">Controls</h3>
          <Button onClick={handleSelectAll} variant="secondary" className="w-full">
            Select All
          </Button>
          <Button onClick={handleDeselectAll} variant="secondary" className="w-full">
            Deselect All
          </Button>
        </div>

        <div className="pt-4 border-t">
          <Button
            onClick={handleGeneratePdf}
            disabled={!resume}
            className="w-full"
          >
            Generate PDF
          </Button>
        </div>
      </div>
    </div>
  );
};
