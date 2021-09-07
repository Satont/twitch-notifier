import { Telegraf, Context, Telegram as TG } from 'telegraf'
import { Update } from 'telegraf/typings/telegram-types'
import { Telegram } from '../../src/services/telegram'

const createRandomNumber = () => Math.floor(Math.random() * (123456789 - 1) + 1);

export const createTelegramClient = (token = 'fake-token') => {
  const telegraf = new Telegraf(token);

  return new Telegram(telegraf)
}

export const createContext = ({ update, telegram }: { update: Update, telegram: TG }) => {
  return new Context(update, telegram)
}

export const createUpdate = ({ text }: { text?: string }) => {
  const from = {
    id: 1,
    is_bot: false,
    first_name: 'Satont',
  }
  const chat =  {
    type: 'private',
    id: from.id,
    first_name: from.first_name,
  }

  return {
    update_id: createRandomNumber(),
    message: {
      chat,
      from,
      text,
      message_id: createRandomNumber(),
      date: Date.now(),
    },
    chat,
  };
}