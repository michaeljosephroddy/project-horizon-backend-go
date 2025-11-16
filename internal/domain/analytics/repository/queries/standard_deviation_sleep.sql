select stddev_pop(hours_slept) as std_dev
from sleep_log
where user_id = ?
    and sleep_date between ? and ?;
