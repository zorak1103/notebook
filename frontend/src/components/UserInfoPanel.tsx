import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { fetchWhoAmI } from '../api/client';
import type { UserInfo } from '../api/types';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './UserInfoPanel.css';

function UserInfoPanel(): React.JSX.Element {
  const { t } = useTranslation();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null);

  useEffect(() => {
    loadUserInfo();
  }, []);

  const loadUserInfo = async () => {
    try {
      setLoading(true);
      setError(null);
      const info = await fetchWhoAmI();
      setUserInfo(info);
    } catch (err) {
      setError(err instanceof Error ? err.message : t('info.loadError'));
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <LoadingSpinner />;
  }

  return (
    <div className="user-info-panel">
      <h1>{t('info.title')}</h1>

      {error && <ErrorMessage message={error} />}

      {userInfo && (
        <section className="info-section">
          <h2>{t('info.sectionTailscale')}</h2>

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
            <span className="info-label">{t('info.displayName')}:</span>
            <span className="info-value">{userInfo.displayName}</span>
          </div>

          <div className="info-row">
            <span className="info-label">{t('info.loginName')}:</span>
            <span className="info-value">{userInfo.loginName}</span>
          </div>

          <div className="info-row">
            <span className="info-label">{t('info.nodeName')}:</span>
            <span className="info-value">{userInfo.nodeName}</span>
          </div>

          <div className="info-row">
            <span className="info-label">{t('info.nodeID')}:</span>
            <span className="info-value code">{userInfo.nodeID}</span>
          </div>
        </section>
      )}
    </div>
  );
}

export default UserInfoPanel;
