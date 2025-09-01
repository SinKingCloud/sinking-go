import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取网站用户信息  *****/
export async function getWebUserInfo(params: API.RequestParams = {}) {
    return get("/user/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}