import { useEffect } from 'react';
import { ResumeProvider } from './contexts/ResumeContext';
import { useResume } from './hooks/useResume';
import { Header } from './components/layout/Header';
import { Sidebar } from './components/layout/Sidebar';
import { ResumeEditor } from './components/resume/ResumeEditor';
import { fetchResume } from './services/api';

const AppContent = () => {
  const { dispatch } = useResume();

  useEffect(() => {
    const loadResume = async () => {
      dispatch({ type: 'SET_LOADING', payload: true });
      try {
        const data = await fetchResume();
        dispatch({ type: 'SET_RESUME', payload: data });
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to load resume';
        dispatch({ type: 'SET_ERROR', payload: message });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    };

    loadResume();
  }, [dispatch]);

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <div className="flex h-[calc(100vh-72px)]">
        <div className="w-1/3 min-w-[400px] max-w-[500px]">
          <Sidebar />
        </div>
        <div className="flex-1 overflow-y-auto">
          <main className="max-w-5xl mx-auto px-6 py-8">
            <ResumeEditor />
          </main>
        </div>
      </div>
    </div>
  );
};

function App() {
  return (
    <ResumeProvider>
      <AppContent />
    </ResumeProvider>
  );
}

export default App;
