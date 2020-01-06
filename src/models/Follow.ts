import { Table, Column, Model, PrimaryKey, AllowNull, ForeignKey } from 'sequelize-typescript';
import Channel from './Channel';
 
@Table({
  tableName: 'follows',
  timestamps: false
})
export default class Follow extends Model<Follow> {
  @AllowNull(false)
  @PrimaryKey
  @Column
  public user_id: number;

  @ForeignKey(() => Channel)
  @Column
  public channel_id: number;
}