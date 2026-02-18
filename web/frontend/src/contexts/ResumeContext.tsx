import { createContext, useReducer, ReactNode } from 'react';
import type { Resume, JobAnalysisResponse, SectionName } from '../types/resume';

interface ResumeState {
  resume: Resume | null;
  jobAnalysis: JobAnalysisResponse | null;
  isLoading: boolean;
  error: string | null;
}

type ResumeAction =
  | { type: 'SET_RESUME'; payload: Resume }
  | { type: 'TOGGLE_SKILL'; payload: { category: string; skillName: string } }
  | { type: 'TOGGLE_SKILL_CATEGORY'; payload: { category: string; selected: boolean } }
  | { type: 'TOGGLE_BULLET'; payload: { entryId: string; bulletId: string; entryType: 'experience' | 'project' } }
  | { type: 'TOGGLE_EXPERIENCE'; payload: string }
  | { type: 'TOGGLE_PROJECT'; payload: string }
  | { type: 'TOGGLE_LEADERSHIP'; payload: string }
  | { type: 'SELECT_ALL' }
  | { type: 'DESELECT_ALL' }
  | { type: 'SET_JOB_ANALYSIS'; payload: JobAnalysisResponse }
  | { type: 'APPLY_SUGGESTIONS'; payload: { threshold: number } }
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'REORDER_SECTION'; payload: { section: SectionName; order: string[] } };

const initialState: ResumeState = {
  resume: null,
  jobAnalysis: null,
  isLoading: false,
  error: null,
};

function reorderSection<T extends { id: string }>(items: T[], order: string[]): T[] {
  const orderMap = new Map(order.map((id, idx) => [id, idx]));
  return [...items].sort(
    (a, b) => (orderMap.get(a.id) ?? items.length) - (orderMap.get(b.id) ?? items.length)
  );
}

const resumeReducer = (state: ResumeState, action: ResumeAction): ResumeState => {
  switch (action.type) {
    case 'SET_RESUME':
      return { ...state, resume: action.payload, error: null };

    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };

    case 'SET_ERROR':
      return { ...state, error: action.payload, isLoading: false };

    case 'SET_JOB_ANALYSIS':
      return { ...state, jobAnalysis: action.payload };

    case 'TOGGLE_SKILL': {
      if (!state.resume) return state;
      const skills = state.resume.skills.map((category) => {
        if (category.category === action.payload.category) {
          return {
            ...category,
            items: category.items.map((item) =>
              item.name === action.payload.skillName
                ? { ...item, selected: !item.selected }
                : item
            ),
          };
        }
        return category;
      });
      return { ...state, resume: { ...state.resume, skills } };
    }

    case 'TOGGLE_SKILL_CATEGORY': {
      if (!state.resume) return state;
      const skills = state.resume.skills.map((category) => {
        if (category.category === action.payload.category) {
          return {
            ...category,
            items: category.items.map((item) => ({
              ...item,
              selected: action.payload.selected,
            })),
          };
        }
        return category;
      });
      return { ...state, resume: { ...state.resume, skills } };
    }

    case 'TOGGLE_BULLET': {
      if (!state.resume) return state;
      const { entryId, bulletId, entryType } = action.payload;

      if (entryType === 'experience') {
        const experience = state.resume.experience.map((entry) => {
          if (entry.id === entryId) {
            return {
              ...entry,
              bullets: entry.bullets.map((bullet) =>
                bullet.id === bulletId ? { ...bullet, selected: !bullet.selected } : bullet
              ),
            };
          }
          return entry;
        });
        return { ...state, resume: { ...state.resume, experience } };
      } else {
        const projects = state.resume.projects.map((entry) => {
          if (entry.id === entryId) {
            return {
              ...entry,
              bullets: entry.bullets.map((bullet) =>
                bullet.id === bulletId ? { ...bullet, selected: !bullet.selected } : bullet
              ),
            };
          }
          return entry;
        });
        return { ...state, resume: { ...state.resume, projects } };
      }
    }

    case 'TOGGLE_EXPERIENCE': {
      if (!state.resume) return state;
      const experience = state.resume.experience.map((entry) => {
        if (entry.id === action.payload) {
          const newSelected = !entry.selected;
          return {
            ...entry,
            selected: newSelected,
            bullets: entry.bullets.map((bullet) => ({ ...bullet, selected: newSelected })),
          };
        }
        return entry;
      });
      return { ...state, resume: { ...state.resume, experience } };
    }

    case 'TOGGLE_PROJECT': {
      if (!state.resume) return state;
      const projects = state.resume.projects.map((entry) => {
        if (entry.id === action.payload) {
          const newSelected = !entry.selected;
          return {
            ...entry,
            selected: newSelected,
            bullets: entry.bullets.map((bullet) => ({ ...bullet, selected: newSelected })),
          };
        }
        return entry;
      });
      return { ...state, resume: { ...state.resume, projects } };
    }

    case 'TOGGLE_LEADERSHIP': {
      if (!state.resume) return state;
      const leadership = state.resume.leadership.map((entry) =>
        entry.id === action.payload ? { ...entry, selected: !entry.selected } : entry
      );
      return { ...state, resume: { ...state.resume, leadership } };
    }

    case 'SELECT_ALL': {
      if (!state.resume) return state;
      const skills = state.resume.skills.map((category) => ({
        ...category,
        items: category.items.map((item) => ({ ...item, selected: true })),
      }));
      const experience = state.resume.experience.map((entry) => ({
        ...entry,
        selected: true,
        bullets: entry.bullets.map((bullet) => ({ ...bullet, selected: true })),
      }));
      const projects = state.resume.projects.map((entry) => ({
        ...entry,
        selected: true,
        bullets: entry.bullets.map((bullet) => ({ ...bullet, selected: true })),
      }));
      const leadership = state.resume.leadership.map((entry) => ({ ...entry, selected: true }));
      return {
        ...state,
        resume: { ...state.resume, skills, experience, projects, leadership },
      };
    }

    case 'DESELECT_ALL': {
      if (!state.resume) return state;
      const skills = state.resume.skills.map((category) => ({
        ...category,
        items: category.items.map((item) => ({ ...item, selected: false })),
      }));
      const experience = state.resume.experience.map((entry) => ({
        ...entry,
        selected: false,
        bullets: entry.bullets.map((bullet) => ({ ...bullet, selected: false })),
      }));
      const projects = state.resume.projects.map((entry) => ({
        ...entry,
        selected: false,
        bullets: entry.bullets.map((bullet) => ({ ...bullet, selected: false })),
      }));
      const leadership = state.resume.leadership.map((entry) => ({ ...entry, selected: false }));
      return {
        ...state,
        resume: { ...state.resume, skills, experience, projects, leadership },
      };
    }

    case 'APPLY_SUGGESTIONS': {
      if (!state.resume || !state.jobAnalysis) return state;
      const { threshold } = action.payload;
      const { scores } = state.jobAnalysis;

      const skills = state.resume.skills.map((category) => ({
        ...category,
        items: category.items.map((item) => {
          const score = scores[item.name] || 0;
          return { ...item, selected: score >= threshold };
        }),
      }));

      const experience = state.resume.experience.map((entry) => {
        const bullets = entry.bullets.map((bullet) => {
          const score = scores[bullet.id] || 0;
          return { ...bullet, selected: score >= threshold, relevanceScore: score };
        });
        const anySelected = bullets.some((b) => b.selected);
        return { ...entry, bullets, selected: anySelected };
      });

      const projects = state.resume.projects.map((entry) => {
        const bullets = entry.bullets.map((bullet) => {
          const score = scores[bullet.id] || 0;
          return { ...bullet, selected: score >= threshold, relevanceScore: score };
        });
        const anySelected = bullets.some((b) => b.selected);
        return { ...entry, bullets, selected: anySelected };
      });

      const leadership = state.resume.leadership.map((entry) => {
        const score = scores[entry.id] || 0;
        return { ...entry, selected: score >= threshold };
      });

      return {
        ...state,
        resume: { ...state.resume, skills, experience, projects, leadership },
      };
    }

    case 'REORDER_SECTION': {
      if (!state.resume) return state;
      const { section, order } = action.payload;
      return {
        ...state,
        resume: {
          ...state.resume,
          [section]: reorderSection(state.resume[section] as { id: string }[], order),
        },
      };
    }

    default:
      return state;
  }
};

interface ResumeContextType {
  state: ResumeState;
  dispatch: React.Dispatch<ResumeAction>;
}

export const ResumeContext = createContext<ResumeContextType | undefined>(undefined);

export const ResumeProvider = ({ children }: { children: ReactNode }) => {
  const [state, dispatch] = useReducer(resumeReducer, initialState);

  return (
    <ResumeContext.Provider value={{ state, dispatch }}>
      {children}
    </ResumeContext.Provider>
  );
};
