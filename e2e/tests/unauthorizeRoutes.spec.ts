import axios from "axios";
import { expect } from "chai";

describe("unauthenticated routes", () => {
    let expectedData: any;

    before(async () => {
        expectedData = { data: "Welcome to public route!" };
    });

    it("should return result", async () => {
        const url = `${process.env.SERVICE_URL}/unauthenticated`;
        const result = await axios.get(url);

        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        // tslint:disable-next-line:no-unused-expression
        expect(result.data).to.exist;
        expect(result.data).to.deep.equal(expectedData);
    });
});
