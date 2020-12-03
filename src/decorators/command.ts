import { ServiceInterface } from '../services/_interface'

export interface CommandDecoratorOptions {
  description: string,
  isVisible?: boolean,
}

export function command(name: string, opts: CommandDecoratorOptions): MethodDecorator {
  return (instance: ServiceInterface, methodName: string): void => {
    if (!instance.commands) instance.commands = []
    const data = { name, fnc: methodName, ...opts }

    instance.commands.push(data)
  }
}
