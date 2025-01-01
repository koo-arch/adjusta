import { Meta, StoryObj } from '@storybook/react';
import DateTimePicker from './DateTImePicker';

const meta: Meta<typeof DateTimePicker> = {
    title: 'Components/DateTimePicker',
    component: DateTimePicker,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        timeIntervals: {
            options: [15, 30, 60],
            control: { type: 'radio' },
        },
        label: {
            control: { type: 'text' },
        },
        error: {
            options: [true, false],
            control: { type: 'radio' },
        },
        helperText: {
            control: { type: 'text' },
        }
    },
};

export default meta;

type Story = StoryObj<typeof DateTimePicker>;

export const Default: Story = {
    args: {
        label: 'Label',
        error: false,
        helperText: 'Helper text',
    },
};