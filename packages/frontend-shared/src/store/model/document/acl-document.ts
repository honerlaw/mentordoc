import {Expose, Type, Exclude} from "class-transformer";
import {Document} from "./document";

export function isAclDocument(data: any): data is AclDocument {
    return data.model.organizationId && Array.isArray(data.model.drafts);
}

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
