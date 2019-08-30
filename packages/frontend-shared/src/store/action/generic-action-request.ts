import {Component} from "react";

export interface IGenericActionRequest {
    options?: {
        alerts: {
            target: Component | string;
        }
    }
}