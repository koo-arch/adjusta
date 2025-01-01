import { Meta, StoryObj } from "@storybook/react";
import IconButton from "./IconButton";
import { PlusIcon } from "@heroicons/react/24/solid";

const meta: Meta<typeof IconButton> = {
    title: "Components/IconButton",
    component: IconButton,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ["autodocs"],
    argTypes: {
        iconSize: {
            options: ["sm", "md", "lg"],
            control: { type: "radio" },
        },
        strokeWidth: {
            options: [1, 2, 3],
            control: { type: "radio" },
        },
    },
}

export default meta;

type Story = StoryObj<typeof IconButton>;

export const Primary: Story = {
    args: {
        children: <PlusIcon />,
    },
};