'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import { useFetchUpcomingEvents } from '@/hooks/event/useFetchUpcomingEvents';
import MiniEventCard from '@/features/events/MiniEventCard';

const UpcomingEvents: React.FC = () => {
    const router = useRouter();
    const { upcomingEvents, isLoading, error } = useFetchUpcomingEvents();

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <section className="bg-white p-4">
            <h2 className="text-lg font-bold mb-4">直近のイベント</h2>
            {upcomingEvents && upcomingEvents.length > 0 ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {upcomingEvents.map((event) => (
                        <MiniEventCard
                            key={event.id}
                            event={event}
                            onClick={() => router.push(`/schedule/draft/${event.id}`)}
                        />
                    ))}
                </div>
            ) : (
                <p>直近のイベントはありません。</p>
            )}
        </section>
    );
};

export default UpcomingEvents;