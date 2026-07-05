'use client'
import React from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import Card from '@/components/Card';
import TextField from '@/components/TextField';
import TextArea from '@/components/TextArea';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';

interface EventBasicFormProps {
    formScope: string;
}

const EventBasicForm: React.FC<EventBasicFormProps> = ({
    formScope,
}) => {
    const [title, setTitle] = useAtom(titleAtomFamily(formScope));
    const [description, setDescription] = useAtom(descriptionAtomFamily(formScope));
    const [location, setLocation] = useAtom(locationAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));

    return (
        <Card variant="outlined" background="inherit" className="w-full">
            <h2 className="text-lg font-bold text-gray-700 dark:text-gray-300">基本情報</h2>
            <p className="text-sm text-gray-500 mb-4">日程調整するイベントのタイトルや詳細を入力してください</p>
            <div className="space-y-6">
                <TextField
                    label="タイトル"
                    value={title}
                    error={!!errors.title}
                    helperText={errors.title}
                    onChange={(e) => {
                        setTitle(e.target.value);
                        clearEditedFieldState('title');
                    }}
                />
                <TextField
                    label="場所"
                    value={location}
                    error={!!errors.location}
                    helperText={errors.location}
                    onChange={(e) => {
                        setLocation(e.target.value);
                        clearEditedFieldState('location');
                    }}
                />
                <TextArea
                    label="説明"
                    value={description}
                    error={!!errors.description}
                    helperText={errors.description}
                    onChange={(e) => {
                        setDescription(e.target.value);
                        clearEditedFieldState('description');
                    }}
                />
            </div>
        </Card>
    );
}

export default EventBasicForm;
