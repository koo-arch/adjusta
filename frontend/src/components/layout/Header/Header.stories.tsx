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

export const Default: Story = {
    args: {
        userMenu: (
            <div className="absolute inset-y-0 right-0 flex items-center pr-2 sm:static sm:inset-auto sm:ml-6 sm:pr-0">
                <div className="ml-3 h-8 w-8 rounded-full bg-gray-300" />
            </div>
        ),
    },
};
