import {MigrationInterface, QueryRunner} from "typeorm";

export class initial1607025517426 implements MigrationInterface {
    name = 'initial1607025517426'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            CREATE TABLE "chats_settings" (
                "id" SERIAL NOT NULL,
                "game_change_notification" boolean NOT NULL DEFAULT false,
                "language" character varying NOT NULL DEFAULT 'english',
                CONSTRAINT "PK_cafeab5c755e94c547405cea61d" PRIMARY KEY ("id")
            )
        `);
        await queryRunner.query(`
            CREATE TABLE "chats" (
                "id" SERIAL NOT NULL,
                "chatId" character varying NOT NULL,
                "service" character varying NOT NULL,
                "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
                "updatedAt" TIMESTAMP NOT NULL DEFAULT now(),
                "settingsId" integer,
                CONSTRAINT "UQ_5b05db8be015cde8647963d41a3" UNIQUE ("chatId", "service"),
                CONSTRAINT "REL_db425de4c28fc6c50cc9c919c2" UNIQUE ("settingsId"),
                CONSTRAINT "PK_0117647b3c4a4e5ff198aeb6206" PRIMARY KEY ("id")
            )
        `);
        await queryRunner.query(`
            CREATE TABLE "follows" (
                "id" SERIAL NOT NULL,
                "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
                "chatId" integer,
                "channelId" character varying,
                CONSTRAINT "UQ_d63caec7a6eee9b38484e88dfef" UNIQUE ("chatId", "channelId"),
                CONSTRAINT "PK_8988f607744e16ff79da3b8a627" PRIMARY KEY ("id")
            )
        `);
        await queryRunner.query(`
            CREATE TABLE "channels" (
                "id" character varying NOT NULL,
                "username" character varying NOT NULL,
                "online" boolean NOT NULL,
                "category" character varying,
                "title" character varying,
                "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
                "updatedAt" TIMESTAMP NOT NULL DEFAULT now(),
                CONSTRAINT "PK_bc603823f3f741359c2339389f9" PRIMARY KEY ("id")
            )
        `);
        await queryRunner.query(`
            ALTER TABLE "chats"
            ADD CONSTRAINT "FK_db425de4c28fc6c50cc9c919c20" FOREIGN KEY ("settingsId") REFERENCES "chats_settings"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
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
            ALTER TABLE "chats" DROP CONSTRAINT "FK_db425de4c28fc6c50cc9c919c20"
        `);
        await queryRunner.query(`
            DROP TABLE "channels"
        `);
        await queryRunner.query(`
            DROP TABLE "follows"
        `);
        await queryRunner.query(`
            DROP TABLE "chats"
        `);
        await queryRunner.query(`
            DROP TABLE "chats_settings"
        `);
    }

}
