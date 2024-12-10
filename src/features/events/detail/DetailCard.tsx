'use client'
import React from 'react';
import Card from '@/components/Card';
import { formatJaDateSpan } from '@/lib/date/format';
import type { EventDraftDetail } from '@/hooks/event/type';
import { MapPinIcon } from '@heroicons/react/20/solid';
import EditButton from './EditButton';
import ConfirmButton from './ConfirmButton';
import DeleteButton from './DeleteButton';

type DetailCardProps = {
    id: string;
    detail: EventDraftDetail;
}

const DetailCard: React.FC<DetailCardProps> = ({ detail, id }) => {
    const confirmedDate = detail.proposed_dates?.find((date) => date.id === detail.confirmed_date_id);
    const isConfirmed = detail.status === 'confirmed' && !!detail.confirmed_date_id && !!confirmedDate;
    return (
        <Card variant="outlined" className="w-full shadow-md">
            <div className="space-y-6">
                <div className="flex justify-between items-center mb-4 border-b pb-2">
                    <h1 className="text-2xl font-bold text-gray-900">{detail.title}</h1>
                    <EditButton to={`/schedule/draft/${id}/edit`} />
                </div>

                {isConfirmed && (
                    <div className='mb-2'>
                        <div className="flex items-center space-x-2 mb-2">
                            <h2 className="text-lg font-semibold text-gray-700">確定日時</h2>
                            <ConfirmButton
                                id={id}
                                detail={detail}
                                isConfirmed={isConfirmed}
                            />
                        </div>
                        <p className="text-lg text-indigo-500">
                            {formatJaDateSpan(confirmedDate.start, confirmedDate.end)}
                        </p>
                    </div>
                )}

                {!isConfirmed && (
                    <div className="space-y-2">
                        <div className="flex items-center space-x-2">
                            <h2 className="text-lg font-semibold text-gray-700">候補日程</h2>
                            <ConfirmButton 
                                id={id}
                                detail={detail}
                                isConfirmed={isConfirmed}
                            />
                        </div>
                        {detail.proposed_dates?.map((date) => (
                            <p key={date.id} className="text-sm text-gray-500">
                                <span className="font-medium">第{date.priority}候補：</span>
                                {formatJaDateSpan(date.start, date.end)}
                            </p>
                        ))}
                    </div>
                )}

                <div className="flex items-center mt-4">
                    <MapPinIcon className="w-6 h-6 text-gray-500 mr-2" />
                    <p className="text-sm text-gray-700">{detail.location || '未設定'}</p>
                </div>

                <div className="mt-4">
                    <h2 className="text-lg font-semibold text-gray-700">説明</h2>
                    <p className="text-sm text-gray-700 mt-2">{detail.description}</p>
                </div>
                <div className="flex justify-end items-end">
                    <DeleteButton id={id} detail={detail} />
                </div>
            </div>
        </Card>
    );
};

export default DetailCard;