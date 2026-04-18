import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchWhoAmI, fetchVersion } from '../api/client';
import type { UserInfo, VersionInfo } from '../api/types';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './UserInfoPanel.css';

function UserInfoPanel(): React.JSX.Element {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);
  const [versionInfo, setVersionInfo] = useState<VersionInfo | null>(null);

  useEffect(() => {
    let cancelled = false;
    Promise.all([fetchWhoAmI(), fetchVersion()])
      .then(([info, ver]) => {
        if (!cancelled) {
          setUserInfo(info);
          setVersionInfo(ver);
          setError(null);
        }
      })
      .catch((err) => {
        if (!cancelled) setError(err instanceof Error ? err.message : t('info.loadError'));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, [t]);

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="user-info-panel page-panel">
      <h1 className="page-heading">{t('info.title')}</h1>

      {error && <ErrorMessage message={error} />}

      {userInfo && (
        <section className="card-section">
          <h2 className="section-heading">{t('info.sectionTailscale')}</h2>

          {userInfo.profilePicURL && (
            <div className="profile-pic-container">
              <img
                src={userInfo.profilePicURL}
                alt={userInfo.displayName}
                className="profile-pic"
              />
            </div>
          )}

          <div className="info-row">
            <span className="data-label">{t('info.displayName')}:</span>
            <span className="info-value">{userInfo.displayName}</span>
          </div>

          <div className="info-row">
            <span className="data-label">{t('info.loginName')}:</span>
            <span className="info-value">{userInfo.loginName}</span>
          </div>

          <div className="info-row">
            <span className="data-label">{t('info.nodeName')}:</span>
            <span className="info-value">{userInfo.nodeName}</span>
          </div>

          <div className="info-row">
            <span className="data-label">{t('info.nodeID')}:</span>
            <span className="info-value">{userInfo.nodeID}</span>
          </div>
        </section>
      )}

      {versionInfo && (
        <section className="card-section">
          <h2 className="section-heading">{t('info.sectionApplication')}</h2>

          <div className="info-row">
            <span className="data-label">{t('info.version')}:</span>
            <span className="info-value">{versionInfo.version}</span>
          </div>

          <div className="info-row">
            <span className="data-label">{t('info.commit')}:</span>
            <span className="info-value">{versionInfo.commit.slice(0, 7)}</span>
          </div>

          <div className="info-row">
            <span className="data-label">{t('info.buildDate')}:</span>
            <span className="info-value">{versionInfo.date}</span>
          </div>
        </section>
      )}
    </div>
  );
}

export default UserInfoPanel;
