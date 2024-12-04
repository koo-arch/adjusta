'use client'
import React, { useEffect } from 'react';
import { useAtom } from 'jotai';
import { titleAtom } from '@/atoms/calendar';
import { useFormContext } from 'react-hook-form';
import Card from '@/components/Card';
import TextField from '@/components/TextField';
import TextArea from '@/components/TextArea';
import type { EventDraftForm } from './type';

interface EventBasicFormProps {
    title?: string;
    description?: string;
    location?: string;
}

const EventBasicForm: React.FC<EventBasicFormProps> =({ title, description, location }) => {
    const { register, formState: { errors } } = useFormContext<EventDraftForm>();
    const [updateTitle, setUpdateTitle] = useAtom(titleAtom);

    useEffect(() => {
       if (title) {
           setUpdateTitle(title);
       }
    }, [updateTitle, setUpdateTitle, title]);

    return (
        <Card variant="outlined" background="inherit" className="w-full">
            <h2 className="text-lg font-bold text-gray-700 dark:text-gray-300">基本情報</h2>
            <p className="text-sm text-gray-500 mb-4">日程調整するイベントのタイトルや詳細を入力してください</p>
            <div className="space-y-6">
                <TextField
                    {...register('title')}
                    label="タイトル"
                    error={!!errors.title}
                    helperText={errors.title?.message}
                    onChange={(e) => setUpdateTitle(e.target.value)}
                />
                <TextField
                    {...register('location')}
                    label="場所"
                    defaultValue={location}
                    error={!!errors.location}
                    helperText={errors.location?.message}
                />
                <TextArea
                    {...register('description')}
                    label="説明"
                    defaultValue={description}
                    error={!!errors.description}
                    helperText={errors.description?.message}
                />
            </div>
        </Card>
    );
}

export default EventBasicForm;