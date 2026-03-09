import { ApiClient } from '@twurple/api';
import { AppTokenAuthProvider } from '@twurple/auth';
import type { Env } from '../types/env';

export class TwitchService {
  private apiClient: ApiClient;
  private authProvider: AppTokenAuthProvider;

  constructor(env: Env) {
    this.authProvider = new AppTokenAuthProvider(
      env.TWITCH_CLIENT_ID,
      env.TWITCH_CLIENT_SECRET
    );
    this.apiClient = new ApiClient({ authProvider: this.authProvider });
  }

  async getUserByLogin(login: string) {
    try {
      return await this.apiClient.users.getUserByName(login);
    } catch (error) {
      return null;
    }
  }

  async getUserById(id: string) {
    try {
      return await this.apiClient.users.getUserById(id);
    } catch (error) {
      return null;
    }
  }

  async getStreamByUserId(userId: string) {
    try {
      return await this.apiClient.streams.getStreamByUserId(userId);
    } catch (error) {
      return null;
    }
  }

  async getGameById(gameId: string) {
    try {
      return await this.apiClient.games.getGameById(gameId);
    } catch (error) {
      return null;
    }
  }

  getApiClient() {
    return this.apiClient;
  }

  getAuthProvider() {
    return this.authProvider;
  }
}
