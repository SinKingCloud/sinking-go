import {API} from "@/../typings";
import {get, post} from "@/utils/request";

/*****  获取node信息  *****/
export async function getNodeList(params: API.RequestParams = {}) {
    return get("/admin/node/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  更新node信息  *****/
export async function updateNode(params: API.RequestParams = {}) {
    return post("/admin/node/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/*****  删除node信息  *****/
export async function deleteNode(params: API.RequestParams = {}) {
    return post("/admin/node/delete", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}