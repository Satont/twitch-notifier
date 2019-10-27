import Telegraf from 'telegraf'
import { Stage, session, BaseScene } from 'telegraf'
import { config } from '../helpers/config'
import { info, error } from '../helpers/logs'
import { User } from '../models/User'
import { Channel } from '../models/Channel'
import { Twitch } from './twitch'
import { remove } from 'lodash'

const service = 'telegram'
const twitch = new Twitch(config.twitch.clientId)
const bot = new Telegraf(config.telegram.token)

bot.use(session())
bot.use((ctx, next) => {
  if (!ctx.message) return
  info(`Telegram | New message from ${ctx.from.username} [${ctx.from.id}], message: ${ctx.message.text}`)
  next()
})

const followScene = new BaseScene('follow', {
  ttl: 60
})
followScene.enter((ctx) => ctx.reply('Enter streamer username you want to follow'))
followScene.on('message', async (ctx) => {
  const [user] = await User.findOrCreate({ where: { id: ctx.from.id, service }, defaults: { follows: [], service } })
  if (/[^a-zA-Z0-9_]/gmu.test(ctx.message.text)) {
    return ctx.reply('Username can cointain only "a-z", "0-9" и "_" symbols')
  }
  try {
    const streamer = await twitch.getChannel(ctx.message.text)
    await Channel.findOrCreate({ where: { id: streamer.id }, defaults: { username: streamer.login, online: false } })
    if (user.follows.includes(streamer.id)) {
      return ctx.reply(`You already followed to ${streamer.displayName}.`)
    } else {
      user.follows.push(streamer.id)
      await user.update({ follows: user.follows })
      await ctx.reply(`You successfully subscribed to ${streamer.displayName}!`)
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
  const user = await User.findOne({ where: { id: ctx.from.id, service } })
  if (!user) {
    return ctx.reply('You are not following to anyone.')
  }
  if (/[^a-zA-Z0-9_]/gmu.test(ctx.message.text)) {
    return ctx.reply('Username can cointain only "a-z", "0-9" и "_" symbols')
  }
  try {
    const streamer = await twitch.getChannel(ctx.message.text)
    if (!user.follows.includes(streamer.id)) {
      return ctx.reply(`You are not followed for ${streamer.displayName}.`)
    } else {
      remove(user.follows, (o: number) => o === streamer.id)
      await user.update({ follows: user.follows })
      await ctx.reply(`You successfully unsubscribed from ${streamer.displayName}.`)
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
  const user = await User.findOne({ where: { id: ctx.from.id, service } })
  if (!user || !user.follows.length) {
    return ctx.reply('You are not followed to anyone.')
  } else {
    const channels = await twitch.getChannelsById(user.follows)
    const follows = channels.map(o => o.displayName).join(', ')
    ctx.reply(`Your list of follows: ${follows}`)
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