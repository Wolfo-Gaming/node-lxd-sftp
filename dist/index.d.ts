/// <reference types="node" />
export type SFTPConfig = {
    serverCert: Buffer;
    clientCert: Buffer;
    clientKey: Buffer;
    host: string;
    instanceName: string;
    listenAddress?: string;
};
declare const _default: (config: SFTPConfig) => Promise<{
    type: string;
    user: string;
    password: string;
    address: string;
}>;
export default _default;
