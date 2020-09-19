declare enum Type {
    POLL = 0,
    SREQ = 1,
    AREQ = 2,
    SRSP = 3
}
declare enum Subsystem {
    RESERVED = 0,
    SYS = 1,
    MAC = 2,
    NWK = 3,
    AF = 4,
    ZDO = 5,
    SAPI = 6,
    UTIL = 7,
    DEBUG = 8,
    APP = 9,
    APP_CNF = 15,
    GREENPOWER = 21
}
declare const DataStart = 4;
declare const SOF = 254;
declare const PositionDataLength = 1;
declare const PositionCmd0 = 2;
declare const PositionCmd1 = 3;
declare const MinMessageLength = 5;
declare const MaxDataSize = 250;
export { Type, Subsystem, DataStart, SOF, PositionDataLength, MinMessageLength, PositionCmd0, PositionCmd1, MaxDataSize, };
//# sourceMappingURL=constants.d.ts.map