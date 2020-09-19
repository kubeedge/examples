interface PortInfoMatch {
    manufacturer: string;
    vendorId: string;
    productId: string;
}
declare function find(matchers: PortInfoMatch[]): Promise<string[]>;
declare function is(path: string, matchers: PortInfoMatch[]): Promise<boolean>;
declare const _default: {
    is: typeof is;
    find: typeof find;
};
export default _default;
//# sourceMappingURL=serialPortUtils.d.ts.map