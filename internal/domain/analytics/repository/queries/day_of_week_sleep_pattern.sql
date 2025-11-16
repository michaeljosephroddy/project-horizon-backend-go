select dayname(sleep_date) as day_of_week,
       dayofweek(sleep_date) as day_number,
       avg(hours_slept) as avg_sleep_hours,
       count(*) as total_entries
from sleep_log
where user_id = ?
    and sleep_date between ? and ?
group by day_of_week,
         day_number
order by day_number;
