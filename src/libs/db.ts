import { Sequelize } from 'sequelize-typescript';
import { readdirSync } from 'fs'
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
    console.log('Succesfuly connected to db.')
    sequelize.sync()
  })
  .catch(err => {
    console.error('Unable to connect to the database:', err)
    process.exit()
  })

/* for (let file of readdirSync(join(__dirname, '../models'))) {
  sequelize.import(`${join(__dirname, '../models')}/${file}`)
}

sequelize.sync() */

export { sequelize }