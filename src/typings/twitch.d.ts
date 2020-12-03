export interface ITwitchStreamChangedChannel {
  user_id: string,
  user_name?: string,
  game_id?: string,
  game_name?: string,
  community_ids?: string[],
  tag_ids?: string[]
  type: string,
  title?: string,
  viewer_count?: number,
  started_at?: string,
  language?: string,
  thumbnail_url?: string
}

export interface ITwitchStreamChangedPayload{
  data: Array<ITwitchStreamChangedChannel>
}
