'use client'
import React from 'react';
import Link from 'next/link';
import { STATUS_TABS, useEventListSearch, type StatusTab } from '@/features/events/hooks/useEventListSearch';
import EventCard from './EventCard';
import EventSearchForm from './EventSearchForm';
import { PaginationControls } from '@/components/common/pagination/PaginationControls';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { EVENT_STATUS_LABELS } from '@/features/events/status';
import { cn } from '@/lib/utils';
import { Plus } from 'lucide-react';

const TAB_LABELS: Record<StatusTab, string> = {
    all: 'すべて',
    ...EVENT_STATUS_LABELS,
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
    const {
        statusTab,
        title,
        selectTab,
        search,
        goToPage,
        searchEvents,
        pagination,
        isPending,
        isPlaceholderData,
        error,
        refetch,
    } = useEventListSearch();

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
            if (title !== '') {
                return (
                    <Card>
                        <CardContent className="flex flex-col items-center gap-4 p-8 text-center">
                            <p className="text-sm text-muted-foreground">
                                「{title}」に一致するイベントはありません。
                            </p>
                            <Button variant="outline" onClick={() => search('')}>
                                検索条件をクリア
                            </Button>
                        </CardContent>
                    </Card>
                );
            }
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
                        onPageChange={goToPage}
                    />
                )}
            </div>
        );
    };

    return (
        <div className="space-y-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
                <Tabs value={statusTab} onValueChange={selectTab}>
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
                <EventSearchForm defaultValue={title} onSearch={search} />
            </div>
            {renderContent()}
        </div>
    );
};

export default EventList;
