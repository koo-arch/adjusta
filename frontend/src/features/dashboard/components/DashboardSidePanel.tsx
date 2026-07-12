'use client'
import React, { useEffect } from 'react';
import { useAtom } from 'jotai';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useFetchNeedsActionDrafts } from '@/features/events/hooks/useFetchNeedsActionDrafts';
import { useFetchUpcomingEvents } from '@/features/events/hooks/useFetchUpcomingEvents';
import { selectedDashboardEventAtom } from '@/features/dashboard/store/selectedEvent';
import NeedsActionDrafts from './NeedsActionDrafts';
import UpcomingEvents from './UpcomingEvents';
import DashboardEventDetail from './DashboardEventDetail';

// 右パネル: 通常はタブ(対応が必要/直近の予定)、カレンダーでイベントを
// クリックしている間はその詳細に切り替わる
const DashboardSidePanel: React.FC = () => {
    const [selectedEvent, setSelectedEvent] = useAtom(selectedDashboardEventAtom);
    // タブの目印用(セクション側と同じクエリなので追加リクエストは発生しない)
    const { needsActionDrafts } = useFetchNeedsActionDrafts();
    const { upcomingEvents } = useFetchUpcomingEvents();
    const hasNeedsAction = (needsActionDrafts?.length ?? 0) > 0;
    const hasUpcoming = (upcomingEvents?.length ?? 0) > 0;

    // 選択状態はグローバル atom のためページ離脱では消えない。
    // 古い選択を持ち越さないよう unmount 時にリセットする(購読解除相当の cleanup)
    useEffect(() => {
        return () => setSelectedEvent(null);
    }, [setSelectedEvent]);

    if (selectedEvent) {
        return <DashboardEventDetail event={selectedEvent} onClose={() => setSelectedEvent(null)} />;
    }

    return (
        <Tabs defaultValue="needs-action">
            <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="needs-action" className="relative">
                    対応が必要
                    {hasNeedsAction && (
                        <span
                            aria-hidden
                            title="対応が必要なイベントがあります"
                            className="absolute right-2 top-1.5 size-2 rounded-full bg-destructive"
                        />
                    )}
                </TabsTrigger>
                <TabsTrigger value="upcoming" className="relative">
                    直近の予定
                    {hasUpcoming && (
                        <span
                            aria-hidden
                            title="直近の予定があります"
                            className="absolute right-2 top-1.5 size-2 rounded-full bg-primary"
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
