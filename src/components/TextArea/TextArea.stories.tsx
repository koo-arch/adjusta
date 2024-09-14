import { Meta, StoryObj } from '@storybook/react';
import TextArea from './TextArea';

const meta: Meta<typeof TextArea> = {
    title: 'Components/TextArea',
    component: TextArea,
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
        areaSize: {
            options: ['sm', 'md', 'lg'],
            control: { type: 'radio' },
        },
        error: {
            options: [true, false],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof TextArea>;

export const Default: Story = {
    args: {
        label: 'Label',
        description: 'Description',
        placeholder: 'Placeholder',
        value: '',
        onChange: () => {},
        helperText: 'Helper text',
    },
};