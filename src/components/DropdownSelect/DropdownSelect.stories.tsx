import { Meta, StoryObj } from '@storybook/react';
import DropdownSelect from './DropdownSelect';

const meta: Meta<typeof DropdownSelect> = {
    title: 'Components/DropdownSelect',
    component: DropdownSelect,
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
        selectSize: {
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

type Story = StoryObj<typeof DropdownSelect>;

export const Default: Story = {
    args: {
        label: 'Label',
        options: [
            { value: '1', label: 'Option 1' },
            { value: '2', label: 'Option 2' },
            { value: '3', label: 'Option 3' },
        ],
        onChange: (item: any) => { console.log(item) },
        renderLabel: (item: any) => item?.label,
    },
};