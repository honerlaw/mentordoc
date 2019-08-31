import {Expose, Exclude} from "class-transformer";
import {Entity} from "../entity";

@Exclude()
export class Organization extends Entity {

    @Expose()
    public name: string;

}