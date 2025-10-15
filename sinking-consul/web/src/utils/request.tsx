import {extend} from 'umi-request';
import defaultSettings from "@/../config/defaultSettings"
import {API} from "../../typings";
import {Modal} from "sinking-antd";

/**
 * request对象
 */
const request = extend({
    prefix: defaultSettings?.gateway || '',
    timeout: 60000,
    errorHandler: errorHandle,
});

/**
 * 异常处理
 */
let lock = false;

function errorHandle(): void {
    if (!lock) {
        lock = true
        const clear = () => {
            lock = false;
            location?.replace(location?.href || "");
        }
        Modal?.error({
            title: "系统错误",
            content: "服务器连接失败,请检查网络配置并稍后刷新重试",
            okText: "刷新",
            onOk: clear,
            onCancel: clear,
            styles: {
                mask: {
                    WebkitBackdropFilter: "blur(10px)",
                },
            },
        } as any);
    }
}


/**
 * 发送post请求
 * @param url 请求地址
 * @param body 请求载荷
 * @param onSuccess code200回调
 * @param onFail code500回调
 * @param onFinally 最后执行
 * @param onError 错误捕获
 * @param options 其他选项
 */
export async function post(url: string = "", body: any = {}, onSuccess: (res) => void = undefined, onFail: (res) => void = undefined, onFinally: () => void = undefined, onError: (error) => void = undefined, options?: {
    [key: string]: any
}) {
    let res = {};
    await request<API.Response>(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        data: body,
        ...(options || {}),
    }).then((r) => {
        res = r;
        if (r?.code == 200) {
            onSuccess?.(r);
        } else {
            onFail?.(r);
        }
    }).catch((error) => {
        onError?.(error);
    }).finally(() => {
        onFinally?.();
    });
    return res as API.Response;
}

/**
 * 发送post请求
 * @param url 请求地址
 * @param params 请求参数
 * @param onSuccess code200回调
 * @param onFail code500回调
 * @param onFinally 最后执行
 * @param onError 错误捕获
 * @param options 其他选项
 */
export async function get(url: string = "", params?: {}, onSuccess: (res) => void = undefined, onFail: (res) => void = undefined, onFinally: () => void = undefined, onError: (error) => void = undefined, options?: {
    [key: string]: any
}) {
    let res = {};
    await request<API.Response>(url, {
        method: 'GET',
        params: {
            ...params,
        },
        ...(options || {}),
    }).then((r) => {
        res = r;
        if (r?.code == 200) {
            onSuccess?.(r);
        } else {
            onFail?.(r);
        }
    }).catch((error) => {
        onError?.(error);
    }).finally(() => {
        onFinally?.();
    });
    return res as API.Response;
}

export default request;
