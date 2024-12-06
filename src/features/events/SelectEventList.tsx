'use client'
import React from 'react';
import { useParams } from 'next/navigation';
import DraggableEventList from '@/features/events/DraggableEventList';
import { useAtomValue, useAtom } from 'jotai';
import { selectedDatesAtom, proposedDatesAtom } from '@/atoms/calendar';
import Card from '@/components/Card';
import ToggleSwitch from '@/components/ToggleSwitch';
import { isConfirmedAtomFamily } from '@/atoms/event';

interface SelectEventListProps {
    formType: 'draft' | 'edit';
}

const SelectEventList: React.FC<SelectEventListProps> = ({ formType }) => {
    const { id } = useParams<{ id?: string }>();
    const [isConfirmed, setIsConfirmed] = useAtom(isConfirmedAtomFamily(id));
    const selectedDates = useAtomValue(selectedDatesAtom);
    const proposedDates = useAtomValue(proposedDatesAtom);

    const dates = formType === 'draft' ? selectedDates : proposedDates;

    return (
        <Card variant="outlined" background="inherit">
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-bold mb-2">選択日程</h2>
                {formType == 'edit' && (
                    <ToggleSwitch
                        checked={isConfirmed}
                        onChange={() => setIsConfirmed(!isConfirmed)}
                        label="候補日程の確定"
                    />
                )}
            </div>
            <p className="text-sm text-gray-500 mb-4">ドラッグで優先順位の入れ替えができます</p>
           {dates.length > 0 ? (
                <DraggableEventList
                    atom={formType === 'draft' ? selectedDatesAtom : proposedDatesAtom}
                    enableTopHighlight={formType === 'edit' && isConfirmed}
                />
            ) : (
                <p className="font-bold py-16 text-center">日程を選択してください</p>
            )}
        </Card>
    )
}

export default SelectEventList;