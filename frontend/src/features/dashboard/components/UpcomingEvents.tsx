'use client'
import React from 'react';
import Link from 'next/link';
import { useFetchUpcomingEvents } from '@/features/events/hooks/useFetchUpcomingEvents';
import MiniEventCard from '@/features/events/components/list/MiniEventCard';
import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';

// パネルに出すのは先頭 5 件まで(残りは一覧へ誘導)
const MAX_ITEMS = 5;

const UpcomingEvents: React.FC = () => {
    const { upcomingEvents, isPending, error, refetch } = useFetchUpcomingEvents();

    return (
        <div className="space-y-3">
            <div className="flex items-center justify-between gap-2">
                <h2 className="text-lg font-bold leading-snug tracking-normal text-gray-900">直近の予定</h2>
                <Link
                    href="/events?status=confirmed"
                    className="shrink-0 text-sm text-primary transition-colors hover:text-primary-dark"
                >
                    すべて見る
                </Link>
            </div>
            {isPending ? (
                <div className="space-y-2">
                    {Array.from({ length: 3 }).map((_, index) => (
                        <Skeleton key={index} className="h-14 w-full rounded-lg" />
                    ))}
                </div>
            ) : error ? (
                <div className="space-y-2 rounded-md border border-border bg-card p-4 text-center">
                    <p className="text-sm text-muted-foreground">イベントの取得に失敗しました。</p>
                    <Button variant="outline" size="sm" onClick={() => refetch()}>
                        再試行
                    </Button>
                </div>
            ) : upcomingEvents && upcomingEvents.length > 0 ? (
                <ul className="space-y-2">
                    {upcomingEvents.slice(0, MAX_ITEMS).map((event) => (
                        <li key={event.id}>
                            <MiniEventCard
                                title={event.title}
                                start={event.start}
                                end={event.end}
                                href={`/events/${event.id}`}
                            />
                        </li>
                    ))}
                </ul>
            ) : (
                <p className="rounded-md border border-dashed border-input px-3 py-4 text-center text-sm text-muted-foreground">
                    直近の予定はありません
                </p>
            )}
        </div>
    );
};

export default UpcomingEvents;
