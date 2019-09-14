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
    public publishedAt: number;

    @Expose()
    public retractedAt: number;

    @Expose()
    @Type(() => DocumentContent)
    public content: DocumentContent | null;

}
