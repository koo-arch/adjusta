'use client'
import React, { useState, useRef } from 'react';
import { toast } from 'react-toastify';
import type { EventClickArg, EventDropArg, DateSelectArg } from '@fullcalendar/core';
import type { EventResizeDoneArg } from '@fullcalendar/interaction';
import Calendar from '@/features/calendar/components/Calendar';
import { EventImpl } from '@fullcalendar/core/internal';
import PopupMenu from '@/components/PopupMenu';
import type { CalendarEvent } from '@/features/calendar/types';
import type { EventDraftDetail } from '@/features/events/types';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';

type DateSelectInfo = {
    id: string;
    start: Date;
    end: Date;
}

interface SelectableCalendarProps<TDate extends DateSelectInfo> {
    dates: TDate[];
    onDatesChange: React.Dispatch<React.SetStateAction<TDate[]>>;
    selectedEvents: CalendarEvent[];
    editingEvent?: EventDraftDetail;
}

const SelectableCalendar = <TDate extends SelectedDate | ProposedDate>({
    dates,
    onDatesChange,
    selectedEvents,
    editingEvent,
}:SelectableCalendarProps<TDate>) => {
    // モバイルは週ビュー(7列)が窮屈なため日ビューで開始する。マウント時に一度だけ判定
    // (FullCalendar は SSR ではほぼ描画されないため、サーバーとの差異は実害なし)
    const [initialView] = useState<'timeGridWeek' | 'timeGridDay'>(() =>
        typeof window !== 'undefined' && window.matchMedia('(max-width: 640px)').matches
            ? 'timeGridDay'
            : 'timeGridWeek',
    );
    // 現在時刻(赤線)が見える位置から開始する。少し手前から表示するため 2 時間戻す
    const [scrollTime] = useState(() => {
        const hour = Math.max(new Date().getHours() - 2, 0);
        return `${String(hour).padStart(2, '0')}:00`;
    });
    const [clickedEvent, setClickedEvent] = useState<EventImpl>();
    const [popupPosition, setPopupPosition] = useState({ top: 0, left: 0 });
    const buttonRef = useRef<HTMLButtonElement>(null);

    const handleDateSelect = (e: DateSelectArg) => {
        // 月ビューの選択は全日(時刻なし)になるため候補にしない
        if (e.allDay) {
            toast.info('月ビューでは時刻を選択できません。週ビューに切り替えるか「日時を追加」を使ってください', {
                toastId: 'all-day-select-info',
            });
            return;
        }

        const newDate = {
            id : new Date().getTime().toString(),
            start: e.start,
            end: e.end,
        } as TDate;
        onDatesChange((prev) => [...prev, newDate]);
    }

    // イベントをクリックした時にポップアップを表示する
    const handleEventClick = (e: EventClickArg) => {
        const event = dates.find((date) => date.id === e.event.id);

        if (event) {
            if (buttonRef.current) {
                buttonRef.current.click();
            }

            setClickedEvent(e.event);
            setPopupPosition({ top: e.jsEvent.pageY, left: e.jsEvent.pageX });
        }
    }

    // イベントのドラッグ＆ドロップ時の処理
    const handleEventDrop = (e: EventDropArg) => {
        const updatedDates = dates.map((date) => {
            if (date.id === e.event.id) {
                return {
                    ...date,
                    start: e.event.start || date.start,
                    end: e.event.end || date.end
                };
            }
            return date;
        });

        onDatesChange(updatedDates);
    }

    // イベントの開始・終了時間の変更時の処理
    const handleEventResize = (e: EventResizeDoneArg) => {
        const updatedDates = dates.map((date) => {
            if (date.id === e.event.id) {
                return {
                    ...date,
                    start: e.event.start || date.start,
                    end: e.event.end || date.end
                };
            }
            return date;
        });

        onDatesChange(updatedDates);
    }

    // イベントの削除時の処理
    const handleDeleteEvent = () => {
        if (clickedEvent) {
            clickedEvent.remove();
            onDatesChange((prev) => prev.filter((date) => date.id !== clickedEvent.id));
        }
    }

    return (
        <div>
            <Calendar
                initialView={initialView}
                // 固定高さで内部スクロールにし、現在時刻基準で開始する
                height="70vh"
                scrollTime={scrollTime}
                headerToolbar={{
                    left: 'prev,next today',
                    center: 'title',
                    // 離れた日付の候補選択のために週/月ビューを切替可能にする(ui-review P2 #7)
                    right: 'timeGridWeek,dayGridMonth',
                }}
                select={handleDateSelect}
                selectedEvents={selectedEvents}
                eventClick={handleEventClick}
                eventDrop={handleEventDrop}
                eventResize={handleEventResize}
                editEvent={editingEvent}
            />
            <PopupMenu
                items={[
                    { label: '削除', onClick: handleDeleteEvent },
                ]}
                position={popupPosition}
                buttonRef={buttonRef}
            />
        </div>
    )
}

export default SelectableCalendar;
