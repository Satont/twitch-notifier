import {MigrationInterface, QueryRunner} from "typeorm";

export class nullableGame1606798115819 implements MigrationInterface {
    name = 'nullableGame1606798115819'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "channels"
            ALTER COLUMN "game" DROP NOT NULL
        `);
        await queryRunner.query(`
            COMMENT ON COLUMN "channels"."game" IS NULL
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            COMMENT ON COLUMN "channels"."game" IS NULL
        `);
        await queryRunner.query(`
            ALTER TABLE "channels"
            ALTER COLUMN "game"
            SET NOT NULL
        `);
    }

}
