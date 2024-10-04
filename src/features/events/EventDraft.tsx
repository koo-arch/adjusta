'use client'
import React, { useEffect } from 'react';
import axios from '@/lib/axios/public';
import { selectedDatesAtom, titleAtom, sendSelectedDatesAtom, selectedEventsAtom } from '@/atoms/calendar';
import { useAtom } from 'jotai';
import { useForm, type SubmitHandler, FormProvider } from 'react-hook-form';
import type { EventDraftForm } from './type';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';
import DraggableEventList from './DraggableEventList';
import TextField from '@/components/TextField';
import TextArea from '@/components/TextArea';

const EventDraft: React.FC = () => {
    const method = useForm<EventDraftForm>();
    const { register, handleSubmit, setValue, formState: { errors } } = method;
    const [selectedDates, setSelectedDates] = useAtom(selectedDatesAtom);
    const [sendSelectedDate] = useAtom(sendSelectedDatesAtom);
    const [title, setTitle] = useAtom(titleAtom);

    useEffect(() => {
        setValue('selected_dates', sendSelectedDate);
    }, [selectedDates, setValue]);

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
            <DraggableEventList atom={selectedDatesAtom} />
            <FormProvider {...method}>
                <form onSubmit={handleSubmit(onSubmit)}>
                    <TextField
                        {...register('title')}
                        label="Title"
                        defaultValue={title}
                        onChange={(e) => setTitle(e.target.value)}
                        error={!!errors.title}
                        helperText={errors.title?.message}
                    />
                    <TextArea
                        {...register('description')}
                        label="Description"
                        error={!!errors.description}
                        helperText={errors.description?.message}
                    />
                    <TextField
                        {...register('location')}
                        label="Location"
                        error={!!errors.location}
                        helperText={errors.location?.message}
                    />
                    <SelectableCalendar
                        dateAtom={selectedDatesAtom}
                        eventAtom={selectedEventsAtom}
                    />
                    <button type="submit">Submit</button>
                </form>
            </FormProvider>
        </div>
    )
}

export default EventDraft;