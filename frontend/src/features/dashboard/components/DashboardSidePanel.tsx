'use client'
import React, { useEffect } from 'react';
import { useAtom } from 'jotai';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useFetchNeedsActionDrafts } from '@/features/events/hooks/useFetchNeedsActionDrafts';
import { useFetchUpcomingEvents } from '@/features/events/hooks/useFetchUpcomingEvents';
import {
    dashboardPanelTabAtom,
    selectedDashboardEventAtom,
    type DashboardPanelTab,
} from '@/features/dashboard/store/selectedEvent';
import NeedsActionDrafts from './NeedsActionDrafts';
import UpcomingEvents from './UpcomingEvents';
import DashboardEventDetail from './DashboardEventDetail';
import { CalendarClock, CircleAlert, Info, MousePointerClick } from 'lucide-react';

// 右パネル: アイコンタブ(対応が必要/直近の予定/詳細)。
// カレンダーでイベントをクリックすると詳細タブへ自動で切り替わる
const DashboardSidePanel: React.FC = () => {
    const [selectedEvent, setSelectedEvent] = useAtom(selectedDashboardEventAtom);
    const [activeTab, setActiveTab] = useAtom(dashboardPanelTabAtom);
    // タブの目印用(セクション側と同じクエリなので追加リクエストは発生しない)
    const { needsActionDrafts } = useFetchNeedsActionDrafts();
    const { upcomingEvents } = useFetchUpcomingEvents();
    const hasNeedsAction = (needsActionDrafts?.length ?? 0) > 0;
    const hasUpcoming = (upcomingEvents?.length ?? 0) > 0;

    // 選択・タブ状態はグローバル atom のためページ離脱では消えない。
    // 古い状態を持ち越さないよう unmount 時にリセットする(購読解除相当の cleanup)
    useEffect(() => {
        return () => {
            setSelectedEvent(null);
            setActiveTab('needs-action');
        };
    }, [setSelectedEvent, setActiveTab]);

    return (
        <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as DashboardPanelTab)}>
            <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger
                    value="needs-action"
                    aria-label="対応が必要"
                    title="対応が必要"
                    className="relative"
                >
                    <CircleAlert className="size-4" />
                    {hasNeedsAction && (
                        <span
                            aria-hidden
                            className="absolute right-2 top-1 size-2 rounded-full bg-destructive"
                        />
                    )}
                </TabsTrigger>
                <TabsTrigger
                    value="upcoming"
                    aria-label="直近の予定"
                    title="直近の予定"
                    className="relative"
                >
                    <CalendarClock className="size-4" />
                    {hasUpcoming && (
                        <span
                            aria-hidden
                            className="absolute right-2 top-1 size-2 rounded-full bg-primary"
                        />
                    )}
                </TabsTrigger>
                <TabsTrigger value="detail" aria-label="イベント詳細" title="イベント詳細">
                    <Info className="size-4" />
                </TabsTrigger>
            </TabsList>
            <TabsContent value="needs-action" className="mt-3">
                <NeedsActionDrafts />
            </TabsContent>
            <TabsContent value="upcoming" className="mt-3">
                <UpcomingEvents />
            </TabsContent>
            <TabsContent value="detail" className="mt-3">
                {selectedEvent ? (
                    <DashboardEventDetail event={selectedEvent} />
                ) : (
                    <div className="flex flex-col items-center gap-2 rounded-md border border-dashed border-input px-3 py-8 text-center text-sm text-muted-foreground">
                        <MousePointerClick className="size-5" aria-hidden />
                        カレンダーのイベントをクリックすると、
                        <br />
                        ここに詳細が表示されます
                    </div>
                )}
            </TabsContent>
        </Tabs>
    );
};

export default DashboardSidePanel;
