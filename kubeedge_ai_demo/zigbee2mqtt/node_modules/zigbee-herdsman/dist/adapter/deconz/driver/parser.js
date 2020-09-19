"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
/* istanbul ignore file */
/* eslint-disable */
const stream = __importStar(require("stream"));
// @ts-ignore
const slip_1 = __importDefault(require("slip"));
const debug_1 = __importDefault(require("debug"));
const debug = debug_1.default('zigbee-herdsman:deconz:driver:parser');
class Parser extends stream.Transform {
    constructor() {
        super();
        this.onMessage = this.onMessage.bind(this);
        this.onError = this.onError.bind(this);
        this.decoder = new slip_1.default.Decoder({
            onMessage: this.onMessage,
            maxMessageSize: 1000000,
            bufferSize: 2048
        });
    }
    onMessage(message) {
        //debug(`message received: ${message}`);
        this.emit('parsed', message);
    }
    onError(_, error) {
        debug(`<-- error '${error}'`);
    }
    _transform(chunk, _, cb) {
        //debug(`<-- [${[...chunk]}]`);
        this.decoder.decode(chunk);
        //debug(`<-- [${[...chunk]}]`);
        cb();
    }
}
exports.default = Parser;
//# sourceMappingURL=parser.js.map