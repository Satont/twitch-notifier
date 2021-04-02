import { resolve } from 'path'
import { error, info } from './logger'
import { promises } from 'fs'

async function* getFiles(dir: string) {
  const dirents = await promises.readdir(dir, { withFileTypes: true })
  for (const dirent of dirents) {
    const res = resolve(dir, dirent.name)
    if (dirent.isDirectory()) {
      yield* getFiles(res)
    } else {
      yield res
    }
  }
}

export default async () => {
  const folders = {
    libs: 'Lib',
    services: 'Service',
    watchers: 'Watcher',
  }
  for (const folder of Object.keys(folders)) {
    try {
      for await (const file of getFiles(resolve(__dirname, '..', folder))) {
        if (!file.endsWith('.js') && !file.endsWith('.ts') || file.endsWith('.d.ts')) continue
        if (file.includes('_interface')) continue

        const loadedFile = (await import(resolve(__dirname, '..', folder, file))).default
        if (!loadedFile) continue
        if (typeof loadedFile.init !== 'undefined') await loadedFile.init()

        info(`${folders[folder]} ${loadedFile.constructor.name.toUpperCase()} loaded`)
      }
    } catch (e) {
      error(e)
      continue
    }
  }

  return true
}

