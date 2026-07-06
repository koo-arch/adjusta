import { Meta, StoryObj } from "@storybook/nextjs";
import Header from "./Header";

const meta: Meta<typeof Header> = {
    title: "Components/Header",
    component: Header,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ["autodocs"],
}

export default meta;

type Story = StoryObj<typeof Header>;

export const Default: Story = {};