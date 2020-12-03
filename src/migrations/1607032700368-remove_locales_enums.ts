import {MigrationInterface, QueryRunner} from "typeorm";

export class removeLocalesEnums1607032700368 implements MigrationInterface {
    name = 'removeLocalesEnums1607032700368'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            COMMENT ON COLUMN "chats_settings"."language" IS NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "chats_settings"
            ALTER COLUMN "language"
            SET DEFAULT 'en'
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "chats_settings"
            ALTER COLUMN "language"
            SET DEFAULT 'english'
        `);
        await queryRunner.query(`
            COMMENT ON COLUMN "chats_settings"."language" IS NULL
        `);
    }

}
