export function ago(string: string) {
    let d;
    const f = string.split(' ', 2);//过滤空格
    if (f[0].search("/") != -1) {//判断是否包含-
        d = (f[0] ? f[0] : '').split('/', 3);//过滤-
    } else {
        d = (f[0] ? f[0] : '').split('-', 3);//过滤-
    }
    const t = (f[1] ? f[1] : '').split(':', 3);//过滤:
    const dateTimeStamp = (new Date(
        parseInt(d[0], 10) || 0,
        (parseInt(d[1], 10) || 1) - 1,
        parseInt(d[2], 10) || 0,
        parseInt(t[0], 10) || 0,
        parseInt(t[1], 10) || 0,
        parseInt(t[2], 10) || 0,
    )).getTime();
    const minute = 1000 * 60;
    const hour = minute * 60;
    const day = hour * 24;
    const week = day * 7;
    const month = day * 30;
    const now = new Date().getTime();
    const diffValue = now - dateTimeStamp;
    if (diffValue < 0) {
        return;
    }
    const minC = diffValue / minute;
    const hourC = diffValue / hour;
    const dayC = diffValue / day;
    const weekC = diffValue / week;
    const monthC = diffValue / month;
    let result;
    if (monthC >= 1 && monthC <= 6) {
        result = " " + parseInt(String(monthC)) + "月前"
    } else if (weekC >= 1 && weekC <= 3) {
        result = " " + parseInt(String(weekC)) + "周前"
    } else if (dayC >= 1 && dayC <= 30) {
        result = " " + parseInt(String(dayC)) + "天前"
    } else if (hourC >= 1 && hourC <= 23) {
        result = " " + parseInt(String(hourC)) + "小时前"
    } else if (minC >= 1 && minC <= 59) {
        result = " " + parseInt(String(minC)) + "分钟前"
    } else if (diffValue >= 0 && diffValue <= minute) {
        result = "刚刚"
    } else {
        const datetime = new Date();
        datetime.setTime(dateTimeStamp);
        const Nyear = datetime.getFullYear();
        const Nmonth = datetime.getMonth() + 1 < 10 ? "0" + (datetime.getMonth() + 1) : datetime.getMonth() + 1;
        const Ndate = datetime.getDate() < 10 ? "0" + datetime.getDate() : datetime.getDate();
        const Nhour = datetime.getHours() < 10 ? "0" + datetime.getHours() : datetime.getHours();
        const Nminute = datetime.getMinutes() < 10 ? "0" + datetime.getMinutes() : datetime.getMinutes();
        const Nsecond = datetime.getSeconds() < 10 ? "0" + datetime.getSeconds() : datetime.getSeconds();
        result = Nyear + "-" + Nmonth + "-" + Ndate + ' ' + Nhour + ":" + Nminute + ":" + Nsecond;
    }
    return result;
}
