declare const UTIL: {
    getNvStatus: {
        SUCCESS: number;
        GET_IEEE_ADDR_FAIL: number;
        GET_SCAN_CHANNEL_FAIL: number;
        GET_PAN_ID_FAIL: number;
        GET_SECURITY_LEVEL_FAIL: number;
        GET_PRECONFIG_KEY_FAIL: number;
    };
    subsystemId: {
        SYS: number;
        MAC: number;
        NWK: number;
        AF: number;
        ZDO: number;
        SAPI: number;
        UTIL: number;
        DBG: number;
        APP: number;
        ALL_SUBSYSTEM: number;
    };
    deviceType: {
        NONE: number;
        COORDINATOR: number;
        ROUTER: number;
        END_DEVICE: number;
    };
    keyEvent: {
        KEY_1: number;
        KEY_2: number;
        KEY_3: number;
        KEY_4: number;
        KEY_5: number;
        KEY_6: number;
        KEY_7: number;
        KEY_8: number;
    };
    keyValue: {
        KEY_1: number;
        KEY_2: number;
        KEY_3: number;
        KEY_4: number;
        KEY_5: number;
        KEY_6: number;
        KEY_7: number;
        KEY_8: number;
    };
    ledMode: {
        OFF: number;
        ON: number;
        BLINK: number;
        FLASH: number;
        TOGGLE: number;
    };
    ledNum: {
        LED_1: number;
        LED_2: number;
        LED_3: number;
        LED_4: number;
        ALL_LEDS: number;
    };
    subsAction: {
        UNSUBSCRIBE: number;
        SUBSCRIBE: number;
    };
    ackPendingOption: {
        ACK_DISABLE: number;
        ACK_ENABLE: number;
    };
    nodeRelation: {
        PARENT: number;
        CHILD_RFD: number;
        CHILD_RFD_RX_IDLE: number;
        CHILD_FFD: number;
        CHILD_FFD_RX_IDLE: number;
        NEIGHBOR: number;
        OTHER: number;
        NOTUSED: number;
    };
};
export default UTIL;
//# sourceMappingURL=util.d.ts.map