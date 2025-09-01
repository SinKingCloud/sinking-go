import {useLocation, useModel, useSelectedRoutes} from "@@/exports";
import defaultSettings from "../../../../config/defaultSettings";
import {useEffect} from "react";

const Title = () => {
    const web = useModel("web");//网站信息
    const match = useSelectedRoutes();
    const location = useLocation();

    const initWeb = () => {
        const route = match?.pop()?.route;
        if (route && route?.title) {
            document.title = route?.title + " - " + (web?.info?.name || defaultSettings?.title);
        }
    }

    useEffect(() => {
        initWeb();
    }, [web?.info?.name, location]);
};

export default Title;