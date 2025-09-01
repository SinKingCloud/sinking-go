import {API} from "../../../typings";
import {post} from "@/utils/request";

export async function getLogList(params: API.RequestParams = {}) {
    return post("/user/person/log", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}