// Domain models - business logic representations
// These are separate from database schema to allow flexibility

export type ChatService = 'telegram';
export type ChannelService = 'twitch';
export type SupportedLanguage = 'en' | 'ru' | 'uk';

export class Chat {
  id: string;
  chatId: string;
  service: ChatService;
  settings?: ChatSettings;
  follows?: Follow[];

  constructor(data: {
    id: string;
    chatId: string;
    service: ChatService;
    settings?: ChatSettings;
    follows?: Follow[];
  }) {
    this.id = data.id;
    this.chatId = data.chatId;
    this.service = data.service;
    this.settings = data.settings;
    this.follows = data.follows;
  }
}

export class ChatSettings {
  id: string;
  chatId: string;
  gameChangeNotification: boolean;
  titleChangeNotification: boolean;
  gameAndTitleChangeNotification: boolean;
  offlineNotification: boolean;
  imageInNotification: boolean;
  language: SupportedLanguage;

  constructor(data: {
    id: string;
    chatId: string;
    gameChangeNotification: boolean;
    titleChangeNotification: boolean;
    gameAndTitleChangeNotification: boolean;
    offlineNotification: boolean;
    imageInNotification: boolean;
    language: SupportedLanguage;
  }) {
    this.id = data.id;
    this.chatId = data.chatId;
    this.gameChangeNotification = data.gameChangeNotification;
    this.titleChangeNotification = data.titleChangeNotification;
    this.gameAndTitleChangeNotification = data.gameAndTitleChangeNotification;
    this.offlineNotification = data.offlineNotification;
    this.imageInNotification = data.imageInNotification;
    this.language = data.language;
  }
}

export class Channel {
  id: string;
  channelId: string;
  service: ChannelService;
  isLive: boolean;
  title?: string;
  category?: string;
  updatedAt?: Date;
  follows?: Follow[];
  streams?: Stream[];

  constructor(data: {
    id: string;
    channelId: string;
    service: ChannelService;
    isLive: boolean;
    title?: string;
    category?: string;
    updatedAt?: Date;
    follows?: Follow[];
    streams?: Stream[];
  }) {
    this.id = data.id;
    this.channelId = data.channelId;
    this.service = data.service;
    this.isLive = data.isLive;
    this.title = data.title;
    this.category = data.category;
    this.updatedAt = data.updatedAt;
    this.follows = data.follows;
    this.streams = data.streams;
  }
}

export class Follow {
  id: string;
  channelId: string;
  chatId: string;
  channel?: Channel;
  chat?: Chat;

  constructor(data: {
    id: string;
    channelId: string;
    chatId: string;
    channel?: Channel;
    chat?: Chat;
  }) {
    this.id = data.id;
    this.channelId = data.channelId;
    this.chatId = data.chatId;
    this.channel = data.channel;
    this.chat = data.chat;
  }
}

export class Stream {
  id: string;
  channelId: string;
  isLive: boolean;
  title?: string;
  category?: string;
  titles: string[];
  categories: string[];
  startedAt: Date;
  updatedAt?: Date;
  endedAt?: Date;

  constructor(data: {
    id: string;
    channelId: string;
    isLive: boolean;
    title?: string;
    category?: string;
    titles: string[];
    categories: string[];
    startedAt: Date;
    updatedAt?: Date;
    endedAt?: Date;
  }) {
    this.id = data.id;
    this.channelId = data.channelId;
    this.isLive = data.isLive;
    this.title = data.title;
    this.category = data.category;
    this.titles = data.titles;
    this.categories = data.categories;
    this.startedAt = data.startedAt;
    this.updatedAt = data.updatedAt;
    this.endedAt = data.endedAt;
  }
}

// Errors
export class FollowAlreadyExistsError extends Error {
  constructor() {
    super('Follow already exists');
    this.name = 'FollowAlreadyExistsError';
  }
}

export class FollowNotFoundError extends Error {
  constructor() {
    super('Follow not found');
    this.name = 'FollowNotFoundError';
  }
}

export class ChannelNotFoundError extends Error {
  constructor() {
    super('Channel not found');
    this.name = 'ChannelNotFoundError';
  }
}
