import {MigrationInterface, QueryRunner} from "typeorm";

export class addStreams1629234153096 implements MigrationInterface {
    name = 'addStreams1629234153096'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            CREATE TABLE "streams" (
                "id" character varying NOT NULL,
                "startedAt" TIMESTAMP NOT NULL,
                "updatedAt" TIMESTAMP NOT NULL DEFAULT now(),
                "channelId" character varying,
                CONSTRAINT "PK_40440b6f569ebc02bc71c25c499" PRIMARY KEY ("id")
            )
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."channels" DROP COLUMN "latestStreamId"
        `);
        await queryRunner.query(`
            ALTER TABLE "streams"
            ADD CONSTRAINT "FK_9acd36eceb50ba13cb7027a4b22" FOREIGN KEY ("channelId") REFERENCES "channels"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "streams" DROP CONSTRAINT "FK_9acd36eceb50ba13cb7027a4b22"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."channels"
            ADD "latestStreamId" character varying
        `);
        await queryRunner.query(`
            DROP TABLE "streams"
        `);
    }

}
