'use client'
import React, { useState } from 'react';
import axios from '@/lib/axios/public';
import { useForm, SubmitHandler } from 'react-hook-form';
import { useRouter } from 'next/navigation';
import type { EventDraftDetail } from '@/hooks/event/type';
import Modal from '@/components/Modal';
import Button from '@/components/Button';
import IconButton from '@/components/IconButton';
import { TrashIcon } from '@heroicons/react/20/solid';

interface EventDeleteProps {
    id: string;
    detail: EventDraftDetail;
}

const DeleteButton: React.FC<EventDeleteProps> = ({ id, detail }) => {
    const [isOpen, setIsOpen] = useState(false);
    const router = useRouter();

    const { handleSubmit } = useForm<EventDraftDetail>({
        defaultValues: detail
    });

    
    const deleteEvent = async (data: EventDraftDetail) => {
        return await axios.delete(`api/calendar/event/draft/${id}`, { data });
    }

    const onSubmit: SubmitHandler<EventDraftDetail> = (data) => {
        deleteEvent(data)
            .then(res => {
                console.log(res);
                setIsOpen(false);
                router.push('/schedule/draft');
            })
            .catch(err => {
                console.log(err);
            })
    }

    return (
        <div>
            <IconButton
                iconColor={'danger'}
                onClick={() => setIsOpen(true)}
            >
                <TrashIcon />
            </IconButton>
            <Modal
                isOpen={isOpen}
                description='このイベントを削除してよろしいですか？'
                onClose={() => setIsOpen(false)}
                actions={
                    <Button
                        type="submit"
                        intent="danger"
                        onClick={handleSubmit(onSubmit)}
                    >
                        削除
                    </Button>
                }
            />
        </div>
    )
}

export default DeleteButton;