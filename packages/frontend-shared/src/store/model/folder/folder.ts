import {Entity} from "../entity";
import {Exclude, Expose} from "class-transformer";

@Exclude()
export class Folder extends Entity {

    @Expose()
    public name: string;

    @Expose()
    public organizationId: string;

    @Expose()
    public parentFolderId: string | null;

    @Expose()
    public childCount: number;

}