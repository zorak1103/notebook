import { useTranslation } from 'react-i18next';
import { updateConfig } from '../api/client';
import './LanguageSwitcher.css';

export function LanguageSwitcher() {
  const { i18n } = useTranslation();

  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const value = event.target.value;
    i18n.changeLanguage(value);

    // Persist language preference to backend
    updateConfig({
      language: value,
      llm_provider_url: '',
      llm_api_key: '',
      llm_model: ''
    }).catch(() => {
      // Silently handle save failures - language still changes locally
    });
  };

  return (
    <div className="language-switcher">
      <select
        value={i18n.language}
        onChange={handleChange}
        className="language-select"
        aria-label="Select language"
      >
        <option value="en">English</option>
        <option value="de">Deutsch</option>
        <option value="fr">Français</option>
        <option value="es">Español</option>
      </select>
    </div>
  );
}
