import {Entity} from "../entity";
import {DocumentContent} from "./document-content";
import {Expose, Exclude, Type} from "class-transformer";

@Exclude()
export class Document extends Entity {

    @Expose()
    public name: string;

    @Expose()
    public organizationId: string;

    @Expose()
    public folderId: string;

    @Expose()
    @Type(() => DocumentContent)
    public content: DocumentContent | null;

}
