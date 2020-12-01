import {MigrationInterface, QueryRunner} from "typeorm";

export class initial1606822190197 implements MigrationInterface {
    name = 'initial1606822190197'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "chats"
            ADD "service" character varying NOT NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "chats"
            ADD CONSTRAINT "UQ_a8783fe5751bf0744cfc23348ab" UNIQUE ("id", "service")
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "chats" DROP CONSTRAINT "UQ_a8783fe5751bf0744cfc23348ab"
        `);
        await queryRunner.query(`
            ALTER TABLE "chats" DROP COLUMN "service"
        `);
    }

}
