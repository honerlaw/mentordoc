import {Entity} from "../entity";
import {Expose, Exclude, Type} from "class-transformer";
import {DocumentDraft} from "./document-draft";

export function isDocument(data: any): data is Document {
    return data.organizationId && data.folderId;
}

@Exclude()
export class Document extends Entity {

    @Expose()
    public organizationId: string;

    @Expose()
    public folderId: string | null;

    @Expose()
    @Type(() => DocumentDraft)
    public drafts: DocumentDraft[];

}
