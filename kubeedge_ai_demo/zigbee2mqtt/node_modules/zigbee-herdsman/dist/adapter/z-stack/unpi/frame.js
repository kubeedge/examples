"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("./constants");
class Frame {
    constructor(type, subsystem, commandID, data, length = null, fcs = null) {
        this.type = type;
        this.subsystem = subsystem;
        this.commandID = commandID;
        this.data = data;
        this.length = length;
        this.fcs = fcs;
    }
    toBuffer() {
        const length = this.data.length;
        const cmd0 = ((this.type << 5) & 0xE0) | (this.subsystem & 0x1F);
        let payload = Buffer.from([constants_1.SOF, length, cmd0, this.commandID]);
        payload = Buffer.concat([payload, this.data]);
        const fcs = Frame.calculateChecksum(payload.slice(1, payload.length));
        return Buffer.concat([payload, Buffer.from([fcs])]);
    }
    static fromBuffer(length, fcsPosition, buffer) {
        const subsystem = buffer.readUInt8(constants_1.PositionCmd0) & 0x1F;
        const type = (buffer.readUInt8(constants_1.PositionCmd0) & 0xE0) >> 5;
        const commandID = buffer.readUInt8(constants_1.PositionCmd1);
        const data = buffer.slice(constants_1.DataStart, fcsPosition);
        const fcs = buffer.readUInt8(fcsPosition);
        // Validate the checksum to see if we fully received the message
        const checksum = this.calculateChecksum(buffer.slice(1, fcsPosition));
        if (checksum === fcs) {
            return new Frame(type, subsystem, commandID, data, length, fcs);
        }
        else {
            throw new Error("Invalid checksum");
        }
    }
    static calculateChecksum(values) {
        let checksum = 0;
        for (const value of values) {
            checksum ^= value;
        }
        return checksum;
    }
    toString() {
        return `${this.length} - ${this.type} - ${this.subsystem} - ${this.commandID} - ` +
            `[${[...this.data]}] - ${this.fcs}`;
    }
}
exports.default = Frame;
//# sourceMappingURL=frame.js.map