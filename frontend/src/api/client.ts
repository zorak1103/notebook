import type { UserInfo, VersionInfo, Meeting, CreateMeetingRequest, Note, CreateNoteRequest, UpdateNoteRequest, ReorderNoteRequest, Config, ConfigUpdateRequest, EnhanceNoteRequest, EnhanceNoteResponse } from './types';

export async function fetchVersion(): Promise<VersionInfo> {
  return apiGet<VersionInfo>('/api/version');
}

/**
 * Fetches the current user's Tailscale information from the backend
 * @returns Promise resolving to UserInfo
 * @throws Error if the request fails
 */
export async function fetchWhoAmI(): Promise<UserInfo> {
  const response = await fetch('/api/whoami');

  if (!response.ok) {
    throw new Error(`Failed to fetch user info: ${response.status} ${response.statusText}`);
  }

  return response.json();
}

// Generic API helpers

async function parseErrorMessage(response: Response): Promise<string> {
  try {
    const data = await response.json();
    return data.error || `${response.status} ${response.statusText}`;
  } catch {
    return `${response.status} ${response.statusText}`;
  }
}

async function apiGet<T>(url: string): Promise<T> {
  const response = await fetch(url);
  if (!response.ok) {
    const message = await parseErrorMessage(response);
    throw new Error(message);
  }
  return response.json();
}

async function apiPost<T>(url: string, data: unknown): Promise<T> {
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    const message = await parseErrorMessage(response);
    throw new Error(message);
  }
  return response.json();
}

async function apiPut<T>(url: string, data: unknown): Promise<T> {
  const response = await fetch(url, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!response.ok) {
    const message = await parseErrorMessage(response);
    throw new Error(message);
  }
  return response.json();
}

async function apiDelete(url: string): Promise<void> {
  const response = await fetch(url, { method: 'DELETE' });
  if (!response.ok) {
    const message = await parseErrorMessage(response);
    throw new Error(message);
  }
}

// Meeting API functions

export async function fetchMeetings(sort?: string, order?: string): Promise<Meeting[]> {
  const params = new URLSearchParams();
  if (sort) params.append('sort', sort);
  if (order) params.append('order', order);
  const query = params.toString() ? `?${params.toString()}` : '';
  return apiGet<Meeting[]>(`/api/meetings${query}`);
}

export async function fetchMeeting(id: number): Promise<Meeting> {
  return apiGet<Meeting>(`/api/meetings/${id}`);
}

export async function createMeeting(data: CreateMeetingRequest): Promise<Meeting> {
  return apiPost<Meeting>('/api/meetings', data);
}

export async function updateMeeting(id: number, data: CreateMeetingRequest): Promise<Meeting> {
  return apiPut<Meeting>(`/api/meetings/${id}`, data);
}

export async function deleteMeeting(id: number): Promise<void> {
  return apiDelete(`/api/meetings/${id}`);
}

export async function searchMeetings(query: string): Promise<Meeting[]> {
  if (!query.trim()) return [];
  return apiGet<Meeting[]>(`/api/search?q=${encodeURIComponent(query)}`);
}

export async function summarizeMeeting(id: number): Promise<Meeting> {
  return apiPost<Meeting>(`/api/meetings/${id}/summarize`, {});
}

// Note API functions

export async function fetchNotes(meetingId: number): Promise<Note[]> {
  return apiGet<Note[]>(`/api/meetings/${meetingId}/notes`);
}

export async function fetchNote(id: number): Promise<Note> {
  return apiGet<Note>(`/api/notes/${id}`);
}

export async function createNote(data: CreateNoteRequest): Promise<Note> {
  return apiPost<Note>('/api/notes', data);
}

export async function updateNote(id: number, data: UpdateNoteRequest): Promise<Note> {
  return apiPut<Note>(`/api/notes/${id}`, data);
}

export async function deleteNote(id: number): Promise<void> {
  return apiDelete(`/api/notes/${id}`);
}

export async function enhanceNote(id: number, content: string): Promise<EnhanceNoteResponse> {
  const req: EnhanceNoteRequest = { content };
  return apiPost<EnhanceNoteResponse>(`/api/notes/${id}/enhance`, req);
}

export async function reorderNote(id: number, direction: 'up' | 'down'): Promise<Note[]> {
  const req: ReorderNoteRequest = { direction };
  return apiPut<Note[]>(`/api/notes/${id}/reorder`, req);
}

// Config API functions

export async function getConfig(): Promise<Config> {
  return apiGet<Config>('/api/config');
}

export async function updateConfig(data: ConfigUpdateRequest): Promise<Config> {
  return apiPost<Config>('/api/config', data);
}
