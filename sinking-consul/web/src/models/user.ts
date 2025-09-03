import {useState} from "react";
import {API} from "@/../typings";
import { getAccountInfo} from "@/service/admin/person";

export default () => {
    const [web, setWeb] = useState<API.UserInfo>();//网站用户
    /**
     * 获取网站用户信息
     */
    const getWebUser = async () => {
        const resp = await getAccountInfo();
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