create table if not exists channels (
    id varchar(255) not null primary key
);

create type chat_service as enum ('telegram');

create table if not exists chats (
    id varchar(255) not null primary key,
    chat_id varchar(255) not null,
    service chat_service not null,
    settings_id varchar(255) not null,
    foreign key (settings_id) references chat_settings(id)
);

create table if not exists chat_settings (
    id varchar(255) not null primary key,
    game_change_notification boolean not null,
    offline_notification boolean not null,
    language varchar(255) not null
);

create table if not exists follows (
    id int not null primary key,
    chat_id varchar(255) not null,
    channel_id varchar(255) not null,
    foreign key (chat_id) references chats(id),
    foreign key (channel_id) references channels(id)
);
