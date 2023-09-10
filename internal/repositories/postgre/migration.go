package postgre

const migrationQuery = `
create table if not exists metrics
(
    id    bigserial not null,
    name  varchar   not null,
    type  varchar   not null,
    delta bigint,
    value double precision,
	PRIMARY KEY (name, type)
);

comment on column metrics.name is 'Имя метрики';

comment on column metrics.type is 'Тип метрики';

comment on column metrics.delta is 'Значение метрики сounter';

comment on column metrics.value is 'Значение метрики gauge';
`
