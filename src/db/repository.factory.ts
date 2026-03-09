import type { 
  IChatRepository,
  IChannelRepository,
  IFollowRepository,
  IStreamRepository
} from './repositories/interfaces';
import {
  ChatDrizzleRepository,
  ChannelDrizzleRepository,
  FollowDrizzleRepository,
  StreamDrizzleRepository
} from './repositories/drizzle';
import type { IDatabaseConnection } from './connection';

export interface IRepositoryFactory {
  createChatRepository(): IChatRepository;
  createChannelRepository(): IChannelRepository;
  createFollowRepository(): IFollowRepository;
  createStreamRepository(): IStreamRepository;
}

/**
 * Factory for creating Drizzle-based repositories
 * Works with any Drizzle-compatible database (D1, PostgreSQL, etc.)
 */
export class DrizzleRepositoryFactory implements IRepositoryFactory {
  constructor(private connection: IDatabaseConnection) {}

  createChatRepository(): IChatRepository {
    return new ChatDrizzleRepository(this.connection.getClient());
  }

  createChannelRepository(): IChannelRepository {
    return new ChannelDrizzleRepository(this.connection.getClient());
  }

  createFollowRepository(): IFollowRepository {
    return new FollowDrizzleRepository(this.connection.getClient());
  }

  createStreamRepository(): IStreamRepository {
    return new StreamDrizzleRepository(this.connection.getClient());
  }
}
