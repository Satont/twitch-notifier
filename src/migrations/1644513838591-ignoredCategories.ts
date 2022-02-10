import {MigrationInterface, QueryRunner} from "typeorm";

export class ignoredCategories1644513838591 implements MigrationInterface {
    name = 'ignoredCategories1644513838591'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            CREATE TABLE "ignored_categories" (
                "id" SERIAL NOT NULL,
                "createdAt" TIMESTAMP NOT NULL DEFAULT now(),
                "categoryId" character varying NOT NULL,
                "chatId" integer NOT NULL,
                CONSTRAINT "UQ_387f7132adcd480068fcb932f3a" UNIQUE ("chatId", "categoryId"),
                CONSTRAINT "PK_bb39c37688d63a593855aeb4930" PRIMARY KEY ("id")
            )
        `);
        await queryRunner.query(`
            ALTER TABLE "ignored_categories"
            ADD CONSTRAINT "FK_6e6f111fde2a848d59354bcaed5" FOREIGN KEY ("chatId") REFERENCES "chats"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "ignored_categories" DROP CONSTRAINT "FK_6e6f111fde2a848d59354bcaed5"
        `);
        await queryRunner.query(`
            DROP TABLE "ignored_categories"
        `);
    }

}
