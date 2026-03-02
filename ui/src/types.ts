export interface Note {
  id: number;
  container_name: string;
  container_id: string;
  compose_project: string;
  compose_service: string;
  note_content: string;
  pinned: boolean;
  tags: string[];
  created_at: string;
  updated_at: string;
}

export interface CreateNoteRequest {
  container_name: string;
  container_id?: string;
  compose_project?: string;
  compose_service?: string;
  note_content?: string;
  tags?: string[];
}

export interface UpdateNoteRequest {
  note_content?: string;
  pinned?: boolean;
  tags?: string[];
  container_id?: string;
}

export interface ContainerInfo {
  id: string;
  name: string;
  image: string;
  state: string;
  status: string;
  composeProject?: string;
  composeService?: string;
}
