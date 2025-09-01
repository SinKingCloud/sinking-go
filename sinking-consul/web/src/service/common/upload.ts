import Settings from "../../../config/defaultSettings";
import request from "@/utils/request";

export function getUploadUrl(){
    return Settings.gateway + "/auth/upload/file"
}

export async function uploadFile(body: any, options?: { [key: string]: any }) {
    return request('/auth/upload/file', {
        method: 'POST',
        data: body,
        ...(options || {}),
    });
}
