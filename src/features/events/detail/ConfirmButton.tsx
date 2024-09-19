'use client'
import React, { useState } from 'react';
import axios from '@/lib/axios/public';
import { useForm, Controller, type SubmitHandler } from 'react-hook-form';
import Button from '@/components/Button';
import ToggleButton from '@/components/ToggleButton';
import Modal from '@/components/Modal';
import DropdownSelect from '@/components/DropdownSelect';
import type { ProposedDate } from '@/hooks/event/type';
import { formatJaDate } from '@/lib/date/format';
import DateTimePicker from '@/components/DateTimePicker';

interface ConfrimForm {
    confirm_date: {
        id: string | null;
        event_id: string;
        start_date: Date | null;
        end_date: Date | null;
        priority: number;
    };
}

interface ConfirmButtonProps {
    id: string;
    selectedDates: ProposedDate[];
}

const ConfirmButton: React.FC<ConfirmButtonProps> = ({ id, selectedDates }) => {
    const [isOpen, setIsOpen] = useState(false);
    const [isDropdownSelected, setIsDropdownSelected] = useState(true); // ドロップダウンが選ばれているかどうか
    const method = useForm<ConfrimForm>({
        defaultValues: {
            confirm_date: {
                id: null,
                event_id: "",
                start_date: null,
                end_date: null,
                priority: 0,
            }
        }
    });
    const { control, handleSubmit, reset, formState: { errors } } = method;

    const handleToggle = (selected: string) => {
        setIsDropdownSelected(selected === '候補日程を選択');
        reset();
    }

    const patchConfirmDate = async (data: ConfrimForm) => {
        return await axios.patch(`api/calendar/event/confirm/${id}`, data);
    }

    const onSubmit: SubmitHandler<ConfrimForm> = (data) => {
        console.log(data);
        patchConfirmDate(data)
            .then(res => {
                console.log(res);
                setIsOpen(false);
            })
            .catch(err => {
                console.log(err);
            }
        )
    }

    return (
        <>
            <Button
                shape='full'
                variant='solid'
                intent='primary'
                size='sm'
                onClick={() => setIsOpen(true)}
            >
                日程確定
            </Button>
            <Modal
                isOpen={isOpen}
                onClose={() => setIsOpen(false)}
                title='日程確定'
                description="確定させる日程を選択してください。"
                actions={
                    <Button
                        variant='solid'
                        intent='primary'
                        size='md'
                        type='submit'
                        onClick={() => handleSubmit(onSubmit)()}
                    >
                        確定
                    </Button>
                }
            >
                {/* 切り替え用ボタン */}
                <div className="mb-4">
                    <ToggleButton
                        options={['候補日程を選択', '手動で入力']}
                        selected={isDropdownSelected ? '候補日程を選択' : '手動で入力'}
                        onToggle={handleToggle}
                        renderLabel={(option) => option}
                    />
                </div>
                <form onSubmit={handleSubmit(onSubmit)}>
                    {isDropdownSelected ? 
                    <Controller
                        control={control}
                        name='confirm_date'
                        render={({ field }) => (
                            <DropdownSelect
                                label='日程'
                                options={selectedDates}
                                renderLabel={(date) => 
                                    date && (
                                        <>
                                            {`第${date.priority}候補: 
                                            ${formatJaDate(date.start_date)} ~ ${formatJaDate(date.end_date)}`}
                                        </>
                                    )
                                }
                                onChange={field.onChange}
                                error={!!errors.confirm_date}
                                helperText={errors.confirm_date?.message}
                            />
                        )}
                    /> : 
                    <div>
                        <Controller
                            control={control}
                            name='confirm_date.start_date'
                            render={({ field }) => (
                                <DateTimePicker
                                    label='開始日時'
                                    onChange={field.onChange}
                                    error={!!errors.confirm_date?.start_date}
                                    helperText={errors.confirm_date?.start_date?.message}
                                />
                            )}
                        />
                        <Controller
                            control={control}
                            name='confirm_date.end_date'
                            render={({ field }) => (
                                <DateTimePicker
                                    label='終了日時'
                                    onChange={field.onChange}
                                    error={!!errors.confirm_date?.end_date}
                                    helperText={errors.confirm_date?.end_date?.message}
                                />
                            )}
                        />
                    </div>}
                </form>
            </Modal>
        </>
    )
}

export default ConfirmButton;