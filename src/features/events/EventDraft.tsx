'use client'
import React, { useEffect } from 'react';
import { useParams } from 'next/navigation';
import axios from '@/lib/axios/public';
import { selectedDatesAtom, sendSelectedDatesAtom, titleAtomFamily } from '@/atoms/calendar';
import { useAtom } from 'jotai';
import { useForm, type SubmitHandler, FormProvider } from 'react-hook-form';
import { type DiscriminatedEventForm, DiscriminatedEventFormResolver } from './zod';
import EventForm from './EventForm';

const EventDraft: React.FC = () => {
    const { id } = useParams<{ id?: string }>()
    const method = useForm<DiscriminatedEventForm>({
        resolver: DiscriminatedEventFormResolver,
        defaultValues: {
            form_type: "draft",
        }
    });
    const { handleSubmit, setValue } = method;
    const [selectedDates, setSelectedDates] = useAtom(selectedDatesAtom);
    const [sendSelectedDate] = useAtom(sendSelectedDatesAtom);
    const [,setTitle] = useAtom(titleAtomFamily(id));

    useEffect(() => {
        setValue('selected_dates', sendSelectedDate);
    }, [selectedDates, setValue, sendSelectedDate]);

    const postEventDraft = async (data: DiscriminatedEventForm) => {
        return await axios.post('api/calendar/event/draft', data);
    }

    const onSubmit: SubmitHandler<DiscriminatedEventForm> = (data) => {
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