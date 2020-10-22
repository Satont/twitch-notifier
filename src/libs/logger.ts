import dayjs from 'dayjs'
import fs from 'fs'
import stripAnsi from 'strip-ansi'
import { inspect } from 'util'
import os from 'os'
import { createStream } from 'rotating-file-stream'
import { getFunctionNameFromStackTrace } from '../commons/stacktrace'

const levelFormat = {
  error: '!!! ERROR !!!',
  debug: 'DEBUG:',
  chatIn: '<<<',
  chatOut: '>>>',
  whisperOut: 'W>>>',
  info: '!!!',
  timeout: '+timeout',
  ban: '+ban',
  unban: '-ban',
  follow: '+follow',
  host: '+host',
  raid: '+raid',
  unfollow: '-follow',
  cheer: '+cheer',
  donate: '+donate',
  sub: '+sub',
  bits: '+bits',
  raided: '+raided',
  hostring: '?hosting',
  hosted: '+hosted',
  subgift: '+subgift',
  resub: '+resub',
  moded: '+moderator',
  unmoded: '-moderator',
  redemption: '+redemption',
  highlight: '+highlight',
}

const logDir = './logs'

if (!fs.existsSync(logDir)) fs.mkdirSync(logDir)

const file = createStream('./logs/bot.log', {
  maxFiles: 10,
  size: '512M',
  compress: 'gzip',
})

function format(level: string, message: any, category?: string) {
  const timestamp = dayjs().format('YYYY-MM-DD[T]HH:mm:ss.SSS')

  if (typeof message === 'object') {
    message = inspect(message)
  }
  return [timestamp, levelFormat[level], category, message].filter(Boolean).join(' ')
}

function log(message: any) {
  const level = getFunctionNameFromStackTrace()

  const formattedMessage = format(level, message)
  process.stdout.write(formattedMessage + '\n')
  file.write(stripAnsi(formattedMessage) + os.EOL)
}

export function isDebugEnabled(category: string) {
  if (!process.env.DEBUG) return false
  const categories = category.split('.')
  let bEnabled = false
  bEnabled = process.env.DEBUG.includes(category) || process.env.DEBUG.includes(categories[0] + '.*')
  bEnabled = process.env.DEBUG === '*' || bEnabled
  return bEnabled
}

export function debug(category: string, message: any) {
  const categories = category.split('.')
  if (categories.length > 2 && category !== '*') {
    throw Error('For debug use only <main>.<sub> or bot*')
  }

  if (isDebugEnabled(category) || category == '*') {
    const formattedMessage = format('debug', message, category)
    process.stdout.write(formattedMessage + '\n')
    file.write(formattedMessage + os.EOL)
  }
}

export function error(message: any) {
  log(message)
}

export function chatIn(message: any) {
  log(message)
}

export function chatOut(message: any) {
  log(message)
}

export function info(message: any) {
  log(message)
}

export function warning(message: any) {
  log(message)
}
