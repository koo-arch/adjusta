import axios, { AxiosInstance } from "axios";
import { getDefaultStore } from "jotai/vanilla";
import { authErrorAtom } from "@/atoms/error";

export const axiosConfig = {
    baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
    withCredentials: true,
    headers: {
        "Content-Type": "application/json",
    },
};

const instance: AxiosInstance = axios.create(axiosConfig);

// Jotaiのデフォルトストアを取得
const store = getDefaultStore();

// インターセプターを設定
instance.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            // 認証エラーを検知してJotaiの状態を更新
            store.set(authErrorAtom, {
                isOpen: true,
                message: "認証エラーが発生しました。再ログインしてください。",
            });
        }
        return Promise.reject(error);
    }
);

export default instance;