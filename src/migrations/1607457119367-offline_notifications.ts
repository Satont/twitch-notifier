import {MigrationInterface, QueryRunner} from "typeorm";

export class offlineNotifications1607457119367 implements MigrationInterface {
    name = 'offlineNotifications1607457119367'

    public async up(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "chats_settings"
            ADD "offline_notification" boolean NOT NULL DEFAULT false
        `);
    }

    public async down(queryRunner: QueryRunner): Promise<void> {
        await queryRunner.query(`
            ALTER TABLE "chats_settings" DROP COLUMN "offline_notification"
        `);
    }

}
