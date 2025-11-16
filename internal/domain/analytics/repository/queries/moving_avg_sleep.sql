with first_query as
    (select sleep_date as DATE,
            avg(hours_slept) as avg_sleep_hours
     from sleep_log
     where user_id = ?
         and sleep_date between ? and ?),
     second_query as
    (select DATE, avg(avg_sleep_hours) over(
                                            order by DATE rows between %s preceding and current row) as moving_avg
     from first_query)
select *
from second_query
order  by DATE;
