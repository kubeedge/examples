declare const MAC: {
    assocStatus: {
        SUCCESSFUL_ASSOCIATION: number;
        PAN_AT_CAPACITY: number;
        PAN_ACCESS_DENIED: number;
    };
    channelPage: {
        PAGE_0: number;
        PAGE_1: number;
        PAGE_2: number;
    };
    txOpt: {
        UNDEFINED: number;
        ACK_TRANS: number;
        GTS_TRANS: number;
        IND_TRANS: number;
        SEC_ENABLED_TRANS: number;
        NO_RE_TRANS: number;
        NO_CONFIRM_TRANS: number;
        USE_PIB_VALUE: number;
        USE_POWER_CHANNEL_VALUES: number;
    };
    commReason: {
        ASSOCIATE_RSP: number;
        ORPHAN_RSP: number;
        RX_SECURE: number;
    };
    disassocReason: {
        RESERVED: number;
        COOR_WISHES_DEV_LEAVE: number;
        DEV_WISHES_LEAVE: number;
    };
    keyIdMode: {
        MODE_NONE_OR_IMPLICIT: number;
        MODE_1: number;
        MODE_4: number;
        MODE_8: number;
    };
    beaconOrder: {
        ORDER_NO_BEACONS: number;
        ORDER_4_MINUTES: number;
        ORDER_2_MINUTES: number;
        ORDER_1_MINUTE: number;
        ORDER_31_SECONDS: number;
        ORDER_15_SECONDS: number;
        ORDER_7_5_SECONDS: number;
        ORDER_4_SECONDS: number;
        ORDER_2_SECONDS: number;
        ORDER_1_SECOND: number;
        ORDER_480_MSEC: number;
        ORDER_240_MSEC: number;
        ORDER_120_MSEC: number;
        ORDER_60_MSEC: number;
        ORDER_30_MSEC: number;
        ORDER_15_MSEC: number;
    };
    scanType: {
        ENERGY_DETECT: number;
        ACTIVE: number;
        PASSIVE: number;
        ORPHAN: number;
        ENHANCED: number;
    };
    frontEndMode: {
        PA_LNA_OFF: number;
        PA_LNA_ON: number;
    };
    pidAttr: {
        ACK_WAIT_DURATION: number;
        ASSOCIATION_PERMIT: number;
        AUTO_REQUEST: number;
        BATT_LIFE_EXT: number;
        BATT_LIFE_EXT_PERIODS: number;
        BEACON_PAYLOAD: number;
        BEACON_PAYLOAD_LENGTH: number;
        BEACON_ORDER: number;
        BEACON_TX_TIME: number;
        BSN: number;
        COORD_EXTENDED_ADDRESS: number;
        COORD_SHORT_ADDRESS: number;
        DSN: number;
        GTS_PERMIT: number;
        MAX_CSMA_BACKOFFS: number;
        MIN_BE: number;
        PAN_ID: number;
        PROMISCUOUS_MODE: number;
        RX_ON_WHEN_IDLE: number;
        SHORT_ADDRESS: number;
        SUPERFRAME_ORDER: number;
        TRANSACTION_PERSISTENCE_TIME: number;
        ASSOCIATED_PAN_COORD: number;
        MAX_BE: number;
        MAX_FRAME_TOTAL_WAIT_TIME: number;
        MAX_FRAME_RETRIES: number;
        RESPONSE_WAIT_TIME: number;
        SYNC_SYMBOL_OFFSET: number;
        TIMESTAMP_SUPPORTED: number;
        SECURITY_ENABLED: number;
        KEY_TABLE: number;
        KEY_TABLE_ENTRIES: number;
        DEVICE_TABLE: number;
        DEVICE_TABLE_ENTRIES: number;
        SECURITY_LEVEL_TABLE: number;
        SECURITY_LEVEL_TABLE_ENTRIES: number;
        FRAME_COUNTER: number;
        AUTO_REQUEST_SECURITY_LEVEL: number;
        AUTO_REQUEST_KEY_ID_MODE: number;
        AUTO_REQUEST_KEY_SOURCE: number;
        AUTO_REQUEST_KEY_INDEX: number;
        DEFAULT_KEY_SOURCE: number;
        PAN_COORD_EXTENDED_ADDRESS: number;
        PAN_COORD_SHORT_ADDRESS: number;
        KEY_ID_LOOKUP_ENTRY: number;
        KEY_DEVICE_ENTRY: number;
        KEY_USAGE_ENTRY: number;
        KEY_ENTRY: number;
        DEVICE_ENTRY: number;
        SECURITY_LEVEL_ENTRY: number;
        PHY_TRANSMIT_POWER: number;
        LOGICAL_CHANNEL: number;
        EXTENDED_ADDRESS: number;
        ALT_BE: number;
        DEVICE_BEACON_ORDER: number;
        PHY_TRANSMIT_POWER_SIGNED: number;
    };
};
export default MAC;
//# sourceMappingURL=mac.d.ts.map