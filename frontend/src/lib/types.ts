export type AuthResponse = {
  token: string;
};

export type LinkItem = {
  id: string;
  name: string;
  url: string;
  description?: string | null;
  createdAt?: string;
  updatedAt?: string;
  clicks?: number;
};

export type CreateLinkPayload = {
  name: string;
  url: string;
  description?: string;
};

export type UpdateLinkPayload = Partial<CreateLinkPayload>;
