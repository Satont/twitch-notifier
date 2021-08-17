import {MigrationInterface, QueryRunner} from "typeorm";

export class storeChannelAndTitleInStreams1629237450026 implements MigrationInterface {
    name = 'storeChannelAndTitleInStreams1629237450026'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "public"."channels" DROP COLUMN "category"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."channels" DROP COLUMN "title"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."streams"
            ADD "category" character varying
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."streams"
            ADD "title" character varying
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "public"."streams" DROP COLUMN "title"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."streams" DROP COLUMN "category"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."channels"
            ADD "title" character varying
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."channels"
            ADD "category" character varying
        `);
    }

}
