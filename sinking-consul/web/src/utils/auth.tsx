/**
 * 设置登录token
 * @param token
 */
export function setLoginToken(token: string) {
    localStorage.setItem('token', token);
}

/**
 * 获取登录token
 * @returns {string}
 */
export function getLoginToken() {
    return localStorage.getItem('token');
}

/**
 * 删除token
 */
export function deleteHeader() {
    localStorage.removeItem('token');
}

/**
 * 返回请求header
 * @returns {{}}
 */
export function getHeaders() {
    return {
        'token': getLoginToken(),
    };
}
