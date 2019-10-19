import { Sequelize } from 'sequelize-typescript'
import { info } from '../helpers/logs'
let connected: boolean = false
import { User } from '../models/User'
import { Channel } from '../models/Channel'
import { config } from '../helpers/config'

const sequelize = new Sequelize(config.db.name, config.db.user, config.db.password, {
  host: config.db.host,
  port: config.db.port,
  dialect: 'postgres',
  pool: {
    max: 10,
    min: 0,
    acquire: 30000,
    idle: 10000
  },
  models: [User, Channel]
})

sequelize.authenticate()
  .then(() => {
    sequelize.sync().then(() => connected = true).then(() => info('Succesfuly connected to db.'))
  })
  .catch(err => {
    console.log(err)
    process.exit()
  })


export { sequelize, connected }