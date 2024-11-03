export type VideoCardSkeletonProps = {
  orientation?: "horizontal" | "vertical";
};

export default function VideoCardSkeleton(props: VideoCardSkeletonProps) {
  const { orientation = "vertical" } = props;
  return orientation === "vertical" ? (
    <div className="animate-pulse shadow-md rounded-lg overflow-hidden">
      <div className="bg-secondary h-44 w-full rounded"></div>
      <div className="mt-3">
        <div className="h-4 bg-secondary rounded w-3/4 mb-1"></div>
        <div className="h-3 bg-secondary rounded w-1/2 mb-2"></div>
      </div>
    </div>
  ) : (
    <div className="flex space-x-2 animate-pulse">
      <span className="bg-secondary w-5/12 h-24 rounded-md"></span>
      <div className="flex-1 space-y-2 py-1 items-center">
        <span className="block bg-secondary h-4 w-3/4 rounded"></span>
        <span className="block bg-secondary h-3 w-1/2 rounded"></span>
      </div>
    </div>
  );
}
