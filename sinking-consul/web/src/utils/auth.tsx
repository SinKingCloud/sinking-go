/**
 * 设置登录token
 * @param type
 * @param token
 */
export function setLoginToken(type: string, token: string) {
    localStorage.setItem('token', token);
    localStorage.setItem('device', type);
}

/**
 * 获取登录token
 * @returns {string}
 */
export function getLoginToken() {
    return localStorage.getItem('token');
}

/**
 * 获取登录方式
 * @returns {string}
 */
export function getLoginType() {
    return localStorage.getItem('device');
}

/**
 * 删除token
 */
export function deleteHeader() {
    localStorage.removeItem('token');
    localStorage.removeItem('device');
}

/**
 * 返回请求header
 * @returns {{}}
 */
export function getHeaders() {
    return {
        'token': getLoginToken(),
        'device': getLoginType()
    };
}
