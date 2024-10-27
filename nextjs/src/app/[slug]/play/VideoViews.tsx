
const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export async function getViews(videoId: number): Promise<number> {
  await sleep(2000);
  const response = await fetch(
    `http://localhost:8000/api/videos/${videoId}/views`,
    {
      next: {
        revalidate: 60,
      },
    }
  );

  return (await response.json()).views;
}

export type VideoViewsProps = {
  videoId: number;
  views?: number;
}

export async function VideoViews(props: VideoViewsProps){
  const {videoId, views: propViews} = props;
  const views = propViews ? propViews : await getViews(videoId);
  return (<span>{views} visualizações</span>);
}
