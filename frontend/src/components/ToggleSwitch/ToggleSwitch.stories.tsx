import { Meta, StoryObj } from '@storybook/react';
import ToggleSwitch from './ToggleSwitch';

const meta: Meta<typeof ToggleSwitch> = {
    title: 'Components/ToggleSwitch',
    component: ToggleSwitch,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        size: {
            options: ['sm', 'md', 'lg'],
            control: { type: 'radio' },
        },
        color: {
            options: ['primary', 'secondary', 'danger', 'warning', 'success', 'indigo'],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof ToggleSwitch>;

export const Primary: Story = {
    args: {
        onChange: (checked) => console.log(checked),
        label: 'Toggle Switch',
    },
};