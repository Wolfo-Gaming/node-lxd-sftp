export type SFTPConfig = {
 serverCert: Buffer,
 clientCert: Buffer,
 clientKey: Buffer,
 host: string,
 instanceName: string,
 listenAddress?: string
}
declare module "node-lxd-sftp" {
    function sftp(config: SFTPConfig): Promise<{type: string,user: string,password: string,address: string}>;
    export default sftp;
}