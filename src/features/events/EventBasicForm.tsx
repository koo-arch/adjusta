'use client'
import React, { useEffect } from 'react';
import { useAtom } from 'jotai';
import { updateTitleAtom } from '@/atoms/calendar';
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
    const [updateTitle, setUpdateTitle] = useAtom(updateTitleAtom);

    useEffect(() => {
       if (title) {
           setUpdateTitle(title);
       }
    }, [updateTitle, setUpdateTitle]);

    return (
        <Card variant="outlined" background="inherit" className="w-full">
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