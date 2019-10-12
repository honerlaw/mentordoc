import {Entity} from "../entity";
import {DocumentContent} from "./document-content";
import {Exclude, Expose, Type} from "class-transformer";

@Exclude()
export class DocumentDraft extends Entity {

    @Expose()
    public name: string;

    @Expose()
    public documentId: string;

    @Expose()
    public publishedAt: number | null;

    @Expose()
    public retractedAt: number | null;

    @Expose()
    @Type(() => DocumentContent)
    public content: DocumentContent | null;

}
