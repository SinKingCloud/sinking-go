import React from "react";
import {Button} from "antd";

export default ({styles, saving, onReset}: any): React.ReactNode => (
    <div className={styles.formActions}>
        <Button onClick={onReset}>重置</Button>
        <Button type="primary" htmlType="submit" loading={saving}>保存</Button>
    </div>
);
