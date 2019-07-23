import axios from "axios";
import * as querystring from "querystring";

export const getToken = async (): Promise<string> => {

    const { CLIENT_CREDENTIALS_CLIENT_ID, CLIENT_CREDENTIALS_CLIENT_SECRET, API_RESOURCE } = process.env;
    const params = {
        client_id: CLIENT_CREDENTIALS_CLIENT_ID,
        client_secret: CLIENT_CREDENTIALS_CLIENT_SECRET,
        grant_type: "client_credentials",
        scope: API_RESOURCE,
    };

    const response = await axios.post(process.env.TOKEN_ENDPOINT, querystring.stringify(params));
    const accessToken = response.data.access_token;
    return accessToken as string;
};
