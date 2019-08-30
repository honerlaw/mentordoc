import {Exclude, Expose} from "class-transformer";
import {Component} from "react";

export enum AlertType {
    SUCCESS = "success",
    ERROR = "error"
}

@Exclude()
export class Alert {

    @Expose()
    public type: AlertType;

    @Expose()
    public message: string;

    @Expose()
    public lifespan: number | undefined;

    @Expose()
    public target: Component | string;

    public getKey(): string {
        return `${this.type}-${this.message}`;
    }

}
