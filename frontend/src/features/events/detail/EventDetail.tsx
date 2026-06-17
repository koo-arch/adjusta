'use client'
import React from 'react';
import Link from 'next/link';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useParams } from 'next/navigation';
import DetailCard from './DetailCard';

const EventDetail = () => {
    const params = useParams<{id: string}>();
    const { eventDetail, isLoading, error } = useFetchEventDetail(params.id);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (error || !eventDetail) {
        return (
            <div className="py-8 text-center">
                <p className="mb-4 text-sm text-gray-500">イベントが見つかりませんでした。</p>
                <Link href="/schedule/draft" className="text-sm text-indigo-600 hover:underline">
                    イベント一覧へ戻る
                </Link>
            </div>
        );
    }

    return (
        <div>
            <DetailCard detail={eventDetail} eventID={params.id} />
        </div>
    );
};

export default EventDetail;
