select stddev_pop(mood_rating) as std_dev
from mood_log
where user_id = ?
    and date(created_at) between ? and ?;
