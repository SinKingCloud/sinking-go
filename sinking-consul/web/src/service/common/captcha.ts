import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取验证码信息 GET /auth/verify/captcha  *****/
export async function getCaptcha(params: API.RequestParams = {}) {
    return get("/auth/verify/captcha", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}