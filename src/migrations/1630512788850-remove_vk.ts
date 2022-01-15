import {MigrationInterface, QueryRunner} from "typeorm";
import { Chat } from '../entities/Chat'

export class removeVk1630511837166 implements MigrationInterface {

    public async up(queryRunner: QueryRunner): Promise<void> {
        const repository = queryRunner.manager.getRepository<Chat>('chats')
        const chats = await repository.find({ relations: ['follows', 'settings'] })

        for (const chat of chats.filter(c => c.service === 'vk' as any)) {
            for (const follow of chat.follows) {
                await follow.remove()
            }
            await chat.remove()
        }
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
    }

}
