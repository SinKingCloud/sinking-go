import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取网站用户信息  *****/
export async function getAccountInfo(params: API.RequestParams = {}) {
    return get("/admin/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}