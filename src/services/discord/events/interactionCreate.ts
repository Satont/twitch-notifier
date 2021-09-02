import { Interaction } from 'discord.js'
import { Client, Discord, On } from 'discordx'

@Discord()
export abstract class AppDiscord {
  @On('interactionCreate')
  interactionCreate(interaction: Interaction, client: Client) {
    return client.executeInteraction(interaction)
  }
}