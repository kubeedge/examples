declare const ZDO: {
    status: {
        SUCCESS: number;
        INVALID_REQTYPE: number;
        DEVICE_NOT_FOUND: number;
        INVALID_EP: number;
        NOT_ACTIVE: number;
        NOT_SUPPORTED: number;
        TIMEOUT: number;
        NO_MATCH: number;
        NO_ENTRY: number;
        NO_DESCRIPTOR: number;
        INSUFFICIENT_SPACE: number;
        NOT_PERMITTED: number;
        TABLE_FULL: number;
        NOT_AUTHORIZED: number;
        BINDING_TABLE_FULL: number;
    };
    initDev: {
        RESTORED_NETWORK_STATE: number;
        NEW_NETWORK_STATE: number;
        LEAVE_NOT_STARTED: number;
    };
    serverCapability: {
        NOT_SUPPORTED: number;
        PRIM_TRUST_CENTER: number;
        BKUP_TRUST_CENTER: number;
        PRIM_BIND_TABLE: number;
        BKUP_BIND_TABLE: number;
        PRIM_DISC_TABLE: number;
        BKUP_DISC_TABLE: number;
        NETWORK_MANAGER: number;
    };
    appDevVer: {
        VER_100: number;
        RESERVE01: number;
        RESERVE02: number;
        RESERVE03: number;
        RESERVE04: number;
        RESERVE05: number;
        RESERVE06: number;
        RESERVE07: number;
        RESERVE08: number;
        RESERVE09: number;
        RESERVE10: number;
        RESERVE11: number;
        RESERVE12: number;
        RESERVE13: number;
        RESERVE14: number;
        RESERVE15: number;
    };
    stackProfileId: {
        NETWORK_SPECIFIC: number;
        HOME_CONTROLS: number;
        ZIGBEEPRO_PROFILE: number;
        GENERIC_STAR: number;
        GENERIC_TREE: number;
    };
    deviceLogicalType: {
        COORDINATOR: number;
        ROUTER: number;
        ENDDEVICE: number;
        COMPLEX_DESC_AVAIL: number;
        USER_DESC_AVAIL: number;
        RESERVED1: number;
        RESERVED2: number;
        RESERVED3: number;
        RESERVED4: number;
    };
    addrReqType: {
        SINGLE: number;
        EXTENDED: number;
    };
    leaveAndRemoveChild: {
        NONE: number;
        LEAVE_REMOVE_CHILDREN: number;
    };
    leaveIndRequest: {
        INDICATION: number;
        REQUEST: number;
    };
    leaveIndRemove: {
        NONE: number;
        REMOVE_CHILDREN: number;
    };
    leaveIndRejoin: {
        NONE: number;
        REJOIN: number;
    };
    descCapability: {
        EXT_LIST_NOT_SUPPORTED: number;
        EXT_ACTIVE_EP_LIST_AVAIL: number;
        EXT_SIMPLE_DESC_LIST_AVAIL: number;
        RESERVED1: number;
        RESERVED2: number;
        RESERVED3: number;
        RESERVED4: number;
        RESERVED5: number;
        RESERVED6: number;
    };
};
export default ZDO;
//# sourceMappingURL=zdo.d.ts.map