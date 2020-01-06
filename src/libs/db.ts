import { Sequelize } from 'sequelize-typescript'
import { info } from '../helpers/logs'
import { db } from '../helpers/config'
import { join } from 'path'

let connected: boolean = false

const sequelize = new Sequelize(db.database, db.username, db.password, {
  host: db.host,
  port: db.port,
  dialect: 'postgres',
  pool: {
    max: 10,
    min: 0,
    acquire: 30000,
    idle: 10000
  },
  models: [join(__dirname, '../models')],
  logging: false
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