import { Meta, StoryObj } from "@storybook/react";
import ThemeButton from './ThemeButton';

const meta: Meta<typeof ThemeButton> = {
    title: "Components/ThemeButton",
    component: ThemeButton,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ["autodocs"],
}

export default meta;

type Story = StoryObj<typeof ThemeButton>;

export const Default: Story = {};