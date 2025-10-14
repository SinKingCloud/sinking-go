import React, {createContext, useContext, useState, ReactNode} from 'react';
import {ThemeProvider, ThemeMode} from 'antd-style';
import {theme} from 'antd';
import defaultSettings from '../../../config/defaultSettings';

// 获取风格
const getDefaultTheme = (color: string, radius: number): any => {
    return {
        token: {
            colorPrimary: color,
            colorInfo: color,
            borderRadius: radius
        },
    }
}

// 获取紧凑风格
const getCompactTheme = (color: string, radius: number): any => {
    let temp = getDefaultTheme(color, radius);
    temp.algorithm = [theme.compactAlgorithm];
    return temp;
}

/**
 * 获取主题模式
 */
const getMode = (): ThemeMode => {
    const mode = localStorage?.getItem("theme")
    if (mode == "light" || mode == "dark") {
        return mode as ThemeMode;
    }
    return "auto" as ThemeMode;
}

/**
 * 设置主题模式到 localStorage
 * @param mode 模式
 */
const setModeToStorage = (mode: string): void => {
    localStorage?.setItem("theme", mode);
}

// 定义主题上下文类型
export interface ThemeContextType {
    themes: any;
    setColor: (color: string) => void;
    setRadius: (radius: number) => void;
    setThemes: (theme: any) => void;
    setDefaultTheme: () => void;
    setCompactTheme: () => void;
    appearance: any;
    setAppearance: (appearance: any) => void;
    mode: ThemeMode;
    getModeName: (mode: ThemeMode) => string;
    lightMode: ThemeMode;
    darkMode: ThemeMode;
    autoMode: ThemeMode;
    setLightMode: () => void;
    setDarkMode: () => void;
    setAutoMode: () => void;
    isLightMode: () => boolean;
    isDarkMode: () => boolean;
    isAutoMode: () => boolean;
    isLightTheme: () => boolean;
    isDarkTheme: () => boolean;
    isCompactTheme: () => boolean;
    toggle: () => void;
}

// 创建主题上下文
const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

// 导出 useTheme hook
export const useTheme = (): ThemeContextType | undefined => {
    const context = useContext(ThemeContext);
    if (!context) {
        return undefined;
    }
    return context;
};

export type ThemeProps = {
    theme?: any; // 自定义主题（可选）
    mode?: ThemeMode; // 自定义模式（可选）
    appearance?: any; // 自定义外观（可选）
    onAppearanceChange?: (appearance: any) => void; // 监听主题改变（可选）
    children?: ReactNode;
};

// 主题 Provider 组件（合并后的统一组件）
const Theme: React.FC<ThemeProps> = ({
                                         children,
                                         theme: customTheme,
                                         mode: customMode,
                                         appearance: customAppearance,
                                         onAppearanceChange
                                     }) => {
    // 检查是否已在 Theme 上下文中（避免重复嵌套）
    const existingContext = useContext(ThemeContext);

    // 如果已存在上下文且没有自定义属性，直接返回子组件（避免不必要的嵌套）
    if (existingContext && !customTheme && !customMode && !customAppearance && !onAppearanceChange) {
        return <>{children}</>;
    }

    // 如果已存在上下文且有自定义属性，只需要覆盖 ThemeProvider 的属性
    if (existingContext) {
        const finalTheme = customTheme || existingContext.themes;
        const finalMode = customMode || existingContext.mode;
        const finalAppearance = customAppearance !== undefined ? customAppearance : existingContext.appearance;

        return (
            <ThemeProvider
                theme={finalTheme}
                themeMode={finalMode}
                appearance={finalAppearance}
                onAppearanceChange={(newAppearance) => {
                    onAppearanceChange?.(newAppearance);
                    existingContext?.setAppearance(newAppearance);
                }}>
                {children}
            </ThemeProvider>
        );
    }

    // 以下是创建新的主题上下文的逻辑（仅在不存在上下文时执行）
    return <ThemeProviderInner
        theme={customTheme}
        mode={customMode}
        appearance={customAppearance}
        onAppearanceChange={onAppearanceChange}>
        {children}
    </ThemeProviderInner>;
};

// 内部 Provider 组件（包含完整的状态管理逻辑）
const ThemeProviderInner: React.FC<ThemeProps> = ({
                                                      children,
                                                      theme: customTheme,
                                                      mode: customMode,
                                                      appearance: customAppearance,
                                                      onAppearanceChange
                                                  }) => {
    const [themes, setThemes] = useState<any>(customTheme || getDefaultTheme(defaultSettings?.color || "", defaultSettings?.radius || 0));
    const [mode, setMode2] = useState<ThemeMode>(customMode || getMode());
    const [appearance, setAppearance] = useState<any>(customAppearance !== undefined ? customAppearance : null);

    const lightMode: ThemeMode = "light";
    const darkMode: ThemeMode = "dark";
    const autoMode: ThemeMode = "auto";

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
    const setColor = (color: string) => {
        let temp = {...themes}
        temp.token.colorPrimary = color;
        temp.token.colorInfo = color;
        setThemes(temp);
    }

    /**
     * 设置主题圆角
     */
    const setRadius = (radius: number) => {
        let temp = {...themes}
        temp.token.borderRadius = radius;
        setThemes(temp);
    }

    /**
     * 获取模式名称
     */
    const getModeName = (modeValue: string) => {
        if (modeValue == lightMode) {
            return "亮色风格";
        }
        if (modeValue == darkMode) {
            return "暗色风格";
        }
        return "跟随系统";
    }

    /**
     * 设置亮色模式
     */
    const setLightMode = () => {
        setModeToStorage(lightMode);
        setMode2(lightMode)
        setAppearance(lightMode);
    }

    /**
     * 设置暗色模式
     */
    const setDarkMode = () => {
        setModeToStorage(darkMode);
        setMode2(darkMode);
        setAppearance(darkMode);
    }

    /**
     * 设置跟随系统模式
     */
    const setAutoMode = () => {
        setModeToStorage(autoMode);
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
     * 当前是否为紧凑主题
     */
    const isCompactTheme = () => {
        return themes?.algorithm && Array.isArray(themes.algorithm) && themes.algorithm.includes(theme.compactAlgorithm);
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

    const value: ThemeContextType = {
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
        isCompactTheme,
        toggle
    };

    return (
        <ThemeContext.Provider value={value}>
            <ThemeProvider
                theme={themes}
                themeMode={mode}
                appearance={appearance}
                onAppearanceChange={(newAppearance) => {
                    onAppearanceChange?.(newAppearance);
                    setAppearance(newAppearance);
                }}>
                {children}
            </ThemeProvider>
        </ThemeContext.Provider>
    );
};

export default Theme;
