import React from "react";
import { Meta, StoryObj } from "@storybook/react";
import { PlusIcon } from "@heroicons/react/20/solid";
import Button from "./Button";

const meta: Meta<typeof Button> = {
    title: "Components/Button",
    component: Button,
    parameters: {
        nextjs: {
            appDirectory: true,
        },
    },
    tags: ["autodocs"],
    argTypes: {
        intent: {
            options: ["primary", "secondary", "danger", "warning", "success", "clear"],
            control: { type: "select" },
        },
        variant: {
            options: ["solid", "outline"],
            control: { type: "radio" },
        },
        size: {
            options: ["sm", "md", "lg"],
            control: { type: "radio" },
        },
        shape: {
            options: ["rounded", "md", "lg", "full"],
            control: { type: "radio" },
        },
        iconSize: {
            options: ["sm", "md", "lg"],
            control: { type: "radio" },
        },
    },
}

export default meta;

type Story = StoryObj<typeof Button>;

export const Primary: Story = {
    args: {
        children: "Button",
    },
};

export const WithStartIcon: Story = {
    args: {
        startIcon: <PlusIcon />,
        children: "Button",
    },
};

export const WithEndIcon: Story = {
    args: {
        endIcon: <PlusIcon />,
        children: "Button",
    },
};