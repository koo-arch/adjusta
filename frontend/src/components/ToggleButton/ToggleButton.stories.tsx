import { Meta, StoryObj } from '@storybook/react';
import ToggleButton from './ToggleButton';

const meta: Meta<typeof ToggleButton> = {
    title: 'Components/ToggleButton',
    component: ToggleButton,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ['autodocs'],
    argTypes: {
        size: {
            options: ['sm', 'md', 'lg', 'xl'],
            control: { type: 'radio' },
        },
    },
};

export default meta;

type Story = StoryObj<typeof ToggleButton>;

export const Primary: Story = {
    args: {
        options: ['Option 1', 'Option 2'],
        selected: 'Option 1',
        onToggle: (selected) => console.log(selected),
        renderLabel: (option: any) => option,
    },
};