import { API } from "../../../typings";
import { get, post } from "@/utils/request";

/** 获取个人信息 GET /user/person/info */
export async function getPersonInfo(params: API.RequestParams = {}) {
    return get("/user/person/info", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 更新个人信息 POST /user/person/update */
export async function updatePersonInfo(params: API.RequestParams = {}) {
    return post("/user/person/update", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改密码 POST /user/person/password */
export async function updatePassword(params: API.RequestParams = {}) {
    return post("/user/person/password", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
}

/** 修改手机号 POST /user/person/phone */
export async function updatePhone(params: API.RequestParams = {}) {
    return post("/user/person/phone", params?.body, params?.onSuccess, params?.onFail, params?.onFinally);
} 