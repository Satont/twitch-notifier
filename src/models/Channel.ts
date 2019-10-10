import { Table, Column, Model, Unique, PrimaryKey, AllowNull, DataType } from 'sequelize-typescript';
 
@Table({
  tableName: 'channels',
  timestamps: false
})
export class Channel extends Model<Channel> {
  @AllowNull(false)
  @Unique
  @PrimaryKey
  @Column
  public id: number;
 
  @Column
  public online: boolean;
}