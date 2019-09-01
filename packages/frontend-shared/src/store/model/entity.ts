import {Exclude, Expose} from "class-transformer";

@Exclude()
export abstract class Entity {

    @Expose()
    public id: string;

    @Expose()
    public updatedAt: number;

    @Expose()
    public createdAt: number;

    @Expose()
    public deletedAt: number | null;

}
