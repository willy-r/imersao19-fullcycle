'use client';

import { useEffect, useRef } from 'react';
import * as shaka from 'shaka-player/dist/shaka-player.compiled';

export type VideoPlayerProps = {
  url: string;
  poster: string;
};

export function VideoPlayer(props: VideoPlayerProps) {
    
  const { url, poster } = props;
  console.log(url);
  
  const videoNodeRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<shaka.Player>();

  useEffect(() => {
    if (!playerRef.current) {
      shaka.polyfill.installAll();
      playerRef.current = new shaka.Player();
      playerRef.current.attach(videoNodeRef.current!).then(() => {
        console.log('The player has been attached');
        playerRef.current!.load(url).then(() => {
          console.log('The video has now been loaded!');
        });
      });
    } 
  }, [url]);

  return (
    <video
      ref={videoNodeRef}
      className="w-full h-full"
      controls
      //   autoPlay
      poster={poster}
    />
  );
}
