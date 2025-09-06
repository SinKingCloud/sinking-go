import React, {useState, useEffect} from 'react';
import {Menu, Card} from 'antd';
import PasswordView from './components/password';
import WebView from './components/web';
import UiView from './components/ui';
import {createStyles} from "antd-style";
import {Body, Title} from "@/components";

type SettingsStateKeys = 'password' | 'web' | 'ui';
type SettingsState = {
    mode: 'inline' | 'horizontal';
    selectKey: SettingsStateKeys;
};

const useStyles = createStyles(({css, responsive, isDarkMode, token}): any => {
    const border = isDarkMode ? "1px solid rgb(50, 50, 50)" : "1px solid #f0f0f0"
    const menus = isDarkMode ? "#dad9d9" : "#4b4b4b"
    const primaryColor = token.colorPrimary || '#1890ff'
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

            .ant-menu-light.ant-menu-root.ant-menu-inline {
                border-right: none;
            }
            
            .ant-menu-horizontal .ant-menu-item::after {
                display: none !important;
            }
            
            .ant-menu-horizontal .ant-menu-item-selected::after {
                display: none !important;
            }

            ${responsive.sm} {
                flex-direction: column;
            }
        `,
        leftMenu: css`
            width: 224px;
            min-width: 224px;
            flex-shrink: 0;
            
            .ant-menu-item {
                position: relative;
                font-weight: bold;
                font-size: 14px;
                border-radius: 10px;
                margin: 4px 0;
                height: 40px;
                line-height: 40px;
            }
            
            ${responsive.sm} {
                width: 100%;
                min-width: 100%;
                border: none;
                margin-bottom: 16px;
                
                .ant-menu {
                    border-bottom: ${border};
                }
                
                .ant-menu-item {
                    display: inline-block;
                    width: auto;
                    min-width: 80px;
                    text-align: center;

                    margin: 0 4px 0 0;
                    border-radius: 6px;
                    height: 36px;
                    line-height: 36px;
                }
            }
            
            ${responsive.xs} {
                .ant-menu-item {
                    min-width: 60px;
                    font-size: 12px;
                    padding: 0 6px;
                }
            }
        `,
        right: css`
            flex: 1;
            min-width: 0;
            padding: 8px 40px;
            margin-left: 1px;
            border-left: ${border};
            
            ${responsive.sm} {
                padding: 16px 20px;
                border-left: none;
                border-top: ${border};
            }
            
            ${responsive.xs} {
                padding: 8px 12px;
            }
        `,
        title: css`
            margin-bottom: 20px;
            font-weight: bolder;
            font-size: 25px;
            line-height: 28px;
            
            ${responsive.md} {
                font-size: 22px;
                margin-bottom: 16px;
            }
            
            ${responsive.sm} {
                font-size: 20px;
                margin-bottom: 12px;
            }
            
            ${responsive.xs} {
                font-size: 18px;
                margin-bottom: 10px;
            }
        `,
        card: css`
            .ant-card-body {
                padding: 24px;
            }
            
            ${responsive.md} {
                .ant-card-body {
                    padding: 16px;
                }
            }
            
            ${responsive.sm} {
                .ant-card-body {
                    padding: 12px;
                }
            }
            
            ${responsive.xs} {
                .ant-card-body {
                    padding: 8px;
                }
            }
        `,
        menu: css`
            .ant-menu-title-content {
                color: ${menus};
            }
            
            ${responsive.sm} {
                &.ant-menu-horizontal {
                    border-bottom: none;
                    
                    .ant-menu-item {
                        border-bottom: 2px solid transparent;
                        position: relative;
                        
                        &.ant-menu-item-selected {
                            border-bottom-color: ${primaryColor};
                            background: ${primaryColor}1a;
                            
                            &::after {
                                display: none !important;
                            }
                        }
                        
                        &::after {
                            display: none !important;
                        }
                    }
                }
            }
        `
    }
})

const Settings: React.FC = () => {
    const menuMap: Record<string, React.ReactNode> = {
        password: '密码设置',
        web: '网站配置',
        ui: '界面配置',
    };
    const {styles: {main, leftMenu, right, title, card, menu}} = useStyles()
    const [initConfig, setInitConfig] = useState<SettingsState>({
        mode: 'inline',
        selectKey: 'password',
    });
    const [isMobile, setIsMobile] = useState(false);

    // 检测屏幕尺寸变化
    useEffect(() => {
        const checkIsMobile = () => {
            const width = window.innerWidth;
            setIsMobile(width < 800); // 只在真正的移动设备(小于576px)才切换
            setInitConfig(prev => ({
                ...prev,
                mode: width < 800 ? 'horizontal' : 'inline'
            }));
        };

        // 初始检查
        checkIsMobile();

        // 监听窗口大小变化
        window.addEventListener('resize', checkIsMobile);
        return () => window.removeEventListener('resize', checkIsMobile);
    }, []);

    const getMenu = () => {
        return Object.keys(menuMap).map((item) => {
            return {
                key: item,
                label: menuMap[item]
            }
        });
    };

    const renderChildren = () => {
        const {selectKey} = initConfig;
        switch (selectKey) {
            case 'password':
                return <PasswordView/>;
            case 'web':
                return <WebView/>;
            case 'ui':
                return <UiView/>;
            default:
                return null;
        }
    };

    return (
        <Body>
            <Card title={<Title>系统设置</Title>} className={card}>
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
