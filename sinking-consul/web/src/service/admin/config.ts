import {API} from "@/../typings";
import {get, post} from "@/utils/request";

/*****  获取config列表  *****/
export async function getConfigList(params: API.RequestParams = {}) {
    return get("/admin/config/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  获取config信息  *****/
export async function getConfigInfo(params: API.RequestParams = {}) {
    return get("/admin/config/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  更新node信息  *****/
export async function updateConfig(params: API.RequestParams = {}) {
    return post("/admin/config/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  创建node信息  *****/
export async function createConfig(params: API.RequestParams = {}) {
    return post("/admin/config/create", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  删除node信息  *****/
export async function deleteConfig(params: API.RequestParams = {}) {
    return post("/admin/config/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}