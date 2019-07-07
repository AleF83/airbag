import axios from "axios";
import { expect } from "chai";
import * as fs from "fs";

import { getToken } from "./utils/token";

describe("authenticated routes", () => {

    let accessToken: string;
    let expectedData: any;

    before(async () => {
        const fileContent = await fs.promises.readFile("../data/server.json", { encoding: "utf8" });
        expectedData = JSON.parse(fileContent).authenticated;

        axios.interceptors.response.use((response) => response, (error) => error.response);
    });

    beforeEach(async () => {
        accessToken = await getToken();
    });

    it("should return result", async () => {
        const result = await axios.get(`${process.env.AIRBAG_URL}/authenticated`, {
            headers: {
                Authorization: `Bearer ${accessToken}`,
            },
        });

        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        expect(result.status).to.equal(200);
        // tslint:disable-next-line:no-unused-expression
        expect(result.data).to.exist;
        expect(result.data).to.deep.equal(expectedData);
    });

    it("should return status Unauthorized", async () => {
        const result = await axios.get(`${process.env.AIRBAG_URL}/authenticated`);

        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        expect(result.status).to.equal(401);
    });
});
