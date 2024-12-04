'use client'
import React, { useEffect } from 'react';
import axios from '@/lib/axios/public';
import { selectedDatesAtom, titleAtom, sendSelectedDatesAtom } from '@/atoms/calendar';
import { useAtom } from 'jotai';
import { useForm, type SubmitHandler, FormProvider } from 'react-hook-form';
import type { EventDraftForm } from './type';
import EventForm from './EventForm';

const EventDraft: React.FC = () => {
    const method = useForm<EventDraftForm>();
    const { handleSubmit, setValue, formState: { errors } } = method;
    const [selectedDates, setSelectedDates] = useAtom(selectedDatesAtom);
    const [sendSelectedDate] = useAtom(sendSelectedDatesAtom);
    const [title, setTitle] = useAtom(titleAtom);

    useEffect(() => {
        setValue('selected_dates', sendSelectedDate);
    }, [selectedDates, setValue, sendSelectedDate]);

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
            <FormProvider {...method}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <EventForm />
                </form>
            </FormProvider>
        </div>
    )
}

export default EventDraft;