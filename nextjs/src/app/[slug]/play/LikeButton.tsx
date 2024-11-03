"use client";

import { useEffect, useState } from "react";
import { likeAction, unlikeAction } from "../../../actions/like.action"; //server action

export type NumLikeProps = {
  videoId: number;
  likes: number;
};

export function LikeButton(props: NumLikeProps) {
  const { videoId, likes: initialLikes } = props;
  const [isLiked, setIsLiked] = useState(false);
  const [currentLikes, setCurrentLikes] = useState(initialLikes);

  async function onSubmit() {
    const likeChange = isLiked ? -1 : +1;
    const updatedLikes = currentLikes + likeChange;

    window.localStorage.setItem("isLiked", isLiked ? "false" : "true");
    window.localStorage.setItem("currentLikes", updatedLikes.toString());

    setCurrentLikes(updatedLikes);
    setIsLiked(!isLiked);

    const formData = new FormData();
    formData.append("videoId", videoId.toString());
    !isLiked ? await likeAction(formData): await unlikeAction(formData);
  }

  useEffect(() => {
    const storedIsLiked = window.localStorage.getItem("isLiked");
    const storedLikesRaw = window.localStorage.getItem("currentLikes");
    setIsLiked(storedIsLiked === "true");
    const storedLikes = storedLikesRaw ? parseInt(storedLikesRaw) : 0;
    setCurrentLikes(storedLikes > initialLikes ? storedLikes : initialLikes);
  }, [initialLikes]);

  return (
    <form className="flex items-center space-x-2" action={onSubmit}>
      <button className="flex items-center text-secondary">
        <svg
          className={`h-8 w-8 text-primary ${isLiked ? "fill-current" : ""}`}
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          id="like-icon"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"
          />
        </svg>
      </button>
      <span className="ml-1 text-primary" id="num-likes">
        {currentLikes}
      </span>
    </form>
  );
}
