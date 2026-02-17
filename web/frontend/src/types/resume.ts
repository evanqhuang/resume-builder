export interface SkillItem {
  name: string;
  tags: string[];
  selected: boolean;
}

export interface SkillCategory {
  category: string;
  items: SkillItem[];
}

export interface Bullet {
  id: string;
  text: string;
  tags: string[];
  selected: boolean;
  relevanceScore?: number;
}

export interface ExperienceEntry {
  id: string;
  title: string;
  company: string;
  location: string;
  start_date: string;
  end_date: string | null;
  tags: string[];
  bullets: Bullet[];
  selected: boolean;
}

export interface ProjectEntry {
  id: string;
  title: string;
  technologies: string;
  github?: string;
  tags: string[];
  bullets: Bullet[];
  selected: boolean;
}

export interface EducationEntry {
  institution: string;
  location: string;
  degree: string;
  focus?: string;
  minor?: string;
  gpa?: string;
  honors?: string;
  program?: string;
}

export interface LeadershipEntry {
  id: string;
  text: string;
  tags: string[];
  selected: boolean;
}

export interface ContactInfo {
  name: string;
  email: string;
  phone: string;
  linkedin?: string;
  github?: string;
}

export interface Resume {
  contact: ContactInfo;
  education: EducationEntry;
  skills: SkillCategory[];
  experience: ExperienceEntry[];
  projects: ProjectEntry[];
  leadership: LeadershipEntry[];
}

export interface JobAnalysisResponse {
  keywords: string[];
  scores: Record<string, number>;
  suggested_items: string[];
}
