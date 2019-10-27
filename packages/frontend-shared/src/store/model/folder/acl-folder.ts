import {Folder} from "./folder";
import {Type, Exclude, Expose} from "class-transformer";

export function isAclFolder(data: any): data is AclFolder {
    return data.model.organizationId && typeof data.model.childCount === "number";
}

@Exclude()
export class AclFolder {

    @Expose()
    @Type(() => Folder)
    public model: Folder;

    @Expose()
    public actions: string[];

    public hasAction(action: string): boolean {
        return this.actions.indexOf(action) !== -1;
    }

}