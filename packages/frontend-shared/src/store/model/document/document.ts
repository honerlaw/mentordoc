import {Entity} from "../entity";
import {Expose, Exclude, Type} from "class-transformer";
import {DocumentDraft} from "./document-draft";

@Exclude()
export class Document extends Entity {

    @Expose()
    public organizationId: string;

    @Expose()
    public folderId: string;

    @Expose()
    @Type(() => DocumentDraft)
    public drafts: DocumentDraft[];

}
