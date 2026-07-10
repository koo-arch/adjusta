'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { proposedDatesAtomFamily, selectedDatesAtomFamily } from '@/features/events/store/calendar';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';
import DraggableDateList from '@/features/events/components/form/DraggableDateList';
import AddDateDialog from '@/features/events/components/form/AddDateDialog';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';

type DraftSelectEventListProps = {
    formType: 'draft';
    formScope: string;
};

type EditSelectEventListProps = {
    formType: 'edit';
    formScope: string;
};

type SelectEventListProps = DraftSelectEventListProps | EditSelectEventListProps;

const SectionHeader: React.FC<{ onAdd: (date: { start: Date; end: Date }) => void }> = ({ onAdd }) => (
    <div className="flex flex-wrap items-start justify-between gap-2">
        <div>
            <h2 className="text-lg font-bold leading-snug tracking-normal text-gray-900">選択日程</h2>
            <p className="mt-1 text-sm text-muted-foreground">
                ドラッグまたは矢印ボタンで優先順位を入れ替えられます
            </p>
        </div>
        <AddDateDialog onAdd={onAdd} />
    </div>
);

const EmptyDates = () => (
    <div className="rounded-md border border-dashed border-input py-10 text-center text-sm text-muted-foreground">
        候補日程がまだありません。
        <br />
        カレンダーで範囲を選択するか、「日時を追加」から登録してください。
    </div>
);

const DraftSelectEventList: React.FC<DraftSelectEventListProps> = (props) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(props.formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(props.formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(props.formScope));

    const handleAdd = ({ start, end }: { start: Date; end: Date }) => {
        const newDate: SelectedDate = { id: new Date().getTime().toString(), start, end };
        setDates((prev) => [...prev, newDate]);
        clearEditedFieldState('selected_dates');
    };

    return (
        <section className="space-y-4">
            <SectionHeader onAdd={handleAdd} />
            {errors.selected_dates && (
                <p className="text-sm text-destructive">{errors.selected_dates}</p>
            )}
            {dates.length > 0 ? (
                <DraggableDateList
                    dates={dates}
                    onDatesChange={(value) => {
                        setDates(value);
                        clearEditedFieldState('selected_dates');
                    }}
                />
            ) : (
                <EmptyDates />
            )}
        </section>
    );
};

const EditSelectEventList: React.FC<EditSelectEventListProps> = (props) => {
    const [dates, setDates] = useAtom(proposedDatesAtomFamily(props.formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(props.formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(props.formScope));

    const handleAdd = ({ start, end }: { start: Date; end: Date }) => {
        setDates((prev) => {
            const newDate: ProposedDate = {
                id: new Date().getTime().toString(),
                start,
                end,
                priority: prev.length + 1,
            };
            return [...prev, newDate];
        });
        clearEditedFieldState('proposed_dates');
    };

    return (
        <section className="space-y-4">
            <SectionHeader onAdd={handleAdd} />
            {errors.proposed_dates && (
                <p className="text-sm text-destructive">{errors.proposed_dates}</p>
            )}
            {dates.length > 0 ? (
                <DraggableDateList
                    dates={dates}
                    onDatesChange={(value) => {
                        setDates(value);
                        clearEditedFieldState('proposed_dates');
                    }}
                />
            ) : (
                <EmptyDates />
            )}
        </section>
    );
};

const SelectEventList: React.FC<SelectEventListProps> = (props) => {
    if (props.formType === 'draft') {
        return <DraftSelectEventList {...props} />;
    }

    return <EditSelectEventList {...props} />;
};

export default SelectEventList;
