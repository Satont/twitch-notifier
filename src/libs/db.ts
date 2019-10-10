import { Sequelize } from 'sequelize-typescript'
import { info, error } from '../helpers/logs'
let connected: boolean = false
import { join } from 'path'

const sequelize = new Sequelize(process.env.DB_NAME, process.env.DB_USER, process.env.DB_PASSWORD, {
  host: process.env.DB_HOST,
  port: Number(process.env.DB_PORT),
  dialect: 'postgres',
  pool: {
    max: 10,
    min: 0,
    acquire: 30000,
    idle: 10000
  },
  models: [join(__dirname, '../models')]
})

sequelize.authenticate()
  .then(() => {
    sequelize.sync().then(() => connected = true).then(() => info('Succesfuly connected to db.')).then(() => console.log(sequelize.models))
  })
  .catch(err => {
    console.log(err)
    process.exit()
  })


export { sequelize, connected }