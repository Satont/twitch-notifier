import { log } from 'console'
import chalk from 'chalk'


export function info (message: string) {
  log(chalk.blue(message))
}

export function error (message: string) {
  log(chalk.red(message))
}

export function warning (message: string) {
  log(chalk.yellow(message))
}
