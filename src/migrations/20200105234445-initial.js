'use strict';

module.exports = {
  up: (queryInterface, Sequelize) => {
    return Promise.all([
      queryInterface.createTable('users', {
        id: Sequelize.INTEGER,
        follows: Sequelize.ARRAY(Sequelize.INTEGER),
        service: {
          allowNull: false,
          type: Sequelize.ENUM('vk', 'telegram')
        }
      }),
      queryInterface.createTable('channels', {
        id: {
          type: Sequelize.INTEGER,
          primaryKey: true,
        },
        username: Sequelize.STRING,
        online: Sequelize.BOOLEAN
      }),
    ])
  },

  down: (queryInterface, Sequelize) => {
    return Promise.all([
      queryInterface.dropTable('users'),
      queryInterface.dropTable('channels'),
    ])
  }
}
