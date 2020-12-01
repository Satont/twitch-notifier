import {MigrationInterface, QueryRunner} from "typeorm";

export class uniqueFollow1606798311968 implements MigrationInterface {
    name = 'uniqueFollow1606798311968'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef" UNIQUE ("chatId", "channelId")
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef"
        `);
    }

}
