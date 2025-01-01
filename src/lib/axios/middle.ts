import axios, { AxiosInstance } from "axios";

export const axiosConfig = {
    baseURL: "http://localhost:3000",
    withCredentials: true,
    headers: {
        "Content-Type": "application/json",
    },
};

const instance: AxiosInstance = axios.create(axiosConfig);

export default instance;