'use strict';

module.exports = {
  up: (queryInterface, Sequelize) => {
    return Promise.all([
      queryInterface.addColumn('channels', 'game', Sequelize.STRING),
      queryInterface.addColumn('users', 'follow_game_change', {
        type: Sequelize.BOOLEAN,
        allowNull: false,
        defaultValue: false
      })
    ])
  },

  down: (queryInterface, Sequelize) => {
    return Promise.all([
      queryInterface.removeColumn('channels', 'game'),
      queryInterface.removeColumn('users', 'follow_game_change')
    ])
  }
};
