import {API} from "@/../typings";
import {get} from "@/utils/request";
/*****  获取当前网站信息 GET /auth/web/info  *****/
export async function getWebInfo(params: API.RequestParams = {}) {
    return get("/auth/web/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}