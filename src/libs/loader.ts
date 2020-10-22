import { resolve } from 'path'
import { error, info } from './logger'
import { promises } from 'fs'

export default async function* getFiles(dir: string) {
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

const loader = async () => {
  const folders = {
    services: 'Service',
  }
  for (const folder of Object.keys(folders)) {
    try {
      for await (const file of getFiles(resolve(__dirname, '..', folder))) {
        const loadedFile = (await import(resolve(__dirname, '..', folder, file))).default
        if (typeof loadedFile.init !== 'undefined') loadedFile.init()

        info(`${folders[folder]} ${loadedFile.constructor.name.toUpperCase()} loaded`)
      }
    } catch (e) {
      error(e)
      continue
    } 
  }
}

loader()