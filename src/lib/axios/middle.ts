import axios, { AxiosInstance } from "axios";

export const axiosConfig = {
    withCredentials: true,
    headers: {
        "Content-Type": "application/json",
    },
};

const instance: AxiosInstance = axios.create(axiosConfig);

export default instance;