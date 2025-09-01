import React from 'react';
import {ThemeProvider} from 'antd-style';
import {useModel} from "umi";

export type ThemeProps = {
    theme?: any;//主题
    mode?: any;//模式
    appearance?: any;//主题
    onAppearanceChange?: (appearance: any) => void;//监听主题改变
    children?: any;//子内容
};

const Theme: React.FC<ThemeProps> = ({...props}) => {
    const systemTheme = useModel("theme");
    return (
        <ThemeProvider theme={props?.theme ? props?.theme : systemTheme?.themes}
                       themeMode={props?.mode ? props?.mode : systemTheme?.mode}
                       appearance={props?.appearance ? props?.appearance : systemTheme?.appearance}
                       onAppearanceChange={(appearance) => {
                           props?.onAppearanceChange?.(appearance);
                           systemTheme?.setAppearance?.(appearance);
                       }}
        >
            {props?.children}
        </ThemeProvider>
    );
};

export default Theme;