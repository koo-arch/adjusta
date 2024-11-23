'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { useFetchNeedsActionDrafts } from '@/hooks/event/useFetchNeedsActionDrafts';
import MiniEventCard from '@/features/events/MiniEventCard';

const NeedsActionDrafts: React.FC = () => {
    const router = useRouter();
    const { needsActionDrafts, isLoading, error } = useFetchNeedsActionDrafts();

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <section className="bg-white">
            <h2 className="text-lg font-bold mb-4">対応が必要なイベント</h2>
            {needsActionDrafts && needsActionDrafts.length > 0 ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {needsActionDrafts.map((event) => (
                        <MiniEventCard
                            key={event.id}
                            onClick={() => router.push(`/schedule/draft/${event.id}`)}
                            {...event}
                        />
                    ))}
                </div>
            ) : (
                <p>イベントはありません。</p>
            )}
        </section>
    )
}

export default NeedsActionDrafts;