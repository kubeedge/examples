"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.attributeList = exports.attributeKeyValue = void 0;
function attributeKeyValue(frame) {
    const payload = {};
    for (const item of frame.Payload) {
        try {
            const attribute = frame.Cluster.getAttribute(item.attrId);
            payload[attribute.name] = item.attrData;
        }
        catch (error) {
            payload[item.attrId] = item.attrData;
        }
    }
    return payload;
}
exports.attributeKeyValue = attributeKeyValue;
function attributeList(frame) {
    const payload = [];
    for (const item of frame.Payload) {
        try {
            const attribute = frame.Cluster.getAttribute(item.attrId);
            payload.push(attribute.name);
        }
        catch (error) {
            payload.push(item.attrId);
        }
    }
    return payload;
}
exports.attributeList = attributeList;
//# sourceMappingURL=zclFrameConverter.js.map