"use client";

import { useEffect, useRef } from "react";
import * as shaka from "shaka-player/dist/shaka-player.compiled";

export type VideoPlayerProps = {
  url: string;
  poster: string;
};

export function VideoPlayer(props: VideoPlayerProps) {
  let { url, poster } = props;
  
  const videoNodeRef = useRef<HTMLVideoElement>(null);
  const playerRef = useRef<shaka.Player>();
  //url = 'http://localhost:9000/media/uploads11/mpeg-dash/output.mpd' // TODO: Arrumar url para localhost e uploads{id}/mpeg-dash/output.mpd
  url = 'http://localhost:9000/media/uploads/1/test.mpd' // TODO: Arrumar url para localhost e uploads{id}/mpeg-dash/output.mpd
  useEffect(() => {
    if (!playerRef.current) {
      shaka.polyfill.installAll();
      playerRef.current = new shaka.Player();
      playerRef.current.attach(videoNodeRef.current!).then(() => {
        console.log("The player has been attached");
        playerRef.current!.load(url).then(() => {
          console.log("The video has now been loaded!");
        }).catch((e) => console.log('e',e));
      });
    } 
  }, [url]);

  return (
    <video
      ref={videoNodeRef}
      className="w-full h-full"
      controls
      autoPlay
      poster={poster}
    />
  );
}