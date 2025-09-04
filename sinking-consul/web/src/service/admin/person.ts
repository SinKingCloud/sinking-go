import {API} from "@/../typings";
import {get, post} from "@/utils/request";

/*****  获取网站用户信息  *****/
export async function getAccountInfo(params: API.RequestParams = {}) {
    return get("/admin/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  更新密码信息  *****/
export async function updatePassword(params: API.RequestParams = {}) {
    return post("/admin/person/password", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}