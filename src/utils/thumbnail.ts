export class ThumbnailBuilder {
  /**
   * Build thumbnail URL from Twitch template URL
   * @param thumbnailUrl - Twitch thumbnail URL with {width} and {height} placeholders
   * @returns Final thumbnail URL
   */
  build(thumbnailUrl: string): string {
    return thumbnailUrl
      .replace('{width}', '1920')
      .replace('{height}', '1080');
  }
}
