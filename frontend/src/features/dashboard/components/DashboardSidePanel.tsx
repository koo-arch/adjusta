'use client'
import React from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useFetchNeedsActionDrafts } from '@/features/events/hooks/useFetchNeedsActionDrafts';
import { useFetchUpcomingEvents } from '@/features/events/hooks/useFetchUpcomingEvents';
import NeedsActionDrafts from './NeedsActionDrafts';
import UpcomingEvents from './UpcomingEvents';
import { CalendarClock, CircleAlert } from 'lucide-react';

// 右パネル: タブ(対応が必要 / 直近の予定)。項目があるタブには目印ドットを出す
const DashboardSidePanel: React.FC = () => {
    // タブの目印用(セクション側と同じクエリなので追加リクエストは発生しない)
    const { needsActionDrafts } = useFetchNeedsActionDrafts();
    const { upcomingEvents } = useFetchUpcomingEvents();
    const hasNeedsAction = (needsActionDrafts?.length ?? 0) > 0;
    const hasUpcoming = (upcomingEvents?.length ?? 0) > 0;

    return (
        <Tabs defaultValue="needs-action">
            <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="needs-action" className="relative gap-1.5">
                    <CircleAlert className="size-4" />
                    対応が必要
                    {hasNeedsAction && (
                        <span
                            aria-hidden
                            className="absolute right-2 top-1 size-2 rounded-full bg-destructive"
                        />
                    )}
                </TabsTrigger>
                <TabsTrigger value="upcoming" className="relative gap-1.5">
                    <CalendarClock className="size-4" />
                    直近の予定
                    {hasUpcoming && (
                        <span
                            aria-hidden
                            className="absolute right-2 top-1 size-2 rounded-full bg-primary"
                        />
                    )}
                </TabsTrigger>
            </TabsList>
            <TabsContent value="needs-action" className="mt-3">
                <NeedsActionDrafts />
            </TabsContent>
            <TabsContent value="upcoming" className="mt-3">
                <UpcomingEvents />
            </TabsContent>
        </Tabs>
    );
};

export default DashboardSidePanel;
