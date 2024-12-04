import React from 'react';
import DraggableEventList from '@/features/events/DraggableEventList';
import { useAtomValue } from 'jotai';
import { selectedDatesAtom } from '@/atoms/calendar';
import Card from '@/components/Card';

const SelectEventList: React.FC = () => {
    const selectedDates = useAtomValue(selectedDatesAtom);

    return (
        <Card variant="outlined" background="inherit">
            <h2 className="text-lg font-bold mb-2">選択日程</h2>
            <p className="text-sm text-gray-500 mb-4">ドラッグで優先順位の入れ替えができます</p>
           {selectedDates.length > 0 ? (
               <DraggableEventList atom={selectedDatesAtom}/>
            ) : (
                <p className="font-bold py-16 text-center">日程を選択してください</p>
            )}
        </Card>
    )
}

export default SelectEventList;