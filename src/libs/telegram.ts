import Telegraf from 'telegraf'
import { Stage, session, BaseScene } from 'telegraf'
import { config } from '../helpers/config'
import { info, error } from '../helpers/logs'
import { User } from '../models/User'
import follow from '../commands/follow'
import unfollow from '../commands/unfollow'
import live from '../commands/live'
import follows from '../commands/follows'
import { isBoolean } from 'util'

const service = 'telegram'
const bot = new Telegraf(config.telegram.token)

bot.use(session())
bot.use(async (ctx, next) => {
  if (!ctx.message) return
  info(`Telegram | New message from ${ctx.from.username} [${ctx.from.id}], message: ${ctx.message.text}`)
  const [user] = await User.findOrCreate({ where: { id: ctx.from.id, service }, defaults: { follows: [], service } })
  ctx.userDb = user
  next()
})

const followScene = new BaseScene('follow', {
  ttl: 60
})
followScene.enter((ctx) => ctx.reply('Enter streamer username you want to follow'))
followScene.on('message', async (ctx) => {
  const channel = ctx.message.text
  try {
    const followed = await follow({ userId: ctx.from.id, service, channel })
    if (followed) {
      ctx.reply(`You successfuly followed to ${channel}`)
    }
  } catch (e) {
    error(e)
    ctx.reply(e.message)
  }
  ctx.scene.leave()
})

const unFollowScene = new BaseScene('unfollow', {
  ttl: 60
})
unFollowScene.enter((ctx) => ctx.reply('Enter streamer username you want to unfollow'))
unFollowScene.on('message', async (ctx) => {
  const channel = ctx.message.text
  try {
    const unfollowed = await unfollow({ service, userId: ctx.from.id, channel })
    if (!unfollowed) {
      ctx.reply(`You aren't followed to ${channel}.`)
    } else {
      ctx.reply(`You was successfuly unfollowed from ${channel}.`)
    }
  } catch (e) {
    error(e)
    ctx.reply(e.message)
  }
  ctx.scene.leave()
})


const stage = new Stage([followScene, unFollowScene])
stage.command('cancel', Stage.leave())
bot.use(stage.middleware())


bot.command(['start', 'help'], ({ reply }) => {
  reply(`Hi! I will notify you about the beginning of the broadcasts on Twitch.`)
})
bot.command('follow', Stage.enter('follow'))
bot.command('unfollow', Stage.enter('unfollow'))
bot.command('follows', async (ctx) => {
  const channels = await follows({ userId: ctx.from.id, service })
  if (isBoolean(channels)) {
    ctx.reply(`You aren't followed to anyone`)
  } else {
    ctx.reply(`You are followed to ${channels.join(', ')}`)
  }
})

bot.command('live', async (ctx) => {
  const channels = await live({ userId: ctx.from.id, service })
  if (isBoolean(channels)) {
    ctx.reply(`There is no channels currently online`)
  } else {
    const links = channels.map((o) => `https://twitch.tv/${o}`)
    ctx.reply(`Currently online: \n ${links.join('\n')}`)
  }
})

export function say (chatId: number | number[], message: string, imgUrl: string) {
  const targets = Array.isArray(chatId) ? chatId : [chatId]
  for (const target of targets) {
    bot.telegram.sendPhoto(target, imgUrl, {
      caption: message
    })
  }
}

bot.launch().then(() => info('Telegram bot connected'))