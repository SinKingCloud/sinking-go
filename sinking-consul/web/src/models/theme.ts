import {useState} from "react";
import {theme} from "antd";
import defaultSettings from "../../config/defaultSettings";
// 获取风格
const getDefaultTheme = (color, radius): any => {
    return {
        token: {
            colorPrimary: color,
            colorInfo: color,
            borderRadius: radius
        },
    }
}
// 获取风格
const getCompactTheme = (color, radius): any => {
    let temp = getDefaultTheme(color, radius);
    temp.algorithm = [theme.compactAlgorithm];
    return temp;
}
/**
 * 获取主题模式
 */
const getMode = (): string => {
    const mode = localStorage?.getItem("theme")
    if (mode == "light" || mode == "dark") {
        return mode;
    }
    return "auto"
}
/**
 * 设置主题模式
 * @param mode 模式
 */
const setMode = (mode): void => {
    localStorage?.setItem("theme", mode);
}

export default () => {
    const [themes, setThemes] = useState<any>(getDefaultTheme(defaultSettings?.color, defaultSettings?.radius));//系统主题
    const [mode, setMode2] = useState<any>(getMode());//系统主题模式
    const [appearance, setAppearance] = useState<any>(null);//当前主题风格

    const lightMode = "light";//亮色

    const darkMode = "dark";//暗色

    const autoMode = "auto";//跟随系统

    /**
     * 设置默认主题
     */
    const setDefaultTheme = () => {
        setThemes(getDefaultTheme(themes?.token?.colorPrimary, themes?.token?.borderRadius))
    }

    /**
     * 设置紧凑主题
     */
    const setCompactTheme = () => {
        setThemes(getCompactTheme(themes?.token?.colorPrimary, themes?.token?.borderRadius))
    }

    /**
     * 设置主题颜色
     */
    const setColor = (color) => {
        let temp = {...themes}
        temp.token.colorPrimary = color;
        temp.token.colorInfo = color;
        setThemes(temp);
    }

    /**
     * 设置主题圆角
     */
    const setRadius = (radius) => {
        let temp = {...themes}
        temp.token.borderRadius = radius;
        setThemes(temp);
    }

    /**
     * 设置模式名称
     */
    const getModeName = (mode) => {
        if (mode == lightMode) {
            return "亮色风格";
        }
        if (mode == darkMode) {
            return "暗色风格";
        }
        return "跟随系统";
    }

    /**
     * 设置亮色模式
     */
    const setLightMode = () => {
        setMode(lightMode);
        setMode2(lightMode)
        setAppearance(lightMode);
    }

    /**
     * 设置暗色模式
     */
    const setDarkMode = () => {
        setMode(darkMode);
        setMode2(darkMode);
        setAppearance(darkMode);
    }

    /**
     * 设置跟随系统模式
     */
    const setAutoMode = () => {
        setMode(autoMode);
        setMode2(autoMode);
        setAppearance(null);
    }

    /**
     * 是否为亮色模式
     */
    const isLightMode = () => {
        return mode == lightMode;
    }

    /**
     * 是否为暗色模式
     */
    const isDarkMode = () => {
        return mode == darkMode;
    }

    /**
     * 是否为跟随系统模式
     */
    const isAutoMode = () => {
        return mode == autoMode;
    }

    /**
     * 当前是否为亮色风格
     */
    const isLightTheme = () => {
        return appearance == lightMode;
    }

    /**
     * 当前是否为暗色风格
     */
    const isDarkTheme = () => {
        return appearance == darkMode;
    }

    /**
     * 模式切换
     */
    const toggle = () => {
        if (isAutoMode()) {
            setLightMode();
        } else if (isLightMode()) {
            setDarkMode();
        } else {
            setAutoMode();
        }
    }

    return {
        themes,
        setColor,
        setRadius,
        setThemes,
        setDefaultTheme,
        setCompactTheme,
        appearance,
        setAppearance,
        mode,
        getModeName,
        lightMode,
        darkMode,
        autoMode,
        setLightMode,
        setDarkMode,
        setAutoMode,
        isLightMode,
        isDarkMode,
        isAutoMode,
        isLightTheme,
        isDarkTheme,
        toggle
    };
};
