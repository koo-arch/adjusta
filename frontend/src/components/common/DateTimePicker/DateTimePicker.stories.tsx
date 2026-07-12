import { useState } from 'react';
import { Meta, StoryObj } from '@storybook/nextjs';
import { DateTimePicker } from './DateTimePicker';

const meta: Meta<typeof DateTimePicker> = {
    title: 'Common/DateTimePicker',
    component: DateTimePicker,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
};

export default meta;

type Story = StoryObj<typeof DateTimePicker>;

const InteractiveTemplate = ({ error, helperText }: { error?: boolean; helperText?: string }) => {
    const [value, setValue] = useState<Date | null>(null);
    return (
        <div className="max-w-xs">
            <DateTimePicker
                label="開始日時"
                selected={value}
                onChange={setValue}
                error={error}
                helperText={helperText}
            />
        </div>
    );
};

export const Default: Story = {
    render: () => <InteractiveTemplate />,
};

export const WithError: Story = {
    render: () => <InteractiveTemplate error helperText="開始日時は必須です" />,
};
