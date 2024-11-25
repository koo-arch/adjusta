'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { useFetchNeedsActionDrafts } from '@/hooks/event/useFetchNeedsActionDrafts';
import MiniEventCard from '@/features/events/MiniEventCard';
import BoardSlider from './BoardSlider';
import EmptyStateCard from './EmptyStateCard';

const NeedsActionDrafts: React.FC = () => {
    const router = useRouter();
    const { needsActionDrafts, isLoading, error } = useFetchNeedsActionDrafts();

    if (isLoading) {
        return <p>Loading...</p>;
    }


    return (
        <section className="bg-inherit">
            <h2 className="text-lg font-bold mb-4">対応が必要なイベント</h2>
            {needsActionDrafts && needsActionDrafts.length > 0 ? (
                <BoardSlider>
                    {needsActionDrafts.map((event) => (
                        <div key={event.id}>
                            <MiniEventCard
                            onClick={() => router.push(`/schedule/draft/${event.id}`)}
                            {...event}
                            />
                        </div>
                    ))}
                </BoardSlider>
            ) : (
                <EmptyStateCard>対応が必要なイベントはありません。</EmptyStateCard>
            )}
        </section>
    )
}

export default NeedsActionDrafts;