
create table public.events
(
    id             uuid                     default uuidv7() not null primary key,
    name           varchar(255)                              not null,
    max_seats      integer                                   not null,
    waitlist_limit integer                                   not null,
    booked_count   integer                  default 0        not null,
    waitlist_count integer                  default 0        not null,
    price          numeric(10, 2)                            not null,
    start_time     timestamp with time zone                  not null,
    end_time       timestamp with time zone                  not null,
    created_at     timestamp with time zone default now()    not null,
    updated_at     timestamp with time zone default now()
);
create index idx_events_created_at
    on public.events (created_at);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
create index idx_events_name_trgm
    on public.events using gin (name gin_trgm_ops);

create table public.bookings
(
    id         uuid                     default uuidv7() not null primary key,
    event_id   uuid                                      not null
        constraint bookings_events_id_fk
            references public.events
            on delete cascade,
    user_id    uuid                                      not null,
    status     booking_status                            not null,
    created_at timestamp with time zone default now()    not null,
    updated_at timestamp with time zone default now()
);

create unique index uniq_user_event
    on public.bookings (event_id, user_id);

INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f5-7376-a505-61718ec959f5', 'Go Workshop 101', 50, 5, 2, 0, 499.00, '2026-03-01 02:00:00.000000 +00:00', '2028-03-01 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-768d-b1a2-8120244bdebd', 'Go Advanced Concurrency', 40, 10, 10, 2, 799.00, '2026-03-02 03:00:00.000000 +00:00', '2026-03-02 11:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-76c9-8c68-848bd943cdbc', 'Microservices with Go', 60, 8, 25, 3, 899.00, '2026-03-03 02:00:00.000000 +00:00', '2026-03-03 09:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-76f9-9dc8-09b29b08a48d', 'REST API Design', 30, 5, 15, 1, 599.00, '2026-03-04 04:00:00.000000 +00:00', '2026-03-04 12:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7713-8e56-316ccbd9b8f6', 'Clean Architecture Go', 45, 5, 20, 0, 699.00, '2026-03-05 02:00:00.000000 +00:00', '2026-03-05 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-772f-b778-0633ea7d4ce1', 'Go Testing Workshop', 35, 5, 12, 0, 499.00, '2026-03-06 03:00:00.000000 +00:00', '2026-03-06 09:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-774b-840c-f11fa4ab6b1c', 'Docker for Go Devs', 50, 10, 30, 4, 799.00, '2026-03-07 02:00:00.000000 +00:00', '2026-03-07 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7765-851d-6f683aad94f6', 'Kubernetes Basics', 60, 10, 45, 6, 999.00, '2026-03-08 02:00:00.000000 +00:00', '2026-03-08 11:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7cdf-b88a-b9eeb7a8fd0d', 'Event Driven Go', 40, 5, 18, 2, 699.00, '2026-03-09 03:00:00.000000 +00:00', '2026-03-09 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7d3d-9a9a-95c7882f7b1e', 'Kafka Integration', 50, 10, 35, 5, 899.00, '2026-03-10 02:00:00.000000 +00:00', '2026-03-10 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7d5c-9d70-cc9925e8a5da', 'Go Performance Tuning', 30, 5, 10, 1, 599.00, '2026-03-11 02:00:00.000000 +00:00', '2026-03-11 09:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7d78-b05e-8ca49d22a2e6', 'gRPC with Go', 45, 8, 22, 2, 799.00, '2026-03-12 03:00:00.000000 +00:00', '2026-03-12 11:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7d8f-a94c-4a3997a4e39d', 'GraphQL in Go', 35, 5, 14, 1, 699.00, '2026-03-13 02:00:00.000000 +00:00', '2026-03-13 09:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7da7-b8fb-030b14de7f10', 'Authentication & JWT', 50, 10, 28, 3, 799.00, '2026-03-14 02:00:00.000000 +00:00', '2026-03-14 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7e22-ad9c-f404b203fc57', 'CI/CD for Go Projects', 40, 5, 16, 0, 699.00, '2026-03-15 03:00:00.000000 +00:00', '2026-03-15 10:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7e4f-8adc-5969b5095145', 'Cloud Native Go', 60, 10, 40, 5, 999.00, '2026-03-16 02:00:00.000000 +00:00', '2026-03-16 11:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7e69-8a59-f4c8b73b9534', 'Go Security Best Practices', 30, 5, 12, 0, 599.00, '2026-03-17 02:00:00.000000 +00:00', '2026-03-17 09:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7e89-a83b-2972808abed1', 'Building CLI with Go', 25, 5, 8, 0, 499.00, '2026-03-18 03:00:00.000000 +00:00', '2026-03-18 08:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7ea5-a096-d6c7984dde9c', 'Go WebSocket Workshop', 35, 5, 15, 2, 699.00, '2026-03-19 02:00:00.000000 +00:00', '2026-03-19 09:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
INSERT INTO public.events (id, name, max_seats, waitlist_limit, booked_count, waitlist_count, price, start_time, end_time, created_at, updated_at) VALUES ('019d1a7c-67f6-7ec0-b8dd-ba6634d1e1d5', 'Scaling Go Systems', 50, 10, 38, 4, 899.00, '2026-03-20 02:00:00.000000 +00:00', '2026-03-20 11:00:00.000000 +00:00', '2026-03-23 11:37:33.939518 +00:00', '2026-03-23 11:37:33.939518 +00:00');
