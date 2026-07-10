'use client'
import React from 'react';
import Link from 'next/link';
import { useFetchEventDetail } from '@/features/events/hooks/useFetchEventDetail';
import { isNotFoundAPIError } from '@/lib/api/errors';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import EventDetailHeader from './EventDetailHeader';
import ProposedDatesSection from './ProposedDatesSection';
import EventInfoSection from './EventInfoSection';

interface EventDetailProps {
    eventID: string;
}

const EventDetailSkeleton = () => (
    <div className="space-y-6">
        <div className="flex items-start justify-between gap-4">
            <div className="space-y-2">
                <Skeleton className="h-8 w-64 max-w-full" />
                <Skeleton className="h-5 w-20" />
            </div>
            <Skeleton className="h-10 w-40" />
        </div>
        {Array.from({ length: 2 }).map((_, index) => (
            <Card key={index}>
                <CardHeader>
                    <Skeleton className="h-6 w-28" />
                </CardHeader>
                <CardContent className="space-y-3">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-3/4" />
                    <Skeleton className="h-4 w-2/3" />
                </CardContent>
            </Card>
        ))}
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
            <Card>
                <CardContent className="flex flex-col items-center gap-4 p-8 text-center">
                    <p className="text-sm text-muted-foreground">イベントが見つかりませんでした。</p>
                    <Button variant="outline" asChild>
                        <Link href="/events">イベント一覧へ戻る</Link>
                    </Button>
                </CardContent>
            </Card>
        );
    }

    if (error || !eventDetail) {
        return (
            <Card>
                <CardContent className="flex flex-col items-center gap-4 p-8 text-center">
                    <p className="text-sm text-muted-foreground">
                        イベントの取得に失敗しました。時間をおいて再度お試しください。
                    </p>
                    <Button variant="outline" onClick={() => refetch()}>
                        再試行
                    </Button>
                </CardContent>
            </Card>
        );
    }

    return (
        <div className="space-y-6">
            <EventDetailHeader eventID={eventID} detail={eventDetail} />
            <ProposedDatesSection eventID={eventID} detail={eventDetail} />
            <EventInfoSection detail={eventDetail} />
        </div>
    );
};

export default EventDetail;
