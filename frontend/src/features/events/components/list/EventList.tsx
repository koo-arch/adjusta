'use client'
import React from 'react';
import { useFetchEventList } from '@/features/events/hooks/useFetchEventList';
import EventCard from './EventCard';
import { useRouter } from 'next/navigation';

const EventList: React.FC = () => {
    const { eventList, isLoading, error } = useFetchEventList();
    const router = useRouter();

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {eventList?.map((event) => (
                <EventCard
                    key={event.id}
                    event={event}
                    onClick={() => router.push(`/events/${event.id}`)}
                />
            ))}
        </div>
    );
};

export default EventList;
