import axios from "axios";
import { expect } from "chai";

describe("metrics route", () => {
    it("should return metrics data", async () => {
        const result = await axios.get(`${process.env.SERVICE_URL}/metrics`);
        // tslint:disable-next-line:no-unused-expression
        expect(result).to.exist;
        // tslint:disable-next-line:no-unused-expression
        expect(result.data).to.exist;
    });
});
