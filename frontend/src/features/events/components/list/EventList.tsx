'use client'
import React from 'react';
import Link from 'next/link';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useSearchEvents } from '@/features/events/hooks/useSearchEvents';
import EventCard from './EventCard';
import { PaginationControls } from '@/components/common/pagination/PaginationControls';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { EVENT_STATUS_LABELS } from '@/features/events/status';
import { cn } from '@/lib/utils';
import { Plus } from 'lucide-react';

const STATUS_TABS = ['all', 'active', 'confirmed', 'draft', 'cancelled'] as const;

type StatusTab = (typeof STATUS_TABS)[number];

const TAB_LABELS: Record<StatusTab, string> = {
    all: 'すべて',
    ...EVENT_STATUS_LABELS,
};

const parseStatusTab = (value: string | null): StatusTab =>
    STATUS_TABS.includes(value as StatusTab) ? (value as StatusTab) : 'all';

const parsePage = (value: string | null): number => {
    const page = Number(value);
    return Number.isInteger(page) && page >= 1 ? page : 1;
};

const EventCardSkeleton = () => (
    <Card className="h-full">
        <CardHeader className="flex-row items-start justify-between gap-2 space-y-0">
            <Skeleton className="h-6 w-2/3" />
            <Skeleton className="h-5 w-16" />
        </CardHeader>
        <CardContent className="space-y-2">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-3/4" />
        </CardContent>
    </Card>
);

export const EventListSkeleton = () => (
    <div className="space-y-4">
        <Skeleton className="h-9 w-80 max-w-full" />
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, index) => (
                <EventCardSkeleton key={index} />
            ))}
        </div>
    </div>
);

const EventList: React.FC = () => {
    const router = useRouter();
    const pathname = usePathname();
    const searchParams = useSearchParams();
    const statusTab = parseStatusTab(searchParams.get('status'));
    const page = parsePage(searchParams.get('page'));

    const { searchEvents, pagination, isPending, isPlaceholderData, error, refetch } = useSearchEvents({
        ...(statusTab === 'all' ? {} : { status: statusTab }),
        page,
    });

    const buildHref = (tab: StatusTab, targetPage: number) => {
        const params = new URLSearchParams();
        if (tab !== 'all') {
            params.set('status', tab);
        }
        if (targetPage > 1) {
            params.set('page', String(targetPage));
        }
        const query = params.toString();
        return query ? `${pathname}?${query}` : pathname;
    };

    const handleTabChange = (value: string) => {
        // タブ切替はページを 1 に戻す。履歴は絞り込み条件で汚さない
        router.replace(buildHref(parseStatusTab(value), 1), { scroll: false });
    };

    const handlePageChange = (nextPage: number) => {
        router.push(buildHref(statusTab, nextPage));
    };

    const renderContent = () => {
        if (isPending) {
            return (
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {Array.from({ length: 6 }).map((_, index) => (
                        <EventCardSkeleton key={index} />
                    ))}
                </div>
            );
        }

        if (error) {
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

        if (!searchEvents || searchEvents.length === 0) {
            if (statusTab !== 'all') {
                return (
                    <Card>
                        <CardContent className="p-8 text-center">
                            <p className="text-sm text-muted-foreground">
                                「{TAB_LABELS[statusTab]}」のイベントはありません。
                            </p>
                        </CardContent>
                    </Card>
                );
            }
            return (
                <Card>
                    <CardContent className="flex flex-col items-center gap-4 p-8 text-center">
                        <div className="space-y-1">
                            <p className="font-medium">イベントがまだありません</p>
                            <p className="text-sm text-muted-foreground">
                                候補日程を登録して、日程調整を始めましょう。
                            </p>
                        </div>
                        <Button asChild>
                            <Link href="/events/new">
                                <Plus className="size-4" />
                                イベントを作成
                            </Link>
                        </Button>
                    </CardContent>
                </Card>
            );
        }

        return (
            <div>
                <ul
                    className={cn(
                        'grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3',
                        // 次ページ取得中は前ページの内容を薄めて表示する
                        isPlaceholderData && 'opacity-60',
                    )}
                >
                    {searchEvents.map((event) => (
                        <li key={event.id}>
                            <EventCard event={event} />
                        </li>
                    ))}
                </ul>
                {pagination && (
                    <PaginationControls
                        page={pagination.page}
                        total={pagination.total_items}
                        limit={pagination.per_page}
                        onPageChange={handlePageChange}
                    />
                )}
            </div>
        );
    };

    return (
        <div className="space-y-4">
            <Tabs value={statusTab} onValueChange={handleTabChange}>
                <div className="overflow-x-auto">
                    <TabsList>
                        {STATUS_TABS.map((tab) => (
                            <TabsTrigger key={tab} value={tab}>
                                {TAB_LABELS[tab]}
                            </TabsTrigger>
                        ))}
                    </TabsList>
                </div>
            </Tabs>
            {renderContent()}
        </div>
    );
};

export default EventList;
