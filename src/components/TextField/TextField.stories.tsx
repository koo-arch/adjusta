import { Meta, StoryObj } from '@storybook/react';
import TextField from './TextField';

const meta: Meta<typeof TextField> = {
    title: 'Components/TextField',
    component: TextField,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        shape: {
            options: ['rounded', 'md', 'lg', 'full'],
            control: { type: 'radio' },
        },
        inputSize: {
            options: ['sm', 'md', 'lg', 'xl'],
            control: { type: 'radio' },
        },
        error: {
            options: [true, false],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof TextField>;

export const Default: Story = {
    args: {
        label: 'Label',
        description: 'Description',
        placeholder: 'Placeholder',
        value: '',
        onChange: () => {},
    },
};