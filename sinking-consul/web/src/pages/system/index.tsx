import React, {useState} from 'react';
import {Card, Menu} from 'antd';
import UiView from "./components/ui";
import {Body, Title} from "@/components";
import {createStyles} from "antd-style";
import WebView from "@/pages/system/components/web";
import PasswordView from "@/pages/system/components/password";

type SettingsStateKeys = 'web' | 'pwd' | 'ui';
type SettingsState = {
    mode: 'inline' | 'horizontal';
    selectKey: SettingsStateKeys;
};
const useStyles = createStyles(({css, responsive, isDarkMode}): any => {
    const border = isDarkMode ? "1px solid rgb(50, 50, 50)" : "1px solid #f0f0f0"
    const menus = isDarkMode ? "#dad9d9" : "#4b4b4b"
    return {
        main: css`
            display: flex;
            width: 100%;
            height: 100%;
            padding-bottom: 16px;

            .ant-menu-light.ant-menu-inline .ant-menu-item::after {
                top: 19%;
                right: 6px;
                border-top-left-radius: 15px;
                border-bottom-left-radius: 15px;
                height: 60%;
            }


            ${responsive.md} {
                flex-direction: column;
            }
        `,
        leftMenu: css`
            width: 224px;

            .ant-menu-item {
                position: absolute;
                font-weight: bold;
                font-size: 14px;
                border-radius: 10px;
            }

            ${responsive.md} {
                width: 100%;
                border: none;
                margin-bottom: 10px;
            }
        `,
        right: css`
            flex: 1;
            padding: 8px 40px;
            margin-left: 1px;
            border-left: ${border};

            ${responsive.md} {
                padding: 10px;
                border-left: none;
            }
        `,
        title: css`
            margin-bottom: 20px;
            font-weight: bolder;
            font-size: 25px;
            line-height: 28px;
        `,
        menu: css`
            .ant-menu-title-content {
                color: ${menus}
            }
        `
    }
})
const Settings: React.FC = () => {
    const menuMap: Record<string, React.ReactNode> = {
        web: "网站设置",
        ui: '界面设置',
        pwd: "密码设置"
    };
    const {styles: {main, leftMenu, right, title, menu}} = useStyles()
    const [initConfig, setInitConfig] = useState<SettingsState>({
        mode: 'inline',
        selectKey: 'web',
    });

    const getMenu = () => {
        return Object.keys(menuMap).map((item) => {
            return {
                key: item, label: menuMap[item]
            }
        });
    };
    const renderChildren = () => {
        const {selectKey} = initConfig;
        switch (selectKey) {
            case 'web':
                return <WebView/>;
            case 'ui':
                return <UiView/>;
            case 'pwd':
                return <PasswordView/>;
            default:
                return null;
        }
    };

    return (
        <Body>
            <Card title={<Title>网站设置</Title>}>
                <div className={main}>
                    <div className={leftMenu}>
                        <Menu
                            className={menu}
                            mode={initConfig.mode}
                            selectedKeys={[initConfig.selectKey]}
                            onClick={({key}) => {
                                setInitConfig({
                                    ...initConfig,
                                    selectKey: key as SettingsStateKeys,
                                });
                            }}
                            items={getMenu()}
                        />
                    </div>
                    <div className={right}>
                        <div className={title}>{menuMap[initConfig.selectKey]}</div>
                        {renderChildren()}
                    </div>
                </div>
            </Card>
        </Body>
    );
};
export default Settings;
