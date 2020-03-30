CREATE OR REPLACE FUNCTION get_next_proxy_for_check(min_interval_minutes int, max_limit int)
    RETURNS TABLE
            (
                proxy_id   int,
                proxy_ip   inet,
                proxy_port int
            )
AS
$$
select proxy.proxy_id,
       proxy_ip,
       proxy_port
from proxy
         left join (
    select max(created_at) as last_check, proxy_id
    from stat
    group by proxy_id
) t2 on proxy.proxy_id = t2.proxy_id
where now() - coalesce(last_check, '2000.01.01'::timestamp) > make_interval(mins => min_interval_minutes)
order by last_check is not null, last_check
limit max_limit;
$$ LANGUAGE sql;

-- На всякий случай аналог этой функции в виде view
create or replace view get_next_proxy_for_check as
select proxy.proxy_id,
       proxy_ip,
       proxy_port,
       last_check
from proxy
         left join (
    select max(created_at) as last_check, proxy_id
    from stat
    group by proxy_id
) t2 on proxy.proxy_id = t2.proxy_id
where now() - coalesce(last_check, '2000.01.01'::timestamp) > make_interval(mins => 60)
order by last_check is not null, last_check
limit 100;
---- create above / drop below ----

drop function get_next_proxy_for_check;
drop view get_next_proxy_for_check;
