import {MigrationInterface, QueryRunner} from "typeorm";

export class noNullableFollows1618419645268 implements MigrationInterface {
    name = 'noNullableFollows1618419645268'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "FK_dc1c79f16fa7bf1ef500a0305fb"
        `);
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "FK_af75a484b4316e01b2248a7d244"
        `);
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef"
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ALTER COLUMN "chatId"
            SET NOT NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ALTER COLUMN "channelId"
            SET NOT NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef" UNIQUE ("chatId", "channelId")
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "FK_dc1c79f16fa7bf1ef500a0305fb" FOREIGN KEY ("chatId") REFERENCES "chats"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "FK_af75a484b4316e01b2248a7d244" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "FK_af75a484b4316e01b2248a7d244"
        `);
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "FK_dc1c79f16fa7bf1ef500a0305fb"
        `);
        await queryRunner.query(`
            ALTER TABLE "follows" DROP CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef"
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ALTER COLUMN "channelId" DROP NOT NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ALTER COLUMN "chatId" DROP NOT NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef" UNIQUE ("chatId", "channelId")
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "FK_af75a484b4316e01b2248a7d244" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD CONSTRAINT "FK_dc1c79f16fa7bf1ef500a0305fb" FOREIGN KEY ("chatId") REFERENCES "chats"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
    }

}
