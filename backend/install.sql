\set
ON_ERROR_STOP 1

DROP
DATABASE IF EXISTS mdb;
DROP
user IF EXISTS mdb;
CREATE
DATABASE mdb;
CREATE
user mdb WITH PASSWORD 'mdb';

\connect
mdb

CREATE SCHEMA mdb;
GRANT
usage
ON
SCHEMA
mdb TO mdb;


CREATE TABLE mdb.users
(
    login             VARCHAR(100) PRIMARY KEY,
    password          VARCHAR(256) NOT NULL,
    img_src           text                  DEFAULT 'http://localhost:8085/avatars/default.jpeg',

    is_admin          integer INTEGER DEFAULT 0,

    firstname         VARCHAR(100),
    lastname          VARCHAR(100),
    sex               INTEGER
        CONSTRAINT sex_t CHECK (sex = 1 OR sex = 0),
    email             VARCHAR(100) NOT NULL UNIQUE,
    registration_date timestamp    NOT NULL DEFAULT NOW(),
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.users TO mdb;
COMMENT
ON TABLE mdb.users IS 'Пользователи';

drop table mdb.jobs;
create table if not exists mdb.jobs
(
    id             bigserial primary key,
    uni_id         text not null,
    status_id      smallint                          default 0,
    queue          text,
    action         text                     not null,
    entity_type    text                     not null,
    payload        text,
    result         text,
    error          text,
    attempts       integer                  not null default 0,
    reserved_until timestamp with time zone          default now() - '1 second'::interval,
    create_txtime  timestamp with time zone not null default now()
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.jobs TO mdb;


drop table mdb.subdivision_org;
create table if not exists mdb.subdivision_org
(
    id   bigserial primary key,
    uid  text not null,
    name text not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.subdivision_org TO mdb;


drop table mdb.education_program;
create table if not exists mdb.education_program
(
    id                bigserial primary key,
    uid               text    not null,
    name              text    not null,
    id_education_form integer not null,
    idocso            integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.education_program TO mdb;


drop table mdb.campaign;
create table if not exists mdb.campaign
(
    id                 bigserial primary key,
    uid                text    not null,
    name               text    not null,
    year_start         text    not null,
    year_end           text    not null,
    education_forms    integer not null,
    education_levels   integer not null,
    id_campaign_type   integer not null,
    id_campaign_status integer not null,
    number_agree       integer not null,
    end_date           text    not null,
    count_directions   integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.campaign TO mdb;


drop table mdb.cmp_achievement;
create table if not exists mdb.cmp_achievement
(
    id           bigserial primary key,
    uid_campaign text    not null,
    uid          text    not null,
    id_category  integer not null,
    name         text    not null,
    max_value    integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.cmp_achievement TO mdb;


drop table mdb.admission_volume;
create table if not exists mdb.admission_volume
(
    id                 bigserial primary key,
    uid                text    not null,
    uid_campaign       text    not null,
    id_direction       integer not null,
    id_education_level integer not null,
    budget_o           integer not null,
    budget_oz          integer not null,
    budget_z           integer not null,
    quota_o            integer not null,
    quota_oz           integer not null,
    quota_z            integer not null,
    paid_o             integer not null,
    paid_oz            integer not null,
    paid_z             integer not null,
    target_o           integer not null,
    target_oz          integer not null,
    target_z           integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.admission_volume TO mdb;

drop table mdb.distibuted_admission_volume;
create table if not exists mdb.distibuted_admission_volume
(
    id                   bigserial primary key,
    uid                  text    not null,
    uid_admission_volume text    not null,
    id_direction         integer not null,
    id_level_budget      integer not null,
    budget_o             integer not null,
    budget_oz            integer not null,
    budget_z             integer not null,
    quota_o              integer not null,
    quota_oz             integer not null,
    quota_z              integer not null,
    paid_o               integer not null,
    paid_oz              integer not null,
    paid_z               integer not null,
    target_o             integer not null,
    target_oz            integer not null,
    target_z             integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.distibuted_admission_volume TO mdb;

drop table mdb.competitive_group;
create table if not exists mdb.competitive_group
(
    id                  bigserial primary key,
    uid                 text    not null,
    uid_campaign        text    not null,
    name                text    not null,
    id_level_budget     integer not null,
    id_education_level  integer not null,
    id_education_source integer not null,
    id_education_form   integer not null,
    admission_number    integer not null,
    comment             text    not null,
    idocso              integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.competitive_group TO mdb;


drop table mdb.competitive_group_program;
create table if not exists mdb.competitive_group_program
(
    id                    bigserial primary key,
    uid_competitive_group text not null,
    uid_subdivision_org   text not null,
    uid                   text not null,
    uid_education_program text not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.competitive_group_program TO mdb;


drop table mdb.competitive_benefit;
create table if not exists mdb.competitive_benefit
(
    id                      bigserial primary key,
    uid_competitive_group   text    not null,
    id_olimpic_type         integer not null,
    id_olimpic_levels       integer not null,
    id_benefit              integer not null,
    id_olimpic_diploma_type integer not null,
    ege_min_value           integer not null,
    olimpic_profile         integer not null,
    uid                     text    not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.competitive_benefit TO mdb;


drop table mdb.entrance_test;
create table if not exists mdb.entrance_test
(
    id                        bigserial primary key,
    uid                       text    not null,
    uid_competitive_group     text    not null,
    id_entrance_test_type     integer not null,
    test_name                 text    not null,
    is_ege                    bool    not null,
    min_score                 integer not null,
    priority                  integer not null,
    id_subject                integer not null,
    uid_replace_entrance_test text    not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.entrance_test TO mdb;


drop table mdb.entrance_test_benefit;
create table if not exists mdb.entrance_test_benefit
(
    id                      bigserial primary key,
    uid                     text    not null,
    uid_entrance_test       text    not null,
    id_benefit              integer not null,
    id_olimpic_diploma_type integer not null,
    id_olympic_classes      integer not null,
    id_olympic_level        integer not null,
    id_olympic_profiles     integer not null,
    ege_min_value           integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.entrance_test_benefit TO mdb;


drop table mdb.entrance_test_location;
create table if not exists mdb.entrance_test_location
(
    id                bigserial primary key,
    id_choice         text    not null,
    uid_choice_epgu   text    not null,
    uid_entrance_test text    not null,
    test_date         text    not null,
    test_location     text    not null,
    entrance_count    integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.entrance_test_location TO mdb;


drop table mdb.service_entrant_photo_file;
create table if not exists mdb.service_entrant_photo_file
(
    id        bigserial primary key,
    uid_oovo  text not null,
    uid_epgu  text not null,
    snils     text not null,
    file_name text not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.service_entrant_photo_file TO mdb;


drop table mdb.identification;
create table if not exists mdb.identification
(
    id        bigserial primary key,
    uid_oovo  text not null,
    uid_epgu  text not null,
    snils     text not null,
    file_name text not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.identification TO mdb;


drop table mdb.document;
create table if not exists mdb.document
(
    id        bigserial primary key,
    uid_oovo  text not null,
    uid_epgu  text not null,
    snils     text not null,
    file_name text not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.document TO mdb;


drop table mdb.service_application;
create table if not exists mdb.service_application
(
    id                                bigserial primary key,
    guid                              text    not null,
    snils                             text    not null,
    passport                          text    not null,
    name                              text    not null,
    surname                           text    not null,
    patronomic                        text    not null,
    series                            text    not null,
    number                            text    not null,
    birthday                          text    not null,
    uid_compititive_group             text    not null,
    uid_compititive_group_target      text    not null,
    app_number                        integer not null,
    registration_date                 text    not null,
    rating                            integer not null,
    priority                          integer not null,
    first_higher_education            bool    not null,
    need_hostel                       bool    not null,
    id_disabeled_document_choice_epgu integer not null,
    id_disabeled_document_choice_oovo integer not null,
    special_conditions                bool    not null,
    agree                             bool    not null,
    agree_date                        text    not null,
    id_return_type                    integer not null,
    return_date                       text    not null,
    orig_doc                          bool    not null,
    id_benefit                        integer not null,
    single_statement                  bool    not null,
    id_aplication_choice_oovo         integer not null,
    id_aplication_choice_epgu         integer not null,
    status_coment                     text    not null,
    id_document                       integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.service_application TO mdb;


drop table mdb.app_achiviment;
create table if not exists mdb.app_achiviment
(
    id                 bigserial primary key,
    id_status          integer not null,
    id_status_action   integer not null,
    id_doc_choice_epgu integer not null,
    id_doc_choice_oovo integer not null,
    comment            text    not null,
    agree              bool    not null,
    agree_date         text    not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.app_achiviment TO mdb;


drop table mdb.terms_admission;
create table if not exists mdb.terms_admission
(
    id                   bigserial primary key,
    uid                  text    not null,
    uid_campaign         text    not null,
    terms_admission_name text    not null,
    id_terms_lfs         integer not null,
    start_date           text    not null,
    end_date             text    not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.terms_admission TO mdb;


drop table mdb.entrance_test_sheet;
create table if not exists mdb.entrance_test_sheet
(
    id                bigserial primary key,
    uid               text    not null,
    uid_entrance_test text    not null,
    name              text    not null,
    id_education_form integer not null,
    doc_date          text    not null,
    file_name         text    not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.entrance_test_sheet TO mdb;


drop table mdb.competitive_group_speciality;
create table if not exists mdb.competitive_group_speciality
(
    id                    bigserial primary key,
    uid                   text    not null,
    idocso                integer not null,
    uid_competitive_group text    not null,
    admission_number      integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.competitive_group_speciality TO mdb;


drop table mdb.entarance_test_agreed_list;
create table if not exists mdb.entarance_test_agreed_list
(
    id                                         bigserial primary key,
    uid_entarance_test_agreed_list_choice_oovo text not null,
    uid_entarance_test_agreed_list_choice_epgu text not null,
    id_aplication_choice_oovo                  text not null,
    id_aplication_choice_epgu                  text not null,
    uid_entarance_test_location_oovo           text not null,
    uid_entarance_test                         text not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.entarance_test_agreed_list TO mdb;


drop table mdb.entrance_test_list_result;
create table if not exists mdb.entrance_test_list_result
(
    id                        bigserial primary key,
    id_aplication_choice_oovo text    not null,
    id_aplication_choice_epgu text    not null,
    uid                       text    not null,
    uid_entarance_test        text    not null,
    id_doc_choice_epgu        integer not null,
    id_doc_choice_oovo        integer not null,
    issue_date                text    not null,
    id_result_source          integer not null,
    result_value              integer not null,
    id_benefit                integer not null,
    id_apeal_status           integer not null,
    has_ege                   bool    not null,
    ege_value                 integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.entrance_test_list_result TO mdb;


drop table mdb.org_direction;
create table if not exists mdb.org_direction
(
    id           bigserial primary key,
    uid          text    not null,
    id_direction integer not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.org_direction TO mdb;


drop table mdb.service_entrant;
create table if not exists mdb.service_entrant
(
    id         bigserial primary key,
    guid       text    not null,
    snils      text    not null,
    passport   text    not null,
    name       text    not null,
    surname    text    not null,
    patronomic text    not null,
    id_gender  integer not null,
    birthday   text    not null,
    birthplace text    not null,
    phone      text    not null,
    email      text    not null,
    idoccm     integer not null,
    full_adr   text    not null,
    index_adr  text    not null,
    id_region  integer not null,
    area       text    not null,
    city       text    not null,
    city_area  text    not null,
    street     text    not null,
    house      text    not null,
    apartment  text    not null
);
GRANT SELECT, INSERT, UPDATE, DELETE ON mdb.service_entrant TO mdb;


GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA mdb TO mdb;
