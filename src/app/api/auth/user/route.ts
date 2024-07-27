import axios from "@/lib/axios/private";
import type { AxiosError } from "axios";

export async function GET(request: Request) {
    console.log('request', request.headers);
    const result = await axios.get('/api/users/me')
        .then((res) => {
            console.log('res', res);
            return { data: res.data, status: res.status }
        })
        .catch((err: AxiosError) => {
            console.log('err', err);
            const status: number = err?.response?.status ? err.response.status : 500;
            return { data: err?.response?.data, status: status }
        });
    return new Response(JSON.stringify(result.data), { status: result.status })
}