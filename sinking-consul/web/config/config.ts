import {defineConfig} from "umi";
import defaultSettings from './defaultSettings';
import routes from "./routes";

export default defineConfig({
    title: defaultSettings?.title,
    routes: routes,
    base: defaultSettings.basePath,
    publicPath: defaultSettings.basePath,  //静态资源基本路径
    outputPath: './dist' + defaultSettings.basePath,//资源输出路径
    plugins: [
        '@umijs/plugins/dist/model',
    ],
    model: {},
    favicons: [
        defaultSettings?.favicons?.toString(),
    ],
    history: {
        type: 'browser'
    },
    hash: true,
    exportStatic: {},
} as any);