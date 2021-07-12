import {MigrationInterface, QueryRunner} from "typeorm";

export class latestStreamIdColumn1626048220050 implements MigrationInterface {
    name = 'latestStreamIdColumn1626048220050'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "channels"
            ADD "latestStreamId" character varying
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "channels" DROP COLUMN "latestStreamId"
        `);
    }

}
