import { Table, Column, Model, Unique, PrimaryKey, AllowNull, DataType } from 'sequelize-typescript';
 
@Table({
  tableName: 'users',
  timestamps: false
})
export class User extends Model<User> {
  @AllowNull(false)
  @Unique
  @PrimaryKey
  @Column
  public id: number;
 
  @Column(DataType.ARRAY(DataType.BIGINT))
  public follows: number[];
}