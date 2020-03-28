CREATE OR REPLACE FUNCTION increment(i integer) RETURNS integer AS
$$
BEGIN
    select proxy_id
    into i
    from proxy;
    RETURN i;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_next_proxy_for_check(min_interval_minutes int, max_limit int)
    RETURNS TABLE
            (
                p_id   int,
                p_ip   inet,
                p_port int
            )
AS
$$
BEGIN
    RETURN QUERY
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
END;
$$ LANGUAGE 'plpgsql';

select *
from get_next_proxy_for_check(1, 100);
select *
from get_next_proxy_for_check
limit 100

select increment(1);


CREATE FUNCTION sum_n_product_with_tab(x int)
    RETURNS TABLE
            (
                p_id   int,
                p_ip   inet,
                p_port int
            )
AS
$$
SELECT proxy_id, proxy_ip, proxy_port
FROM proxy;
$$ LANGUAGE SQL;
select * from sum_n_product_with_tab(1)