import React from "react";
import {Spin} from "antd";

export default ({styles}: any): React.ReactNode => (
    <div className={styles.loading}><Spin/></div>
);
