/**
 * 设置支付跳转URL
 * @param url
 */
export function setPayJumpUrl(url: string = "") {
    if (url == "") {
        url = window.location.href;
    }
    localStorage.setItem('payJumpUrl', url);
}

/**
 * 获取支付跳转URL
 */
export function getPayJumpUrl(): string {
    return localStorage.getItem('payJumpUrl') || '';
}
