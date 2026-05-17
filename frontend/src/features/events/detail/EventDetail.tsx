'use client'
import React, { useEffect } from 'react';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useParams, useRouter } from 'next/navigation';
import DetailCard from './DetailCard';

const EventDetail = () => {
    const params = useParams<{slug: string}>();
    const router = useRouter();

    const { eventDetail, isLoading, error } = useFetchEventDetail(params.slug);

    useEffect(() => {
        if (!isLoading && (!eventDetail || error)) {
            router.replace('/schedule/draft');
        }
    },[isLoading, eventDetail, error, router]);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (error || !eventDetail) {
        return null;
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