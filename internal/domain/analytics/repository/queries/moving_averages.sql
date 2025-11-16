with first_query as
    (select date(created_at) as DATE,
            avg(mood_rating) as daily_avg
     from mood_log
     where user_id = ?
         and date(created_at) between ? and ? group  by date(created_at)),
     second_query as
    (select DATE, avg(daily_avg) over(
                                      order by DATE rows between %s preceding and current row) as moving_avg
     from first_query)
select *
from second_query
order  by DATE;
