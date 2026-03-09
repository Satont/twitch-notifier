export class ThumbnailBuilder {
  /**
   * Build thumbnail URL from Twitch template URL
   * @param thumbnailUrl - Twitch thumbnail URL with {width} and {height} placeholders
   * @param checkValidity - Whether to check if the URL is accessible (with retry logic)
   * @returns Final thumbnail URL
   */
  async build(thumbnailUrl: string, checkValidity = false): Promise<string> {
    let thumbnail = thumbnailUrl
      .replace('{width}', '1920')
      .replace('{height}', '1080');

    if (!checkValidity) {
      return thumbnail;
    }

    const isValid = await this.checkValidity(thumbnail, 0);

    if (!isValid) {
      // Fallback to lower resolution
      thumbnail = thumbnail
        .replace('1920', '1280')
        .replace('1080', '720');
    }

    return thumbnail;
  }

  /**
   * Check if thumbnail URL is accessible with retry logic
   * @param url - URL to check
   * @param attempt - Current attempt number (max 5)
   * @returns Whether the URL is valid
   */
  private async checkValidity(url: string, attempt: number): Promise<boolean> {
    try {
      const response = await fetch(url, {
        method: 'HEAD',
        redirect: 'manual',
      });

      if (response.status === 200) {
        return true;
      }

      if (attempt >= 5) {
        return false;
      }

      // Wait 5 seconds before retry
      await new Promise(resolve => setTimeout(resolve, 5000));
      return this.checkValidity(url, attempt + 1);
    } catch (error) {
      if (attempt >= 5) {
        return false;
      }

      await new Promise(resolve => setTimeout(resolve, 5000));
      return this.checkValidity(url, attempt + 1);
    }
  }
}
