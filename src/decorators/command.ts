import { ServiceInterface } from '../services/_interface'

export function command(name: string): MethodDecorator {
  return (instance: ServiceInterface, methodName: string): void => {
    if (!instance.commands) instance.commands = []
    const data = { name, fnc: methodName }

    instance.commands.push(data)
  }
}
