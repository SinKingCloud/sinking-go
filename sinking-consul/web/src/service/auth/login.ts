// @ts-ignore
import {request} from 'umi';
import {get, post} from "@/utils/request";
import {API} from "../../../typings";

/** 账号登录 POST /auth/login */
export async function login(params: API.RequestParams = {}) {
    return post("/auth/login", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 退出登录 GET /auth/out */
export async function logout(params: API.RequestParams = {}) {
    return get("/auth/logout", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}