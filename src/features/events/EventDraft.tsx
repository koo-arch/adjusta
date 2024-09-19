'use client'
import React, { useEffect } from 'react';
import axios from '@/lib/axios/public';
import { selectedDatesAtom, titleAtom, prioritizedSelectedDatesAtom } from '@/atoms/calendar';
import { useAtom } from 'jotai';
import { useForm, type SubmitHandler, FormProvider } from 'react-hook-form';
import type { EventDraftForm } from './type';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';
import DraggableEventList from './DraggableEventList';
import TextField from '@/components/TextField';

const EventDraft: React.FC = () => {
    const method = useForm<EventDraftForm>();
    const { register, handleSubmit, setValue } = method;
    const [selectedDates, setSelectedDates] = useAtom(selectedDatesAtom);
    const [prioritizedSelectedDate] = useAtom(prioritizedSelectedDatesAtom);
    const [title, setTitle] = useAtom(titleAtom);

    useEffect(() => {
        setValue('selected_dates', prioritizedSelectedDate);
    },[selectedDates, setValue]);

    const postEventDraft = async (data: EventDraftForm) => {
        return await axios.post('api/calendar/event/draft', data);
    }

    const onSubmit: SubmitHandler<EventDraftForm> = (data) => {
        postEventDraft(data)
            .then(res => {
                console.log(res);
                setSelectedDates([]);
                setTitle('');
            })
            .catch(err => {
                console.log(err);
                // Todo: エラーメッセージを表示する
            })
    }

    return (
        <div>
            <DraggableEventList />
            <FormProvider {...method}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <TextField
                        {...register('title')}
                        label="Title"
                        defaultValue={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />
                    <TextField
                        {...register('description')}
                        label="Description"
                    />
                    <TextField
                        {...register('location')}
                        label="Location"
                    />
                    <SelectableCalendar />
                    <button type="submit">Submit</button>
                </form>
            </FormProvider>
        </div>
    )
}

export default EventDraft;