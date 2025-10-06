import {createFromIconfontCN} from '@ant-design/icons';

/**
 * 图标组件(生产环境需使用本地资源)
 */
const Icon = createFromIconfontCN({
    scriptUrl: '//at.alicdn.com/t/c/font_5013902_1i6dqszic9d.js',
});
/**
 * 图标数据
 */
const Exit = "icon-exit";
const Right = "icon-right";
const Log = "icon-log";
const Set = "icon-set";
const Bottom = "icon-bottom";
const Light = "icon-light";//亮色
const Dark = "icon-dark";
const Auto = "icon-auto";
const Home = "icon-home";
const Logo = "icon-logo";
export {
    Icon,
    Exit,
    Right,
    Log,
    Set,
    Bottom,
    Light,
    Dark,
    Auto,
    Home,
    Logo,
}