'use client'
import React from 'react';
import { useParams } from 'next/navigation';
import { useAtom } from 'jotai';
import { titleAtomFamily } from '@/atoms/calendar';
import { useFormContext } from 'react-hook-form';
import Card from '@/components/Card';
import TextField from '@/components/TextField';
import TextArea from '@/components/TextArea';
import type { DiscriminatedEventForm } from './zod';

interface EventBasicFormProps {
    description?: string;
    location?: string;
}

const EventBasicForm: React.FC<EventBasicFormProps> =({ description, location }) => {
    const { id } = useParams<{ id?: string }>();
    const { register, formState: { errors } } = useFormContext<DiscriminatedEventForm>();
    const [title, setUpdateTitle] = useAtom(titleAtomFamily(id));

    return (
        <Card variant="outlined" background="inherit" className="w-full">
            <h2 className="text-lg font-bold text-gray-700 dark:text-gray-300">基本情報</h2>
            <p className="text-sm text-gray-500 mb-4">日程調整するイベントのタイトルや詳細を入力してください</p>
            <div className="space-y-6">
                <TextField
                    {...register('title')}
                    label="タイトル"
                    defaultValue={title}
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