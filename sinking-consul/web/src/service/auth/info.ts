import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取当前网站信息 GET /auth/info  *****/
export async function getWebInfo(params: API.RequestParams = {}) {
    return get("/auth/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}