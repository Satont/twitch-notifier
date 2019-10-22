import express from 'express'
import { config } from '../helpers/config'
import { info } from '../helpers/logs'
const app = express()
import { router as api } from '../routes/api'

app.use(api)

app.listen(config.panel.port, () => {
  info(`Api listening on ${config.panel.port} port`)
})