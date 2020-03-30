select *
from get_next_proxy_for_check;
-- ACCESS (Subquery Scan)		100	100	212.06	2.056	210.81	2.042

select *
from get_next_proxy_for_check
limit 1;
-- ACCESS (Subquery Scan)		100	1	212.06	2.132	210.81	2.132

select * from get_next_proxy_for_check(60, 100);
-- TABLE_FUNCTION (Function Scan)		1000	1	10.25	1.753	0.25	1.752

select *
from get_next_proxy_for_check_sql(60, 100);
-- TABLE_FUNCTION (Function Scan)		1000	1	10.25	1.715	0.25	1.714

