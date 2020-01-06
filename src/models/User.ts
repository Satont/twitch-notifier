import { Table, Column, Model, PrimaryKey, AllowNull, IsIn, Default } from 'sequelize-typescript';
 
@Table({
  tableName: 'users',
  timestamps: false
})
export default class User extends Model<User> {
  @AllowNull(false)
  @PrimaryKey
  @Column
  public id: number;
 
  @IsIn({ args: [['vk', 'telegram']], msg: 'Service must be vk or telegram' })
  @Default('vk')
  @Column
  public service: string;
}