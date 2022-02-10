import { execFile } from 'child_process';
import { resolve } from 'path';
import dotenv from 'dotenv';
dotenv.config({ path: resolve(process.cwd(), '.env') });

const eventName: string = process.argv[2] ?? 'streamup';
const channelId: string = process.argv[3] ?? '128644134';

const eventNamesMapping: {
  [x: string]: string;
} = {
  'stream-change': `channel.update.${channelId}`,
  streamup: `stream.online.${channelId}`,
  streamdown: `stream.offline.${channelId}`,
};

// event trigger stream-change -F http://localhost:3001/eventsub/event/channel.update.46571894 -f 46571894 -s channel.update.46571894.secret
const command = [
  'event',
  'trigger',
  eventName,
  '-F',
  `https://${process.env.SITE_URL.replace('http://', '').replace('https://', '')}/twitch/eventsub/event/${eventNamesMapping[eventName]}`,
  '-f',
  channelId,
  '-t',
  channelId,
  '-s',
  `${eventNamesMapping[eventName]}.${process.env.TWITCH_CLIENT_ID}`,
];
execFile(resolve(__dirname, 'twitch-cli'), command, (error, stdout, stderr) => {
  console.log('Executing: ' + command.join(' '));

  if (error) console.error('Error:', error);
  if (stdout) console.log('stdout', stdout);
  if (stderr) console.log('stderr', stderr);
});