import {Exclude, Expose} from "class-transformer";

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

    public getKey(): string {
        return `${this.type}-${this.message}`;
    }

}
