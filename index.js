const { spawn } = require("child_process")

function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min)
}
/**
 * @param {import("./types").SFTPConfig} config 
 * @returns {Promise<{type: string,user: string,password: string,address: string}>}
 */
module.exports.default = (config) => {
    return new Promise((resolve, reject) => {
        var file = `sftp-${process.platform}-${process.arch}${process.platform == "win32" ? ".exe" : ""}`
        if (existsSync("./bin/" + file)) {
            var listenAddr = config.listenAddress ? config.listenAddress : "0.0.0.0:" + randomInt(3000, 4000);
            var s = spawn(`./bin/${file} ${config.host} "${config.serverCert}" "${config.clientCert}" "${config.clientKey}" ${config.instanceName} ${listenAddr}`, { shell: process.platform == "win32" ? "powershell" : "bash" })
            s.stdout.on("data", (data) => {
                try {
                    if (data.toString().startsWith("{")) {
                        var msg = JSON.parse(data.toString())
                        if (msg.type == "auth") {
                            resolve(msg)
                        } else if (msg.type == "error") {
                            reject(new Error(msg))
                        } else if (msg.type == "error-withclient") {
                            reject(new Error(msg))
                        }
                    }
                } catch (error) {
                    s.kill()
                }
            })
            s.stderr.on("data", (data) => {
                reject(new Error(data.toString()))
            })
            return s;
        } else {
            reject(new Error(`Platform ${process.platform}-${process.arch} not supported`))
        }
    })
}
