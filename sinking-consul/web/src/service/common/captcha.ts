import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取验证码信息 GET /auth/captcha  *****/
export async function getCaptcha(params: API.RequestParams = {}) {
    return get("/auth/captcha", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}