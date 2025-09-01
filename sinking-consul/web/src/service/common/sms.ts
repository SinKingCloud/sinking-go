
import {post} from "@/utils/request";
import {API} from "../../../typings";

/** 发送短信验证码 POST /auth/verify/sms */
export async function sendSms(params: API.RequestParams = {}) {
    return post("/auth/verify/sms", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}