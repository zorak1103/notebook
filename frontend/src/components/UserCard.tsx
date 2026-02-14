import type { UserInfo } from '../api/types';
import './UserCard.css';

interface UserCardProps {
  user: UserInfo;
}

export function UserCard({ user }: UserCardProps) {
  return (
    <div className="user-card">
      <img
        src={user.profilePicURL}
        alt={`${user.displayName}'s profile`}
        className="user-avatar"
      />
      <div className="user-info">
        <h2 className="user-name">{user.displayName}</h2>
        <p className="user-detail">
          <span className="label">Email:</span> {user.loginName}
        </p>
        <p className="user-detail">
          <span className="label">Node:</span> {user.nodeName}
        </p>
        <p className="user-detail">
          <span className="label">Node ID:</span> {user.nodeID}
        </p>
      </div>
    </div>
  );
}
