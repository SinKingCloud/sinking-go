import React, {useMemo} from "react";
import {Breadcrumb} from "antd";
import {useLocation, useSelectedRoutes} from "umi";
import {getAllMenuItems, getFirstMenuWithoutChildren, getParentList, historyPush} from "@/utils/route";
import {createStyles} from "antd-style";

const useBreadCrumbStyles = createStyles(({token}) => {
    return {
        bread: {
            backgroundColor: token?.colorBgContainer,
            padding: "5px 15px 5px 15px",
            fontSize: "12px",
        },
        breadStyle: {
            color: "rgb(156, 156, 156)",
            cursor: "pointer"
        },
    };
});

export type BreadCrumbProps = {
    style?: React.CSSProperties; // 样式
    className?: string; // 样式名
    enabled?: boolean; // 是否启用面包屑
    homeTitle?: string; // 首页标题
};

/**
 * 面包屑组件
 * @param props
 * @constructor
 */
const BreadCrumb: React.FC<BreadCrumbProps> = React.memo((props) => {
    const {
        style,
        className,
        enabled = true,
        homeTitle = "首页"
    } = props;

    const {styles: {bread, breadStyle}} = useBreadCrumbStyles();
    const location = useLocation();
    const match = useSelectedRoutes();

    const breadCrumbData = useMemo(() => {
        if (!enabled) return [];

        const items = getParentList(getAllMenuItems(false), match?.at(-1)?.route?.name);
        const temp = [{
            title: homeTitle,
            onClick: () => {
                historyPush(getFirstMenuWithoutChildren(getAllMenuItems(location?.pathname))?.name || "");
            },
            className: breadStyle,
        }];

        const handleItemClick = (x: any) => {
            if (x?.children && x?.children?.length > 0) {
                historyPush(getFirstMenuWithoutChildren(x?.children)?.name || "");
            } else {
                historyPush(x?.name);
            }
        };

        items.forEach((x) => {
            temp.push({
                title: x?.label,
                onClick: () => handleItemClick(x),
                className: breadStyle,
            });
        });

        return temp;
    }, [enabled, location?.pathname, homeTitle, breadStyle]);

    if (!enabled || match?.at(-1)?.route?.hideBreadCrumb || breadCrumbData?.length === 0) {
        return null;
    }

    return (
        <Breadcrumb
            className={`${bread} ${className || ''}`}
            style={style}
            items={breadCrumbData}
        />
    );
});

export default BreadCrumb;