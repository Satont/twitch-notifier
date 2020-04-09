import { config } from './config'
import { SocksProxyAgent } from 'socks-proxy-agent';

export const agent = new SocksProxyAgent({
  host: config.proxy.host,
  userId: config.proxy.username,
  password: config.proxy.password,
  port:  config.proxy.port,
  timeout: 1 * 60 * 1000,
}) || null
