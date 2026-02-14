// UserInfo represents the Tailscale user information returned by the API
export interface UserInfo {
  displayName: string;
  loginName: string;
  profilePicURL: string;
  nodeName: string;
  nodeID: string;
}

// Meeting represents a meeting record
export interface Meeting {
  id: number;
  created_by: string;
  subject: string;
  meeting_date: string;
  start_time: string;
  end_time: string | null;
  participants: string | null;
  summary: string | null;
  keywords: string | null;
  created_at: string;
  updated_at: string;
}

// CreateMeetingRequest represents the request body for creating a meeting
export interface CreateMeetingRequest {
  subject: string;
  meeting_date: string;
  start_time: string;
  end_time?: string | null;
  participants?: string | null;
  summary?: string | null;
  keywords?: string | null;
}

// UpdateMeetingRequest represents the request body for updating a meeting
export interface UpdateMeetingRequest extends CreateMeetingRequest {
  id: number;
}
