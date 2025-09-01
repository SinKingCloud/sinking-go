import 'umi/typings';

declare namespace API {

    /**
     * 快捷请求结构
     */
    type RequestParams = {
        body?: any,
        onSuccess?: (res) => void,
        onFail?: (res) => void,
        onFinally?: () => void
    }

    /**
     * 通用响应
     */
    type Response = {
        code?: number;
        message?: string;
        data?: any;
        request_id?: string;
    }
    /**
     * 站点信息
     */
    type WebInfo = {
        name?: string;
        contact?: string;
    }

    /**
     * 网站用户信息
     */
    /**
     * 用户信息
     */
    type UserInfo = {
        id?: number;
        money?: any;
        nick_name?: string;
        avatar?: string;
        phone?: string;
        login_ip?: string;
        login_time?: string;
        is_admin?: boolean;
    }
    type Ui = {
        logo?: string;
        water_mark?: string;
        layout?: string;
        theme?: string;
        color?: string;
        compact?: boolean;
    }
    /**
     * 网站信息
     */
    type WebInfo = {
        name?: string;
        title?: string;
        keywords?: string;
        describe?: string;
        contact?: string;
        url?: string;
        ui?: Ui,
    }
}

