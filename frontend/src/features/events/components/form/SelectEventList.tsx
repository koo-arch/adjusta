'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { proposedDatesAtomFamily, selectedDatesAtomFamily } from '@/features/events/store/calendar';
import { isConfirmedAtomFamily } from '@/features/events/store/confirmation';
import DraggableDateList from '@/features/events/components/form/DraggableDateList';
import Card from '@/components/Card';
import ToggleSwitch from '@/components/ToggleSwitch';
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

const DraftSelectEventList: React.FC<DraftSelectEventListProps> = (props) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(props.formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(props.formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(props.formScope));

    return (
        <Card variant="outlined" background="inherit">
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-bold mb-2">選択日程</h2>
            </div>
            <p className="text-sm text-gray-500 mb-4">ドラッグで優先順位の入れ替えができます</p>

            {errors.selected_dates && (
                <p className="text-sm text-red-500 mb-4">{errors.selected_dates}</p>
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
                <p className="font-bold py-16 text-center">日程を選択してください</p>
            )}
        </Card>
    );
};

const EditSelectEventList: React.FC<EditSelectEventListProps> = (props) => {
    const [dates, setDates] = useAtom(proposedDatesAtomFamily(props.formScope));
    const [isConfirmed, setIsConfirmed] = useAtom(isConfirmedAtomFamily(props.formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(props.formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(props.formScope));

    return (
        <Card variant="outlined" background="inherit">
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-bold mb-2">選択日程</h2>
                <ToggleSwitch
                    checked={isConfirmed}
                    onChange={(checked) => {
                        setIsConfirmed(checked);
                        clearEditedFieldState('confirmed');
                    }}
                    label="候補日程の確定"
                />
            </div>
            <p className="text-sm text-gray-500 mb-4">ドラッグで優先順位の入れ替えができます</p>

            {errors.proposed_dates && (
                <p className="text-sm text-red-500 mb-4">{errors.proposed_dates}</p>
            )}
            {dates.length > 0 ? (
                <DraggableDateList
                    dates={dates}
                    onDatesChange={(value) => {
                        setDates(value);
                        clearEditedFieldState('proposed_dates');
                    }}
                    enableTopHighlight={isConfirmed}
                />
            ) : (
                <p className="font-bold py-16 text-center">日程を選択してください</p>
            )}
        </Card>
    );
};

const SelectEventList: React.FC<SelectEventListProps> = (props) => {
    if (props.formType === 'draft') {
        return <DraftSelectEventList {...props} />;
    }

    return <EditSelectEventList {...props} />;
};

export default SelectEventList;
