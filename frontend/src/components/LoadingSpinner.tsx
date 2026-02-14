import './LoadingSpinner.css';

export function LoadingSpinner() {
  return (
    <div className="loading-spinner" role="status" aria-label="Loading">
      <div className="spinner"></div>
    </div>
  );
}
