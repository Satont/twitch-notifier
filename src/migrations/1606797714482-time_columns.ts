import {MigrationInterface, QueryRunner} from "typeorm";

export class timeColumns1606797714482 implements MigrationInterface {
    name = 'timeColumns1606797714482'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "chats"
            ADD "createdAt" TIMESTAMP NOT NULL DEFAULT now()
        `);
        await queryRunner.query(`
            ALTER TABLE "chats"
            ADD "updatedAt" TIMESTAMP NOT NULL DEFAULT now()
        `);
        await queryRunner.query(`
            ALTER TABLE "follows"
            ADD "createdAt" TIMESTAMP NOT NULL DEFAULT now()
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "follows" DROP COLUMN "createdAt"
        `);
        await queryRunner.query(`
            ALTER TABLE "chats" DROP COLUMN "updatedAt"
        `);
        await queryRunner.query(`
            ALTER TABLE "chats" DROP COLUMN "createdAt"
        `);
    }

}
