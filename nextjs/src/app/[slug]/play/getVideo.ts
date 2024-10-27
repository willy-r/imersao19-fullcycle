import { VideoModel } from '@/models';

export async function getVideo(slug: string): Promise<VideoModel> {
  const response = await fetch(`http://localhost:8000/api/videos/${slug}`, {
    next: {
      tags: [`video-${slug}`]
    }
  });
  return response.json();
}
