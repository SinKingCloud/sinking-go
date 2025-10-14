import {useEffect, useState} from "react";
import {API} from "@/../typings";
import {getWebInfo} from "@/service/auth/info";
import defaultSettings from "../../config/defaultSettings";

export default () => {
    const [info, setInfo] = useState<API.WebInfo>({
        ...defaultSettings as API.WebInfo
    });
    /**
     * 获取站点信息
     */
    const getInfo = async () => {
        const resp = await getWebInfo();
        return resp?.data || {};
    }
    /**
     * 刷新站点信息
     */
    const refreshInfo = () => {
        getInfo().then((d) => {
            setInfo(d);
        });
    }

    useEffect(() => {
        refreshInfo();
    }, []);

    return {
        info,
        setInfo,
        getInfo,
        refreshInfo
    };
};