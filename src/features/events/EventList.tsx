'use client'
import React from 'react';
import { useFetchEventList } from '@/hooks/event/useFetchEventList';
import EventCard from './EventCard';
import { useRouter } from 'next/navigation';

const EventList: React.FC = () => {
    const { eventList, isLoading, error } = useFetchEventList();
    const router = useRouter();

    console.log(eventList);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-1 lg:grid-cols-3">
            {eventList?.map((event) => (
                <EventCard
                    key={event.id}
                    event={event}
                    onClick={() => router.push(`/schedule/draft/${event.id}`)}
                />
            ))}
        </div>
    );
};

export default EventList;