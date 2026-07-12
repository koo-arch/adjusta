'use client'
import React, { useEffect } from 'react';
import { useAtom } from 'jotai';
import { toast } from 'react-toastify';
import { allEventsAtom } from '@/features/events/store/calendar';
import { StyleWrapper } from './style';
import FullCalendar from '@fullcalendar/react';
import type { ToolbarInput, DateRangeInput, EventClickArg, EventDropArg, DateSelectArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import timeGridPlugin from '@fullcalendar/timegrid';
import jaLocale from '@fullcalendar/core/locales/ja';
import momentPlugin from '@fullcalendar/moment';
import { useFetchGoogleCalendarEvents } from '@/features/calendar/hooks/useFetchGoogleCalendarEvents';
import { useSearchEvents } from '@/features/events/hooks/useSearchEvents';
import { renderDayCell, renderDayHeader, renderSlotLabel } from './render';
import type { CalendarEvent } from '../types';
import type { EventDraftDetail } from '@/features/events/types';


interface CalendarProps<T extends CalendarEvent> {
    initialView?: "dayGridMonth" | "dayGridWeek" | "timeGridWeek" | "timeGridDay";
    firstDay?: number;
    headerToolbar?: ToolbarInput;
    select?: (arg: DateSelectArg) => void;
    selectedEvents?: T[];
    visibleRange?: DateRangeInput;
    eventClick?: (e: EventClickArg) => void;
    eventDrop?: (e: EventDropArg) => void;
    eventResize?: (e: EventResizeDoneArg) => void;
    editEvent?: EventDraftDetail;
    // 'auto' は全展開(内部スクロールなし)。固定値を渡すと timeGrid が内部スクロールになり scrollTime が効く
    height?: number | string;
    scrollTime?: string;
}

const Calendar = <T extends CalendarEvent>({
    initialView = 'dayGridMonth',
    firstDay,
    headerToolbar,
    select,
    selectedEvents,
    visibleRange,
    eventClick,
    eventDrop,
    eventResize,
    editEvent,
    height = 'auto',
    scrollTime,
}: CalendarProps<T>) => {
    const { events, isLoading: isGoogleEventLoading } = useFetchGoogleCalendarEvents();
    const { searchEvents, isPending: isSearchLoading } = useSearchEvents({ status: "active" });
    const [allEvents, setAllEvents] = useAtom(allEventsAtom);
    const confirmedGoogleEventID = editEvent?.confirmed_google_event_id ?? editEvent?.google_event_id;

    const warningToastId = 'google-calendar-warning';

    // FullCalendar へ渡す表示イベント一覧を、取得済みデータと編集中のローカル状態から合成する。
    useEffect(() => {
        if (isGoogleEventLoading || isSearchLoading) return;

        // Google 由来の予定はニュートラル系で区別する(DESIGN.md「Third-Party Components」)
        const googleEventList: CalendarEvent[]  = events?.events
            .filter(ge => !confirmedGoogleEventID || confirmedGoogleEventID !== ge.id)
            .map(event => ({
                id: event.id,
                title: event.summary,
                start: event.start,
                end: event.end,
                location: event.location,
                description: event.description,
                origin: "google",
                local_event_id: null,
                backgroundColor: '#e5e7eb',
                borderColor: '#e5e7eb',
                textColor: '#374151',
            })) ?? [];

        // 他イベントの調整中候補は薄いインディゴ。編集中イベント自身の候補
        // (selectedEvents 経由)は既定の Primary 塗りのままにして区別する
        const searchEventList: CalendarEvent[] = searchEvents?.flatMap(event =>
            event.proposed_dates
                .filter(date => !editEvent?.proposed_dates?.some(edit => edit.id === date.id))
                .map(date => ({
                    id: date.id,
                    title: event.title,
                    start: date.start,
                    end: date.end,
                    location: event.location,
                    description: event.description,
                    origin: "local",
                    local_event_id: event.id,
                    backgroundColor: '#e0e7ff',
                    borderColor: '#c7d2fe',
                    textColor: '#4338ca',
                }))
        ) ?? [];

        const allEvents: CalendarEvent[] = [...googleEventList, ...searchEventList, ...(selectedEvents || [])];
        setAllEvents(allEvents);
    }, [events, searchEvents, isGoogleEventLoading, isSearchLoading, editEvent, confirmedGoogleEventID, setAllEvents, selectedEvents]);

    // Google Calendar API から部分的な取得失敗が返ったときだけ警告トーストを出す。
    useEffect(() => {
        if (events?.warning?.failed_calendars) {
            toast.warn(`取得に失敗したカレンダーがあります: ${events.warning.failed_calendars.join(', ')}`,{
                toastId: warningToastId,
            });
        }
    }, [events?.warning])
    

    return (
        <div>
            <StyleWrapper>
                <FullCalendar
                    plugins={[dayGridPlugin, timeGridPlugin, interactionPlugin, momentPlugin]}
                    initialView={initialView}
                    firstDay={firstDay}
                    headerToolbar={headerToolbar || {
                        left: 'prev,next today',
                        center: 'title',
                        right: 'dayGridMonth,timeGridWeek,timeGridDay'
                    }}
                    businessHours={{ daysOfWeek: [1, 2, 3, 4, 5] }}
                    eventClick={eventClick || (() => {})}
                    snapDuration={'00:10:00'}
                    height={height}
                    nowIndicator={true}
                    scrollTime={scrollTime}
                    scrollTimeReset={false}
                    selectable={true}
                    selectMirror={true}
                    editable={true} // イベントのドラッグ＆ドロップを可能に
                    eventDrop={eventDrop || (() => {})}
                    eventResizableFromStart={true} // イベントの開始時間もリサイズ可能にする
                    eventResize={eventResize || (() => {})}
                    select={select || (() => {})}
                    visibleRange={visibleRange}
                    dayCellContent={renderDayCell}
                    dayHeaderContent={renderDayHeader}
                    slotLabelContent={renderSlotLabel}
                    aspectRatio={1.6}
                    locale={jaLocale}
                    events={allEvents}
                />
            </StyleWrapper>
        </div>
    );
}

export default Calendar;
