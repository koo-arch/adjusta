'use client'
import React from 'react';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useParams } from 'next/navigation';
import DetailCard from './DetailCard';

const EventDetail = () => {
    const params = useParams<{id: string}>();

    const { eventDetail, isLoading, error } = useFetchEventDetail(params.id);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            {eventDetail &&
                <DetailCard detail={eventDetail} id={params.id} />
            }
        </div>
    );
};

export default EventDetail;