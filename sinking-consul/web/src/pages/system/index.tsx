import React from "react";
import {Col, Menu, Row} from "antd";
import {Body, Icon, useTheme} from "sinking-antd";
import {history, useLocation, useModel} from "umi";
import Ui from "./components/ui";
import Web from "./components/web";
import Password from "./components/password";
import HeroGraphic from "./components/hero-graphic";
import useStyles from "./styles";

const items = [
    {key: "web", title: "网站设置", icon: "GlobalOutlined", component: Web},
    {key: "ui", title: "界面设置", icon: "LayoutOutlined", component: Ui},
    {key: "pwd", title: "密码设置", icon: "SafetyCertificateOutlined", component: Password},
];

export default (): React.ReactNode => {
    const theme = useTheme();
    const web = useModel("web");
    const location = useLocation();
    const isCompactMode = theme?.isCompactTheme?.() || false;
    const isDarkMode = Boolean(theme?.isDarkMode?.() || theme?.isDarkTheme?.());
    const {styles} = useStyles({isCompactMode, isDarkMode});
    const selected = new URLSearchParams(location.search).get("tab") || "web";
    const activeItem = items.find((item) => item.key === selected) || items[0];
    const Component = activeItem.component;

    const changeActive = (key: string) => {
        const search = new URLSearchParams(location.search);
        search.set("tab", key);
        history.replace(`${location.pathname}?${search.toString()}`);
    };

    return (
        <Body>
            <Row className={styles.page} gutter={[0, isCompactMode ? 10 : 12]}>
                <Col span={24}>
                    <section className={styles.workspace}>
                        <section className={styles.hero}>
                            <div className="hero-copy">
                                <div className="eyebrow"><span className="status-dot"/>SYSTEM SETTINGS</div>
                                <h1>系统设置</h1>
                            </div>
                            <div className="hero-visual"><HeroGraphic variant="task"/></div>
                        </section>
                        <div className={styles.settingsMain}>
                            <div className={styles.leftMenu}>
                                <div className={styles.menuScroll}>
                                    <Menu
                                        className={styles.menu}
                                        mode="inline"
                                        selectedKeys={[activeItem.key]}
                                        onClick={({key}) => changeActive(key)}
                                        items={items.map((item) => ({
                                            key: item.key,
                                            label: item.title,
                                            icon: <Icon type={item.icon}/>,
                                        }))}
                                    />
                                </div>
                            </div>
                            <section className={styles.right}>
                                <div className={styles.title}>{activeItem.title}</div>
                                <Component styles={styles} info={web?.info}/>
                            </section>
                        </div>
                    </section>
                </Col>
            </Row>
        </Body>
    );
};
