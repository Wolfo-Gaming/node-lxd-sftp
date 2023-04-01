# Node LXD SFTP
![npm](https://badges.aleen42.com/src/npm.svg) [![npm version](https://badge.fury.io/js/node-lxd-sftp.svg)](https://github.com/Wolfo-Gaming/node-lxd-sftp) ![typescript](https://badges.aleen42.com/src/javascript.svg)

Optional dependency for [node-lxd](github.com/Wolfo-Gaming/node-lxd) wich adds SFTP support for better filesystem performance.

# Installing

```bash
$ npm install --save node-lxd-sftp
```

## Getting Started ##

You can include this module in to your own project, but it is made for my other module [node-lxd](github.com/Wolfo-Gaming/node-lxd) as an optional dependency because of its increased size due to precompiled binaries.

```js
var sftp = require("node-lxd-sftp");
sftp({
 serverCert: fs.readFileSync("./server.crt"),
 clientCert: fs.readFileSync("./cert.crt"),
 clientKey: fs.readFileSync("./key.key"),
 host: "https://192.168.2.63:8443/",
 instanceName: "ubuntu2204"
}).then(auth => {
    // this includes the listening address + credentials
    console.log(auth);
});
```