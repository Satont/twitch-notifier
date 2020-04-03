/* do not change everything in this file until you don't know what are you doing */

module.exports = {
  development: {
    username: 'database_dev',
    password: 'database_dev',
    database: 'database_dev',
    host: '127.0.0.1',
    port: Number(process.env.DB_PORT),
    dialect: 'postgres',
    migrationStorageTableName: 'migrations',
  },
  production: {
    username: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    database: process.env.DB_NAME,
    port: Number(process.env.DB_PORT),
    host: process.env.DB_HOST,
    dialect: 'postgres',
    migrationStorageTableName: 'migrations',
  },
}
