'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { Button } from '@/components/ui/button';
import {
    locationAtomFamily,
    selectedDatesAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { formStepAtomFamily } from '@/features/events/store/formStep';
import { buildDefaultCandidateSpan, type SelectedDate } from '@/features/events/store/dates';
import DatesPanelView from '@/features/events/components/form/DatesPanelView';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';
import { useCandidateSyncSetting } from '@/features/auth/hooks/useCandidateSyncSetting';

interface DraftDatesPanelProps {
    formScope: string;
}

const DraftDatesPanel: React.FC<DraftDatesPanelProps> = ({ formScope }) => {
    const [dates, setDates] = useAtom(selectedDatesAtomFamily(formScope));
    const title = useAtomValue(titleAtomFamily(formScope));
    const location = useAtomValue(locationAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));
    const setStep = useSetAtom(formStepAtomFamily(formScope));
    const candidateSync = useCandidateSyncSetting();

    const handleDatesChange: React.Dispatch<React.SetStateAction<SelectedDate[]>> = (value) => {
        setDates(value);
        clearEditedFieldState('selected_dates');
    };

    const handleAdd = () => {
        const newDate: SelectedDate = { id: new Date().getTime().toString(), ...buildDefaultCandidateSpan() };
        handleDatesChange((prev) => [...prev, newDate]);
        return newDate.id;
    };

    return (
        <div className="space-y-4">
            <DatesPanelView
                title={title}
                location={location}
                dates={dates}
                onDatesChange={handleDatesChange}
                onAdd={handleAdd}
                onEditBasic={() => setStep('basic')}
                error={errors.selected_dates}
            />
            {dates.length > 0 && !candidateSync.isLoading && !candidateSync.setting?.enabled && (
                <div className="rounded-md border border-primary/30 bg-primary/5 p-3 text-sm">
                    <p className="font-medium text-foreground">候補日程をカレンダーにも表示</p>
                    <p className="mt-1 text-muted-foreground">専用カレンダーに仮予定として追加できます。</p>
                    <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        className="mt-3"
                        disabled={candidateSync.isUpdating}
                        onClick={() => candidateSync.setEnabled(true)}
                    >
                        {candidateSync.isUpdating ? '有効化しています…' : 'この場で有効にする'}
                    </Button>
                </div>
            )}
        </div>
    );
};

export default DraftDatesPanel;
