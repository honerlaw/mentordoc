import {Entity} from "../entity";
import {Expose, Exclude} from "class-transformer";

@Exclude()
export class DocumentContent extends Entity {

    @Expose()
    public documentId: string;

    @Expose()
    public content: string;

}