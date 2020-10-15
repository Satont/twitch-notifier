import {
  Table, Column, Model, PrimaryKey, AllowNull, DataType, IsIn, Default
} from 'sequelize-typescript';

@Table({
  tableName: 'users',
  timestamps: false
})
export class User extends Model<User> {
  @AllowNull(false)
  @PrimaryKey
  @Column
  public id: number;

  @Column(DataType.ARRAY(DataType.INTEGER))
  public follows: number[];

  @IsIn({ args: [['vk', 'telegram']], msg: 'Service must be vk or telegram' })
  @Default('vk')
  @Column
  public service: string;

  @Default(false)
  @Column(DataType.BOOLEAN)
  public follow_game_change: boolean;
}
