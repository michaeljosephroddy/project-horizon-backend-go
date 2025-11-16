with first_query as
    (select date(created_at) as date,
            avg(mood_rating) as daily_avg_rating
     from mood_log
     where user_id = ?
         and date(created_at) between ? and ? group  by date(created_at)),
     second_query as
    (select avg(daily_avg_rating) as period_mood_rating_avg
     from first_query)
select period_mood_rating_avg
from second_query;
