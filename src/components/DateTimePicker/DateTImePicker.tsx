'use client'
import React, { useState } from 'react';
import DatePicker, { DatePickerProps } from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';
import { ja } from 'date-fns/locale';

type DateTimePickerProps = DatePickerProps & {
    label?: string;
    error?: boolean;
    timeIntervals?: number;
    helperText?: string;
    initialDate?: Date;
} & {
    selectsRange?: never;
    selectsMultiple?: never;
};

const DateTimePicker: React.FC<DateTimePickerProps> = ({ timeIntervals, label, error, helperText, initialDate, onChange, ...props }) => {
    const [selectedDate, setSelectedDate] = useState<Date | null>(initialDate || null);

    const handleDateChange = (date: Date | null, event?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>) => {
        setSelectedDate(date);
        onChange && onChange(date, event);
    }

    return (
        <div>
            {label && <label className="block text-base font-medium text-gray-700 mb-2">{label}</label>}
            <DatePicker
                locale={ja}
                showTimeSelect
                selected={selectedDate}
                onChange={handleDateChange}
                timeFormat="HH:mm"
                timeIntervals={timeIntervals || 15}
                timeCaption="time"
                dateFormat="Pp"
                {...props}
                wrapperClassName='w-full'
                className={`block w-full mt-1 border px-3 py-1.5 text-base rounded-md focus:outline-none focus:ring-2 ${error ? 'border-red-500 focus:ring-red-500' : 'focus:ring-indigo-500'}`}
            />
            {helperText && (
                <p className={`mt-1 text-sm ${error ? 'text-red-500' : 'text-gray-500'}`}>
                    {helperText}
                </p>
            )}
        </div>
    )
}

export default DateTimePicker;