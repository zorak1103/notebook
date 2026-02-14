import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { WelcomePage } from './components/WelcomePage';
import { LanguageSwitcher } from './components/LanguageSwitcher';
import { MeetingList } from './components/MeetingList';
import { MeetingForm } from './components/MeetingForm';
import { MeetingDetail } from './components/MeetingDetail';
import { SearchPanel } from './components/SearchPanel';
import ConfigPanel from './components/ConfigPanel';
import './App.css';

type View = 'welcome' | 'list' | 'create' | 'edit' | 'detail' | 'search' | 'config';

function App() {
  const { t } = useTranslation();
  const [view, setView] = useState<View>('welcome');
  const [selectedId, setSelectedId] = useState<number | undefined>();

  const handleNewMeeting = () => {
    setSelectedId(undefined);
    setView('create');
  };

  const handleMeetingList = () => {
    setView('list');
  };

  const handleSearch = () => {
    setView('search');
  };

  const handleConfig = () => {
    setView('config');
  };

  const handleSearchSelect = (id: number) => {
    setSelectedId(id);
    setView('detail');
  };

  const handleViewDetail = (id: number) => {
    setSelectedId(id);
    setView('detail');
  };

  const handleEdit = (id: number) => {
    setSelectedId(id);
    setView('edit');
  };

  const handleEditFromDetail = () => {
    setView('edit');
  };

  const handleFormSuccess = () => {
    setView('list');
  };

  const handleFormCancel = () => {
    setView('list');
  };

  const handleDetailBack = () => {
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
            className={`nav-item ${view === 'list' || view === 'edit' || view === 'detail' ? 'nav-item--active' : ''}`}
            onClick={handleMeetingList}
          >
            {t('navigation.meetingList')}
          </button>
          <button
            className={`nav-item ${view === 'search' ? 'nav-item--active' : ''}`}
            onClick={handleSearch}
          >
            {t('navigation.search')}
          </button>
          <button
            className={`nav-item ${view === 'config' ? 'nav-item--active' : ''}`}
            onClick={handleConfig}
          >
            {t('navigation.configuration')}
          </button>
        </nav>

        <div className="sidebar-footer">
          <LanguageSwitcher />
        </div>
      </aside>

      <main className="main-content">
        {view === 'welcome' && <WelcomePage />}
        {view === 'list' && <MeetingList onEdit={handleEdit} onViewDetail={handleViewDetail} />}
        {view === 'create' && (
          <MeetingForm onSuccess={handleFormSuccess} onCancel={handleFormCancel} />
        )}
        {view === 'edit' && selectedId && (
          <MeetingForm
            meetingId={selectedId}
            onSuccess={handleFormSuccess}
            onCancel={handleFormCancel}
          />
        )}
        {view === 'detail' && selectedId && (
          <MeetingDetail
            meetingId={selectedId}
            onBack={handleDetailBack}
            onEdit={handleEditFromDetail}
          />
        )}
        {view === 'search' && <SearchPanel onSelectMeeting={handleSearchSelect} />}
        {view === 'config' && <ConfigPanel />}
      </main>
    </div>
  );
}

export default App;
