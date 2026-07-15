'use client'
import React from 'react';
import { STATUS_TABS, useEventListSearch, type StatusTab } from '@/features/events/hooks/useEventListSearch';
import EventsToolbar from './EventsToolbar';
import EventGrid from './EventGrid';
import { PaginationControls } from '@/components/common/pagination/PaginationControls';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { EVENT_STATUS_LABELS } from '@/features/events/status';

const TAB_LABELS: Record<StatusTab, string> = {
    all: 'すべて',
    ...EVENT_STATUS_LABELS,
};

const TOOLBAR_TABS = STATUS_TABS.map((tab) => ({ value: tab, label: TAB_LABELS[tab] }));

// 「疎な状態」(絞り込みなし・1ページ目・lg の1行未満)でのみ作成プレースホルダを出す
const SPARSE_THRESHOLD = 3;

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
        <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <Skeleton className="h-9 w-80 max-w-full" />
            <div className="flex items-center gap-2">
                <Skeleton className="h-10 w-full md:w-64" />
                <Skeleton className="hidden h-10 w-28 md:block" />
            </div>
        </div>
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
        page,
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

    const isUnfilteredFirstPage = statusTab === 'all' && title === '' && page === 1;

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
                <div className="space-y-4">
                    <div className="space-y-1 pt-4 text-center">
                        <p className="font-medium">イベントがまだありません</p>
                        <p className="text-sm text-muted-foreground">
                            候補日程を登録して、日程調整を始めましょう。
                        </p>
                    </div>
                    <EventGrid events={[]} showCreatePlaceholder />
                </div>
            );
        }

        return (
            <div>
                <EventGrid
                    events={searchEvents}
                    showCreatePlaceholder={isUnfilteredFirstPage && searchEvents.length < SPARSE_THRESHOLD}
                    isDimmed={isPlaceholderData}
                />
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
            <EventsToolbar
                tabs={TOOLBAR_TABS}
                activeTab={statusTab}
                onTabChange={selectTab}
                searchValue={title}
                onSearch={search}
            />
            {renderContent()}
        </div>
    );
};

export default EventList;
