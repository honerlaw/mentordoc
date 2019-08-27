import {Exclude, Expose} from "class-transformer";

@Exclude()
export class HttpError {

    @Expose()
    public errors: string[];

    public constructor(...errors: string[]) {
        this.errors = errors;
    }

}
