import {createStyles} from "antd-style";

export default createStyles(({css, token, isDarkMode}: any, props: any = {}) => {
    const compact = Boolean(props?.isCompactMode);
    const dark = typeof props?.isDarkMode === "boolean" ? props.isDarkMode : Boolean(isDarkMode);

    return {
        page: css`
            width: 100%;
            max-width: 1450px;
            margin: 0 auto;

            > .ant-col {
                min-width: 0;
            }
        `,
        workspace: css`
            container-name: setting-workspace;
            container-type: inline-size;
            min-width: 0;
            overflow: hidden;
            border-radius: ${token.borderRadiusLG}px;
            background: ${token.colorBgContainer};
            box-shadow: ${token.boxShadowTertiary};
        `,
        hero: css`
            position: relative;
            min-height: ${compact ? 104 : 118}px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            overflow: hidden;
            isolation: isolate;
            background: linear-gradient(
                118deg,
                color-mix(in srgb, ${token.colorPrimary}, ${token.colorBgContainer} ${dark ? 94 : 93}%) 0%,
                ${token.colorBgContainer} 58%,
                color-mix(in srgb, #13c2c2, ${token.colorBgContainer} ${dark ? 97 : 96}%) 100%
            );

            &::before {
                position: absolute;
                z-index: 0;
                inset: 0;
                background-image:
                    linear-gradient(${dark ? "rgba(255,255,255,0.018)" : "rgba(72,132,202,0.035)"} 1px, transparent 1px),
                    linear-gradient(90deg, ${dark ? "rgba(255,255,255,0.018)" : "rgba(72,132,202,0.035)"} 1px, transparent 1px);
                background-size: 24px 24px;
                content: "";
                -webkit-mask-image: linear-gradient(90deg, transparent 36%, rgba(0, 0, 0, .2) 58%, #000 100%);
                mask-image: linear-gradient(90deg, transparent 36%, rgba(0, 0, 0, .2) 58%, #000 100%);
                pointer-events: none;
            }

            &::after {
                position: absolute;
                z-index: 0;
                inset: -45% -10% -55% 48%;
                background:
                    radial-gradient(ellipse at 48% 32%, color-mix(in srgb, ${token.colorPrimary}, transparent 84%) 0%, transparent 60%),
                    radial-gradient(ellipse at 76% 72%, color-mix(in srgb, #13c2c2, transparent 89%) 0%, transparent 58%);
                content: "";
                filter: blur(28px);
                opacity: ${dark ? .34 : .52};
                pointer-events: none;
            }

            .hero-copy {
                position: relative;
                z-index: 2;
                min-width: 0;
                max-width: calc(100% - 178px);
                padding: ${compact ? "14px 18px" : "18px 22px"};
            }

            .eyebrow {
                display: flex;
                align-items: center;
                gap: 7px;
                color: color-mix(in srgb, ${token.colorPrimary}, ${token.colorTextTertiary} 62%);
                font-size: 9px;
                font-weight: 600;
                line-height: 16px;
                letter-spacing: 0;
            }

            .eyebrow .status-dot {
                width: 5px;
                height: 5px;
                flex: none;
                border-radius: 50%;
                background: color-mix(in srgb, ${token.colorPrimary}, #fff 26%);
                box-shadow: 0 0 0 4px color-mix(in srgb, ${token.colorPrimary}, transparent 90%);
            }

            h1 {
                margin: 3px 0 0;
                color: ${token.colorTextHeading};
                font-size: ${compact ? 19 : 20}px;
                font-weight: 600;
                line-height: ${compact ? 26 : 28}px;
                letter-spacing: 0;
            }

            .hero-visual {
                position: absolute;
                z-index: 1;
                top: 0;
                right: 82px;
                width: 52%;
                height: 100%;
                opacity: ${dark ? .12 : .25};
                filter: ${dark ? "blur(.8px) saturate(.78)" : "blur(.25px)"};
                -webkit-mask-image: linear-gradient(90deg, transparent 0%, rgba(0, 0, 0, .38) 28%, #000 52%);
                mask-image: linear-gradient(90deg, transparent 0%, rgba(0, 0, 0, .38) 28%, #000 52%);
                pointer-events: none;
            }

            @container setting-workspace (max-width: 600px) {
                min-height: ${compact ? 104 : 116}px;

                .hero-copy {
                    max-width: calc(100% - 145px);
                    padding: 16px 14px;
                }

                h1 {
                    font-size: 18px;
                    line-height: 25px;
                }

                .hero-visual {
                    right: 42px;
                    width: 68%;
                    opacity: ${dark ? .07 : .14};
                }

                .hero-visual .flow-secondary,
                .hero-visual .flow-detail {
                    display: none;
                }
            }

            @container setting-workspace (max-width: 430px) {
                .hero-visual {
                    right: 16px;
                    width: 76%;
                    opacity: ${dark ? .045 : .09};
                }

                .hero-visual .flow-node-a {
                    display: none;
                }
            }

            @container setting-workspace (max-width: 360px) {
                .hero-visual {
                    display: none;
                }
            }
        `,
        settingsMain: css`
            width: 100%;
            height: 100%;
            padding: ${token.paddingLG}px;
            display: flex;
            box-sizing: border-box;

            .ant-menu-light.ant-menu-inline .ant-menu-item::after {
                top: 19%;
                right: 6px;
                height: 60%;
                border-top-left-radius: 15px;
                border-bottom-left-radius: 15px;
            }

            @container setting-workspace (max-width: 760px) {
                flex-direction: column;
                padding: ${compact ? 10 : 12}px;
            }
        `,
        leftMenu: css`
            width: 224px;
            flex: none;

            .ant-menu-item {
                border-radius: ${token.borderRadius}px;
                font-size: ${token.fontSize}px;
                font-weight: bold;
            }

            @container setting-workspace (max-width: 760px) {
                width: 100%;
                margin-bottom: 10px;
                overflow: hidden;
                border: none;
            }
        `,
        menuScroll: css`
            min-width: 0;

            @container setting-workspace (max-width: 760px) {
                width: 100%;
                overflow-x: auto;
                overflow-y: hidden;
                overscroll-behavior-x: contain;
                scrollbar-width: none;
                touch-action: pan-x;
                -webkit-overflow-scrolling: touch;

                &::-webkit-scrollbar {
                    display: none;
                }
            }
        `,
        menu: css`
            border-inline-end: none !important;

            &:focus,
            &:focus-visible {
                outline: none !important;
            }

            @container setting-workspace (max-width: 760px) {
                display: flex !important;
                flex-direction: row;
                flex-wrap: nowrap;
                gap: 6px;
                width: max-content;
                min-width: 100%;
                padding: 0;
                background: transparent;

                .ant-menu-item {
                    width: auto !important;
                    min-width: 96px;
                    height: 34px;
                    min-height: 34px;
                    margin: 0 !important;
                    padding-inline: 12px !important;
                    display: flex;
                    flex: 0 0 auto;
                    align-items: center;
                    justify-content: center;
                    background: ${token.colorFillQuaternary};
                    color: ${token.colorTextSecondary};
                    font-weight: 400;
                    line-height: 34px;
                    white-space: nowrap;
                }

                .ant-menu-item::after {
                    display: none !important;
                }

                .ant-menu-item-selected {
                    background: ${token.colorPrimaryBg} !important;
                    color: ${token.colorTextHeading} !important;
                    font-weight: 600;
                }
            }
        `,
        right: css`
            min-width: 0;
            flex: 1;
            margin-left: 1px;
            padding: 8px 40px;
            border-left: 1px solid ${token.colorSplit};

            @container setting-workspace (max-width: 760px) {
                margin-left: 0;
                padding: 10px;
                border-left: none;
            }
        `,
        title: css`
            margin-bottom: 20px;
            color: ${token.colorTextHeading};
            font-size: ${token.fontSizeHeading4}px;
            font-weight: bolder;
            line-height: 28px;
        `,
        loading: css`
            width: 100%;
            min-height: 240px;
            display: flex;
            align-items: center;
            justify-content: center;
        `,
        uiWrapper: css`
            width: 100%;
            max-width: 720px;
        `,
        formField: css`
            width: 100%;
            max-width: ${compact ? 300 : 350}px;
        `,
        formActions: css`
            margin-top: ${token.marginLG}px;
            display: flex;
            align-items: center;
            gap: ${token.marginSM}px;
        `,
        errorState: css`
            min-height: 240px;
            padding: 24px;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            gap: ${token.marginSM}px;
            color: ${token.colorTextSecondary};
            font-size: ${token.fontSize}px;
            text-align: center;

            > .anticon {
                color: ${token.colorError};
                font-size: ${token.fontSizeHeading3}px;
            }
        `,
    };
});
