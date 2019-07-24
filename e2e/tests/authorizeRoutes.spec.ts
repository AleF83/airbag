import axios from "axios";
import { expect } from "chai";

import { getToken } from "./utils/token";

describe("authenticated routes", () => {

    let accessToken: string;
    let expectedData: any;

    before(async () => {
        expectedData = { data: "Welcome to secured route!" };

        axios.interceptors.response.use((response) => response, (error) => error.response);
    });

    beforeEach(async () => {
        accessToken = await getToken();
    });

    it("should return result", async () => {
        const url = `${process.env.SERVICE_URL}/authenticated`;
        const options = {
            headers: {
                Authorization: `Bearer ${accessToken}`,
            },
        };

        const result = await axios.get(url, options);

        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        expect(result.status).to.equal(200);
        // tslint:disable-next-line:no-unused-expression
        expect(result.data).to.exist;
        expect(result.data).to.deep.equal(expectedData);
    });

    it("should return status Unauthorized", async () => {
        const url = `${process.env.SERVICE_URL}/authenticated`;
        const result = await axios.get(url);

        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        expect(result.status).to.equal(401);
    });
});
