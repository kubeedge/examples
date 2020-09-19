"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const constants_1 = require("../unpi/constants");
const unpi_1 = require("../unpi");
const definition_1 = __importDefault(require("./definition"));
const buffaloZnp_1 = __importDefault(require("./buffaloZnp"));
const parameterType_1 = __importDefault(require("./parameterType"));
const BufferAndListTypes = [
    parameterType_1.default.BUFFER, parameterType_1.default.BUFFER8, parameterType_1.default.BUFFER16,
    parameterType_1.default.BUFFER18, parameterType_1.default.BUFFER32, parameterType_1.default.BUFFER42,
    parameterType_1.default.BUFFER100, parameterType_1.default.LIST_UINT16, parameterType_1.default.LIST_ROUTING_TABLE,
    parameterType_1.default.LIST_BIND_TABLE, parameterType_1.default.LIST_NEIGHBOR_LQI, parameterType_1.default.LIST_NETWORK,
    parameterType_1.default.LIST_ASSOC_DEV, parameterType_1.default.LIST_UINT8,
];
class ZpiObject {
    constructor(type, subsystem, command, commandID, payload, parameters) {
        this.subsystem = subsystem;
        this.command = command;
        this.commandID = commandID;
        this.payload = payload;
        this.type = type;
        this.parameters = parameters;
    }
    static createRequest(subsystem, command, payload) {
        if (!definition_1.default[subsystem]) {
            throw new Error(`Subsystem '${subsystem}' does not exist`);
        }
        const cmd = definition_1.default[subsystem].find((c) => c.name === command);
        if (!cmd) {
            throw new Error(`Command '${command}' from subsystem '${subsystem}' not found`);
        }
        return new ZpiObject(cmd.type, subsystem, command, cmd.ID, payload, cmd.request);
    }
    toUnpiFrame() {
        const buffer = this.createPayloadBuffer();
        return new unpi_1.Frame(this.type, this.subsystem, this.commandID, buffer);
    }
    static fromUnpiFrame(frame) {
        const cmd = definition_1.default[frame.subsystem].find((c) => c.ID === frame.commandID);
        if (!cmd) {
            throw new Error(`CommandID '${frame.commandID}' from subsystem '${frame.subsystem}' not found`);
        }
        const parameters = frame.type === constants_1.Type.SRSP ? cmd.response : cmd.request;
        if (parameters === undefined) {
            /* istanbul ignore next */
            throw new Error(`CommandID '${frame.commandID}' from subsystem '${frame.subsystem}' cannot be a ` +
                `${frame.type === constants_1.Type.SRSP ? 'response' : 'request'}`);
        }
        const payload = this.readParameters(frame.data, parameters);
        return new ZpiObject(frame.type, frame.subsystem, cmd.name, cmd.ID, payload, parameters);
    }
    static readParameters(buffer, parameters) {
        const buffalo = new buffaloZnp_1.default(buffer);
        const result = {};
        for (const parameter of parameters) {
            const options = {};
            if (BufferAndListTypes.includes(parameter.parameterType)) {
                // When reading a buffer, assume that the previous parsed parameter contains
                // the length of the buffer
                const lengthParameter = parameters[parameters.indexOf(parameter) - 1];
                const length = result[lengthParameter.name];
                /* istanbul ignore else */
                if (typeof length === 'number') {
                    options.length = length;
                }
                if (parameter.parameterType === parameterType_1.default.LIST_ASSOC_DEV) {
                    // For LIST_ASSOC_DEV, we also need to grab the startindex which is right before the length
                    const startIndexParameter = parameters[parameters.indexOf(parameter) - 2];
                    const startIndex = result[startIndexParameter.name];
                    /* istanbul ignore else */
                    if (typeof startIndex === 'number') {
                        options.startIndex = startIndex;
                    }
                }
            }
            result[parameter.name] = buffalo.read(parameterType_1.default[parameter.parameterType], options);
        }
        return result;
    }
    createPayloadBuffer() {
        const buffalo = new buffaloZnp_1.default(Buffer.alloc(constants_1.MaxDataSize));
        for (const parameter of this.parameters) {
            const value = this.payload[parameter.name];
            buffalo.write(parameterType_1.default[parameter.parameterType], value, {});
        }
        return buffalo.getBuffer().slice(0, buffalo.getPosition());
    }
    isResetCommand() {
        return (this.command === 'resetReq' && this.subsystem === constants_1.Subsystem.SYS) ||
            (this.command === 'systemReset' && this.subsystem === constants_1.Subsystem.SAPI);
    }
}
exports.default = ZpiObject;
//# sourceMappingURL=zpiObject.js.map