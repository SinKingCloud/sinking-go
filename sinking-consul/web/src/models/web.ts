import {useEffect, useState} from "react";
import {API} from "@/../typings";
import {getWebInfo} from "@/service/auth/info";
import defaultSettings from "../../config/defaultSettings";
import {useModel} from "umi";

export default () => {
    const [info, setInfo] = useState<API.WebInfo>({
        ...defaultSettings as API.WebInfo
    });
    const theme = useModel("theme");
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
            if (d?.ui?.color) {
                theme?.setColor(d?.ui?.color);
            }
            if (d?.ui?.radius >= 0) {
                theme?.setRadius(d?.ui?.radius <= 15 ? d?.ui?.radius : 0);
            }
            if (d?.ui?.compact) {
                theme?.setCompactTheme();
            } else {
                theme?.setDefaultTheme();
            }
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