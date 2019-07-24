import { before } from "mocha";
import { URL } from "url";
import waitOn = require("wait-on");

before("waiting for oidc server & airbag", async function() {
    const timeout = Number(process.env.TIMEOUT);
    this.timeout(timeout);

    const { hostname: tknHostname, port: tknPort, protocol: tknProtocol } = new URL(process.env.TOKEN_ENDPOINT);
    await waitFor(tknHostname, tknPort, tknProtocol, timeout);

    const { hostname: svcHostName, port: svcPort, protocol: svcProtocol } = new URL(process.env.SERVICE_URL);
    await waitFor(svcHostName, svcPort, svcProtocol, timeout);
});

const getImplicitPort = (protocol) => {
    switch (protocol) {
        case "http:":
            return 80;
        case "https:":
            return 443;
        default:
    }
};

const waitFor = async (hostname: string, port: string, protocol: string, timeout: number): Promise<void> => {
    const explicitPort = port || getImplicitPort(protocol);
    await waitOn({
        delay: 1000,
        interval: 100,
        log: true,
        resources: [
            `tcp:${hostname}:${explicitPort}`,
        ],
        timeout,
        verbose: true,
        window: 1000,
    });
};
