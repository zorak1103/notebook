import { useTranslation } from 'react-i18next';
import { WelcomePage } from './components/WelcomePage';
import { LanguageSwitcher } from './components/LanguageSwitcher';
import './App.css';

function App() {
  const { t } = useTranslation();

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h1 className="app-title">{t('app.title')}</h1>
        </div>

        <nav className="sidebar-nav">
          <button className="nav-item" disabled>
            {t('navigation.newMeeting')}
          </button>
          <button className="nav-item" disabled>
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
        <WelcomePage />
      </main>
    </div>
  );
}

export default App;
