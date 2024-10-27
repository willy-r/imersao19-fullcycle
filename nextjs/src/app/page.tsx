import { Suspense } from 'react';
import { VideosList } from '@/components/VideosList';
import VideoCardSkeleton from '@/components/VideoCardSkeleton';

export default async function Home({ searchParams }: { searchParams: { search: string } }) {
  const { search } = await searchParams;
  return (
    <main className="container mx-auto px-4 py-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <Suspense
          fallback={new Array(15).fill(null).map((_, index) => (
            <VideoCardSkeleton key={index} />
          ))}
        >
          <VideosList search={search} />
        </Suspense>
      </div>
    </main>
  );
}
