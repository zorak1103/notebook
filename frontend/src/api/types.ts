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

// Note represents a note record
export interface Note {
  id: number;
  meeting_id: number;
  note_number: number;
  content: string;
  created_at: string;
  updated_at: string;
}

// CreateNoteRequest represents the request body for creating a note
export interface CreateNoteRequest {
  meeting_id: number;
  content: string;
}

// UpdateNoteRequest represents the request body for updating a note
export interface UpdateNoteRequest {
  content: string;
}

// Config represents the application configuration
export interface Config {
  llm_provider_url: string;
  llm_api_key: string;
  llm_model: string;
  language: string;
  llm_prompt_summary: string;
  llm_prompt_enhance: string;
}

// ConfigUpdateRequest represents the request body for updating configuration
export interface ConfigUpdateRequest {
  llm_provider_url: string;
  llm_api_key: string;
  llm_model: string;
  language?: string;
  llm_prompt_summary: string;
  llm_prompt_enhance: string;
}
