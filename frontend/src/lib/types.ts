export type AuthResponse = {
  token: string;
};

export type BackendLinkItem = {
  id: number;
  title: string;
  url: string;
  description?: string | null;
  created_at: string;
  updated_at: string;
  click_count: number;
  folder_id?: number | null;
};

export type BackendFolder = {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
};

export type LinkItem = {
  id: string;
  name: string;
  url: string;
  description?: string | null;
  createdAt?: string;
  updatedAt?: string;
  clicks?: number;
  folderId?: string | null;
};

export type FolderItem = {
  id: string;
  name: string;
  createdAt?: string;
  updatedAt?: string;
};

export type CreateLinkPayload = {
  name: string;
  url: string;
  description?: string;
  folderId?: number | null;
};

export type UpdateLinkPayload = Partial<CreateLinkPayload>;

export type CreateFolderPayload = {
  name: string;
};

export type UpdateFolderPayload = Partial<CreateFolderPayload>;
