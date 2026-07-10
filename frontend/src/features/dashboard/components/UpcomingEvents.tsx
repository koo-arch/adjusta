'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { useFetchUpcomingEvents } from '@/features/events/hooks/useFetchUpcomingEvents';
import MiniEventCard from '@/features/events/components/list/MiniEventCard';
import BoardSlider from '@/features/dashboard/components/BoardSlider';
import EmptyStateCard from '@/features/dashboard/components/EmptyStateCard';

const UpcomingEvents: React.FC = () => {
    const router = useRouter();
    const { upcomingEvents, isLoading } = useFetchUpcomingEvents();

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <section className="bg-inherit">
            <h2 className="text-lg font-bold mb-4">直近のイベント</h2>
            {upcomingEvents && upcomingEvents.length > 0 ? (
                <BoardSlider>
                    {upcomingEvents.map((event) => (
                        <MiniEventCard
                        key={event.id}
                        onClick={() => router.push(`/events/${event.id}`)}
                        {...event}
                        />
                    ))}
                </BoardSlider>
            ) : (
                <EmptyStateCard>直近のイベントはありません。</EmptyStateCard>
            )}
        </section>
    );
};

export default UpcomingEvents;
