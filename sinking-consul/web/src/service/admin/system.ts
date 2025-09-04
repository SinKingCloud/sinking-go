import {API} from "@/../typings";
import {get, post} from "@/utils/request";

/*****  获取网站统计信息  *****/
export async function getOverviewInfo(params: API.RequestParams = {}) {
    return get("/admin/system/overview", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  获取系统配置信息  *****/
export async function getConfig(params: API.RequestParams = {}) {
    return get("/admin/system/config", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  设置系统配置信息  *****/
export async function setConfig(params: API.RequestParams = {}) {
    return post("/admin/system/config", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}