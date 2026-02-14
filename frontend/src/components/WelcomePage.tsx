import { useWhoAmI } from '../hooks/useWhoAmI';
import { UserCard } from './UserCard';
import { LoadingSpinner } from './LoadingSpinner';
import { ErrorMessage } from './ErrorMessage';
import './WelcomePage.css';

export function WelcomePage() {
  const { user, loading, error, fetchUser } = useWhoAmI();

  return (
    <div className="welcome-page">
      <div className="welcome-container">
        <h1 className="welcome-title">Welcome to Notebook</h1>
        <p className="welcome-subtitle">
          A Tailscale-powered application for managing meetings and notes
        </p>

        {!user && !loading && !error && (
          <button onClick={fetchUser} className="who-am-i-button">
            Who Am I?
          </button>
        )}

        {loading && <LoadingSpinner />}

        {error && <ErrorMessage message={error} onRetry={fetchUser} />}

        {user && <UserCard user={user} />}
      </div>
    </div>
  );
}
