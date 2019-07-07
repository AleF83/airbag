import axios from "axios";
import { expect } from "chai";
import * as fs from "fs";

describe("unauthenticated routes", () => {
    let expectedData: any;

    before(async () => {
        const fileContent = await fs.promises.readFile("../data/server.json", { encoding: "utf8" });
        expectedData = JSON.parse(fileContent).unauthenticated;
    });

    it("should return result", async () => {
        const result = await axios.get(`${process.env.AIRBAG_URL}/unauthenticated`);

        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        // tslint:disable-next-line:no-unused-expression
        expect(result.data).to.exist;
        expect(result.data).to.deep.equal(expectedData);
    });
});
