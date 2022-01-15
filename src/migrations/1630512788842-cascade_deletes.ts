import {MigrationInterface, QueryRunner} from "typeorm";

export class cascadeDeletes1630512788842 implements MigrationInterface {
    name = 'cascadeDeletes1630512788842'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "public"."chats" DROP CONSTRAINT "FK_db425de4c28fc6c50cc9c919c20"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."chats"
            ADD CONSTRAINT "FK_db425de4c28fc6c50cc9c919c20" FOREIGN KEY ("settingsId") REFERENCES "chats_settings"("id") ON DELETE CASCADE ON UPDATE NO ACTION
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "public"."chats" DROP CONSTRAINT "FK_db425de4c28fc6c50cc9c919c20"
        `);
        await queryRunner.query(`
            ALTER TABLE "public"."chats"
            ADD CONSTRAINT "FK_db425de4c28fc6c50cc9c919c20" FOREIGN KEY ("settingsId") REFERENCES "chats_settings"("id") ON DELETE NO ACTION ON UPDATE NO ACTION
        `);
    }

}
