// @ts-ignore
import {request} from 'umi';
import {get, post} from "@/utils/request";
import {API} from "../../../typings";

/** 账号登录 POST /auth/web/login */
export async function login(params: API.RequestParams = {}) {
    return post("/auth/web/login", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 退出登录 GET /auth/login/out */
export async function outLogin(params: API.RequestParams = {}) {
    return get("/auth/web/out_login", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}