import { Router } from 'express'
const router = Router()
import { Channel } from '../models/Channel'
import { User } from '../models/User'
import { sequelize as db} from '../libs/db'
import axios from 'axios'
import { config } from '../helpers/config'

router.get('/', async (req, res) => {
  res.send('Nothing to do here :)')
})

router.get('/counts', async (req, res) => {
  const users = await User.count()
  const streamers = await Channel.count()
  res.json({ users, streamers })
})

router.get('/top10', async (req, res) => {
  const query = `select channels.username, count(users.id) as followers_count, channels.id as channel_id
  from channels 
  JOIN users ON channels.id = ANY(users.follows)
  group by channels.username, channel_id
  order by followers_count desc
  limit 10`
  const [top]: any = await db.query(query)
  const request = await axios.get(`https://api.twitch.tv/helix/users?id=${top.map(o => o.channel_id).join('&id=')}`, { headers: {
    'Client-ID': config.twitch.clientId
  }})
  for (const channel of request.data.data) {
    const index = request.data.data.indexOf(channel)
    top[index].imageUrl = channel.profile_image_url
  }
  res.json(top)
})

export { router }
