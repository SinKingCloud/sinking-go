import {useState} from "react";
import {API} from "@/../typings";
import { getWebUserInfo} from "@/service/person/info";

export default () => {
    const [web, setWeb] = useState<API.UserInfo>();//网站用户
    /**
     * 获取网站用户信息
     */
    const getWebUser = async () => {
        const resp = await getWebUserInfo();
        if (resp?.code != 200) {
            return undefined;
        }
        return resp?.data || undefined;
    }

    /**
     * 刷新网站用户信息
     */
    const refreshWebUser = (callback: () => void = undefined) => {
        getWebUser().then((d) => {
            if (d) {
                setWeb(d);
                callback?.();
            }
        });
    }

    return {
        web,
        setWeb,
        getWebUser,
        refreshWebUser,
    };
};