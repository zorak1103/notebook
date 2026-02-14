import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { WelcomePage } from './components/WelcomePage';
import { LanguageSwitcher } from './components/LanguageSwitcher';
import { MeetingList } from './components/MeetingList';
import { MeetingForm } from './components/MeetingForm';
import './App.css';

type View = 'welcome' | 'list' | 'create' | 'edit';

function App() {
  const { t } = useTranslation();
  const [view, setView] = useState<View>('welcome');
  const [editingId, setEditingId] = useState<number | undefined>();

  const handleNewMeeting = () => {
    setEditingId(undefined);
    setView('create');
  };

  const handleMeetingList = () => {
    setView('list');
  };

  const handleEdit = (id: number) => {
    setEditingId(id);
    setView('edit');
  };

  const handleFormSuccess = () => {
    setView('list');
  };

  const handleFormCancel = () => {
    setView('list');
  };

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h1 className="app-title">{t('app.title')}</h1>
        </div>

        <nav className="sidebar-nav">
          <button
            className={`nav-item ${view === 'create' ? 'nav-item--active' : ''}`}
            onClick={handleNewMeeting}
          >
            {t('navigation.newMeeting')}
          </button>
          <button
            className={`nav-item ${view === 'list' || view === 'edit' ? 'nav-item--active' : ''}`}
            onClick={handleMeetingList}
          >
            {t('navigation.meetingList')}
          </button>
          <button className="nav-item" disabled>
            {t('navigation.search')}
          </button>
          <button className="nav-item" disabled>
            {t('navigation.configuration')}
          </button>
        </nav>

        <div className="sidebar-footer">
          <LanguageSwitcher />
        </div>
      </aside>

      <main className="main-content">
        {view === 'welcome' && <WelcomePage />}
        {view === 'list' && <MeetingList onEdit={handleEdit} />}
        {view === 'create' && (
          <MeetingForm onSuccess={handleFormSuccess} onCancel={handleFormCancel} />
        )}
        {view === 'edit' && (
          <MeetingForm
            meetingId={editingId}
            onSuccess={handleFormSuccess}
            onCancel={handleFormCancel}
          />
        )}
      </main>
    </div>
  );
}

export default App;
