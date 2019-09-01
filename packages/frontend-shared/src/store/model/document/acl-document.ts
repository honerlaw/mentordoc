import {Expose, Type, Exclude} from "class-transformer";
import {Document} from "./document";

@Exclude()
export class AclDocument {

    @Expose()
    @Type(() => Document)
    public model: Document;

    @Expose()
    public actions: string[];

    public hasAction(action: string): boolean {
        return this.actions.indexOf(action) !== -1;
    }

}
