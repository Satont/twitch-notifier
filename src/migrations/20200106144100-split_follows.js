'use strict';

module.exports = {
  up: (queryInterface, Sequelize) => {
    return Promise.all([
      queryInterface.createTable('follows', {
        user_id: Sequelize.INTEGER,
        channel_id: {
          type: Sequelize.INTEGER,
          references: {
            model: {
              tableName: 'channels',
              schema: 'public'
            },
            key: 'id'
          },
        },
      }),
      queryInterface.sequelize.query('SELECT * FROM public.users').then(data => {
        data = data[0]
        return data
      }).then(data => {
        // data = [ { id, follows[], service } ]
        const follows = []
        for (const user of data) {
          if (!user.follows || !user.follows.length) continue
          for (const channel of user.follows) {
            follows.push({ user_id: user.id, channel_id: channel })
          }
        }
        return follows
      }).then(data => {
        if (!data.length) return true
        return queryInterface.bulkInsert('follows', data)
      }),
      queryInterface.removeColumn('users', 'follows')
    ])
  },

  down: (queryInterface, Sequelize) => {
    return Promise.all([
      queryInterface.dropTable('follows')
    ])
  }
};
