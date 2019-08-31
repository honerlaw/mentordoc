import {plainToClass} from "class-transformer";
import {HttpError} from "../store/model/request-status/http-error";
import {IRootState} from "../store/model/root-state";
import {MiddlewareAPI} from "redux";
import {AuthenticationData} from "../store/model/user/authentication-data";
import {SetAuthenticationData} from "../store/action/user/set-authentication-data";
import {Logout} from "../store/action/user/logout";
import {isArray} from "lodash";

export interface IRequestOptions<T> {
    method: "POST" | "GET" | "PUT" | "DELETE";
    path: string,
    model?: new () => any,
    body?: any;
    api?: MiddlewareAPI,
    useRefreshToken?: boolean;
    headers?: Record<string, string>;
}

export async function request<T>(options: IRequestOptions<T>): Promise<T | null> {
    const headers: Record<string, string> = {
        "Content-Type": "application/json",
        ...options.headers
    };

    setAuthorizationHeader(headers, options.api, options.useRefreshToken);

    const init: RequestInit = {
        method: options.method,
        headers: headers
    };

    if (options.body) {
        init.body = JSON.stringify(options.body);
    }

    const resp: Response = await fetch(`http://${process.env.API_HOST}:${process.env.API_PORT}/v1${options.path}`, init);

    if (resp.status < 200 || resp.status >= 300) {

        // 401 on the token, so attempt to use the refresh token to get a new access token
        if (resp.status === 401) {
            const retry: boolean = await handleUnauthorized(options.api, options.useRefreshToken);
            if (retry) {
                return request(options);
            }
        }

        throw await parseResponse<HttpError>(resp, HttpError);
    }

    if (options.model) {
        return await parseResponse<T>(resp, options.model);
    }
    return null;
}

async function parseResponse<T>(resp: Response, model: new () => T): Promise<T> {
    try {
        const json: any = await resp.json();
        if (isArray(json)) {
            const items: unknown = json.map((item: any): any => plainToClass(model, item));

            return items as T;
        }
        return plainToClass(model, json);
    } catch (err) {
        throw new HttpError("something went wrong");
    }
}

function setAuthorizationHeader(headers: Record<string, string>, api?: MiddlewareAPI, useRefreshToken?: boolean): void {
    if (!api) {
        return;
    }

    const state: IRootState = api.getState();
    const data: AuthenticationData | null = state.user.authenticationData;
    if (!data) {
        return;
    }

    if (data.accessToken) {
        headers["Authorization"] = `Bearer ${data.accessToken}`;
    }

    if (useRefreshToken && data.refreshToken) {
        headers["Authorization"] = `Bearer ${data.accessToken}`;
    }
}

/**
 * returns true to denote we need to retry the parent request (e.g. we successfully got a new auth token)
 */
async function handleUnauthorized(api?: MiddlewareAPI, useRefreshToken?: boolean): Promise<boolean> {
    if (!api) {
        return false;
    }

    // so the last request wasn't a refresh token, so we can do the refresh token request
    if (!useRefreshToken) {
        try {
            await refreshToken(api);

            return true;
        } catch (err) {
            // failed, but we don't care, going to logout anyways
        }
    }

    api.dispatch(Logout.action());

    return false;
}

async function refreshToken(api: MiddlewareAPI): Promise<void> {
    const data: AuthenticationData | null = await request({
        method: "POST",
        path: "/user/auth/refresh",
        model: AuthenticationData,
        useRefreshToken: true
    });

    api.dispatch(SetAuthenticationData.action({
        data
    }));
}