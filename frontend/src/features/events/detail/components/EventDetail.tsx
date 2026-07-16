'use client'
import React from 'react';
import Link from 'next/link';
import { useFetchEventDetail } from '@/features/events/hooks/useFetchEventDetail';
import { isNotFoundAPIError } from '@/lib/api/errors';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import EventDetailHeader from './EventDetailHeader';
import ProposedDatesSection from './ProposedDatesSection';

interface EventDetailProps {
    eventID: string;
}

export const EventDetailSkeleton = () => (
    <div className="space-y-6">
        <div className="space-y-3">
            <div className="flex items-start justify-between gap-4">
                <Skeleton className="h-8 w-64 max-w-full" />
                <Skeleton className="h-10 w-24" />
            </div>
            <Skeleton className="h-5 w-32" />
            <Skeleton className="h-4 w-full" />
        </div>
        <div className="space-y-3 border-t border-border pt-6">
            <Skeleton className="h-6 w-24" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-9 w-full" />
            <Skeleton className="h-9 w-3/4" />
        </div>
    </div>
);

const EventDetail: React.FC<EventDetailProps> = ({ eventID }) => {
    const { eventDetail, isPending, error, refetch } = useFetchEventDetail(eventID);

    if (isPending) {
        return <EventDetailSkeleton />;
    }

    // 存在しない ID・他ユーザーのイベントは 404 で返る(screen-design 5.6)
    if (isNotFoundAPIError(error)) {
        return (
            <div className="flex flex-col items-center gap-4 py-16 text-center">
                <p className="text-sm text-muted-foreground">イベントが見つかりませんでした。</p>
                <Button variant="outline" asChild>
                    <Link href="/events">イベント一覧へ戻る</Link>
                </Button>
            </div>
        );
    }

    if (error || !eventDetail) {
        return (
            <div className="flex flex-col items-center gap-4 py-16 text-center">
                <p className="text-sm text-muted-foreground">
                    イベントの取得に失敗しました。時間をおいて再度お試しください。
                </p>
                <Button variant="outline" onClick={() => refetch()}>
                    再試行
                </Button>
            </div>
        );
    }

    return (
        <article className="space-y-6">
            <EventDetailHeader eventID={eventID} detail={eventDetail} />
            <ProposedDatesSection eventID={eventID} detail={eventDetail} />
        </article>
    );
};

export default EventDetail;
