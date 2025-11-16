select avg(hours_slept) as avg_sleep_hours
from sleep_log
where user_id = ?
    and sleep_date between ? and ?;
