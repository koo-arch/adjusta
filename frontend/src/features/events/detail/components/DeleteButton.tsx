'use client'
import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import Modal from '@/components/Modal';
import Button from '@/components/Button';
import IconButton from '@/components/IconButton';
import { TrashIcon } from '@heroicons/react/20/solid';
import { useDeleteDraftMutation } from '@/features/events/detail/hooks/useDeleteDraftMutation';

interface EventDeleteProps {
    eventID: string;
}

const DeleteButton: React.FC<EventDeleteProps> = ({ eventID }) => {
    const [isOpen, setIsOpen] = useState(false);
    const router = useRouter();
    const deleteDraftMutation = useDeleteDraftMutation(eventID);

    const onSubmit = async () => {
        const deleted = await deleteDraftMutation.submit();
        if (!deleted) {
            return;
        }

        setIsOpen(false);
        router.push('/schedule/draft');
    };

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
                        disabled={deleteDraftMutation.isPending}
                    >
                        削除
                    </Button>
                }
            />
        </div>
    )
}

export default DeleteButton;
