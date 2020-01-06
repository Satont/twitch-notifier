import { Table, Column, Model, Unique, PrimaryKey, AllowNull, HasMany } from 'sequelize-typescript';
import Follow from './Follow'

@Table({
  tableName: 'channels',
  timestamps: false
})
export default class Channel extends Model<Channel> {
  @AllowNull(false)
  @Unique
  @PrimaryKey
  @Column
  public id: number;
 
  @Column
  public username: string;

  @Column
  public online: boolean;

  @HasMany(() => Follow)
  public follows: Follow[]
}