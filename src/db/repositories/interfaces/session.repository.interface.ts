export interface ISessionRepository {
  get(key: string): Promise<string | undefined>;
  set(key: string, value: string, expiresAt?: number): Promise<void>;
  delete(key: string): Promise<void>;
  cleanup(): Promise<void>; // Remove expired sessions
}
