'use client'
import React from 'react';
import { useRouter } from 'next/navigation';
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';
import { Trash2 } from 'lucide-react';
import { useDeleteDraftMutation } from '@/features/events/detail/hooks/useDeleteDraftMutation';

interface DeleteButtonProps {
    eventID: string;
    title: string;
}

const DeleteButton: React.FC<DeleteButtonProps> = ({ eventID, title }) => {
    const router = useRouter();
    const { submit, isPending } = useDeleteDraftMutation(eventID);

    const handleDelete = async () => {
        const deleted = await submit();
        if (deleted) {
            router.push('/events');
        }
    };

    return (
        <AlertDialog>
            <AlertDialogTrigger asChild>
                <Button
                    variant="ghost"
                    size="icon"
                    aria-label="削除"
                    title="削除"
                    className="text-muted-foreground hover:text-destructive [&_svg]:size-5"
                    disabled={isPending}
                >
                    <Trash2 />
                </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>イベントを削除しますか?</AlertDialogTitle>
                    <AlertDialogDescription>
                        「{title}」を候補日程ごと削除します。この操作は取り消せません。
                    </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>キャンセル</AlertDialogCancel>
                    <AlertDialogAction
                        onClick={handleDelete}
                        className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                    >
                        削除する
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
};

export default DeleteButton;
