'use client'
import React from 'react';
import { useSetAtom, useAtomValue } from 'jotai';
import { allEventsAtom } from '@/features/events/store/calendar';
import { dashboardPanelTabAtom, selectedDashboardEventAtom } from '@/features/dashboard/store/selectedEvent';
import Calendar from '@/features/calendar/components/Calendar';
import type { EventClickArg } from '@fullcalendar/core';

const DashboardCalendar: React.FC = () => {
    const allEvents = useAtomValue(allEventsAtom);
    const setSelectedEvent = useSetAtom(selectedDashboardEventAtom);
    const setPanelTab = useSetAtom(dashboardPanelTabAtom);

    const handleEventClick = (e: EventClickArg) => {
        const event = allEvents.find(event => event.id === e.event.id);
        if (!event) {
            return;
        }
        setSelectedEvent(event);
        setPanelTab('detail');
        // 縦積み時(lg 未満)はパネルが画面外にあるため、詳細まで滑らかに移動する
        document.getElementById('dashboard-side-panel')?.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }

    return (
        /* ダッシュボードのカレンダーは閲覧専用(選択・ドラッグ・リサイズ不可。ui-review P1 #3) */
        <Calendar
            selectable={false}
            editable={false}
            eventClick={handleEventClick}
        />
    );
}

export default DashboardCalendar;
