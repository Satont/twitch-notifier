import dayjs from 'dayjs'
import fs from 'fs'
import stripAnsi from 'strip-ansi'
import { inspect } from 'util'
import os from 'os'
import { createStream } from 'rotating-file-stream'
import { getFunctionNameFromStackTrace } from '../commons/stacktrace'

const levelFormat = {
  error: '!!! ERROR !!!',
  chatIn: '<<<',
  chatOut: '>>>',
  info: '!!!',
  warning: 'WARNING',
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
