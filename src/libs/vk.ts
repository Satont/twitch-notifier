import { VK } from 'vk-io'
import { info, error } from '../helpers/logs'
import { User } from '../models/User'
import follow from '../commands/follow'
import unfollow from '../commands/unfollow'
import follows from '../commands/follows'
import live from '../commands/live'
import { config } from '../helpers/config'
import { chunk, isBoolean } from 'lodash'

const bot = new VK({
  token: config.vk.token
})
const service = 'vk'

bot.updates.on('message', async (ctx, next) => {
  if (ctx.senderId !== -187752469) {
    info(`Vk | New message from: ${ctx.senderId}, message: ${ctx.text}`)
  }
  const [user] = await User.findOrCreate({ where: { id: ctx.senderId, service: 'vk' }, defaults: { follows: [], service: 'vk' } })
  ctx.dbUser = user
  return next()
})

bot.updates.hear(value => (value.startsWith('!подписка')), async (ctx) => {
  const channel: string = ctx.text.split(' ')[1]
  if (!channel) {
    return ctx.reply('Вы должны указать на кого подписаться.')
  }
  try {
    const followed = await follow({ service, userId: ctx.senderId, channel })
    if (followed) {
      ctx.reply(`Вы успешно подписались на ${channel}.`)
    }
  } catch (e) {
    error(e)
    ctx.reply(e.message)
  }
})

bot.updates.hear(value => (value.startsWith('!отписка')), async (ctx) => {
  const channel: string = ctx.text.split(' ')[1]
  if (!channel) {
    return ctx.reply('Вы должны указать от кого отписаться.')
  }
  try {
    const unfollowed = await unfollow({ service, userId: ctx.senderId, channel })
    if (!unfollowed) {
      ctx.reply(`Вы не подписаны на ${channel}.`)
    } else {
      ctx.reply(`Вы успешно отписались от ${channel}.`)
    }
  } catch (e) {
    error(e)
    ctx.reply(e.message)
  }
})

bot.updates.hear(value => (value.startsWith('!подписки')), async (ctx) => {
  const followed = await follows({ userId: ctx.senderId, service })
  if (isBoolean(followed)) {
    ctx.reply(`Вы ни на кого не подписаны.`)
  } else {
    ctx.reply(`Вы подписаны на ${followed.join(', ')}.`)
  }
})

bot.updates.hear(value => (value.startsWith('!онлайн')), async (ctx) => {
  const channels = await live({ userId: ctx.senderId, service })
  if (isBoolean(channels)) {
    ctx.reply(`Сейчас нет ни одного канала онлайн из ваших подписок.`)
  } else {
    const links = channels.map((o) => `https://twitch.tv/${o}`)
    ctx.reply(`Сейчас онлайн: \n${links.join('\n')}`)
  }
})

bot.updates.hear(value => (value.startsWith('!команды')), async (ctx) => {
  ctx.reply(`
  На данный момент доступны следующие команды: 
  !подписка username
  !отписка username
  !подписки 
  !команды
  !онлайн
  `)
})

export function say(userId: number | number[], message: string, attachment?: string) {
  info(`Send message to ${Array.isArray(userId) ? userId.join(', ') : userId}. message: ${message}`)
  const targets = Array.isArray(userId) ? userId : [userId]
  const chunks = chunk(targets, 100)
  for (const chunk of chunks) {
    bot.api.messages.send({
      random_id: Math.random() * (1000000000 - 9) + 10,
      user_ids: chunk,
      message,
      dont_parse_links: true,
      attachment: attachment
    })
  }
} 

bot.updates.start().then(() => info('VK bot connected.')).catch(console.error)

export { bot }
