'use client'
import React from 'react';
import Link from 'next/link';
import { useFetchEventDetail } from '@/features/events/hooks/useFetchEventDetail';
import DetailCard from './DetailCard';

interface EventDetailProps {
    eventID: string;
}

const EventDetail: React.FC<EventDetailProps> = ({ eventID }) => {
    const { eventDetail, isLoading, error } = useFetchEventDetail(eventID);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (error || !eventDetail) {
        return (
            <div className="py-8 text-center">
                <p className="mb-4 text-sm text-gray-500">イベントが見つかりませんでした。</p>
                <Link href="/events" className="text-sm text-indigo-600 hover:underline">
                    イベント一覧へ戻る
                </Link>
            </div>
        );
    }

    return (
        <div>
            <DetailCard detail={eventDetail} eventID={eventID} />
        </div>
    );
};

export default EventDetail;
