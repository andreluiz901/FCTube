"use server";

export async function likeAction(formData: FormData) {
  const videoId = formData.get("videoId");
  await fetch(`${process.env.DJANGO_API_URL}/videos/${videoId}/like`, {
    method: "POST",
  });
}

export async function unlikeAction(formData: FormData) {
  const videoId = formData.get("videoId");
  await fetch(`${process.env.DJANGO_API_URL}/videos/${videoId}/unlike`, {
    method: "POST",
  });
}
