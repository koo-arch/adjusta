'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    proposedDatesAtomFamily,
    proposedEventsAtomFamily,
    selectedDatesAtomFamily,
    selectedEventsAtomFamily,
} from '@/features/events/store/calendar';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';
import SelectableCalendar from '@/features/calendar/components/SelectableCalendar';
import DraggableDateList from '@/features/events/components/form/DraggableDateList';
import AddDateDialog from '@/features/events/components/form/AddDateDialog';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';
import type { EventDraftDetail } from '@/features/events/types';

type DraftDateSelectionSectionProps = {
    formType: 'draft';
    formScope: string;
    editingEvent?: EventDraftDetail;
};

type EditDateSelectionSectionProps = {
    formType: 'edit';
    formScope: string;
    editingEvent?: EventDraftDetail;
};

type DateSelectionSectionProps = DraftDateSelectionSectionProps | EditDateSelectionSectionProps;

const EmptyDates = () => (
    <div className="rounded-md border border-dashed border-input py-10 text-center text-sm text-muted-foreground">
        候補日程がまだありません。
        <br />
        カレンダーで範囲を選択するか、「日時を追加」から登録してください。
    </div>
);

// カレンダー(操作)と選択済みリスト(結果)を「候補日程」という 1 つのセクションにまとめる
const SectionLayout: React.FC<{
    onAdd: (date: { start: Date; end: Date }) => void;
    error?: string;
    calendar: React.ReactNode;
    list: React.ReactNode;
}> = ({ onAdd, error, calendar, list }) => (
    <section className="space-y-4">
        <div className="flex flex-wrap items-start justify-between gap-2">
            <div>
                <h2 className="text-lg font-bold leading-snug tracking-normal text-gray-900">候補日程</h2>
                <p className="mt-1 text-sm text-muted-foreground">
                    カレンダーから選ぶか、「日時を追加」で登録します。ドラッグまたは矢印ボタンで優先順位を変更できます
                </p>
            </div>
            <AddDateDialog onAdd={onAdd} />
        </div>
        {error && <p className="text-sm text-destructive">{error}</p>}
        <div className="grid grid-cols-1 gap-6 md:grid-cols-10">
            <div className="md:col-span-6">{calendar}</div>
            <div className="md:col-span-4">{list}</div>
        </div>
    </section>
);

const DraftDateSelectionSection: React.FC<DraftDateSelectionSectionProps> = ({ formScope, editingEvent }) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(formScope));
    const selectedEvents = useAtomValue(selectedEventsAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    const handleDatesChange: React.Dispatch<React.SetStateAction<SelectedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('selected_dates');
    };

    const handleAdd = ({ start, end }: { start: Date; end: Date }) => {
        const newDate: SelectedDate = { id: new Date().getTime().toString(), start, end };
        handleDatesChange((prev) => [...prev, newDate]);
    };

    return (
        <SectionLayout
            onAdd={handleAdd}
            error={errors.selected_dates}
            calendar={
                <SelectableCalendar
                    dates={dates}
                    onDatesChange={handleDatesChange}
                    selectedEvents={selectedEvents}
                    editingEvent={editingEvent}
                />
            }
            list={
                dates.length > 0 ? (
                    <DraggableDateList dates={dates} onDatesChange={handleDatesChange} />
                ) : (
                    <EmptyDates />
                )
            }
        />
    );
};

const EditDateSelectionSection: React.FC<EditDateSelectionSectionProps> = ({ formScope, editingEvent }) => {
    const [dates, setDates] = useAtom(proposedDatesAtomFamily(formScope));
    const proposedEvents = useAtomValue(proposedEventsAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    const handleDatesChange: React.Dispatch<React.SetStateAction<ProposedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('proposed_dates');
    };

    const handleAdd = ({ start, end }: { start: Date; end: Date }) => {
        handleDatesChange((prev) => {
            const newDate: ProposedDate = {
                id: new Date().getTime().toString(),
                start,
                end,
                priority: prev.length + 1,
            };
            return [...prev, newDate];
        });
    };

    return (
        <SectionLayout
            onAdd={handleAdd}
            error={errors.proposed_dates}
            calendar={
                <SelectableCalendar
                    dates={dates}
                    onDatesChange={handleDatesChange}
                    selectedEvents={proposedEvents}
                    editingEvent={editingEvent}
                />
            }
            list={
                dates.length > 0 ? (
                    <DraggableDateList dates={dates} onDatesChange={handleDatesChange} />
                ) : (
                    <EmptyDates />
                )
            }
        />
    );
};

const DateSelectionSection: React.FC<DateSelectionSectionProps> = (props) => {
    if (props.formType === 'draft') {
        return <DraftDateSelectionSection {...props} />;
    }

    return <EditDateSelectionSection {...props} />;
};

export default DateSelectionSection;
