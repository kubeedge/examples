declare const SAPI: {
    zbDeviceInfo: {
        DEV_STATE: number;
        IEEE_ADDR: number;
        SHORT_ADDR: number;
        PARENT_SHORT_ADDR: number;
        PARENT_IEEE_ADDR: number;
        CHANNEL: number;
        PAN_ID: number;
        EXT_PAN_ID: number;
    };
    bindAction: {
        REMOVE_BIND: number;
        CREATE_BIND: number;
    };
    searchType: {
        ZB_IEEE_SEARCH: number;
    };
    txOptAck: {
        NONE: number;
        END_TO_END_ACK: number;
    };
};
export default SAPI;
//# sourceMappingURL=sapi.d.ts.map