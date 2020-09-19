declare const SYS: {
    resetType: {
        HARD: number;
        SOFT: number;
    };
    capabilities: {
        SYS: number;
        MAC: number;
        NWK: number;
        AF: number;
        ZDO: number;
        SAPI: number;
        UTIL: number;
        DEBUG: number;
        APP: number;
        ZOAD: number;
    };
    osalTimerEvent: {
        EVENT_0: number;
        EVENT_1: number;
        EVENT_2: number;
        EVENT_3: number;
    };
    adcChannels: {
        AIN0: number;
        AIN1: number;
        AIN2: number;
        AIN3: number;
        AIN4: number;
        AIN5: number;
        AIN6: number;
        AIN7: number;
        TEMP_SENSOR: number;
        VOLT_READ: number;
    };
    adcResolution: {
        BIT_8: number;
        BIT_10: number;
        BIT_12: number;
        BIT_14: number;
    };
    gpioOperation: {
        SET_DIRECTION: number;
        SET_INPUT_MODE: number;
        SET: number;
        CLEAR: number;
        TOGGLE: number;
        READ: number;
    };
    sysStkTune: {
        TX_PWR: number;
        RX_ON_IDLE: number;
    };
    resetReason: {
        POWER_UP: number;
        EXTERNAL: number;
        WATCH_DOG: number;
    };
    nvItemInitStatus: {
        ALREADY_EXISTS: number;
        SUCCESS: number;
        FAILED: number;
    };
    nvItemDeleteStatus: {
        SUCCESS: number;
        NOT_EXISTS: number;
        FAILED: number;
        BAD_LENGTH: number;
    };
};
export default SYS;
//# sourceMappingURL=sys.d.ts.map