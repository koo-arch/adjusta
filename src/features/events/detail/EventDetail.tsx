'use client'
import React from 'react';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useParams } from 'next/navigation';
import DetailCard from './DetailCard';

const EventDetail = () => {
    const params = useParams<{slug: string}>();

    const { eventDetail, isLoading, error } = useFetchEventDetail(params.slug);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            {eventDetail &&
                <DetailCard detail={eventDetail} slug={params.slug} />
            }
        </div>
    );
};

export default EventDetail;