import {API} from "@/../typings";
import {get} from "@/utils/request";

/*****  获取cluster信息  *****/
export async function getClusterList(params: API.RequestParams = {}) {
    return get("/admin/cluster/list", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}