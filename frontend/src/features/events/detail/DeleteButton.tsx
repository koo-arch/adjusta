'use client'
import React, { useState } from 'react';
import axios from '@/lib/axios/public';
import { useRouter } from 'next/navigation';
import Modal from '@/components/Modal';
import Button from '@/components/Button';
import IconButton from '@/components/IconButton';
import { TrashIcon } from '@heroicons/react/20/solid';

interface EventDeleteProps {
    eventID: string;
}

const DeleteButton: React.FC<EventDeleteProps> = ({ eventID }) => {
    const [isOpen, setIsOpen] = useState(false);
    const router = useRouter();

    const deleteEvent = async () => {
        return await axios.delete(`api/calendar/event/draft/${eventID}`);
    }

    const onSubmit = () => {
        deleteEvent()
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
                        onClick={onSubmit}
                    >
                        削除
                    </Button>
                }
            />
        </div>
    )
}

export default DeleteButton;
