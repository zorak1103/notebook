import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { getConfig, updateConfig } from '../api/client';
import type { ConfigUpdateRequest } from '../api/types';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import { LanguageSwitcher } from './LanguageSwitcher';
import './ConfigPanel.css';

function ConfigPanel(): React.JSX.Element {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const [originalKey, setOriginalKey] = useState<string>('');
  const [formData, setFormData] = useState<ConfigUpdateRequest>({
    llm_provider_url: '',
    llm_api_key: '',
    llm_model: '',
    llm_prompt_summary: '',
    llm_prompt_enhance: '',
  });

  useEffect(() => {
    const loadConfig = async () => {
      try {
        setLoading(true);
        setError(null);
        const config = await getConfig();
        setFormData({
          llm_provider_url: config.llm_provider_url || '',
          llm_api_key: config.llm_api_key || '',
          llm_model: config.llm_model || '',
          llm_prompt_summary: config.llm_prompt_summary || '',
          llm_prompt_enhance: config.llm_prompt_enhance || '',
        });
        setOriginalKey(config.llm_api_key || '');
      } catch (err) {
        setError(err instanceof Error ? err.message : t('config.loadError'));
      } finally {
        setLoading(false);
      }
    };
    loadConfig();
  }, [t]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);
    setSuccess(false);

    try {
      // If API key hasn't changed (still masked), send empty string
      // LanguageSwitcher handles language persistence independently
      const dataToSend = {
        ...formData,
        llm_api_key: formData.llm_api_key === originalKey ? '' : formData.llm_api_key,
        language: '',
      };

      const result = await updateConfig(dataToSend);
      setFormData({
        llm_provider_url: result.llm_provider_url || '',
        llm_api_key: result.llm_api_key || '',
        llm_model: result.llm_model || '',
        llm_prompt_summary: result.llm_prompt_summary || '',
        llm_prompt_enhance: result.llm_prompt_enhance || '',
      });
      setOriginalKey(result.llm_api_key || '');
      setSuccess(true);

      // Auto-dismiss success message after 3 seconds
      setTimeout(() => {
        setSuccess(false);
      }, 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : t('config.saveError'));
    } finally {
      setSaving(false);
    }
  };

  const handleChange = (field: keyof ConfigUpdateRequest, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="config-panel page-panel">
      <h1 className="page-heading">{t('config.title')}</h1>

      {error && <ErrorMessage message={error} />}

      {success && (
        <div className="config-success">
          {t('config.saveSuccess')}
        </div>
      )}

      <section className="card-section">
        <h2 className="section-heading">{t('config.sectionLanguage')}</h2>
        <div className="form-group">
          <label>{t('config.language')}</label>
          <LanguageSwitcher />
          <small className="hint">{t('config.languageHint')}</small>
        </div>
      </section>

      <form onSubmit={handleSubmit}>
        <section className="card-section">
          <h2 className="section-heading">{t('config.sectionLlm')}</h2>

          <div className="form-group">
            <label htmlFor="provider-url">{t('config.providerUrl')}</label>
            <input
              type="url"
              id="provider-url"
              value={formData.llm_provider_url}
              onChange={(e) => handleChange('llm_provider_url', e.target.value)}
              placeholder={t('config.providerUrlPlaceholder')}
            />
            <small className="hint">{t('config.providerUrlHint')}</small>
          </div>

          <div className="form-group">
            <label htmlFor="api-key">{t('config.apiKey')}</label>
            <input
              type="password"
              id="api-key"
              value={formData.llm_api_key}
              onChange={(e) => handleChange('llm_api_key', e.target.value)}
              placeholder={t('config.apiKeyPlaceholder')}
            />
            <small className="hint">{t('config.apiKeyHint')}</small>
          </div>

          <div className="form-group">
            <label htmlFor="model">{t('config.model')}</label>
            <input
              type="text"
              id="model"
              value={formData.llm_model}
              onChange={(e) => handleChange('llm_model', e.target.value)}
              placeholder={t('config.modelPlaceholder')}
            />
            <small className="hint">{t('config.modelHint')}</small>
          </div>
        </section>

        <section className="card-section">
          <h2 className="section-heading">{t('config.sectionPrompts')}</h2>

          <div className="form-group">
            <label htmlFor="prompt-summary">{t('config.promptSummary')}</label>
            <textarea
              id="prompt-summary"
              value={formData.llm_prompt_summary}
              onChange={(e) => handleChange('llm_prompt_summary', e.target.value)}
              rows={6}
              placeholder={t('config.promptSummaryPlaceholder')}
            />
            <small className="hint">{t('config.promptSummaryHint')}</small>
          </div>

          <div className="form-group">
            <label htmlFor="prompt-enhance">{t('config.promptEnhance')}</label>
            <textarea
              id="prompt-enhance"
              value={formData.llm_prompt_enhance}
              onChange={(e) => handleChange('llm_prompt_enhance', e.target.value)}
              rows={6}
              placeholder={t('config.promptEnhancePlaceholder')}
            />
            <small className="hint">{t('config.promptEnhanceHint')}</small>
          </div>
        </section>

        <div className="form-actions">
          <button type="submit" className="btn btn-submit" disabled={saving}>
            {saving ? t('config.saving') : t('config.save')}
          </button>
        </div>
      </form>
    </div>
  );
}

export default ConfigPanel;
