const config = require('../../config.js')
const db = require('../../database.js')[process.env.NODE_ENV]

export default { config }
export { config, db }