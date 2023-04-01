export type SFTPConfig = {
 serverCert: Buffer,
 clientCert: Buffer,
 clientKey: Buffer,
 host: string,
 instanceName: string,
 listenAddress?: string
}