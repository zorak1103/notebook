// UserInfo represents the Tailscale user information returned by the API
export interface UserInfo {
  displayName: string;
  loginName: string;
  profilePicURL: string;
  nodeName: string;
  nodeID: string;
}
